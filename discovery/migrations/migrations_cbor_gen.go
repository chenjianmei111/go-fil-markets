// Code generated by github.com/whyrusleeping/cbor-gen. DO NOT EDIT.

package migrations

import (
	"fmt"
	"io"
	"sort"

	migrations "github.com/chenjianmei111/go-fil-markets/retrievalmarket/migrations"
	cid "github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"
)

var _ = xerrors.Errorf
var _ = cid.Undef
var _ = sort.Sort

var lengthBufRetrievalPeers0 = []byte{129}

func (t *RetrievalPeers0) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufRetrievalPeers0); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.Peers ([]migrations.RetrievalPeer0) (slice)
	if len(t.Peers) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.Peers was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajArray, uint64(len(t.Peers))); err != nil {
		return err
	}
	for _, v := range t.Peers {
		if err := v.MarshalCBOR(w); err != nil {
			return err
		}
	}
	return nil
}

func (t *RetrievalPeers0) UnmarshalCBOR(r io.Reader) error {
	*t = RetrievalPeers0{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Peers ([]migrations.RetrievalPeer0) (slice)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("t.Peers: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Peers = make([]migrations.RetrievalPeer0, extra)
	}

	for i := 0; i < int(extra); i++ {

		var v migrations.RetrievalPeer0
		if err := v.UnmarshalCBOR(br); err != nil {
			return err
		}

		t.Peers[i] = v
	}

	return nil
}
