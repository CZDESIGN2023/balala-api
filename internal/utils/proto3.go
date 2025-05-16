package utils

import (
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func ConvertInterfaceToAnyProtoBean(v interface{}) (*any.Any, error) {
	anyValue := &any.Any{}
	bytes, _ := json.Marshal(v)
	bytesValue := &wrappers.BytesValue{
		Value: bytes,
	}
	err := anypb.MarshalFrom(anyValue, bytesValue, proto.MarshalOptions{})
	return anyValue, err
}

func Pbany(v interface{}) (*anypb.Any, error) {
	pv, ok := v.(proto.Message)
	if !ok {
		return &anypb.Any{}, fmt.Errorf("%v is not proto.Message", pv)
	}
	return anypb.New(pv)
}
