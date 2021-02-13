// Code generated by SQLBoiler 3.7.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import "testing"

func TestUpsert(t *testing.T) {
	t.Run("Heartbeats", testHeartbeatsUpsert)

	t.Run("NetworkSnapshots", testNetworkSnapshotsUpsert)

	t.Run("NetworkSnapshotBins", testNetworkSnapshotBinsUpsert)

	t.Run("Nodes", testNodesUpsert)

	t.Run("NodeLocations", testNodeLocationsUpsert)

	t.Run("NodeVersions", testNodeVersionsUpsert)
}
