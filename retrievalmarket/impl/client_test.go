package retrievalimpl_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	dss "github.com/ipfs/go-datastore/sync"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	cbg "github.com/whyrusleeping/cbor-gen"

	"github.com/chenjianmei111/go-address"
	datatransfer "github.com/chenjianmei111/go-data-transfer"
	"github.com/chenjianmei111/go-multistore"
	"github.com/chenjianmei111/go-state-types/abi"
	"github.com/chenjianmei111/go-state-types/big"
	"github.com/chenjianmei111/go-storedcounter"

	"github.com/chenjianmei111/go-fil-markets/retrievalmarket"
	retrievalimpl "github.com/chenjianmei111/go-fil-markets/retrievalmarket/impl"
	"github.com/chenjianmei111/go-fil-markets/retrievalmarket/impl/testnodes"
	"github.com/chenjianmei111/go-fil-markets/retrievalmarket/migrations"
	"github.com/chenjianmei111/go-fil-markets/retrievalmarket/network"
	rmnet "github.com/chenjianmei111/go-fil-markets/retrievalmarket/network"
	"github.com/chenjianmei111/go-fil-markets/shared"
	"github.com/chenjianmei111/go-fil-markets/shared_testutil"
	tut "github.com/chenjianmei111/go-fil-markets/shared_testutil"
)

func TestClient_Construction(t *testing.T) {

	ds := dss.MutexWrap(datastore.NewMapDatastore())
	storedCounter := storedcounter.New(ds, datastore.NewKey("nextDealID"))
	multiStore, err := multistore.NewMultiDstore(ds)
	require.NoError(t, err)
	dt := tut.NewTestDataTransfer()
	net := tut.NewTestRetrievalMarketNetwork(tut.TestNetworkParams{})
	_, err = retrievalimpl.NewClient(
		net,
		multiStore,
		dt,
		testnodes.NewTestRetrievalClientNode(testnodes.TestRetrievalClientNodeParams{}),
		&tut.TestPeerResolver{},
		ds,
		storedCounter)
	require.NoError(t, err)

	require.Len(t, dt.Subscribers, 1)
	require.Len(t, dt.RegisteredVoucherResultTypes, 2)
	_, ok := dt.RegisteredVoucherResultTypes[0].(*retrievalmarket.DealResponse)
	require.True(t, ok)
	_, ok = dt.RegisteredVoucherResultTypes[1].(*migrations.DealResponse0)
	require.True(t, ok)
	require.Len(t, dt.RegisteredVoucherTypes, 4)
	_, ok = dt.RegisteredVoucherTypes[0].VoucherType.(*retrievalmarket.DealProposal)
	require.True(t, ok)
	_, ok = dt.RegisteredVoucherTypes[1].VoucherType.(*migrations.DealProposal0)
	require.True(t, ok)
	_, ok = dt.RegisteredVoucherTypes[2].VoucherType.(*retrievalmarket.DealPayment)
	require.True(t, ok)
	_, ok = dt.RegisteredVoucherTypes[3].VoucherType.(*migrations.DealPayment0)
	require.True(t, ok)
	require.Len(t, dt.RegisteredTransportConfigurers, 2)
	_, ok = dt.RegisteredTransportConfigurers[0].VoucherType.(*retrievalmarket.DealProposal)
	_, ok = dt.RegisteredTransportConfigurers[1].VoucherType.(*migrations.DealProposal0)
	require.True(t, ok)
}

