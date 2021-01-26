// TODO:
// - Create the handlers
// Create the views
// Adapts the layout
// Adapt the js controllers

package mempool

import (
	"context"
	"database/sql"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/chaincfg/v2"
	dcrjson "github.com/decred/dcrd/rpc/jsonrpc/types/v2"
	"github.com/decred/dcrd/rpcclient/v5"
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

	tmpls := []string{"mempool", "status"}

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
	webServer.AddRoute("/api/charts/mempool/{chartDataType}", web.GET, c.chart, chartDataTypeCtx)

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

	err := GetResponse(ctx, &http.Client{Timeout: 10 * time.Second}, explorerUrl, &bestBlock)
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
			Time:                 NowUTC(),
			FirstSeenTime:        NowUTC(), //todo: use the time of the first tx in the mempool
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
				mempoolDto.FirstSeenTime = UnixTime(tx.Time)
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
