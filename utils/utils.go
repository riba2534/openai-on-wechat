package utils

import (
	"encoding/base64"

	jsoniter "github.com/json-iterator/go"
)

func MarshalAnyToString(param interface{}) string {
	s, err := jsoniter.MarshalToString(param)
	if err != nil {
		return "{}"
	}
	return s
}

func MarshalAnyToByte(param interface{}) []byte {
	s, err := jsoniter.Marshal(param)
	if err != nil {
		return []byte{}
	}
	return s
}

func DecodeBase64(s string) []byte {
	decodeBytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return []byte{}
	}
	return decodeBytes
}
