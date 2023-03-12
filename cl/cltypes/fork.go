package cltypes

import (
	"github.com/ledgerwatch/erigon/cl/clparams"
	"github.com/ledgerwatch/erigon/cl/cltypes/ssz"
	"github.com/ledgerwatch/erigon/cl/merkle_tree"
)

// Fork data, contains if we were on bellatrix/alteir/phase0 and transition epoch.
type Fork struct {
	PreviousVersion [4]byte `ssz:"true"`
	CurrentVersion  [4]byte `ssz:"true"`
	Epoch           uint64  `ssz:"true"`
}

func (f *Fork) EncodeSSZ(dst []byte) ([]byte, error) {
	return ssz.Encode(f, dst)
}

func (f *Fork) DecodeSSZ(buf []byte) error {
	if len(buf) < f.EncodingSizeSSZ() {
		return ssz.ErrLowBufferSize
	}
	copy(f.PreviousVersion[:], buf)
	copy(f.CurrentVersion[:], buf[clparams.VersionLength:])
	f.Epoch = ssz.UnmarshalUint64SSZ(buf[clparams.VersionLength*2:])
	return nil
}

func (f *Fork) EncodingSizeSSZ() int {
	return clparams.VersionLength*2 + 8
}

func (f *Fork) HashSSZ() ([32]byte, error) {
	return merkle_tree.HashTreeRoot(f)
}
