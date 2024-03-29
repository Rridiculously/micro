package net

import (
	"encoding/binary"
	"net"
	"time"
)

//func Connect(network, address string) error {
//	conn, err := net.DialTimeout(network, address, 10*time.Second)
//	if err != nil {
//		return err
//	}
//	defer func() {
//		_ = conn.Close()
//	}()
//	for {
//		_, err := conn.Write([]byte("hello"))
//		if err != nil {
//			return err
//		}
//		time.Sleep(time.Second)
//		buf := make([]byte, 1024)
//		_, err = conn.Read(buf)
//		if err != nil {
//			return err
//		}
//		fmt.Println(string(buf))
//	}
//}

type Client struct {
	network string
	address string
}

func (c *Client) Send(data string) (string, error) {
	conn, err := net.DialTimeout(c.network, c.address, 10*time.Second)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = conn.Close()
	}()
	reqLen := len(data)
	//构建响应数据
	//data=respLen + respData
	req := make([]byte, numOfLengthBytes+reqLen)

	//把数据写进去前八个字节
	binary.BigEndian.PutUint64(req[:numOfLengthBytes], uint64(reqLen))
	copy(req[numOfLengthBytes:], data)
	_, err = conn.Write(req)

	lenBs := make([]byte, numOfLengthBytes)
	_, err = conn.Read(lenBs)

	if err != nil {
		return "", err
	}
	//响应有多长？
	length := binary.BigEndian.Uint64(lenBs)
	respBs := make([]byte, length)
	_, err = conn.Read(respBs)

	if err != nil {
		return "", err
	}
	return string(respBs), nil
}
