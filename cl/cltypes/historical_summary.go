package cltypes

import (
	libcommon "github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon-lib/common/length"
	"github.com/ledgerwatch/erigon/cl/cltypes/ssz"
	"github.com/ledgerwatch/erigon/cl/merkle_tree"
)

type HistoricalSummary struct {
	BlockSummaryRoot libcommon.Hash `ssz:"true"`
	StateSummaryRoot libcommon.Hash `ssz:"true"`
}

func (h *HistoricalSummary) EncodeSSZ(buf []byte) ([]byte, error) {
	return ssz.Encode(h, buf)
}

func (h *HistoricalSummary) DecodeSSZ(buf []byte) error {
	if len(buf) < h.EncodingSizeSSZ() {
		return ssz.ErrLowBufferSize
	}
	copy(h.BlockSummaryRoot[:], buf)
	copy(h.StateSummaryRoot[:], buf[length.Hash:])
	return nil
}

func (h *HistoricalSummary) DecodeSSZWithVersion(buf []byte, _ int) error {
	return h.DecodeSSZ(buf)
}

func (h *HistoricalSummary) HashSSZ() ([32]byte, error) {
	return merkle_tree.HashTreeRoot(h)
}

func (*HistoricalSummary) EncodingSizeSSZ() int {
	return length.Hash * 2
}
