package main

import netchan "../"
import "log"
import "io"

func main() {
	ln, err := netchan.Listen("tcp", ":9999")

	if err != nil {
		log.Fatalf("Error binding to socket: %q", err)
	}

	log.Println("Listening on :9999. Stop with [CTRL] + [c]")

	for {
		select {
		case c, ok := <-ln.Accept:
			if !ok {
				log.Println("Channel was closed")
				return
			}

			io.Copy(c, c)
			c.Close()
		}
	}
}
