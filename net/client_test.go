package net

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

//	func TestConnect(t *testing.T) {
//		go func() {
//			err := Serve("tcp", ":8080")
//			t.Log(err)
//		}()
//		time.Sleep(time.Second * 5)
//		err := Connect("tcp", ":8080")
//		if err != nil {
//			t.Error(err)
//		}
//	}
func TestClient_Send(t *testing.T) {
	server := &Server{}
	go func() {
		err := server.Start("tcp", ":8080")
		t.Log(err)
	}()
	time.Sleep(time.Second * 3)
	client := &Client{
		network: "tcp",
		address: ":8080",
	}
	resp, err := client.Send("hello")
	require.NoError(t, err)
	assert.Equal(t, "hellohello", resp)
}
