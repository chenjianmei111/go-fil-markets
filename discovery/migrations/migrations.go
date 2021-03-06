package migrations

import (
	versioning "github.com/chenjianmei111/go-ds-versioning/pkg"
	"github.com/chenjianmei111/go-ds-versioning/pkg/versioned"

	"github.com/chenjianmei111/go-fil-markets/discovery"
	"github.com/chenjianmei111/go-fil-markets/retrievalmarket"
	"github.com/chenjianmei111/go-fil-markets/retrievalmarket/migrations"
)

//go:generate cbor-gen-for RetrievalPeers0

// RetrievalPeers0 is version 0 of RetrievalPeers
type RetrievalPeers0 struct {
	Peers []migrations.RetrievalPeer0
}

// MigrateRetrievalPeers0To1 migrates a tuple encoded list of retrieval peers to a map encoded list
func MigrateRetrievalPeers0To1(oldRps *RetrievalPeers0) (*discovery.RetrievalPeers, error) {
	peers := make([]retrievalmarket.RetrievalPeer, 0, len(oldRps.Peers))
	for _, oldRp := range oldRps.Peers {
		peers = append(peers, retrievalmarket.RetrievalPeer{
			Address:  oldRp.Address,
			ID:       oldRp.ID,
			PieceCID: oldRp.PieceCID,
		})
	}
	return &discovery.RetrievalPeers{
		Peers: peers,
	}, nil
}

// RetrievalPeersMigrations are migrations for the store local discovery list of peers we can retrieve from
var RetrievalPeersMigrations = versioned.BuilderList{
	versioned.NewVersionedBuilder(MigrateRetrievalPeers0To1, versioning.VersionKey("1")),
}
