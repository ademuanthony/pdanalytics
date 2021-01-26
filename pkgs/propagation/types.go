package propagation

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/decred/dcrd/blockchain/stake"
	"github.com/decred/dcrd/chaincfg/v2"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/rpcclient/v5"
	"github.com/decred/dcrd/wire"
	"github.com/planetdecred/dcrextdata/datasync"
	"github.com/planetdecred/pdanalytics/web"
)

type propagation struct {
	ctx             context.Context
	dcrClient       *rpcclient.Client
	dataStore       store
	ticketInds      BlockValidatorIndex
	bestBlockHeight uint32
	syncIsDone      bool
	ticketIndsMutex sync.Mutex

	activeChain      *chaincfg.Params
	Version          string
	NetName          string
	MeanVotingBlocks int64

	templates web.Templates
	webServer *web.Server
	pageData  *web.PageData
}

type store interface {
	BlockTableName() string
	VoteTableName() string
	SaveBlock(context.Context, Block) error
	UpdateBlockBinData(context.Context) error
	FetchBlockForSync(ctx context.Context, blockHeight int64, offtset int, limit int) ([]Block, int64, error)
	SaveVote(ctx context.Context, vote Vote) error
	UpdateVoteTimeDeviationData(context.Context) error
	FetchVoteForSync(ctx context.Context, date time.Time, offtset int, limit int) ([]Vote, int64, error)

	datasync.Store
}

// TicketIndex is used to assign an index to a ticket hash.
type TicketIndex map[string]int

// BlockValidatorIndex keeps a list of arbitrary indexes for unique combinations
// of block hash and the ticket being spent to validate the block, i.e.
// map[validatedBlockHash]map[ticketHash]index.
type BlockValidatorIndex map[string]TicketIndex

type Dto struct {
	Time                 string  `json:"time"`
	FirstSeenTime        string  `json:"first_seen_time"`
	NumberOfTransactions int     `json:"number_of_transactions"`
	Voters               int     `json:"voters"`
	Tickets              int     `json:"tickets"`
	Revocations          int     `json:"revocations"`
	Size                 int32   `json:"size"`
	TotalFee             float64 `json:"total_fee"`
	Total                float64 `json:"total"`
}

type Block struct {
	BlockReceiveTime  time.Time
	BlockInternalTime time.Time
	BlockHeight       uint32
	BlockHash         string
}

type BlockDto struct {
	BlockReceiveTime  string    `json:"block_receive_time"`
	BlockInternalTime string    `json:"block_internal_time"`
	Delay             string    `json:"delay"`
	BlockHeight       uint32    `json:"block_height"`
	BlockHash         string    `json:"block_hash"`
	Votes             []VoteDto `json:"votes"`
}

type PropagationChartData struct {
	BlockHeight    int64     `json:"block_height"`
	TimeDifference float64   `json:"time_difference"`
	BlockTime      time.Time `json:"block_time"`
}

type BlockReceiveTime struct {
	BlockHeight int64 `json:"block_height"`
	ReceiveTime time.Time
}

type Vote struct {
	Hash              string
	ReceiveTime       time.Time
	TargetedBlockTime time.Time
	BlockReceiveTime  time.Time
	VotingOn          int64
	BlockHash         string
	ValidatorId       int
	Validity          string
}

type VoteDto struct {
	Hash                  string `json:"hash"`
	ReceiveTime           string `json:"receive_time"`
	TargetedBlockTimeDiff string `json:"block_time_diff"`
	BlockReceiveTimeDiff  string `json:"block_receive_time_diff"`
	VotingOn              int64  `json:"voting_on"`
	BlockHash             string `json:"block_hash"`
	ShortBlockHash        string `json:"short_block_hash"`
	ValidatorId           int    `json:"validator_id"`
	Validity              string `json:"validity"`
}

// VoteInfo models data about a SSGen transaction (vote)
type VoteInfo struct {
	Validation         BlockValidation `json:"block_validation"`
	Version            uint32          `json:"vote_version"`
	Bits               uint16          `json:"vote_bits"`
	Choices            []*VoteChoice   `json:"vote_choices"`
	TicketSpent        string          `json:"ticket_spent"`
	MempoolTicketIndex int             `json:"mempool_ticket_index"`
	ForLastBlock       bool            `json:"last_block"`
}

// SetTicketIndex assigns the VoteInfo an index based on the block that the vote
// is (in)validating and the spent ticket hash. The ticketSpendInds tracks
// known combinations of target block and spent ticket hash. This index is used
// for sorting in views and counting total unique votes for a block.
func (vi *VoteInfo) SetTicketIndex(ticketSpendInds BlockValidatorIndex) {
	// One-based indexing
	startInd := 1
	// Reference the sub-index for the block being (in)validated by this vote.
	if idxs, ok := ticketSpendInds[vi.Validation.Hash]; ok {
		// If this ticket has been seen before voting on this block, set the
		// known index. Otherwise, assign the next index in the series.
		if idx, ok := idxs[vi.TicketSpent]; ok {
			vi.MempoolTicketIndex = idx
		} else {
			idx := len(idxs) + startInd
			idxs[vi.TicketSpent] = idx
			vi.MempoolTicketIndex = idx
		}
	} else {
		// First vote encountered for this block. Create new ticket sub-index.
		ticketSpendInds[vi.Validation.Hash] = TicketIndex{
			vi.TicketSpent: startInd,
		}
		vi.MempoolTicketIndex = startInd
	}
}

