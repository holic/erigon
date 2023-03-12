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
	var err error
	size := uint64(len(buf))
	if size < uint64(c.EncodingSizeSSZ()) {
		return ssz.ErrLowBufferSize
	}
	c.Epoch = ssz.UnmarshalUint64SSZ(buf[0:8])
	copy(c.Root[:], buf[8:40])

	return err
}

func (c *Checkpoint) EncodingSizeSSZ() int {
	return length.BlockNum + length.Hash
}

func (c *Checkpoint) HashSSZ() ([32]byte, error) {
	return merkle_tree.HashTreeRoot(c)
}
