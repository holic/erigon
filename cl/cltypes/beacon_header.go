package cltypes

import (
	libcommon "github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon-lib/common/length"

	"github.com/ledgerwatch/erigon/cl/cltypes/ssz"
	"github.com/ledgerwatch/erigon/cl/merkle_tree"
)

/*
 * BeaconBlockHeader is the message we validate in the lightclient.
 * It contains the hash of the block body, and state root data.
 */
type BeaconBlockHeader struct {
	Slot          uint64         `ssz:"true"`
	ProposerIndex uint64         `ssz:"true"`
	ParentRoot    libcommon.Hash `ssz:"true"`
	Root          libcommon.Hash `ssz:"true"`
	BodyRoot      libcommon.Hash `ssz:"true"`
}

func (b *BeaconBlockHeader) EncodeSSZ(dst []byte) ([]byte, error) {
	return ssz.Encode(b, dst)
}

func (b *BeaconBlockHeader) DecodeSSZ(buf []byte) error {
	return ssz.Decode(b, buf)
}

func (b *BeaconBlockHeader) HashSSZ() ([32]byte, error) {
	return merkle_tree.HashTreeRoot(b)
}

func (b *BeaconBlockHeader) EncodingSizeSSZ() int {
	return length.Hash*3 + length.BlockNum*2
}

/*
 * SignedBeaconBlockHeader is a beacon block header + validator signature.
 */
type SignedBeaconBlockHeader struct {
	Header    *BeaconBlockHeader `ssz:"true"`
	Signature [96]byte           `ssz:"true"`
}

func (b *SignedBeaconBlockHeader) EncodeSSZ(dst []byte) ([]byte, error) {
	return ssz.Encode(b, dst)
}

func (b *SignedBeaconBlockHeader) DecodeSSZ(buf []byte) error {
	return ssz.Decode(b, buf)
}

func (b *SignedBeaconBlockHeader) HashSSZ() ([32]byte, error) {
	return merkle_tree.HashTreeRoot(b)
}

func (b *SignedBeaconBlockHeader) EncodingSizeSSZ() int {
	return b.Header.EncodingSizeSSZ() + 96
}
