package merkle_tree

import (
	"reflect"

	libcommon "github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon/cl/cltypes/ssz"
)

// HashTreeRoot computes the hash of an object according to ETH 2.0 specs.
func HashTreeRoot(obj any) ([32]byte, error) {
	// Get values in struct.
	reflValue := reflect.ValueOf(obj)
	reflType := reflect.TypeOf(obj)
	if reflValue.Kind() == reflect.Ptr {
		reflValue = reflValue.Elem()
		reflType = reflType.Elem()
	}
	// We need to accumulate
	var hashes [][32]byte
	// Iterate over all the fields.
	for i := 0; i < reflValue.NumField(); i++ {
		// Process each field.
		field := reflValue.Field(i)
		if !field.CanInterface() || reflType.Field(i).Tag.Get(ssz.TagSSZ) != ssz.TagTrueFlag {
			continue
		}
		switch fieldVal := field.Interface().(type) {
		// Base field can just be appended
		case uint64:
			hashes = append(hashes, Uint64Root(fieldVal))
		case libcommon.Hash:
			hashes = append(hashes, fieldVal)
		// Will be fixed in the future.
		case [32]byte:
			hashes = append(hashes, fieldVal)
		case [48]byte:
			root, err := PublicKeyRoot(fieldVal)
			if err != nil {
				return [32]byte{}, err
			}
			hashes = append(hashes, root)
		case [96]byte:
			root, err := SignatureRoot(fieldVal)
			if err != nil {
				return [32]byte{}, err
			}
			hashes = append(hashes, root)
		case [4]byte:
			var root [32]byte
			copy(root[:], fieldVal[:])
			hashes = append(hashes, root)
		case libcommon.Address:
			var root [32]byte
			copy(root[:], fieldVal[:])
			hashes = append(hashes, root)
		case bool:
			hashes = append(hashes, BoolRoot(fieldVal))
		default:
			root, err := HashTreeRoot(fieldVal)
			if err != nil {
				return [32]byte{}, err
			}
			hashes = append(hashes, root)
		}
	}

	return ArraysRoot(hashes, extendToPowerOf2(len(hashes)))
}

func extendToPowerOf2(x int) (v uint64) {
	v = uint64(x)
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return
}