func TestClient_Query(t *testing.T) {
	ctx := context.Background()

	ds := dss.MutexWrap(datastore.NewMapDatastore())
	storedCounter := storedcounter.New(ds, datastore.NewKey("nextDealID"))
	multiStore, err := multistore.NewMultiDstore(ds)
	require.NoError(t, err)
	dt := tut.NewTestDataTransfer()

	pcid := tut.GenerateCids(1)[0]
	expectedPeer := peer.ID("somevalue")
	rpeer := retrievalmarket.RetrievalPeer{
		Address: address.TestAddress2,
		ID:      expectedPeer,
	}

	expectedQuery := retrievalmarket.Query{
		PayloadCID: pcid,
	}

	expectedQueryResponse := retrievalmarket.QueryResponse{
		Status:                     retrievalmarket.QueryResponseAvailable,
		Size:                       1234,
		PaymentAddress:             address.TestAddress,
		MinPricePerByte:            abi.NewTokenAmount(5678),
		MaxPaymentInterval:         4321,
		MaxPaymentIntervalIncrease: 0,
	}

	t.Run("it works", func(t *testing.T) {
		var qsb tut.QueryStreamBuilder = func(p peer.ID) (rmnet.RetrievalQueryStream, error) {
			return tut.NewTestRetrievalQueryStream(tut.TestQueryStreamParams{
				Writer:     tut.ExpectQueryWriter(t, expectedQuery, "queries should match"),
				RespReader: tut.StubbedQueryResponseReader(expectedQueryResponse),
			}), nil
		}
		net := tut.NewTestRetrievalMarketNetwork(tut.TestNetworkParams{
			QueryStreamBuilder: tut.ExpectPeerOnQueryStreamBuilder(t, expectedPeer, qsb, "Peers should match"),
		})
		node := testnodes.NewTestRetrievalClientNode(testnodes.TestRetrievalClientNodeParams{})
		node.ExpectKnownAddresses(rpeer, nil)
		c, err := retrievalimpl.NewClient(
			net,
			multiStore,
			dt,
			node,
			&tut.TestPeerResolver{},
			ds,
			storedCounter)
		require.NoError(t, err)

		resp, err := c.Query(ctx, rpeer, pcid, retrievalmarket.QueryParams{})
		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, expectedQueryResponse, resp)
		node.VerifyExpectations(t)
	})

	t.Run("when the stream returns error, returns error", func(t *testing.T) {
		net := tut.NewTestRetrievalMarketNetwork(tut.TestNetworkParams{
			QueryStreamBuilder: tut.FailNewQueryStream,
		})
		node := testnodes.NewTestRetrievalClientNode(testnodes.TestRetrievalClientNodeParams{})
		node.ExpectKnownAddresses(rpeer, nil)
		c, err := retrievalimpl.NewClient(
			net,
			multiStore,
			dt,
			node,
			&tut.TestPeerResolver{},
			ds,
			storedCounter)
		require.NoError(t, err)

		_, err = c.Query(ctx, rpeer, pcid, retrievalmarket.QueryParams{})
		assert.EqualError(t, err, "new query stream failed")
		node.VerifyExpectations(t)
	})

	t.Run("when WriteDealStatusRequest fails, returns error", func(t *testing.T) {

		qsbuilder := func(p peer.ID) (network.RetrievalQueryStream, error) {
			newStream := tut.NewTestRetrievalQueryStream(tut.TestQueryStreamParams{
				PeerID: p,
				Writer: tut.FailQueryWriter,
			})
			return newStream, nil
		}

		net := tut.NewTestRetrievalMarketNetwork(tut.TestNetworkParams{
			QueryStreamBuilder: qsbuilder,
		})
		node := testnodes.NewTestRetrievalClientNode(testnodes.TestRetrievalClientNodeParams{})
		node.ExpectKnownAddresses(rpeer, nil)
		c, err := retrievalimpl.NewClient(
			net,
			multiStore,
			dt,
			node,
			&tut.TestPeerResolver{},
			ds,
			storedCounter)
		require.NoError(t, err)

		statusCode, err := c.Query(ctx, rpeer, pcid, retrievalmarket.QueryParams{})
		assert.EqualError(t, err, "write query failed")
		assert.Equal(t, retrievalmarket.QueryResponseUndefined, statusCode)
		node.VerifyExpectations(t)
	})

	t.Run("when ReadDealStatusResponse fails, returns error", func(t *testing.T) {
		qsbuilder := func(p peer.ID) (network.RetrievalQueryStream, error) {
			newStream := tut.NewTestRetrievalQueryStream(tut.TestQueryStreamParams{
				PeerID:     p,
				RespReader: tut.FailResponseReader,
			})
			return newStream, nil
		}
		net := tut.NewTestRetrievalMarketNetwork(tut.TestNetworkParams{
			QueryStreamBuilder: qsbuilder,
		})
		node := testnodes.NewTestRetrievalClientNode(testnodes.TestRetrievalClientNodeParams{})
		node.ExpectKnownAddresses(rpeer, nil)
		c, err := retrievalimpl.NewClient(
			net,
			multiStore,
			dt,
			node,
			&tut.TestPeerResolver{},
			ds,
			storedCounter)
		require.NoError(t, err)

		statusCode, err := c.Query(ctx, rpeer, pcid, retrievalmarket.QueryParams{})
		assert.EqualError(t, err, "query response failed")
		assert.Equal(t, retrievalmarket.QueryResponseUndefined, statusCode)
		node.VerifyExpectations(t)
	})
}

