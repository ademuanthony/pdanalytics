package mempool

import (
	"context"
	"time"

	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/rpcclient"
)

type Mempool struct {
	Time                 time.Time `json:"time"`
	FirstSeenTime        time.Time `json:"first_seen_time"`
	NumberOfTransactions int       `json:"number_of_transactions"`
	Voters               int       `json:"voters"`
	Tickets              int       `json:"tickets"`
	Revocations          int       `json:"revocations"`
	Size                 int32     `json:"size"`
	TotalFee             float64   `json:"total_fee"`
	Total                float64   `json:"total"`
}

type DataStore interface {
	CreateTables(ctx context.Context) error
	DropTables() error
	ClearCache() error
	MempoolTableName() string
	StoreMempool(context.Context, Mempool) error
	UpdateMempoolAggregateData(ctx context.Context) error
	LastMempoolTime() (entryTime time.Time, err error)
	LastMempoolBlockHeight() (height int64, err error)
	MempoolCount(ctx context.Context) (int64, error)
	Mempools(ctx context.Context, offtset int, limit int) ([]Dto, error)
}

type Collector struct {
	ctx                context.Context
	collectionInterval float64
	dcrClient          *rpcclient.Client
	dataStore          DataStore
	activeChain        *chaincfg.Params
	syncIsDone         bool
	bestBlockHeight    uint32
}
