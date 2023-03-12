package ssz

import (
	"reflect"

	libcommon "github.com/ledgerwatch/erigon-lib/common"
)

var (
	TagSSZ      = "ssz"
	TagTrueFlag = "true"
)

// This package is a working progress. only base functionality is supported. no more no less.

// Encode just encodes a specific struct. it appends the encoding to the given bytes buffer.
func Encode(x any, buf []byte) (dst []byte, err error) {
	dst = buf
	// Get values in struct.
	reflValue := reflect.ValueOf(x)
	reflType := reflect.TypeOf(x)
	if reflValue.Kind() == reflect.Ptr {
		reflValue = reflValue.Elem()
		reflType = reflType.Elem()
	}
	// Iterate over all the fields.
	for i := 0; i < reflValue.NumField(); i++ {
		// Process each field.
		field := reflValue.Field(i)
		if !field.CanInterface() || reflType.Field(i).Tag.Get(TagSSZ) != TagTrueFlag {
			continue
		}
		switch fieldVal := field.Interface().(type) {
		// Base field can just be appended
		case uint64:
			dst = append(dst, Uint64SSZ(fieldVal)...)
		case libcommon.Hash:
			dst = append(dst, fieldVal[:]...)
		// Will be fixed in the future.
		case [32]byte:
			dst = append(dst, fieldVal[:]...)
		case [48]byte:
			dst = append(dst, fieldVal[:]...)
		case [96]byte:
			dst = append(dst, fieldVal[:]...)
		case [4]byte:
			dst = append(dst, fieldVal[:]...)
		case libcommon.Address:
			dst = append(dst, fieldVal[:]...)
		case bool:
			dst = append(dst, BoolSSZ(fieldVal))
		default:
			dst, err = Encode(fieldVal, dst)
			if err != nil {
				return nil, err
			}
		}
	}
	return dst, nil
}