func TestClient_FindProviders(t *testing.T) {
	ds := dss.MutexWrap(datastore.NewMapDatastore())
	storedCounter := storedcounter.New(ds, datastore.NewKey("nextDealID"))
	multiStore, err := multistore.NewMultiDstore(ds)
	require.NoError(t, err)
	dt := tut.NewTestDataTransfer()
	expectedPeer := peer.ID("somevalue")

	var qsb tut.QueryStreamBuilder = func(p peer.ID) (rmnet.RetrievalQueryStream, error) {
		return tut.NewTestRetrievalQueryStream(tut.TestQueryStreamParams{
			Writer:     tut.TrivialQueryWriter,
			RespReader: tut.TrivialQueryResponseReader,
		}), nil
	}
	net := tut.NewTestRetrievalMarketNetwork(tut.TestNetworkParams{
		QueryStreamBuilder: tut.ExpectPeerOnQueryStreamBuilder(t, expectedPeer, qsb, "Peers should match"),
	})

	t.Run("when providers are found, returns providers", func(t *testing.T) {
		peers := tut.RequireGenerateRetrievalPeers(t, 3)
		testResolver := tut.TestPeerResolver{Peers: peers}

		c, err := retrievalimpl.NewClient(net, multiStore, dt, &testnodes.TestRetrievalClientNode{}, &testResolver, ds, storedCounter)
		require.NoError(t, err)

		testCid := tut.GenerateCids(1)[0]
		assert.Len(t, c.FindProviders(testCid), 3)
	})

	t.Run("when there is an error, returns empty provider list", func(t *testing.T) {
		testResolver := tut.TestPeerResolver{Peers: []retrievalmarket.RetrievalPeer{}, ResolverError: errors.New("boom")}
		c, err := retrievalimpl.NewClient(net, multiStore, dt, &testnodes.TestRetrievalClientNode{}, &testResolver, ds, storedCounter)
		require.NoError(t, err)

		badCid := tut.GenerateCids(1)[0]
		assert.Len(t, c.FindProviders(badCid), 0)
	})

	t.Run("when there are no providers", func(t *testing.T) {
		testResolver := tut.TestPeerResolver{Peers: []retrievalmarket.RetrievalPeer{}}
		c, err := retrievalimpl.NewClient(net, multiStore, dt, &testnodes.TestRetrievalClientNode{}, &testResolver, ds, storedCounter)
		require.NoError(t, err)

		testCid := tut.GenerateCids(1)[0]
		assert.Len(t, c.FindProviders(testCid), 0)
	})
}

