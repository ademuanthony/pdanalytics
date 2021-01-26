// TODO:
// - Create the handlers
// Create the views
// Adapts the layout
// Adapt the js controllers

package mempool

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/chaincfg/v2"
	dcrjson "github.com/decred/dcrd/rpc/jsonrpc/types/v2"
	"github.com/decred/dcrd/rpcclient/v5"
	"github.com/planetdecred/dcrextdata/app/helpers"
	"github.com/planetdecred/pdanalytics/web"
)

const (
	chartViewOption             = "chart"
	defaultViewOption           = chartViewOption
	mempoolDefaultChartDataType = "size"
	maxPageSize                 = 250
	defaultPageSize             = 20
	noDataMessage               = "does not have data for the selected query option(s)."
)

var (
	pageSizeSelector = map[int]int{
		20:  20,
		30:  30,
		50:  50,
		100: 100,
		150: 150,
	}
)

func NewCollector(ctx context.Context, client *rpcclient.Client, interval float64,
	activeChain *chaincfg.Params, dataStore DataStore, webServer *web.Server, viewFolder string) (*Collector, error) {

	if viewFolder == "" {
		viewFolder = "./pkgs/mempool/views"
	}

	c := &Collector{
		ctx:                ctx,
		templates:          web.NewTemplates(viewFolder, false, []string{"extras"}, web.MakeTemplateFuncMap(activeChain)),
		webServer:          webServer,
		dcrClient:          client,
		collectionInterval: interval,
		dataStore:          dataStore,
		activeChain:        activeChain,
	}

	tmpls := []string{"mempool"}

	for _, name := range tmpls {
		if err := c.templates.AddTemplate(name); err != nil {
			log.Errorf("Unable to create new html template: %v", err)
			return nil, err
		}
	}

	c.webServer.AddMenuItem(web.MenuItem{
		Href:      "/mempool",
		HyperText: "Mempool",
		Attributes: map[string]string{
			"class": "menu-item",
			"title": "Historic mempool data",
		},
	})

	// Development subsidy address of the current network
	devSubsidyAddress, err := web.DevSubsidyAddress(activeChain)
	if err != nil {
		log.Warnf("mempool.New: %v", err)
		return nil, err
	}
	log.Debugf("Organization address: %s", devSubsidyAddress)

	c.pageData = &web.PageData{
		BlockInfo: new(web.BlockInfo),
		HomeInfo: &web.HomeInfo{
			DevAddress: devSubsidyAddress,
			Params: web.ChainParams{
				WindowSize:       c.activeChain.StakeDiffWindowSize,
				RewardWindowSize: c.activeChain.SubsidyReductionInterval,
				BlockTime:        c.activeChain.TargetTimePerBlock.Nanoseconds(),
				MeanVotingBlocks: c.MeanVotingBlocks,
			},
			PoolInfo: web.TicketPoolInfo{
				Target: uint32(c.activeChain.TicketPoolSize * c.activeChain.TicketsPerBlock),
			},
		},
	}

	webServer.AddRoute("/mempool", web.GET, c.mempoolPage)
	webServer.AddRoute("/getmempool", web.GET, c.getMempool)

	return c, nil
}

func (c *Collector) SetExplorerBestBlock(ctx context.Context) error {
	var explorerUrl string
	switch c.activeChain.Name {
	case chaincfg.MainNetParams().Name:
		explorerUrl = "https://explorer.dcrdata.org/api/block/best"
	case chaincfg.TestNet3Params().Name:
		explorerUrl = "https://testnet.dcrdata.org/api/block/best"
	}

	var bestBlock = struct {
		Height uint32 `json:"height"`
	}{}

	err := helpers.GetResponse(ctx, &http.Client{Timeout: 10 * time.Second}, explorerUrl, &bestBlock)
	if err != nil {
		return err
	}

	log.Infof("Current best block height: %d", bestBlock.Height)
	c.bestBlockHeight = bestBlock.Height
	return nil
}

