// Code generated by SQLBoiler 4.5.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import "testing"

// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("Blocks", testBlocks)
	t.Run("BlockBins", testBlockBins)
	t.Run("Exchanges", testExchanges)
	t.Run("ExchangeTicks", testExchangeTicks)
	t.Run("Githubs", testGithubs)
	t.Run("Heartbeats", testHeartbeats)
	t.Run("Mempools", testMempools)
	t.Run("MempoolBins", testMempoolBins)
	t.Run("NetworkSnapshots", testNetworkSnapshots)
	t.Run("NetworkSnapshotBins", testNetworkSnapshotBins)
	t.Run("Nodes", testNodes)
	t.Run("NodeLocations", testNodeLocations)
	t.Run("NodeVersions", testNodeVersions)
	t.Run("Propagations", testPropagations)
	t.Run("Reddits", testReddits)
	t.Run("Twitters", testTwitters)
	t.Run("Votes", testVotes)
	t.Run("VoteReceiveTimeDeviations", testVoteReceiveTimeDeviations)
	t.Run("Youtubes", testYoutubes)
}

func TestDelete(t *testing.T) {
	t.Run("Blocks", testBlocksDelete)
	t.Run("BlockBins", testBlockBinsDelete)
	t.Run("Exchanges", testExchangesDelete)
	t.Run("ExchangeTicks", testExchangeTicksDelete)
	t.Run("Githubs", testGithubsDelete)
	t.Run("Heartbeats", testHeartbeatsDelete)
	t.Run("Mempools", testMempoolsDelete)
	t.Run("MempoolBins", testMempoolBinsDelete)
	t.Run("NetworkSnapshots", testNetworkSnapshotsDelete)
	t.Run("NetworkSnapshotBins", testNetworkSnapshotBinsDelete)
	t.Run("Nodes", testNodesDelete)
	t.Run("NodeLocations", testNodeLocationsDelete)
	t.Run("NodeVersions", testNodeVersionsDelete)
	t.Run("Propagations", testPropagationsDelete)
	t.Run("Reddits", testRedditsDelete)
	t.Run("Twitters", testTwittersDelete)
	t.Run("Votes", testVotesDelete)
	t.Run("VoteReceiveTimeDeviations", testVoteReceiveTimeDeviationsDelete)
	t.Run("Youtubes", testYoutubesDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("Blocks", testBlocksQueryDeleteAll)
	t.Run("BlockBins", testBlockBinsQueryDeleteAll)
	t.Run("Exchanges", testExchangesQueryDeleteAll)
	t.Run("ExchangeTicks", testExchangeTicksQueryDeleteAll)
	t.Run("Githubs", testGithubsQueryDeleteAll)
	t.Run("Heartbeats", testHeartbeatsQueryDeleteAll)
	t.Run("Mempools", testMempoolsQueryDeleteAll)
	t.Run("MempoolBins", testMempoolBinsQueryDeleteAll)
	t.Run("NetworkSnapshots", testNetworkSnapshotsQueryDeleteAll)
	t.Run("NetworkSnapshotBins", testNetworkSnapshotBinsQueryDeleteAll)
	t.Run("Nodes", testNodesQueryDeleteAll)
	t.Run("NodeLocations", testNodeLocationsQueryDeleteAll)
	t.Run("NodeVersions", testNodeVersionsQueryDeleteAll)
	t.Run("Propagations", testPropagationsQueryDeleteAll)
	t.Run("Reddits", testRedditsQueryDeleteAll)
	t.Run("Twitters", testTwittersQueryDeleteAll)
	t.Run("Votes", testVotesQueryDeleteAll)
	t.Run("VoteReceiveTimeDeviations", testVoteReceiveTimeDeviationsQueryDeleteAll)
	t.Run("Youtubes", testYoutubesQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("Blocks", testBlocksSliceDeleteAll)
	t.Run("BlockBins", testBlockBinsSliceDeleteAll)
	t.Run("Exchanges", testExchangesSliceDeleteAll)
	t.Run("ExchangeTicks", testExchangeTicksSliceDeleteAll)
	t.Run("Githubs", testGithubsSliceDeleteAll)
	t.Run("Heartbeats", testHeartbeatsSliceDeleteAll)
	t.Run("Mempools", testMempoolsSliceDeleteAll)
	t.Run("MempoolBins", testMempoolBinsSliceDeleteAll)
	t.Run("NetworkSnapshots", testNetworkSnapshotsSliceDeleteAll)
	t.Run("NetworkSnapshotBins", testNetworkSnapshotBinsSliceDeleteAll)
	t.Run("Nodes", testNodesSliceDeleteAll)
	t.Run("NodeLocations", testNodeLocationsSliceDeleteAll)
	t.Run("NodeVersions", testNodeVersionsSliceDeleteAll)
	t.Run("Propagations", testPropagationsSliceDeleteAll)
	t.Run("Reddits", testRedditsSliceDeleteAll)
	t.Run("Twitters", testTwittersSliceDeleteAll)
	t.Run("Votes", testVotesSliceDeleteAll)
	t.Run("VoteReceiveTimeDeviations", testVoteReceiveTimeDeviationsSliceDeleteAll)
	t.Run("Youtubes", testYoutubesSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("Blocks", testBlocksExists)
	t.Run("BlockBins", testBlockBinsExists)
	t.Run("Exchanges", testExchangesExists)
	t.Run("ExchangeTicks", testExchangeTicksExists)
	t.Run("Githubs", testGithubsExists)
	t.Run("Heartbeats", testHeartbeatsExists)
	t.Run("Mempools", testMempoolsExists)
	t.Run("MempoolBins", testMempoolBinsExists)
	t.Run("NetworkSnapshots", testNetworkSnapshotsExists)
	t.Run("NetworkSnapshotBins", testNetworkSnapshotBinsExists)
	t.Run("Nodes", testNodesExists)
	t.Run("NodeLocations", testNodeLocationsExists)
	t.Run("NodeVersions", testNodeVersionsExists)
	t.Run("Propagations", testPropagationsExists)
	t.Run("Reddits", testRedditsExists)
	t.Run("Twitters", testTwittersExists)
	t.Run("Votes", testVotesExists)
	t.Run("VoteReceiveTimeDeviations", testVoteReceiveTimeDeviationsExists)
	t.Run("Youtubes", testYoutubesExists)
}

func TestFind(t *testing.T) {
	t.Run("Blocks", testBlocksFind)
	t.Run("BlockBins", testBlockBinsFind)
	t.Run("Exchanges", testExchangesFind)
	t.Run("ExchangeTicks", testExchangeTicksFind)
	t.Run("Githubs", testGithubsFind)
	t.Run("Heartbeats", testHeartbeatsFind)
	t.Run("Mempools", testMempoolsFind)
	t.Run("MempoolBins", testMempoolBinsFind)
	t.Run("NetworkSnapshots", testNetworkSnapshotsFind)
	t.Run("NetworkSnapshotBins", testNetworkSnapshotBinsFind)
	t.Run("Nodes", testNodesFind)
	t.Run("NodeLocations", testNodeLocationsFind)
	t.Run("NodeVersions", testNodeVersionsFind)
	t.Run("Propagations", testPropagationsFind)
	t.Run("Reddits", testRedditsFind)
	t.Run("Twitters", testTwittersFind)
	t.Run("Votes", testVotesFind)
	t.Run("VoteReceiveTimeDeviations", testVoteReceiveTimeDeviationsFind)
	t.Run("Youtubes", testYoutubesFind)
}

func TestBind(t *testing.T) {
	t.Run("Blocks", testBlocksBind)
	t.Run("BlockBins", testBlockBinsBind)
	t.Run("Exchanges", testExchangesBind)
	t.Run("ExchangeTicks", testExchangeTicksBind)
	t.Run("Githubs", testGithubsBind)
	t.Run("Heartbeats", testHeartbeatsBind)
	t.Run("Mempools", testMempoolsBind)
	t.Run("MempoolBins", testMempoolBinsBind)
	t.Run("NetworkSnapshots", testNetworkSnapshotsBind)
	t.Run("NetworkSnapshotBins", testNetworkSnapshotBinsBind)
	t.Run("Nodes", testNodesBind)
	t.Run("NodeLocations", testNodeLocationsBind)
	t.Run("NodeVersions", testNodeVersionsBind)
	t.Run("Propagations", testPropagationsBind)
	t.Run("Reddits", testRedditsBind)
	t.Run("Twitters", testTwittersBind)
	t.Run("Votes", testVotesBind)
	t.Run("VoteReceiveTimeDeviations", testVoteReceiveTimeDeviationsBind)
	t.Run("Youtubes", testYoutubesBind)
}

func TestOne(t *testing.T) {
	t.Run("Blocks", testBlocksOne)
	t.Run("BlockBins", testBlockBinsOne)
	t.Run("Exchanges", testExchangesOne)
	t.Run("ExchangeTicks", testExchangeTicksOne)
	t.Run("Githubs", testGithubsOne)
	t.Run("Heartbeats", testHeartbeatsOne)
	t.Run("Mempools", testMempoolsOne)
	t.Run("MempoolBins", testMempoolBinsOne)
	t.Run("NetworkSnapshots", testNetworkSnapshotsOne)
	t.Run("NetworkSnapshotBins", testNetworkSnapshotBinsOne)
	t.Run("Nodes", testNodesOne)
	t.Run("NodeLocations", testNodeLocationsOne)
	t.Run("NodeVersions", testNodeVersionsOne)
	t.Run("Propagations", testPropagationsOne)
	t.Run("Reddits", testRedditsOne)
	t.Run("Twitters", testTwittersOne)
	t.Run("Votes", testVotesOne)
	t.Run("VoteReceiveTimeDeviations", testVoteReceiveTimeDeviationsOne)
	t.Run("Youtubes", testYoutubesOne)
}

func TestAll(t *testing.T) {
	t.Run("Blocks", testBlocksAll)
	t.Run("BlockBins", testBlockBinsAll)
	t.Run("Exchanges", testExchangesAll)
	t.Run("ExchangeTicks", testExchangeTicksAll)
	t.Run("Githubs", testGithubsAll)
	t.Run("Heartbeats", testHeartbeatsAll)
	t.Run("Mempools", testMempoolsAll)
	t.Run("MempoolBins", testMempoolBinsAll)
	t.Run("NetworkSnapshots", testNetworkSnapshotsAll)
	t.Run("NetworkSnapshotBins", testNetworkSnapshotBinsAll)
	t.Run("Nodes", testNodesAll)
	t.Run("NodeLocations", testNodeLocationsAll)
	t.Run("NodeVersions", testNodeVersionsAll)
	t.Run("Propagations", testPropagationsAll)
	t.Run("Reddits", testRedditsAll)
	t.Run("Twitters", testTwittersAll)
	t.Run("Votes", testVotesAll)
	t.Run("VoteReceiveTimeDeviations", testVoteReceiveTimeDeviationsAll)
	t.Run("Youtubes", testYoutubesAll)
}

func TestCount(t *testing.T) {
	t.Run("Blocks", testBlocksCount)
	t.Run("BlockBins", testBlockBinsCount)
	t.Run("Exchanges", testExchangesCount)
	t.Run("ExchangeTicks", testExchangeTicksCount)
	t.Run("Githubs", testGithubsCount)
	t.Run("Heartbeats", testHeartbeatsCount)
	t.Run("Mempools", testMempoolsCount)
	t.Run("MempoolBins", testMempoolBinsCount)
	t.Run("NetworkSnapshots", testNetworkSnapshotsCount)
	t.Run("NetworkSnapshotBins", testNetworkSnapshotBinsCount)
	t.Run("Nodes", testNodesCount)
	t.Run("NodeLocations", testNodeLocationsCount)
	t.Run("NodeVersions", testNodeVersionsCount)
	t.Run("Propagations", testPropagationsCount)
	t.Run("Reddits", testRedditsCount)
	t.Run("Twitters", testTwittersCount)
	t.Run("Votes", testVotesCount)
	t.Run("VoteReceiveTimeDeviations", testVoteReceiveTimeDeviationsCount)
	t.Run("Youtubes", testYoutubesCount)
}

func TestInsert(t *testing.T) {
	t.Run("Blocks", testBlocksInsert)
	t.Run("Blocks", testBlocksInsertWhitelist)
	t.Run("BlockBins", testBlockBinsInsert)
	t.Run("BlockBins", testBlockBinsInsertWhitelist)
	t.Run("Exchanges", testExchangesInsert)
	t.Run("Exchanges", testExchangesInsertWhitelist)
	t.Run("ExchangeTicks", testExchangeTicksInsert)
	t.Run("ExchangeTicks", testExchangeTicksInsertWhitelist)
	t.Run("Githubs", testGithubsInsert)
	t.Run("Githubs", testGithubsInsertWhitelist)
	t.Run("Heartbeats", testHeartbeatsInsert)
	t.Run("Heartbeats", testHeartbeatsInsertWhitelist)
	t.Run("Mempools", testMempoolsInsert)
	t.Run("Mempools", testMempoolsInsertWhitelist)
	t.Run("MempoolBins", testMempoolBinsInsert)
	t.Run("MempoolBins", testMempoolBinsInsertWhitelist)
	t.Run("NetworkSnapshots", testNetworkSnapshotsInsert)
	t.Run("NetworkSnapshots", testNetworkSnapshotsInsertWhitelist)
	t.Run("NetworkSnapshotBins", testNetworkSnapshotBinsInsert)
	t.Run("NetworkSnapshotBins", testNetworkSnapshotBinsInsertWhitelist)
	t.Run("Nodes", testNodesInsert)
	t.Run("Nodes", testNodesInsertWhitelist)
	t.Run("NodeLocations", testNodeLocationsInsert)
	t.Run("NodeLocations", testNodeLocationsInsertWhitelist)
	t.Run("NodeVersions", testNodeVersionsInsert)
	t.Run("NodeVersions", testNodeVersionsInsertWhitelist)
	t.Run("Propagations", testPropagationsInsert)
	t.Run("Propagations", testPropagationsInsertWhitelist)
	t.Run("Reddits", testRedditsInsert)
	t.Run("Reddits", testRedditsInsertWhitelist)
	t.Run("Twitters", testTwittersInsert)
	t.Run("Twitters", testTwittersInsertWhitelist)
	t.Run("Votes", testVotesInsert)
	t.Run("Votes", testVotesInsertWhitelist)
	t.Run("VoteReceiveTimeDeviations", testVoteReceiveTimeDeviationsInsert)
	t.Run("VoteReceiveTimeDeviations", testVoteReceiveTimeDeviationsInsertWhitelist)
	t.Run("Youtubes", testYoutubesInsert)
	t.Run("Youtubes", testYoutubesInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
	t.Run("ExchangeTickToExchangeUsingExchange", testExchangeTickToOneExchangeUsingExchange)
	t.Run("HeartbeatToNodeUsingNode", testHeartbeatToOneNodeUsingNode)
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
	t.Run("ExchangeToExchangeTicks", testExchangeToManyExchangeTicks)
	t.Run("NodeToHeartbeats", testNodeToManyHeartbeats)
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
	t.Run("ExchangeTickToExchangeUsingExchangeTicks", testExchangeTickToOneSetOpExchangeUsingExchange)
	t.Run("HeartbeatToNodeUsingHeartbeats", testHeartbeatToOneSetOpNodeUsingNode)
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {
	t.Run("ExchangeToExchangeTicks", testExchangeToManyAddOpExchangeTicks)
	t.Run("NodeToHeartbeats", testNodeToManyAddOpHeartbeats)
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {}

func TestReload(t *testing.T) {
	t.Run("Blocks", testBlocksReload)
	t.Run("BlockBins", testBlockBinsReload)
	t.Run("Exchanges", testExchangesReload)
	t.Run("ExchangeTicks", testExchangeTicksReload)
	t.Run("Githubs", testGithubsReload)
	t.Run("Heartbeats", testHeartbeatsReload)
	t.Run("Mempools", testMempoolsReload)
	t.Run("MempoolBins", testMempoolBinsReload)
	t.Run("NetworkSnapshots", testNetworkSnapshotsReload)
	t.Run("NetworkSnapshotBins", testNetworkSnapshotBinsReload)
	t.Run("Nodes", testNodesReload)
	t.Run("NodeLocations", testNodeLocationsReload)
	t.Run("NodeVersions", testNodeVersionsReload)
	t.Run("Propagations", testPropagationsReload)
	t.Run("Reddits", testRedditsReload)
	t.Run("Twitters", testTwittersReload)
	t.Run("Votes", testVotesReload)
	t.Run("VoteReceiveTimeDeviations", testVoteReceiveTimeDeviationsReload)
	t.Run("Youtubes", testYoutubesReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("Blocks", testBlocksReloadAll)
	t.Run("BlockBins", testBlockBinsReloadAll)
	t.Run("Exchanges", testExchangesReloadAll)
	t.Run("ExchangeTicks", testExchangeTicksReloadAll)
	t.Run("Githubs", testGithubsReloadAll)
	t.Run("Heartbeats", testHeartbeatsReloadAll)
	t.Run("Mempools", testMempoolsReloadAll)
	t.Run("MempoolBins", testMempoolBinsReloadAll)
	t.Run("NetworkSnapshots", testNetworkSnapshotsReloadAll)
	t.Run("NetworkSnapshotBins", testNetworkSnapshotBinsReloadAll)
	t.Run("Nodes", testNodesReloadAll)
	t.Run("NodeLocations", testNodeLocationsReloadAll)
	t.Run("NodeVersions", testNodeVersionsReloadAll)
	t.Run("Propagations", testPropagationsReloadAll)
	t.Run("Reddits", testRedditsReloadAll)
	t.Run("Twitters", testTwittersReloadAll)
	t.Run("Votes", testVotesReloadAll)
	t.Run("VoteReceiveTimeDeviations", testVoteReceiveTimeDeviationsReloadAll)
	t.Run("Youtubes", testYoutubesReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("Blocks", testBlocksSelect)
	t.Run("BlockBins", testBlockBinsSelect)
	t.Run("Exchanges", testExchangesSelect)
	t.Run("ExchangeTicks", testExchangeTicksSelect)
	t.Run("Githubs", testGithubsSelect)
	t.Run("Heartbeats", testHeartbeatsSelect)
	t.Run("Mempools", testMempoolsSelect)
	t.Run("MempoolBins", testMempoolBinsSelect)
	t.Run("NetworkSnapshots", testNetworkSnapshotsSelect)
	t.Run("NetworkSnapshotBins", testNetworkSnapshotBinsSelect)
	t.Run("Nodes", testNodesSelect)
	t.Run("NodeLocations", testNodeLocationsSelect)
	t.Run("NodeVersions", testNodeVersionsSelect)
	t.Run("Propagations", testPropagationsSelect)
	t.Run("Reddits", testRedditsSelect)
	t.Run("Twitters", testTwittersSelect)
	t.Run("Votes", testVotesSelect)
	t.Run("VoteReceiveTimeDeviations", testVoteReceiveTimeDeviationsSelect)
	t.Run("Youtubes", testYoutubesSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("Blocks", testBlocksUpdate)
	t.Run("BlockBins", testBlockBinsUpdate)
	t.Run("Exchanges", testExchangesUpdate)
	t.Run("ExchangeTicks", testExchangeTicksUpdate)
	t.Run("Githubs", testGithubsUpdate)
	t.Run("Heartbeats", testHeartbeatsUpdate)
	t.Run("Mempools", testMempoolsUpdate)
	t.Run("MempoolBins", testMempoolBinsUpdate)
	t.Run("NetworkSnapshots", testNetworkSnapshotsUpdate)
	t.Run("NetworkSnapshotBins", testNetworkSnapshotBinsUpdate)
	t.Run("Nodes", testNodesUpdate)
	t.Run("NodeLocations", testNodeLocationsUpdate)
	t.Run("NodeVersions", testNodeVersionsUpdate)
	t.Run("Propagations", testPropagationsUpdate)
	t.Run("Reddits", testRedditsUpdate)
	t.Run("Twitters", testTwittersUpdate)
	t.Run("Votes", testVotesUpdate)
	t.Run("VoteReceiveTimeDeviations", testVoteReceiveTimeDeviationsUpdate)
	t.Run("Youtubes", testYoutubesUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("Blocks", testBlocksSliceUpdateAll)
	t.Run("BlockBins", testBlockBinsSliceUpdateAll)
	t.Run("Exchanges", testExchangesSliceUpdateAll)
	t.Run("ExchangeTicks", testExchangeTicksSliceUpdateAll)
	t.Run("Githubs", testGithubsSliceUpdateAll)
	t.Run("Heartbeats", testHeartbeatsSliceUpdateAll)
	t.Run("Mempools", testMempoolsSliceUpdateAll)
	t.Run("MempoolBins", testMempoolBinsSliceUpdateAll)
	t.Run("NetworkSnapshots", testNetworkSnapshotsSliceUpdateAll)
	t.Run("NetworkSnapshotBins", testNetworkSnapshotBinsSliceUpdateAll)
	t.Run("Nodes", testNodesSliceUpdateAll)
	t.Run("NodeLocations", testNodeLocationsSliceUpdateAll)
	t.Run("NodeVersions", testNodeVersionsSliceUpdateAll)
	t.Run("Propagations", testPropagationsSliceUpdateAll)
	t.Run("Reddits", testRedditsSliceUpdateAll)
	t.Run("Twitters", testTwittersSliceUpdateAll)
	t.Run("Votes", testVotesSliceUpdateAll)
	t.Run("VoteReceiveTimeDeviations", testVoteReceiveTimeDeviationsSliceUpdateAll)
	t.Run("Youtubes", testYoutubesSliceUpdateAll)
}
