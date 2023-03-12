package cltypes

import (
	"github.com/ledgerwatch/erigon-lib/common/length"
	"github.com/ledgerwatch/erigon/cl/cltypes/clonable"
	"github.com/ledgerwatch/erigon/cl/cltypes/ssz"
	"github.com/ledgerwatch/erigon/common"
)

type Metadata struct {
	SeqNumber uint64  `ssz:"true"`
	Attnets   uint64  `ssz:"true"`
	Syncnets  *uint64 `ssz:"true"`
}

func (m *Metadata) EncodeSSZ(buf []byte) ([]byte, error) {
	ret := buf
	ret = append(ret, ssz.Uint64SSZ(m.SeqNumber)...)
	ret = append(ret, ssz.Uint64SSZ(m.Attnets)...)
	if m.Syncnets == nil {
		return ret, nil
	}
	ret = append(ret, ssz.Uint64SSZ(*m.Syncnets)...)

	return ret, nil
}

func (m *Metadata) DecodeSSZ(buf []byte) error {
	m.SeqNumber = ssz.UnmarshalUint64SSZ(buf)
	m.Attnets = ssz.UnmarshalUint64SSZ(buf[8:])
	if len(buf) < 24 {
		return nil
	}
	m.Syncnets = new(uint64)
	*m.Syncnets = ssz.UnmarshalUint64SSZ(buf[16:])
	return nil
}

func (m *Metadata) EncodingSizeSSZ() (ret int) {
	ret = common.BlockNumberLength * 2
	if m.Syncnets != nil {
		ret += 8
	}
	return
}

func (m *Metadata) DecodeSSZWithVersion(buf []byte, _ int) error {
	return m.DecodeSSZ(buf)
}

// Ping is a test P2P message, used to test out liveness of our peer/signaling disconnection.
type Ping struct {
	Id uint64
}

func (p *Ping) EncodeSSZ(buf []byte) ([]byte, error) {
	return append(buf, ssz.Uint64SSZ(p.Id)...), nil
}

func (p *Ping) DecodeSSZ(buf []byte) error {
	p.Id = ssz.UnmarshalUint64SSZ(buf)
	return nil
}

func (p *Ping) EncodingSizeSSZ() int {
	return common.BlockNumberLength
}

func (p *Ping) DecodeSSZWithVersion(buf []byte, _ int) error {
	return p.DecodeSSZ(buf)
}

// P2P Message for bootstrap
type SingleRoot struct {
	Root [32]byte
}

func (s *SingleRoot) EncodeSSZ(buf []byte) ([]byte, error) {
	return append(buf, s.Root[:]...), nil
}

func (s *SingleRoot) DecodeSSZ(buf []byte) error {
	copy(s.Root[:], buf)
	return nil
}

func (s *SingleRoot) EncodingSizeSSZ() int {
	return length.Hash
}

func (s *SingleRoot) DecodeSSZWithVersion(buf []byte, _ int) error {
	return s.DecodeSSZ(buf)
}

func (*SingleRoot) Clone() clonable.Clonable {
	return &SingleRoot{}
}

/*
 * LightClientUpdatesByRangeRequest that helps syncing to chain tip from a past point.
 * It takes the Period of the starting update and the amount of updates we want (MAX: 128).
 */
type LightClientUpdatesByRangeRequest struct {
	Period uint64 `ssz:"true"`
	Count  uint64 `ssz:"true"`
}

func (*LightClientUpdatesByRangeRequest) Clone() clonable.Clonable {
	return &LightClientUpdatesByRangeRequest{}
}

func (l *LightClientUpdatesByRangeRequest) DecodeSSZWithVersion(buf []byte, _ int) error {
	return l.DecodeSSZ(buf)
}

func (l *LightClientUpdatesByRangeRequest) EncodeSSZ(buf []byte) ([]byte, error) {
	return ssz.Encode(l, buf)
}

func (l *LightClientUpdatesByRangeRequest) DecodeSSZ(buf []byte) error {
	l.Period = ssz.UnmarshalUint64SSZ(buf)
	l.Count = ssz.UnmarshalUint64SSZ(buf[8:])
	return nil
}

func (l *LightClientUpdatesByRangeRequest) EncodingSizeSSZ() int {
	return 2 * common.BlockNumberLength
}

/*
 * BeaconBlocksByRangeRequest is the request for getting a range of blocks.
 */
type BeaconBlocksByRangeRequest struct {
	StartSlot uint64 `ssz:"true"`
	Count     uint64 `ssz:"true"`
	Step      uint64 `ssz:"true"`
}

func (b *BeaconBlocksByRangeRequest) EncodeSSZ(buf []byte) ([]byte, error) {
	return ssz.Encode(b, buf)
}

func (b *BeaconBlocksByRangeRequest) DecodeSSZ(buf []byte) error {
	b.StartSlot = ssz.UnmarshalUint64SSZ(buf)
	b.Count = ssz.UnmarshalUint64SSZ(buf[8:])
	b.Step = ssz.UnmarshalUint64SSZ(buf[16:])
	return nil
}

func (b *BeaconBlocksByRangeRequest) DecodeSSZWithVersion(buf []byte, _ int) error {
	return b.DecodeSSZ(buf)
}

func (b *BeaconBlocksByRangeRequest) EncodingSizeSSZ() int {
	return 3 * common.BlockNumberLength
}

func (*BeaconBlocksByRangeRequest) Clone() clonable.Clonable {
	return &BeaconBlocksByRangeRequest{}
}

/*
 * Status is a P2P Message we exchange when connecting to a new Peer.
 * It contains network information about the other peer and if mismatching we drop it.
 */
type Status struct {
	ForkDigest     [4]byte  `ssz:"true"`
	FinalizedRoot  [32]byte `ssz:"true"`
	FinalizedEpoch uint64   `ssz:"true"`
	HeadRoot       [32]byte `ssz:"true"`
	HeadSlot       uint64   `ssz:"true"`
}

func (s *Status) EncodeSSZ(buf []byte) ([]byte, error) {
	return ssz.Encode(s, buf)
}

func (s *Status) DecodeSSZ(buf []byte) error {
	copy(s.ForkDigest[:], buf)
	copy(s.FinalizedRoot[:], buf[4:])
	s.FinalizedEpoch = ssz.UnmarshalUint64SSZ(buf[36:])
	copy(s.HeadRoot[:], buf[44:])
	s.HeadSlot = ssz.UnmarshalUint64SSZ(buf[76:])
	return nil
}

func (s *Status) DecodeSSZWithVersion(buf []byte, _ int) error {
	return s.DecodeSSZ(buf)
}

func (s *Status) EncodingSizeSSZ() int {
	return 84
}
