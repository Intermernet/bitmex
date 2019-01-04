package websocket

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestConnectMaster(t *testing.T) {
	conn, err := Connect("www.bitmex.com")
	if err != nil {
		t.Errorf("%v\n", err)
	}
	if conn == nil {
		t.Error("No connect to ws")
	}
}

func TestConnectDev(t *testing.T) {
	conn, err := Connect("testnet.bitmex.com")
	if err != nil {
		t.Errorf("%v\n", err)
	}
	if conn == nil {
		t.Error("No connect to testnet ws")
	}
}

func TestConnectFail(t *testing.T) {
	_, err := Connect("")
	if err == nil {
		t.Error("expected error, got nil")
	}

}

func TestWorkerReadMessages(t *testing.T) {
	chReaderMessage := make(chan []byte)
	conn, err := Connect("testnet.bitmex.com")
	if err != nil {
		t.Errorf("%v\n", err)
	}
	go ReadFromWSToChannel(conn, chReaderMessage)
	message := <-chReaderMessage
	if message == nil {
		t.Error("Empty message")
	}
	close(chReaderMessage)
}

func TestWorkerWriteMessages(t *testing.T) {

	conn, err := Connect("testnet.bitmex.com")
	if err != nil {
		t.Errorf("%v\n", err)
	}

	// Read
	chReadFromWS := make(chan []byte, 10)
	go ReadFromWSToChannel(conn, chReadFromWS)

	// Write
	chWriteToWS := make(chan interface{}, 10)
	go WriteFromChannelToWS(conn, chWriteToWS)

	// Send ping
	chWriteToWS <- []byte(`ping`)

	// Read first response message
	message := <-chReadFromWS
	if !strings.Contains(string(message), "Welcome to the BitMEX") {
		fmt.Println(string(message))
		t.Error("No welcome message")
	}

	// Read second response message
	message = <-chReadFromWS
	if !strings.Contains(string(message), "pong") {
		fmt.Println(string(message))
		t.Error("No pong message")
	}

	time.Sleep(1 * time.Second)
}
