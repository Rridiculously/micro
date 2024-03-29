package net

import (
	"encoding/binary"
	"net"
)

//	func Serve(network, address string) error {
//		listener, err := net.Listen(network, address)
//		if err != nil {
//			return err
//		}
//		defer func() {
//			_ = listener.Close()
//		}()
//		for {
//			conn, err := listener.Accept()
//			if err != nil {
//				return err
//			}
//			go func() {
//				if err := handleConn(conn); err != nil {
//					_ = conn.Close()
//				}
//			}()
//		}
//	}
const numOfLengthBytes = 8

func (s *Server) handleConn(conn net.Conn) error {
	for {
		//lenBs长度字段的字节表示
		lenBs := make([]byte, numOfLengthBytes)
		_, err := conn.Read(lenBs)

		if err != nil {
			return err
		}
		//消息有多长？
		length := binary.BigEndian.Uint64(lenBs)
		reqBs := make([]byte, length)
		_, err = conn.Read(reqBs)

		if err != nil {
			return err
		}
		respData := handleMsg(reqBs)
		respLen := len(respData)
		//构建响应数据
		//打他=respLen + respData
		res := make([]byte, numOfLengthBytes+respLen)

		//把数据写进去前八个字节
		binary.BigEndian.PutUint64(res[:numOfLengthBytes], uint64(respLen))
		copy(res[numOfLengthBytes:], respData)
		_, err = conn.Write(res)
		if err != nil {
			return err
		}

	}
}

func handleMsg(req []byte) []byte {
	res := make([]byte, 2*len(req))
	copy(res[:len(req)], req)
	copy(res[len(req):], req)
	return res
}

type Server struct {
}

func (s *Server) Start(network, address string) error {
	listener, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	defer func() {
		_ = listener.Close()
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go func() {
			if err := s.handleConn(conn); err != nil {
				_ = conn.Close()
			}
		}()
	}
}
