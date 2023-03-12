package cltypes

import (
	libcommon "github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon/cl/cltypes/ssz"
	"github.com/ledgerwatch/erigon/cl/merkle_tree"
)

// Change to EL engine
type BLSToExecutionChange struct {
	ValidatorIndex uint64            `ssz:"true"`
	From           [48]byte          `ssz:"true"`
	To             libcommon.Address `ssz:"true"`
}

func (b *BLSToExecutionChange) EncodeSSZ(buf []byte) ([]byte, error) {
	return ssz.Encode(b, buf)
}

func (b *BLSToExecutionChange) HashSSZ() ([32]byte, error) {
	return merkle_tree.HashTreeRoot(b)
}

func (b *BLSToExecutionChange) DecodeSSZ(buf []byte) error {
	if len(buf) < b.EncodingSizeSSZ() {
		return ssz.ErrLowBufferSize
	}
	b.ValidatorIndex = ssz.UnmarshalUint64SSZ(buf)
	copy(b.From[:], buf[8:])
	copy(b.To[:], buf[56:])
	return nil
}

func (*BLSToExecutionChange) EncodingSizeSSZ() int {
	return 76
}

type SignedBLSToExecutionChange struct {
	Message   *BLSToExecutionChange `ssz:"true"`
	Signature [96]byte              `ssz:"true"`
}

func (s *SignedBLSToExecutionChange) EncodeSSZ(buf []byte) ([]byte, error) {
	return ssz.Encode(s, buf)
}

func (s *SignedBLSToExecutionChange) DecodeSSZ(buf []byte) error {
	if len(buf) < s.EncodingSizeSSZ() {
		return ssz.ErrLowBufferSize
	}
	s.Message = new(BLSToExecutionChange)
	if err := s.Message.DecodeSSZ(buf); err != nil {
		return err
	}
	copy(s.Signature[:], buf[s.Message.EncodingSizeSSZ():])
	return nil
}

func (s *SignedBLSToExecutionChange) DecodeSSZWithVersion(buf []byte, _ int) error {
	return s.DecodeSSZ(buf)
}

func (s *SignedBLSToExecutionChange) HashSSZ() ([32]byte, error) {
	return merkle_tree.HashTreeRoot(s)
}

func (s *SignedBLSToExecutionChange) EncodingSizeSSZ() int {
	return 96 + s.Message.EncodingSizeSSZ()
}
