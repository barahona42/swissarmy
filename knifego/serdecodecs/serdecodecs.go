/*
serdecodecs provides utilities for serialization and de-serialization of custom types.
*/
package serdecodecs

import (
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TimestampCodec: implements bson.Value{Encoder,Decoder} interface for serde of `*timestamppb.Timestamp`.
type TimestampCodec struct{}

// EncodeValue: value encoder.
func (codec TimestampCodec) EncodeValue(_ bson.EncodeContext, vw bson.ValueWriter, val reflect.Value) error {
	if !val.IsValid() {
		return bson.ValueEncoderError{
			Name:     "TimestampEncodeValue",
			Kinds:    []reflect.Kind{reflect.TypeOf(timestamppb.Timestamp{}).Kind()},
			Received: val,
		}
	}
	t, ok := val.Interface().(*timestamppb.Timestamp)
	if !ok {
		return bson.ValueEncoderError{
			Name:     "TimestampEncodeType",
			Kinds:    []reflect.Kind{reflect.TypeOf(timestamppb.Timestamp{}).Kind()},
			Received: val,
		}
	}
	// timestamppb.Timestamp starts at Unix start
	return vw.WriteDateTime(t.AsTime().UnixMilli())
}

// DecodeValue: value decoder.
func (codec TimestampCodec) DecodeValue(_ bson.DecodeContext, vr bson.ValueReader, val reflect.Value) error {
	if !val.IsValid() {
		return bson.ValueDecoderError{
			Name:     "TimestampDecodeValue",
			Kinds:    []reflect.Kind{reflect.Kind(bson.TypeDateTime)},
			Received: val,
		}
	}
	if vr.Type() != bson.TypeDateTime {
		return bson.ValueDecoderError{
			Name:     "TimestampDecodeType",
			Kinds:    []reflect.Kind{reflect.Kind(bson.TypeDateTime)},
			Received: val,
		}
	}
	t, err := vr.ReadDateTime()
	if err != nil {
		return fmt.Errorf("read datetime: %w", err)
	}
	val.Set(reflect.ValueOf(timestamppb.New(time.UnixMilli(t))))
	return nil
}

// NewRegistry: creates a new bson registry that adds codecs defined in this package.
func NewRegistry() *bson.Registry {
	r := bson.NewRegistry() // not actually empty as docstring says. includes "default" codecs
	RegisterBSON(r)
	return r
}

func RegisterBSON(r *bson.Registry) {
	t := reflect.TypeOf(&timestamppb.Timestamp{}) // pointer is important to prevent copying the sync.Mutex value
	r.RegisterTypeEncoder(t, TimestampCodec{})
	r.RegisterTypeDecoder(t, TimestampCodec{})
}
