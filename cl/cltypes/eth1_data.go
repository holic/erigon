package cltypes

import (
	libcommon "github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon-lib/common/length"
	"github.com/ledgerwatch/erigon/cl/cltypes/ssz"
	"github.com/ledgerwatch/erigon/cl/merkle_tree"
	"github.com/ledgerwatch/erigon/common"
)

type Eth1Data struct {
	Root         libcommon.Hash `ssz:"true"`
	DepositCount uint64         `ssz:"true"`
	BlockHash    libcommon.Hash `ssz:"true"`
}

func (e *Eth1Data) Equal(b *Eth1Data) bool {
	return e.BlockHash == b.BlockHash && e.Root == b.Root && b.DepositCount == e.DepositCount
}

// MarshalSSZTo ssz marshals the Eth1Data object to a target array
func (e *Eth1Data) EncodeSSZ(buf []byte) ([]byte, error) {
	return ssz.Encode(e, buf)
}

// DecodeSSZ ssz unmarshals the Eth1Data object
func (e *Eth1Data) DecodeSSZ(buf []byte) error {
	return ssz.Decode(e, buf)
}

func (e *Eth1Data) DecodeSSZWithVersion(buf []byte, _ int) error {
	return e.DecodeSSZ(buf)
}

// EncodingSizeSSZ returns the ssz encoded size in bytes for the Eth1Data object
func (e *Eth1Data) EncodingSizeSSZ() int {
	return common.BlockNumberLength + length.Hash*2
}

// HashSSZ ssz hashes the Eth1Data object
func (e *Eth1Data) HashSSZ() ([32]byte, error) {
	return merkle_tree.HashTreeRoot(e)
}