func (c *Collector) StartMonitoring(ctx context.Context) {
	var mu sync.Mutex

	collectMempool := func() {
		if !c.syncIsDone {
			return
		}

		mu.Lock()
		defer mu.Unlock()

		mempoolTransactionMap, err := c.dcrClient.GetRawMempoolVerbose(dcrjson.GRMAll)
		if err != nil {
			log.Error(err)
			return
		}

		if len(mempoolTransactionMap) == 0 {
			return
		}

		mempoolDto := Mempool{
			NumberOfTransactions: len(mempoolTransactionMap),
			Time:                 helpers.NowUTC(),
			FirstSeenTime:        helpers.NowUTC(), //todo: use the time of the first tx in the mempool
		}

		for hashString, tx := range mempoolTransactionMap {
			hash, err := chainhash.NewHashFromStr(hashString)
			if err != nil {
				log.Error(err)
				continue
			}
			rawTx, err := c.dcrClient.GetRawTransactionVerbose(hash)
			if err != nil {
				log.Error(err)
				continue
			}

			totalOut := 0.0
			for _, v := range rawTx.Vout {
				totalOut += v.Value
			}

			mempoolDto.Total += totalOut
			mempoolDto.TotalFee += tx.Fee
			mempoolDto.Size += tx.Size
			if mempoolDto.FirstSeenTime.Unix() > tx.Time {
				mempoolDto.FirstSeenTime = helpers.UnixTime(tx.Time)
			}

		}

		votes, err := c.dcrClient.GetRawMempool(dcrjson.GRMVotes)
		if err != nil {
			log.Error(err)
			return
		}
		mempoolDto.Voters = len(votes)

		tickets, err := c.dcrClient.GetRawMempool(dcrjson.GRMTickets)
		if err != nil {
			log.Error(err)
			return
		}
		mempoolDto.Tickets = len(tickets)

		revocations, err := c.dcrClient.GetRawMempool(dcrjson.GRMRevocations)
		if err != nil {
			log.Error(err)
			return
		}
		mempoolDto.Revocations = len(revocations)

		if err = c.dataStore.StoreMempool(ctx, mempoolDto); err != nil {
			log.Error(err)
		}
	}

	lastMempoolTime, err := c.dataStore.LastMempoolTime()
	if err != nil {
		if err != sql.ErrNoRows {
			log.Errorf("Unable to get last mempool entry time: %s", err.Error())
		}
	} else {
		sencodsPassed := math.Abs(time.Since(lastMempoolTime).Seconds())
		if sencodsPassed < c.collectionInterval {
			timeLeft := c.collectionInterval - sencodsPassed
			log.Infof("Fetching mempool every %dm, collected %0.2f ago, will fetch in %0.2f.", 1, sencodsPassed,
				timeLeft)
			time.Sleep(time.Duration(timeLeft) * time.Second)
		}
	}
	collectMempool()
	ticker := time.NewTicker(time.Duration(c.collectionInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			collectMempool()
			break
		case <-ctx.Done():
			return
		}
	}
}

