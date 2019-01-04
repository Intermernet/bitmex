# Bitmex API
Forked from orignal code at github.com/vmpartner/bitmex

Packages for work with bitmex rest and websocket API on golang.  
Target of this packages make easy access to bitmex API including testnet platform.


## Usage
Please see full example in main.go

####  REST
```
// Load config
cfg, err := config.LoadConfig("config.json")
if err != nil {
    log.Fatal(err)
}
ctx := rest.MakeContext(cfg.Key, cfg.Secret, cfg.Host)

// Get wallet
wallet, response, err := rest.GetWallet(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Status: %v, wallet amount: %v\n", response.StatusCode, wallet.Amount)

// Place order
params := map[string]interface{}{
    "side":     "Buy",
    "symbol":   "XBTUSD",
    "ordType":  "Limit",
    "orderQty": 1,
    "price":    9000,
    "clOrdID":  "MyUniqID_123",
    "execInst": "ParticipateDoNotInitiate",
}
order, response, err := rest.NewOrder(ctx, params)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Order: %+v, Response: %+v\n", order, response)
```

#### Websocket
```
// Load config
cfg, err := config.LoadConfig("config.json")
if err != nil {
    log.Fatal(err)
}

// Connect to WS
conn, err := websocket.Connect(cfg.Host)
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

// Listen read WS
chReadFromWS := make(chan []byte, 100)
go func() {
    err != websocket.ReadFromWSToChannel(conn, chReadFromWS)
    if err != nil {
        log.Fatal(err)
    }
}()

// Listen write WS
chWriteToWS := make(chan interface{}, 100)
go func() {
    err != websocket.WriteFromChannelToWS(conn, chWriteToWS)
    if err != nil {
        log.Fatal(err)
    }
}()

// Authorize
auth, err != websocket.GetAuthMessage(cfg.Key, cfg.Secret)
if err != nil {
    log.Fatal(err)
}
chWriteToWS <- auth

// Listen
go func() {
    for {
        message := <-chReadFromWS
        res, err := bitmex.DecodeMessage(message)
        if err != nil {
            log.Fatal(err)
        }

        // Business logic
        switch res.Table {
        case "orderBookL2":
            if res.Action == "partial" {
                // Update table
            } else {
                // Update row
            }
        case "order":
            if res.Action == "partial" {
                // Update table
            } else {
                // Update row
            }
        case "position":
            if res.Action == "partial" {
                // Update table
            } else {
                // Update row
            }
        }
    }
}()

```

## Example
Example of usage look in main.go

## More
This is forked from github.com/vmpartner/bitmex
Support the original author! Thank you!
```
eth: 0x3e9b92625c49Bfd41CCa371D1e4A1f0d4c25B6fC
btc: 35XDoFSA8QeM26EnCyhQPTMBZm4S1DvncE
```
