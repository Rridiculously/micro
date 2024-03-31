package rpc

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
	"reflect"
	"time"
)

// InitClientProxy 初始化客户端代理,为GetById类的函数类型的字段赋值
func InitClientProxy(service Service, addr string) error {
	client := NewClient(addr)
	return setFuncField(service, client)
}

func setFuncField(service Service, proxy Proxy) error {
	if service == nil {
		return errors.New("service is nil")
	}
	val := reflect.ValueOf(service)
	typ := val.Type()
	// 只支持指向结构体的一级指针
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return errors.New("service is not a pointer to a struct")
	}
	val = val.Elem()
	typ = typ.Elem()
	numFiled := typ.NumField()
	for i := 0; i < numFiled; i++ {
		fieldTyp := typ.Field(i)
		fieldVal := val.Field(i)
		if fieldVal.CanSet() {
			//本地调用捕捉到的地方
			fn := func(args []reflect.Value) (results []reflect.Value) {
				retVal := reflect.New(fieldTyp.Type.Out(0).Elem())
				ctx := args[0].Interface().(context.Context)
				reqData, err := json.Marshal(args[1].Interface())
				if err != nil {
					return []reflect.Value{retVal, reflect.ValueOf(err)}
				}

				req := &Request{
					ServiceName: service.Name(),
					MethodName:  fieldTyp.Name,
					Args:        reqData,
				}

				resp, err := proxy.Invoke(ctx, req)
				if err != nil {
					return []reflect.Value{retVal, reflect.ValueOf(err)}
				}
				err = json.Unmarshal(resp.Data, retVal.Interface())
				if err != nil {
					return nil
				}
				return []reflect.Value{retVal, reflect.Zero(reflect.TypeOf(new(error)).Elem())}
				// 创建与字段类型匹配的零值
			}
			//设值给 GetById
			fnVal := reflect.MakeFunc(fieldTyp.Type, fn)
			fieldVal.Set(fnVal)
		}

	}
	return nil
}

const numOfLengthBytes = 8

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
	}
}
func (c *Client) Invoke(ctx context.Context, req *Request) (*Response, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.Send(data)
	if err != nil {
		return nil, err
	}
	return &Response{
		Data: resp,
	}, nil
}
func (c *Client) Send(data []byte) ([]byte, error) {
	conn, err := net.DialTimeout("tcp", c.addr, 10*time.Second)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	//响应有多长？
	length := binary.BigEndian.Uint64(lenBs)
	respBs := make([]byte, length)
	_, err = conn.Read(respBs)

	if err != nil {
		return nil, err
	}
	return respBs, nil
}
