package protohelpers

/*
prov
*/

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
)

var (
	ErrNoOptionsMatched   error = errors.New("no options matched")
	ErrExtensionValueType error = errors.New("extension value is not of expected type")
)

// optionsreceiver: is an interface for retrieving extensions for a given protobuf-generated type.
// as of this writing, any enum type should implement these functions OOB and should be accepted by GetStringNameExtension
// as an argument
type optionsreceiver interface {
	Descriptor() protoreflect.EnumDescriptor
	Number() protoreflect.EnumNumber
}

// getValueOfExtensionType: wraps reflection work to resolve the target option
func getValueOfExtensionType(desc protoreflect.EnumDescriptor, n protoreflect.EnumNumber, ext protoreflect.ExtensionType) (string, error) {
	vd := desc.Values()
	vdn := vd.ByNumber(n)
	if vdn == nil {
		return "", fmt.Errorf("no descriptor for number '%d'", n)
	}
	opts := vdn.Options()
	if opts == nil {
		return "", ErrNoOptionsMatched
	}
	if s, ok := proto.GetExtension(opts, ext).(string); ok {
		return s, nil
	}
	return "", ErrExtensionValueType
}

// GetExtensionValue: retrieves the value of the specified extension for the provided receiver.
func GetExtensionValue(ext *protoimpl.ExtensionInfo, v optionsreceiver) (string, error) {
	return getValueOfExtensionType(v.Descriptor(), v.Number(), ext)
}
