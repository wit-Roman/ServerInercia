package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"sync"
	"time"

	"server/constant"
)

func runListen() {
	listener, err := net.Listen("tcp", constant.Address)
	if err != nil {
		log.Fatal(err)
	}

	serverHandler := &TserverHandler{
		subscribers: make(map[*Tsubscriber]struct{}), // TODO: писать не в ключ или читать из msgs
	}
	//cs.serveMux.HandleFunc("/subscribe", cs.subscribeHandler)
	serverHandler.serveMux.HandleFunc("/", connectHandler)
	log.Println("Запуск сервера ", constant.Address)

	const timePeriod = time.Second * 10
	server := &http.Server{
		Handler:      serverHandler,
		ReadTimeout:  timePeriod,
		WriteTimeout: timePeriod,
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Serve(listener)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	select {
	case err := <-errCh:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigCh:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timePeriod)
	defer cancel()

	log.Println("Выкл сервера ", constant.Address)

	server.Shutdown(ctx)
}

type TserverHandler struct {
	subscribers   map[*Tsubscriber]struct{}
	serveMux      http.ServeMux
	subscribersMu sync.Mutex
}
type Tsubscriber struct {
	msgs      chan []byte
	closeSlow func()
}

func (cs *TserverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cs.serveMux.ServeHTTP(w, r)
}
