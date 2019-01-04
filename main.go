package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Intermernet/bitmex/bitmex"
	"github.com/Intermernet/bitmex/config"
	"github.com/Intermernet/bitmex/rest"
	"github.com/Intermernet/bitmex/websocket"
)

// Usage example
func main() {

	// Load config
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	ctx := rest.MakeContext(cfg.Key, cfg.Secret, cfg.Host, cfg.Timeout)

	// Get wallet
	w, response, err := rest.GetWallet(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Status: %v, wallet amount: %v\n", response.StatusCode, w.Amount)

	// Connect to WS
	conn, err := websocket.Connect(cfg.Host)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Listen read WS
	chReadFromWS := make(chan []byte, 100)
	go func() {
		err := websocket.ReadFromWSToChannel(conn, chReadFromWS)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Listen write WS
	chWriteToWS := make(chan interface{}, 100)

	go func() {
		err := websocket.WriteFromChannelToWS(conn, chWriteToWS)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Authorize
	auth, err := websocket.GetAuthMessage(cfg.Key, cfg.Secret)

	if err != nil {
		log.Fatal(err)
	}
	chWriteToWS <- auth
	// Read first response message
	message := <-chReadFromWS
	if !strings.Contains(string(message), "Welcome to the BitMEX") {
		fmt.Println(string(message))
		panic("No welcome message")
	}

	// Read auth response success
	message = <-chReadFromWS
	res, err := bitmex.DecodeMessage(message)
	if err != nil {
		log.Fatal(err)
	}
	if res.Success != true || res.Request.(map[string]interface{})["op"] != "authKey" {
		log.Fatal("No auth response success")
	}

	// Listen websocket before subscribe
	go func() {
		for {
			message := <-chReadFromWS
			res, err := bitmex.DecodeMessage(message)
			if err != nil {
				log.Fatal(err)
			}

			// Your logic here
			fmt.Printf("%+v\n", res)
		}
	}()

	// Subscribe
	messageWS := websocket.Message{Op: "subscribe"}
	messageWS.AddArgument("orderBookL2:XBTUSD")
	messageWS.AddArgument("order")
	messageWS.AddArgument("position")
	chWriteToWS <- messageWS

	// Loop forever
	select {}
}
