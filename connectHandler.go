package main

import (
	"fmt"
	"net/http"
	"server/world"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func connectHandler(w http.ResponseWriter, r *http.Request) {
	if token, ok := r.Header["Bearer"]; ok {
		match, err := verifyToken("$argon2id$v=" + token[0])
		if !match || err != nil {
			fmt.Println("cannot connect ", err)
			return
		}

		player := world.CreateNewPlayer(token[0])

		ctx := r.Context()

		//upgrade http connection to websocket
		wsconn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true}) // TODO CORS
		if err != nil {
			return
		}

		defer wsconn.Close(websocket.StatusInternalError, "")

		var readContent interface{}
		// Read
		go func() {
			for {
				if err := wsjson.Read(ctx, wsconn, &readContent); err != nil {
					if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
						websocket.CloseStatus(err) == websocket.StatusGoingAway {
						return
					}
					fmt.Println("read error", err) //after 5 seconds I get either of:
					//failed to get reader: received close frame: status = StatusPolicyViolation and reason = "unexpected data message"
					//failed to read JSON message: failed to get reader: failed to read frame header: EOF
					break
				}
				if control, ok := readContent.(int8); ok {
					// TODO Ship{}
					player.InputCh <- control
				}
			}
		}()
		// Write

		//go func() {
		for {
			timeout := time.After(5 * time.Second)
			select {
			case <-timeout:
				if err := wsjson.Write(ctx, wsconn, time.Now().Format("2006-01-02T15:04:05Z07:00")); err != nil {
					fmt.Println("err writing to websocket: ", err)
					break
				}
			case command := <-player.OutputCh:
				fmt.Println(command)
				if err := wsjson.Write(ctx, wsconn, command); err != nil {
					fmt.Println("err writing to websocket: ", err)
					break
				}

			}
		}
		//}()
	}
}
