package websocket

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Op   string        `json:"op,omitempty"`
	Args []interface{} `json:"args,omitempty"`
}

func (m *Message) AddArgument(argument string) {
	m.Args = append(m.Args, argument)
}

func Connect(host string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "wss", Host: host, Path: "/realtime"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func ReadFromWSToChannel(c *websocket.Conn, chRead chan<- []byte) error {
	for {
		_, message, err := c.ReadMessage()
		//fmt.Println("Read", string(message))
		if err != nil {
			return err
		}
		chRead <- message
	}
}

func WriteFromChannelToWS(c *websocket.Conn, chWrite <-chan interface{}) error {
	for {
		message := <-chWrite
		if reflect.TypeOf(message).String() == "websocket.Message" {
			var err error
			message, err = json.Marshal(message)
			if err != nil {
				return err
			}
		}
		err := c.WriteMessage(websocket.TextMessage, message.([]byte))
		if err != nil {
			return err
		}
	}
}

func GetAuthMessage(key string, secret string) (Message, error) {
	nonce := time.Now().Unix() + 412
	req := fmt.Sprintf("GET/realtime%d", nonce)
	sig := hmac.New(sha256.New, []byte(secret))
	_, err := sig.Write([]byte(req))
	if err != nil {
		return Message{}, err
	}
	signature := hex.EncodeToString(sig.Sum(nil))
	var msgKey []interface{}
	msgKey = append(msgKey, key)
	msgKey = append(msgKey, nonce)
	msgKey = append(msgKey, signature)

	return Message{"authKey", msgKey}, nil
}
