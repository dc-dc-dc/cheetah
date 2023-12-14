package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dc-dc-dc/cheetah/cmd/web_server/ws"
	"github.com/dc-dc-dc/cheetah/market"
	"github.com/dc-dc-dc/cheetah/market/csv"
)

func main() {

	wsManager := ws.NewWebSocketManager(context.Background())
	defer wsManager.Close()
	wsManager.RegisterHandler("market:search", func(ctx context.Context, payload interface{}, client *ws.WebSocketClient) error {
		req, ok := payload.(map[string]interface{})
		if !ok {
			return fmt.Errorf("could not cast to dict")
		}
		fmt.Printf("got search request for %+v\n", req["symbol"])
		symbol, ok := req["symbol"]
		if !ok {
			return fmt.Errorf("could not get symbol from req")
		}
		producer := csv.NewYFinanceProducer(symbol.(string), market.Interval1Day, time.Now().Add(time.Hour*24*-365), time.Now())
		out := make(chan market.MarketLine, 1)
		err := producer.Produce(ctx, out)
		for err == nil {
			if err := client.SendMessage(ws.MessageWrapper{Type: "market:receive", Payload: <-out}); err != nil {
				return err
			}
			err = producer.Produce(ctx, out)
		}
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Printf("[error] got error=%s\n", err)
			}
		}
		return nil
	})
	router := http.NewServeMux()
	router.HandleFunc("/ws", wsManager.HttpHandler)
	router.Handle("/", http.FileServer(http.Dir("./static")))
	go func() {
		if err := http.ListenAndServe(":8080", router); err != nil {
			panic(err)
		}
	}()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