func (c *Collector) mempoolPage(w http.ResponseWriter, r *http.Request) {
	mempoolData, err := c.fetchMempoolData(r)
	if err != nil {
		web.RenderErrorfJSON(w, err.Error())
		return
	}

	str, err := c.templates.ExecTemplateToString("mempool", struct {
		*web.CommonPageData
		Mempool   map[string]interface{}
		BlockTime float64
	}{
		CommonPageData: c.commonData(r),
		Mempool:        mempoolData,
		BlockTime:      c.activeChain.MinDiffReductionTime.Seconds(),
	})

	if err != nil {
		log.Errorf("Template execute failure: %v", err)
		c.StatusPage(w, r, web.DefaultErrorCode, web.DefaultErrorMessage, "", web.ExpStatusError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if _, err = io.WriteString(w, str); err != nil {
		log.Error(err)
	}
}

// /getmempool
func (s *Collector) getMempool(res http.ResponseWriter, req *http.Request) {
	data, err := s.fetchMempoolData(req)

	if err != nil {
		web.RenderErrorfJSON(res, err.Error())
		return
	}

	web.RenderJSON(res, data)
}

func (s *Collector) fetchMempoolData(req *http.Request) (map[string]interface{}, error) {
	req.ParseForm()
	page := req.FormValue("page")
	numberOfRows := req.FormValue("records-per-page")
	viewOption := req.FormValue("view-option")
	chartDataType := req.FormValue("chart-data-type")

	if chartDataType == "" {
		chartDataType = mempoolDefaultChartDataType
	}

	if viewOption == "" {
		viewOption = defaultViewOption
	}

	var pageSize int
	numRows, err := strconv.Atoi(numberOfRows)
	switch {
	case err != nil || numRows <= 0:
		pageSize = defaultPageSize
	case numRows > maxPageSize:
		pageSize = maxPageSize
	default:
		pageSize = numRows
	}

	pageToLoad, err := strconv.Atoi(page)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	offset := (pageToLoad - 1) * pageSize

	data := map[string]interface{}{
		"chartView":            true,
		"chartDataType":        chartDataType,
		"selectedViewOption":   viewOption,
		"pageSizeSelector":     pageSizeSelector,
		"selectedNumberOfRows": pageSize,
		"currentPage":          pageToLoad,
		"previousPage":         pageToLoad - 1,
		"totalPages":           0,
	}

	if viewOption == defaultViewOption {
		return data, nil
	}

	ctx := req.Context()

	mempoolSlice, err := s.dataStore.Mempools(ctx, offset, pageSize)
	if err != nil {
		return nil, err
	}

	totalCount, err := s.dataStore.MempoolCount(ctx)
	if err != nil {
		return nil, err
	}

	if len(mempoolSlice) == 0 {
		data["message"] = fmt.Sprintf("Mempool %s", noDataMessage)
		return data, nil
	}

	data["mempoolData"] = mempoolSlice
	data["totalPages"] = int(math.Ceil(float64(totalCount) / float64(pageSize)))

	totalTxLoaded := offset + len(mempoolSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = pageToLoad + 1
	}

	return data, nil
}

// commonData grabs the common page data that is available to every page.
// This is particularly useful for extras.tmpl, parts of which
// are used on every page
func (c *Collector) commonData(r *http.Request) *web.CommonPageData {

	darkMode, err := r.Cookie(web.DarkModeCoookie)
	if err != nil && err != http.ErrNoCookie {
		log.Errorf("Cookie dcrdataDarkBG retrieval error: %v", err)
	}
	return &web.CommonPageData{
		Version:       c.Version,
		ChainParams:   c.activeChain,
		BlockTimeUnix: int64(c.activeChain.TargetTimePerBlock.Seconds()),
		DevAddress:    c.pageData.HomeInfo.DevAddress,
		NetName:       c.NetName,
		Links:         web.ExplorerLinks,
		Cookies: web.Cookies{
			DarkMode: darkMode != nil && darkMode.Value == "1",
		},
		RequestURI: r.URL.RequestURI(),
		MenuItems:  c.webServer.MenuItems,
	}
}

// StatusPage provides a page for displaying status messages and exception
// handling without redirecting. Be sure to return after calling StatusPage if
// this completes the processing of the calling http handler.
func (c *Collector) StatusPage(w http.ResponseWriter, r *http.Request, code, message, additionalInfo string, sType web.ExpStatus) {
	commonPageData := c.commonData(r)
	if commonPageData == nil {
		// exp.blockData.GetTip likely failed due to empty DB.
		http.Error(w, "The database is initializing. Try again later.",
			http.StatusServiceUnavailable)
		return
	}
	str, err := c.templates.Exec("status", struct {
		*web.CommonPageData
		StatusType     web.ExpStatus
		Code           string
		Message        string
		AdditionalInfo string
	}{
		CommonPageData: commonPageData,
		StatusType:     sType,
		Code:           code,
		Message:        message,
		AdditionalInfo: additionalInfo,
	})
	if err != nil {
		log.Errorf("Template execute failure: %v", err)
		str = "Something went very wrong if you can see this, try refreshing"
	}

	w.Header().Set("Content-Type", "text/html")
	switch sType {
	case web.ExpStatusDBTimeout:
		w.WriteHeader(http.StatusServiceUnavailable)
	case web.ExpStatusNotFound:
		w.WriteHeader(http.StatusNotFound)
	case web.ExpStatusFutureBlock:
		w.WriteHeader(http.StatusOK)
	case web.ExpStatusError:
		w.WriteHeader(http.StatusInternalServerError)
	// When blockchain sync is running, status 202 is used to imply that the
	// other requests apart from serving the status sync page have been received
	// and accepted but cannot be processed now till the sync is complete.
	case web.ExpStatusSyncing:
		w.WriteHeader(http.StatusAccepted)
	case web.ExpStatusNotSupported:
		w.WriteHeader(http.StatusUnprocessableEntity)
	case web.ExpStatusBadRequest:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	io.WriteString(w, str)
}
