module github.com/chenjianmei111/go-fil-markets

go 1.13

require (
	github.com/chenjianmei111/go-address v0.0.6
	github.com/chenjianmei111/go-commp-utils v0.2.0
	github.com/chenjianmei111/go-data-transfer v1.5.0
	github.com/chenjianmei111/go-ds-versioning v0.2.0
	github.com/chenjianmei111/go-multistore v0.0.4
	github.com/chenjianmei111/go-padreader v0.0.1
	github.com/chenjianmei111/go-state-types v0.2.0
	github.com/chenjianmei111/go-statemachine v1.0.0
	github.com/chenjianmei111/go-statestore v0.1.2
	github.com/chenjianmei111/specs-actors v0.9.14
	github.com/chenjianmei111/specs-actors/v2 v2.3.5
	github.com/hannahhoward/cbor-gen-for v0.0.0-20200817222906-ea96cece81f1
	github.com/hannahhoward/go-pubsub v0.0.0-20200423002714-8d62886cc36e
	github.com/ipfs/go-block-format v0.0.2
	github.com/ipfs/go-blockservice v0.1.4-0.20200624145336-a978cec6e834
	github.com/ipfs/go-cid v0.0.7
	github.com/ipfs/go-datastore v0.4.5
	github.com/ipfs/go-graphsync v0.6.0
	github.com/ipfs/go-ipfs-blockstore v1.0.3
	github.com/ipfs/go-ipfs-blocksutil v0.0.1
	github.com/ipfs/go-ipfs-chunker v0.0.5
	github.com/ipfs/go-ipfs-ds-help v1.0.0
	github.com/ipfs/go-ipfs-exchange-offline v0.0.1
	github.com/ipfs/go-ipfs-files v0.0.8
	github.com/ipfs/go-ipld-format v0.2.0
	github.com/ipfs/go-log/v2 v2.1.2-0.20200626104915-0016c0b4b3e4
	github.com/ipfs/go-merkledag v0.3.2
	github.com/ipfs/go-unixfs v0.2.4
	github.com/ipld/go-car v0.1.1-0.20201119040415-11b6074b6d4d
	github.com/ipld/go-ipld-prime v0.5.1-0.20201021195245-109253e8a018
	github.com/jbenet/go-random v0.0.0-20190219211222-123a90aedc0c
	github.com/jpillora/backoff v1.0.0
	github.com/libp2p/go-libp2p v0.12.0
	github.com/libp2p/go-libp2p-core v0.7.0
	github.com/multiformats/go-multiaddr v0.3.1
	github.com/multiformats/go-multibase v0.0.3
	github.com/stretchr/testify v1.6.1
	github.com/whyrusleeping/cbor-gen v0.0.0-20210219115102-f37d292932f2
	golang.org/x/exp v0.0.0-20200207192155-f17229e696bd
	golang.org/x/net v0.0.0-20200625001655-4c5254603344
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
)

replace github.com/chenjianmei111/filecoin-ffi => ./extern/filecoin-ffi
