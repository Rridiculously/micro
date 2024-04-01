package rpc

import (
	"encoding/binary"
	"net"
)

func ReadMsg(conn net.Conn) ([]byte, error) {
	lenBs := make([]byte, numOfLengthBytes)
	_, err := conn.Read(lenBs)
	if err != nil {
		return nil, err
	}
	//响应有多长？
	length := binary.BigEndian.Uint64(lenBs)
	data := make([]byte, length)
	_, err = conn.Read(data)
	return data, err
}

func EncodeMsg(data []byte) []byte {
	respLen := len(data)
	//构建响应数据
	//打他=respLen + respData
	res := make([]byte, numOfLengthBytes+respLen)
	//把数据写进去前八个字节
	binary.BigEndian.PutUint64(res[:numOfLengthBytes], uint64(respLen))
	copy(res[numOfLengthBytes:], data)
	return res
}
