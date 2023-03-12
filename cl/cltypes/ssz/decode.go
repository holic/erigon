package ssz

import (
	"reflect"

	"github.com/ledgerwatch/erigon-lib/common"
	libcommon "github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon-lib/common/length"
)

// This package is a working progress. only base functionality is supported. no more no less.

// Decodes just decodes a specific struct.
func Decode(x any, buf []byte) error {
	// Get values in struct.
	reflValue := reflect.ValueOf(x)
	reflType := reflect.TypeOf(x)

	// Iterate over all the fields.
	_, err := decodeValue(reflValue, reflType, buf)
	return err

}

func decodeValue(reflValue reflect.Value, reflType reflect.Type, buf []byte) (int, error) {
	if reflValue.Kind() == reflect.Ptr {
		reflValue = reflValue.Elem()
		reflType = reflType.Elem()
	}
	pos := 0
	for i := 0; i < reflValue.NumField(); i++ {
		// Process each field.
		field := reflValue.Field(i)
		if !field.CanInterface() || reflType.Field(i).Tag.Get(TagSSZ) != TagTrueFlag {
			continue
		}
		switch fieldVal := field.Interface().(type) {
		// Base field can just be appended
		case uint64:
			if len(buf) < 8 {
				return 0, ErrLowBufferSize
			}
			num := UnmarshalUint64SSZ(buf[pos:])
			field.Set(reflect.ValueOf(num))
			pos += 8
		case libcommon.Hash:
			if len(buf) < length.Hash {
				return 0, ErrLowBufferSize
			}
			field.Set(reflect.ValueOf(common.BytesToHash(buf[pos : pos+length.Hash])))
			pos += length.Hash
		// Will be fixed in the future.
		case [32]byte:
			if len(buf) < length.Hash {
				return 0, ErrLowBufferSize
			}
			field.Set(reflect.ValueOf(common.BytesToHash(buf[pos : pos+length.Hash])))
			pos += length.Hash
		case [48]byte:
			var val [48]byte
			if len(buf) < 48 {
				return 0, ErrLowBufferSize
			}
			copy(val[:], buf[pos:])
			field.Set(reflect.ValueOf(val))
			pos += 48
		case [96]byte:
			var val [96]byte
			if len(buf) < 96 {
				return 0, ErrLowBufferSize
			}
			copy(val[:], buf[pos:])
			field.Set(reflect.ValueOf(val))
			pos += 96
		case [4]byte:
			var val [4]byte
			if len(buf) < 4 {
				return 0, ErrLowBufferSize
			}
			copy(val[:], buf[pos:])
			field.Set(reflect.ValueOf(val))
			pos += 4
		case libcommon.Address:
			if len(buf) < length.Addr {
				return 0, ErrLowBufferSize
			}
			field.Set(reflect.ValueOf(common.BytesToAddress(buf[pos : pos+length.Addr])))
			pos += length.Addr
		case bool:
			if len(buf) == 0 {
				return 0, ErrLowBufferSize
			}
			if buf[pos] == 0x01 {
				field.Set(reflect.ValueOf(true))
			} else {
				field.Set(reflect.ValueOf(false))
			}
			pos++
		default:
			// Create pointer to default.
			t := reflect.TypeOf(fieldVal)
			v := reflect.New(t.Elem())
			n, err := decodeValue(v, t, buf[pos:])
			if err != nil {
				return 0, err
			}
			field.Set(v)
			pos += n
		}
	}
	return pos, nil
}
