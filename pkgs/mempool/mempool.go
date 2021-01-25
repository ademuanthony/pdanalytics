// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mempool

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrjson"
	chainjson "github.com/decred/dcrd/rpc/jsonrpc/types/v2"
	"github.com/decred/dcrd/rpcclient"
	"github.com/decred/dcrd/wire"
	"github.com/planetdecred/dcrextdata/app/helpers"
	"github.com/planetdecred/dcrextdata/datasync"
)

func NewCollector(ctx context.Context, client *rpcclient.Client, interval float64, activeChain *chaincfg.Params, dataStore DataStore) *Collector {
	c := &Collector{
		ctx:                ctx,
		dcrClient:          client,
		collectionInterval: interval,
		dataStore:          dataStore,
		activeChain:        activeChain,
		ticketInds:         make(BlockValidatorIndex),
	}
	return c
}

func (c *Collector) SetExplorerBestBlock(ctx context.Context) error {
	var explorerUrl string
	switch c.activeChain.Name {
	case chaincfg.MainNetParams.Name:
		explorerUrl = "https://explorer.dcrdata.org/api/block/best"
	case chaincfg.TestNet3Params.Name:
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

func (c *Collector) ConnectBlock(blockHeader *wire.BlockHeader) error {
	if blockHeader.Height > c.bestBlockHeight {
		c.syncIsDone = true
	}

	if !c.syncIsDone {
		log.Infof("Received a stale block height %d, block dropped", blockHeader.Height)
		return nil
	}

	block := Block{
		BlockInternalTime: blockHeader.Timestamp.UTC(),
		BlockReceiveTime:  helpers.NowUTC(),
		BlockHash:         blockHeader.BlockHash().String(),
		BlockHeight:       blockHeader.Height,
	}
	if err := c.dataStore.SaveBlock(c.ctx, block); err != nil {
		log.Error(err)
		return err
	}
	if err := c.dataStore.UpdateBlockBinData(c.ctx); err != nil {
		log.Errorf("Error in block bin data update, %s", err.Error())
		return err
	}
	return nil
}

func (c *Collector) TxReceived(txDetails *chainjson.TxRawResult) error {
	if !c.syncIsDone {
		return nil
	}
	receiveTime := helpers.NowUTC()

	msgTx, err := MsgTxFromHex(txDetails.Hex)
	if err != nil {
		log.Errorf("Failed to decode transaction hex: %v", err)
		return err
	}

	if txType := DetermineTxTypeString(msgTx); txType != "Vote" {
		return nil
	}

	var voteInfo *VoteInfo
	validation, version, bits, choices, err := SSGenVoteChoices(msgTx, c.activeChain)
	if err != nil {
		log.Errorf("Error in getting vote choice: %s", err.Error())
		return err
	}

	voteInfo = &VoteInfo{
		Validation: BlockValidation{
			Hash:     validation.Hash,
			Height:   validation.Height,
			Validity: validation.Validity,
		},
		Version:     version,
		Bits:        bits,
		Choices:     choices,
		TicketSpent: msgTx.TxIn[1].PreviousOutPoint.Hash.String(),
	}

	c.ticketIndsMutex.Lock()
	voteInfo.SetTicketIndex(c.ticketInds)
	c.ticketIndsMutex.Unlock()

	vote := Vote{
		ReceiveTime: receiveTime,
		VotingOn:    validation.Height,
		Hash:        txDetails.Txid,
		ValidatorId: voteInfo.MempoolTicketIndex,
	}

	if voteInfo.Validation.Validity {
		vote.Validity = "Valid"
	} else {
		vote.Validity = "Invalid"
	}

	var retries = 3
	var targetedBlock *wire.MsgBlock

	// try to get the block from the blockchain until the number of retries has elapsed
	for i := 0; i <= retries; i++ {
		hash, _ := chainhash.NewHashFromStr(validation.Hash)
		targetedBlock, err = c.dcrClient.GetBlock(hash)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	// err is ignored since the vote will be updated when the block becomes available
	if targetedBlock != nil {
		vote.TargetedBlockTime = targetedBlock.Header.Timestamp.UTC()
		vote.BlockHash = targetedBlock.Header.BlockHash().String()
	}

	if err = c.dataStore.SaveVote(c.ctx, vote); err != nil {
		log.Error(err)
	}

	if err = c.dataStore.UpdateVoteTimeDeviationData(c.ctx); err != nil {
		log.Errorf("Error in vote receive time deviation data update, %s", err.Error())
	}
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

func (c *Collector) RegisterSyncer(syncCoordinator *datasync.SyncCoordinator) {
	c.registerBlockSyncer(syncCoordinator)
	c.registerVoteSyncer(syncCoordinator)
}

func (c *Collector) registerBlockSyncer(syncCoordinator *datasync.SyncCoordinator) {
	syncCoordinator.AddSyncer(c.dataStore.BlockTableName(), datasync.Syncer{
		LastEntry: func(ctx context.Context, db datasync.Store) (string, error) {
			var lastHeight int64
			err := db.LastEntry(ctx, c.dataStore.BlockTableName(), &lastHeight)
			if err != nil && err != sql.ErrNoRows {
				return "0", fmt.Errorf("error in fetching last block height, %s", err.Error())
			}
			return strconv.FormatInt(lastHeight, 10), nil
		},
		Collect: func(ctx context.Context, url string) (result *datasync.Result, err error) {
			result = new(datasync.Result)
			result.Records = []Block{}
			err = helpers.GetResponse(ctx, &http.Client{Timeout: 10 * time.Second}, url, result)
			return
		},
		Retrieve: func(ctx context.Context, last string, skip, take int) (result *datasync.Result, err error) {
			blockHeight, _ := strconv.ParseInt(last, 10, 64)
			result = new(datasync.Result)
			blocks, totalCount, err := c.dataStore.FetchBlockForSync(ctx, blockHeight, skip, take)
			if err != nil {
				result.Message = err.Error()
				return
			}
			result.Records = blocks
			result.TotalCount = totalCount
			result.Success = true
			return
		},
		Append: func(ctx context.Context, store datasync.Store, data interface{}) {
			mappedData := data.([]interface{})
			var blocks []Block
			for _, item := range mappedData {
				var block Block
				err := datasync.DecodeSyncObj(item, &block)
				if err != nil {
					log.Errorf("Error in decoding the received block data, %s", err.Error())
					return
				}
				blocks = append(blocks, block)
			}

			for _, block := range blocks {
				err := store.SaveBlockFromSync(ctx, block)
				if err != nil {
					log.Errorf("Error while appending block synced data, %s", err.Error())
				}
			}
			// update propagation data
			if err := store.UpdatePropagationData(ctx); err != nil {
				log.Errorf("Error in initial propagation data update, %s", err.Error())
			}
		},
	})
}

func (c *Collector) registerVoteSyncer(syncCoordinator *datasync.SyncCoordinator) {
	syncCoordinator.AddSyncer(c.dataStore.VoteTableName(), datasync.Syncer{
		LastEntry: func(ctx context.Context, db datasync.Store) (string, error) {
			var receiveTime time.Time
			err := db.LastEntry(ctx, c.dataStore.VoteTableName(), &receiveTime)
			if err != nil && err != sql.ErrNoRows {
				return "0", fmt.Errorf("error in fetching last vote receive time, %s", err.Error())
			}
			return strconv.FormatInt(receiveTime.Unix(), 10), nil
		},
		Collect: func(ctx context.Context, url string) (result *datasync.Result, err error) {
			result = new(datasync.Result)
			result.Records = []Vote{}
			err = helpers.GetResponse(ctx, &http.Client{Timeout: 10 * time.Second}, url, result)
			return
		},
		Retrieve: func(ctx context.Context, last string, skip, take int) (result *datasync.Result, err error) {
			unixDate, _ := strconv.ParseInt(last, 10, 64)
			result = new(datasync.Result)
			votes, totalCount, err := c.dataStore.FetchVoteForSync(ctx, helpers.UnixTime(unixDate), skip, take)
			if err != nil {
				result.Message = err.Error()
				return
			}
			fmt.Println("Total count", totalCount)
			result.Records = votes
			result.TotalCount = totalCount
			result.Success = true
			return
		},
		Append: func(ctx context.Context, store datasync.Store, data interface{}) { //todo: should return an error
			mappedData := data.([]interface{})
			var votes []Vote
			for _, item := range mappedData {
				var vote Vote
				err := datasync.DecodeSyncObj(item, &vote)
				if err != nil {
					log.Errorf("Error in decoding the received vote data, %s", err.Error())
					return
				}
				votes = append(votes, vote)
			}

			for _, vote := range votes {
				err := store.SaveVoteFromSync(ctx, vote)
				if err != nil {
					log.Errorf("Error while appending vote synced data, %s", err.Error())
				}
			}
		},
	})
}