package propagation

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/planetdecred/dcrextdata/datasync"
	"github.com/planetdecred/pdanalytics/web"
)

const (
	chartViewOption   = "chart"
	defaultViewOption = chartViewOption
	maxPageSize       = 250
	defaultPageSize   = 20
	noDataMessage     = "does not have data for the selected query option(s)."
)

var (
	pageSizeSelector = map[int]int{
		20:  20,
		30:  30,
		50:  50,
		100: 100,
		150: 150,
	}

	propagationRecordSet = map[string]string{
		"blocks": "Blocks",
		"votes":  "Votes",
	}
)

// /propagation
func (c *propagation) propagationPage(w http.ResponseWriter, r *http.Request) {

	block, err := c.fetchPropagationData(r)
	if err != nil {
		log.Error(err)
		c.StatusPage(w, r, web.DefaultErrorCode, web.DefaultErrorMessage, "", web.ExpStatusError)
		return
	}

	str, err := c.templates.ExecTemplateToString("propagation", struct {
		*web.CommonPageData
		Propagation map[string]interface{}
		BlockTime   float64
	}{
		CommonPageData: c.commonData(r),
		Propagation:    block,
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

// /getPropagationData
func (c *propagation) getPropagationData(w http.ResponseWriter, r *http.Request) {
	data, err := c.fetchPropagationData(r)
	if err != nil {
		web.RenderErrorfJSON(w, err.Error())
		return
	}
	web.RenderJSON(w, data)
}

func (c *propagation) fetchPropagationData(req *http.Request) (map[string]interface{}, error) {
	req.ParseForm()
	page := req.FormValue("page")
	numberOfRows := req.FormValue("records-per-page")
	viewOption := req.FormValue("view-option")
	recordSet := req.FormValue("record-set")
	chartType := req.FormValue("chart-type")

	if viewOption == "" {
		viewOption = "chart"
	}

	if recordSet == "" {
		recordSet = "both"
	}

	if chartType == "" {
		chartType = "block-propagation"
	}

	var pageSize int
	numRows, err := strconv.Atoi(numberOfRows)
	if err != nil || numRows <= 0 {
		pageSize = defaultPageSize
	} else if numRows > maxPageSize {
		pageSize = maxPageSize
	} else {
		pageSize = numRows
	}

	pageToLoad, err := strconv.Atoi(page)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	offset := (pageToLoad - 1) * pageSize

	ctx := req.Context()

	syncSources, _ := datasync.RegisteredSources()

	data := map[string]interface{}{
		"chartView":            viewOption == "chart",
		"selectedViewOption":   viewOption,
		"chartType":            chartType,
		"currentPage":          pageToLoad,
		"propagationRecordSet": propagationRecordSet,
		"pageSizeSelector":     pageSizeSelector,
		"selectedRecordSet":    recordSet,
		"both":                 true,
		"selectedNum":          pageSize,
		"url":                  "/propagation",
		"previousPage":         pageToLoad - 1,
		"totalPages":           0,
		"syncSources":          strings.Join(syncSources, "|"),
	}

	if viewOption == defaultViewOption {
		return data, nil
	}

	blockSlice, err := c.dataStore.BlocksWithoutVotes(ctx, offset, pageSize)
	if err != nil {
		return nil, err
	}

	for i := 0; i <= 1 && i <= len(blockSlice)-1; i++ {
		votes, err := c.dataStore.VotesByBlock(ctx, blockSlice[i].BlockHash)
		if err != nil {
			return nil, err
		}
		blockSlice[i].Votes = votes
	}

	totalCount, err := c.dataStore.BlockCount(ctx)
	if err != nil {
		return nil, err
	}

	if len(blockSlice) == 0 {
		data["message"] = fmt.Sprintf("%s %s", recordSet, noDataMessage)
		return data, nil
	}

	data["records"] = blockSlice
	data["totalPages"] = int(math.Ceil(float64(totalCount) / float64(pageSize)))

	totalTxLoaded := offset + len(blockSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = pageToLoad + 1
	}

	return data, nil
}

// /getblocks
func (c *propagation) getBlocks(res http.ResponseWriter, req *http.Request) {
	data, err := c.fetchBlockData(req)
	if err != nil {
		web.RenderErrorfJSON(res, err.Error())
		return
	}

	web.RenderJSON(res, data)
}

func (c *propagation) fetchBlockData(req *http.Request) (map[string]interface{}, error) {
	req.ParseForm()
	page := req.FormValue("page")
	numberOfRows := req.FormValue("records-per-page")
	viewOption := req.FormValue("view-option")

	if viewOption == "" {
		viewOption = defaultViewOption
	}

	var pageSize int
	numRows, err := strconv.Atoi(numberOfRows)
	if err != nil || numRows <= 0 {
		pageSize = defaultPageSize
	} else if numRows > maxPageSize {
		pageSize = maxPageSize
	} else {
		pageSize = numRows
	}

	pageToLoad, err := strconv.Atoi(page)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	offset := (pageToLoad - 1) * pageSize

	ctx := req.Context()

	data := map[string]interface{}{
		"chartView":            true,
		"selectedViewOption":   defaultViewOption,
		"currentPage":          pageToLoad,
		"propagationRecordSet": propagationRecordSet,
		"pageSizeSelector":     pageSizeSelector,
		"selectedFilter":       "blocks",
		"blocks":               true,
		"url":                  "/blockdata",
		"selectedNum":          pageSize,
		"previousPage":         pageToLoad - 1,
		"totalPages":           pageToLoad,
	}

	if viewOption == defaultViewOption {
		return data, nil
	}

	blocksSlice, err := c.dataStore.BlocksWithoutVotes(ctx, offset, pageSize)
	if err != nil {
		return nil, err
	}

	if len(blocksSlice) == 0 {
		data["message"] = fmt.Sprintf("Blocks %s", noDataMessage)
		return data, nil
	}

	totalCount, err := c.dataStore.BlockCount(ctx)
	if err != nil {
		return nil, err
	}

	data["records"] = blocksSlice
	data["totalPages"] = int(math.Ceil(float64(totalCount) / float64(pageSize)))

	totalTxLoaded := offset + len(blocksSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = pageToLoad + 1
	}

	return data, nil
}

// /getvotes
func (c *propagation) getVotes(res http.ResponseWriter, req *http.Request) {
	data, err := c.fetchVoteData(req)

	if err != nil {
		web.RenderErrorfJSON(res, err.Error())
		return
	}
	web.RenderJSON(res, data)
}

func (c *propagation) fetchVoteData(req *http.Request) (map[string]interface{}, error) {
	req.ParseForm()
	page := req.FormValue("page")
	numberOfRows := req.FormValue("records-per-page")
	viewOption := req.FormValue("view-option")

	if viewOption == "" {
		viewOption = defaultViewOption
	}

	var pageSize int
	numRows, err := strconv.Atoi(numberOfRows)
	if err != nil || numRows <= 0 {
		pageSize = defaultPageSize
	} else if numRows > maxPageSize {
		pageSize = maxPageSize
	} else {
		pageSize = numRows
	}

	pageToLoad, err := strconv.Atoi(page)
	if err != nil || pageToLoad <= 0 {
		pageToLoad = 1
	}

	offset := (pageToLoad - 1) * pageSize

	ctx := req.Context()

	data := map[string]interface{}{
		"chartView":            true,
		"selectedViewOption":   defaultViewOption,
		"currentPage":          pageToLoad,
		"propagationRecordSet": propagationRecordSet,
		"pageSizeSelector":     pageSizeSelector,
		"selectedFilter":       "votes",
		"votes":                true,
		"selectedNum":          pageSize,
		"url":                  "/votesdata",
		"previousPage":         pageToLoad - 1,
		"totalPages":           pageToLoad,
	}

	if viewOption == defaultViewOption {
		return data, nil
	}

	voteSlice, err := c.dataStore.Votes(ctx, offset, pageSize)
	if err != nil {
		return nil, err
	}

	if len(voteSlice) == 0 {
		data["message"] = fmt.Sprintf("Votes %s", noDataMessage)
		return data, nil
	}

	totalCount, err := c.dataStore.VotesCount(ctx)
	if err != nil {
		return nil, err
	}

	data["voteRecords"] = voteSlice
	data["totalPages"] = int(math.Ceil(float64(totalCount) / float64(pageSize)))

	totalTxLoaded := offset + len(voteSlice)
	if int64(totalTxLoaded) < totalCount {
		data["nextPage"] = pageToLoad + 1
	}

	return data, nil
}

// getVoteByBlock
func (c *propagation) getVoteByBlock(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	hash := req.FormValue("block_hash")
	votes, err := c.dataStore.VotesByBlock(req.Context(), hash)
	if err != nil {
		web.RenderErrorfJSON(res, err.Error())
		return
	}
	web.RenderJSON(res, votes)
}

// commonData grabs the common page data that is available to every page.
// This is particularly useful for extras.tmpl, parts of which
// are used on every page
func (c *propagation) commonData(r *http.Request) *web.CommonPageData {

	darkMode, err := r.Cookie(web.DarkModeCoookie)
	if err != nil && err != http.ErrNoCookie {
		log.Errorf("Cookie pdanalyticsDarkBG retrieval error: %v", err)
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
func (c *propagation) StatusPage(w http.ResponseWriter, r *http.Request, code, message, additionalInfo string, sType web.ExpStatus) {
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