func (vi *VoteInfo) DeepCopy() *VoteInfo {
	if vi == nil {
		return nil
	}
	out := *vi
	out.Choices = make([]*VoteChoice, len(vi.Choices))
	copy(out.Choices, vi.Choices)
	return &out
}

// BlockValidation models data about a vote's decision on a block
type BlockValidation struct {
	Hash     string `json:"hash"`
	Height   int64  `json:"height"`
	Validity bool   `json:"validity"`
}

// VoteChoice represents the choice made by a vote transaction on a single vote
// item in an agenda. The ID, Description, and Mask fields describe the vote
// item for which the choice is being made. Those are the initial fields in
// chaincfg.Params.Deployments[VoteVersion][VoteIndex].
type VoteChoice struct {
	// Single unique word identifying the vote.
	ID string `json:"id"`

	// Longer description of what the vote is about.
	Description string `json:"description"`

	// Usable bits for this vote.
	Mask uint16 `json:"mask"`

	// VoteVersion and VoteIndex specify which vote item is referenced by this
	// VoteChoice (i.e. chaincfg.Params.Deployments[VoteVersion][VoteIndex]).
	VoteVersion uint32 `json:"vote_version"`
	VoteIndex   int    `json:"vote_index"`

	// ChoiceIdx indicates the corresponding element in the vote item's []Choice
	ChoiceIdx int `json:"choice_index"`

	// Choice is the selected choice for the specified vote item
	Choice *chaincfg.Choice `json:"choice"`
}

// SSGenVoteBlockValid determines if a vote transaction is voting yes or no to a
// block, and returns the votebits in case the caller wants to check agenda
// votes. The error return may be ignored if the input transaction is known to
// be a valid ssgen (vote), otherwise it should be checked.
func SSGenVoteBlockValid(msgTx *wire.MsgTx) (BlockValidation, uint16, error) {
	if !stake.IsSSGen(msgTx) {
		return BlockValidation{}, 0, fmt.Errorf("not a vote transaction")
	}

	ssGenVoteBits := stake.SSGenVoteBits(msgTx)
	blockHash, blockHeight := stake.SSGenBlockVotedOn(msgTx)
	blockValid := BlockValidation{
		Hash:     blockHash.String(),
		Height:   int64(blockHeight),
		Validity: dcrutil.IsFlagSet16(ssGenVoteBits, dcrutil.BlockValid),
	}
	return blockValid, ssGenVoteBits, nil
}

// SSGenVoteChoices gets a ssgen's vote choices (block validity and any
// agendas). The vote's stake version, to which the vote choices correspond, and
// vote bits are also returned. Note that []*VoteChoice may be an empty slice if
// there are no consensus deployments for the transaction's vote version. The
// error value may be non-nil if the tx is not a valid ssgen.
func SSGenVoteChoices(tx *wire.MsgTx, params *chaincfg.Params) (BlockValidation, uint32, uint16, []*VoteChoice, error) {
	validBlock, voteBits, err := SSGenVoteBlockValid(tx)
	if err != nil {
		return validBlock, 0, 0, nil, err
	}

	// Determine the ssgen's vote version and get the relevant consensus
	// deployments containing the vote items targeted.
	voteVersion := stake.SSGenVersion(tx)
	deployments := params.Deployments[voteVersion]

	// Allocate space for each choice
	choices := make([]*VoteChoice, 0, len(deployments))

	// For each vote item (consensus deployment), extract the choice from the
	// vote bits and store the vote item's Id, Description and vote bits Mask.
	for d := range deployments {
		voteAgenda := &deployments[d].Vote
		choiceIndex := voteAgenda.VoteIndex(voteBits)
		voteChoice := VoteChoice{
			ID:          voteAgenda.Id,
			Description: voteAgenda.Description,
			Mask:        voteAgenda.Mask,
			VoteVersion: voteVersion,
			VoteIndex:   d,
			ChoiceIdx:   choiceIndex,
			Choice:      &voteAgenda.Choices[choiceIndex],
		}
		choices = append(choices, &voteChoice)
	}

	return validBlock, voteVersion, voteBits, choices, nil
}

// DetermineTxTypeString returns a string representing the transaction type given
// a wire.MsgTx struct
func DetermineTxTypeString(msgTx *wire.MsgTx) string {
	switch stake.DetermineTxType(msgTx) {
	case stake.TxTypeSSGen:
		return "Vote"
	case stake.TxTypeSStx:
		return "Ticket"
	case stake.TxTypeSSRtx:
		return "Revocation"
	default:
		return "Regular"
	}
}

// MsgTxFromHex returns a wire.MsgTx struct built from the transaction hex string.
func MsgTxFromHex(txhex string) (*wire.MsgTx, error) {
	msgTx := wire.NewMsgTx()
	if err := msgTx.Deserialize(hex.NewDecoder(strings.NewReader(txhex))); err != nil {
		return nil, err
	}
	return msgTx, nil
}
