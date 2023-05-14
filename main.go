package main

import (
	"log"
	"server/world"
)

func main() {
	log.SetFlags(0)

	go runListen()

	/*go func() {
		err := runListen()
		if err != nil {
			log.Fatal(err)
		}
	}()*/

	world.Create()
}
