package utils

import (
	"encoding/binary"
	"errors"
	"github.com/golang/protobuf/proto"
	"unsafe"

)

func TypeToBytes(x interface{})([]byte,error){

	switch x.(type) {
	case string:
		return []byte(x.(string)),nil
	case int:
		return intToBytes(x.(int)),nil
	case int64:
	case uint:
	case uint64:
		return uint64ToBytes(x.(uint64)),nil
	case []byte:
		return x.([]byte),nil
	case proto.Message:
		return proto.Marshal(x.(proto.Message))

	}

	return nil,errors.New("type is not exist")
}
func uint64ToBytes(i uint64)[]byte{
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b
}

func intToBytes(data int)(ret []byte){
	 l := unsafe.Sizeof(data)
	ret = make([]byte, l)
	 tmp  := 0xff
	var index uint = 0
	for index=0; index<uint(l); index++{
		ret[index] = byte((tmp<<(index*8) & data)>>(index*8))
	}
	return ret
}