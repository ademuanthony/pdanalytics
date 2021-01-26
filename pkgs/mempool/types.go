// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mempool

import (
	"context"
	"time"

	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/rpcclient"
	"github.com/planetdecred/dcrextdata/datasync"
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
	MempoolTableName() string
	BlockTableName() string
	VoteTableName() string
	StoreMempool(context.Context, Mempool) error
	LastMempoolTime() (entryTime time.Time, err error)
	FetchMempoolForSync(ctx context.Context, date time.Time, offtset int, limit int) ([]Mempool, int64, error)
	SaveBlock(context.Context, Block) error
	UpdateBlockBinData(context.Context) error
	FetchBlockForSync(ctx context.Context, blockHeight int64, offtset int, limit int) ([]Block, int64, error)
	SaveVote(ctx context.Context, vote Vote) error
	UpdateVoteTimeDeviationData(context.Context) error
	FetchVoteForSync(ctx context.Context, date time.Time, offtset int, limit int) ([]Vote, int64, error)

	datasync.Store
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
