package rpc

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
	"reflect"
)

type Server struct {
	services map[string]Service
}

func NewServer() *Server {
	return &Server{
		services: make(map[string]Service, 16),
	}

}
func (s *Server) RegisterService(service Service) {
	s.services[service.Name()] = service
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
		respData, err := s.handleMsg(reqBs)
		if err != nil {
			//业务err
			return err
		}
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
func (s *Server) handleMsg(reqData []byte) ([]byte, error) {
	req := &Request{}
	err := json.Unmarshal(reqData, req)
	if err != nil {
		return nil, err
	}
	//业务调用
	service, ok := s.services[req.ServiceName]
	if !ok {
		return nil, errors.New("service not found")
	}
	val := reflect.ValueOf(service)
	method := val.MethodByName(req.MethodName)
	in := make([]reflect.Value, 2)
	in[0] = reflect.ValueOf(context.Background())
	inReq := reflect.New(method.Type().In(1).Elem())
	err = json.Unmarshal(req.Args, inReq.Interface())
	if err != nil {
		return nil, err
	}
	in[1] = inReq
	results := method.Call(in)
	if results[1].Interface() != nil {
		return nil, results[1].Interface().(error)
	}
	resp, err := json.Marshal(results[0].Interface())
	return resp, err
}
