package core

import (
    "github.com/ugorji/go/codec"
)

var mh codec.MsgpackHandle
var jh codec.JsonHandle

func MsgUnpack(pack []byte, res *[]interface{}) error {
    err := codec.NewDecoderBytes(pack, &mh).Decode(res)
    return err
}

func MsgUnpackScalar(pack []byte, res *interface{}) error {
    err := codec.NewDecoderBytes(pack, &mh).Decode(res)
    return err
}

func MsgPack(pack interface{}) (res []byte) {
    codec.NewEncoderBytes(&res, &mh).Encode(pack)
    return res
}

func JsonUnpack(pack []byte, res interface{}) error {
    err := codec.NewDecoderBytes(pack, &jh).Decode(res)
    return err
}

func JsonPack(pack interface{}) (res []byte) {
    codec.NewEncoderBytes(&res, &jh).Encode(pack)
    return res
}

