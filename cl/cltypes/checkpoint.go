package cltypes

import (
	libcommon "github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon-lib/common/length"

	"github.com/ledgerwatch/erigon/cl/cltypes/ssz"
	"github.com/ledgerwatch/erigon/cl/merkle_tree"
)

type Checkpoint struct {
	Epoch uint64         `ssz:"true"`
	Root  libcommon.Hash `ssz:"true"`
}

func (c *Checkpoint) Equal(other *Checkpoint) bool {
	return c.Epoch == other.Epoch && c.Root == other.Root
}

func (c *Checkpoint) EncodeSSZ(buf []byte) ([]byte, error) {
	return ssz.Encode(c, buf)
}

func (c *Checkpoint) DecodeSSZ(buf []byte) error {
	return ssz.Decode(c, buf)
}

func (c *Checkpoint) EncodingSizeSSZ() int {
	return length.BlockNum + length.Hash
}

func (c *Checkpoint) HashSSZ() ([32]byte, error) {
	return merkle_tree.HashTreeRoot(c)
}