func TestMigrations(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	ds := dss.MutexWrap(datastore.NewMapDatastore())
	storedCounter := storedcounter.New(ds, datastore.NewKey("nextDealID"))
	multiStore, err := multistore.NewMultiDstore(ds)
	require.NoError(t, err)
	dt := tut.NewTestDataTransfer()
	net := tut.NewTestRetrievalMarketNetwork(tut.TestNetworkParams{})
	retrievalDs := namespace.Wrap(ds, datastore.NewKey("/retrievals/client"))

	numDeals := 5
	payloadCIDs := make([]cid.Cid, numDeals)
	iDs := make([]retrievalmarket.DealID, numDeals)
	pieceCIDs := make([]*cid.Cid, numDeals)
	pricePerBytes := make([]abi.TokenAmount, numDeals)
	paymentIntervals := make([]uint64, numDeals)
	paymentIntervalIncreases := make([]uint64, numDeals)
	unsealPrices := make([]abi.TokenAmount, numDeals)
	storeIDs := make([]*multistore.StoreID, numDeals)
	channelIDs := make([]datatransfer.ChannelID, numDeals)
	lastPaymentRequesteds := make([]bool, numDeals)
	allBlocksReceiveds := make([]bool, numDeals)
	totalFundss := make([]abi.TokenAmount, numDeals)
	lanes := make([]uint64, numDeals)
	senders := make([]peer.ID, numDeals)
	totalReceiveds := make([]uint64, numDeals)
	messages := make([]string, numDeals)
	bytesPaidFors := make([]uint64, numDeals)
	currentIntervals := make([]uint64, numDeals)
	paymentRequesteds := make([]abi.TokenAmount, numDeals)
	fundsSpents := make([]abi.TokenAmount, numDeals)
	unsealFundsPaids := make([]abi.TokenAmount, numDeals)
	voucherShortfalls := make([]abi.TokenAmount, numDeals)
	selfPeer := tut.GeneratePeers(1)[0]

	allSelectorBuf := new(bytes.Buffer)
	err = dagcbor.Encoder(shared.AllSelector(), allSelectorBuf)
	require.NoError(t, err)
	allSelectorBytes := allSelectorBuf.Bytes()

	for i := 0; i < numDeals; i++ {
		payloadCIDs[i] = tut.GenerateCids(1)[0]
		iDs[i] = retrievalmarket.DealID(rand.Uint64())
		pieceCID := tut.GenerateCids(1)[0]
		pieceCIDs[i] = &pieceCID
		pricePerBytes[i] = big.NewInt(rand.Int63())
		paymentIntervals[i] = rand.Uint64()
		paymentIntervalIncreases[i] = rand.Uint64()
		unsealPrices[i] = big.NewInt(rand.Int63())
		storeID := multistore.StoreID(rand.Uint64())
		storeIDs[i] = &storeID
		senders[i] = tut.GeneratePeers(1)[0]
		channelIDs[i] = datatransfer.ChannelID{
			Initiator: selfPeer,
			Responder: senders[i],
			ID:        datatransfer.TransferID(rand.Uint64()),
		}
		lastPaymentRequesteds[i] = rand.Intn(2) == 1
		allBlocksReceiveds[i] = rand.Intn(2) == 1
		totalFundss[i] = big.NewInt(rand.Int63())
		lanes[i] = rand.Uint64()
		totalReceiveds[i] = rand.Uint64()
		messages[i] = string(tut.RandomBytes(20))
		bytesPaidFors[i] = rand.Uint64()
		currentIntervals[i] = rand.Uint64()
		paymentRequesteds[i] = big.NewInt(rand.Int63())
		fundsSpents[i] = big.NewInt(rand.Int63())
		unsealFundsPaids[i] = big.NewInt(rand.Int63())
		voucherShortfalls[i] = big.NewInt(rand.Int63())
		deal := migrations.ClientDealState0{
			DealProposal0: migrations.DealProposal0{
				PayloadCID: payloadCIDs[i],
				ID:         iDs[i],
				Params0: migrations.Params0{
					Selector: &cbg.Deferred{
						Raw: allSelectorBytes,
					},
					PieceCID:                pieceCIDs[i],
					PricePerByte:            pricePerBytes[i],
					PaymentInterval:         paymentIntervals[i],
					PaymentIntervalIncrease: paymentIntervalIncreases[i],
					UnsealPrice:             unsealPrices[i],
				},
			},
			StoreID:              storeIDs[i],
			ChannelID:            channelIDs[i],
			LastPaymentRequested: lastPaymentRequesteds[i],
			AllBlocksReceived:    allBlocksReceiveds[i],
			TotalFunds:           totalFundss[i],
			ClientWallet:         address.TestAddress,
			MinerWallet:          address.TestAddress2,
			PaymentInfo: &migrations.PaymentInfo0{
				PayCh: address.TestAddress,
				Lane:  lanes[i],
			},
			Status:           retrievalmarket.DealStatusCompleted,
			Sender:           senders[i],
			TotalReceived:    totalReceiveds[i],
			Message:          messages[i],
			BytesPaidFor:     bytesPaidFors[i],
			CurrentInterval:  currentIntervals[i],
			PaymentRequested: paymentRequesteds[i],
			FundsSpent:       fundsSpents[i],
			UnsealFundsPaid:  unsealFundsPaids[i],
			WaitMsgCID:       nil,
			VoucherShortfall: voucherShortfalls[i],
		}
		buf := new(bytes.Buffer)
		err := deal.MarshalCBOR(buf)
		require.NoError(t, err)
		err = retrievalDs.Put(datastore.NewKey(fmt.Sprint(deal.ID)), buf.Bytes())
		require.NoError(t, err)
	}
	retrievalClient, err := retrievalimpl.NewClient(
		net,
		multiStore,
		dt,
		testnodes.NewTestRetrievalClientNode(testnodes.TestRetrievalClientNodeParams{}),
		&tut.TestPeerResolver{},
		retrievalDs,
		storedCounter)
	require.NoError(t, err)
	shared_testutil.StartAndWaitForReady(ctx, t, retrievalClient)
	deals, err := retrievalClient.ListDeals()
	require.NoError(t, err)
	for i := 0; i < numDeals; i++ {
		deal, ok := deals[iDs[i]]
		require.True(t, ok)
		expectedDeal := retrievalmarket.ClientDealState{
			DealProposal: retrievalmarket.DealProposal{
				PayloadCID: payloadCIDs[i],
				ID:         iDs[i],
				Params: retrievalmarket.Params{
					Selector: &cbg.Deferred{
						Raw: allSelectorBytes,
					},
					PieceCID:                pieceCIDs[i],
					PricePerByte:            pricePerBytes[i],
					PaymentInterval:         paymentIntervals[i],
					PaymentIntervalIncrease: paymentIntervalIncreases[i],
					UnsealPrice:             unsealPrices[i],
				},
			},
			StoreID:              storeIDs[i],
			ChannelID:            channelIDs[i],
			LastPaymentRequested: lastPaymentRequesteds[i],
			AllBlocksReceived:    allBlocksReceiveds[i],
			TotalFunds:           totalFundss[i],
			ClientWallet:         address.TestAddress,
			MinerWallet:          address.TestAddress2,
			PaymentInfo: &retrievalmarket.PaymentInfo{
				PayCh: address.TestAddress,
				Lane:  lanes[i],
			},
			Status:           retrievalmarket.DealStatusCompleted,
			Sender:           senders[i],
			TotalReceived:    totalReceiveds[i],
			Message:          messages[i],
			BytesPaidFor:     bytesPaidFors[i],
			CurrentInterval:  currentIntervals[i],
			PaymentRequested: paymentRequesteds[i],
			FundsSpent:       fundsSpents[i],
			UnsealFundsPaid:  unsealFundsPaids[i],
			WaitMsgCID:       nil,
			VoucherShortfall: voucherShortfalls[i],
			LegacyProtocol:   true,
		}
		require.Equal(t, expectedDeal, deal)
	}
}
