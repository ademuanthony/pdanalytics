package propagation

import (
	"context"
	"net/http"
	"time"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/chaincfg/v2"
	chainjson "github.com/decred/dcrd/rpc/jsonrpc/types/v2"
	"github.com/decred/dcrd/rpcclient/v5"
	"github.com/decred/dcrd/wire"
	"github.com/planetdecred/dcrextdata/app/helpers"
	"github.com/planetdecred/pdanalytics/web"
)

func New(ctx context.Context, dcrdClient *rpcclient.Client, dataStore store,
	webServer *web.Server, viewFolder string,
	params *chaincfg.Params) (*propagation, error) {

	if viewFolder == "" {
		viewFolder = "./pkgs/propagation/views"
	}

	ac := &propagation{
		dataStore:   dataStore,
		templates:   web.NewTemplates(viewFolder, false, []string{"extras"}, web.MakeTemplateFuncMap(params)),
		webServer:   webServer,
		activeChain: params,
		dcrClient:   dcrdClient,
	}

	return ac, nil
}

func (c *propagation) SetExplorerBestBlock(ctx context.Context) error {
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

func (c *propagation) ConnectBlock(blockHeader *wire.BlockHeader) error {
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

func (c *propagation) TxReceived(txDetails *chainjson.TxRawResult) error {
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
