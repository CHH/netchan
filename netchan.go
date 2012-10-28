/*
	Wrapper which yields socket connections on a channel
*/
package netchan

import "net"

type Netchan struct {
	Accept chan net.Conn
	Quit chan bool
	Listener net.Listener
}

const NdefaultBacklog = 15

/*
	Replacement for net.Listen, which yields connections on a channel.
*/
func Listen(network string, addr string) (*Netchan, error) {
	ln, err := net.Listen(network, addr)

	if err != nil {
		return nil, err
	}

	return Chan(ln, NdefaultBacklog), nil
}

/*
	Yields all accepted connections to the net.Listener on the returned Channel.
*/
func Chan(ln net.Listener, backlog int) *Netchan {
	c := &Netchan{}
	c.Accept = make(chan net.Conn, backlog)
	c.Quit = make(chan bool, 1)
	c.Listener = ln

	go func() {
		select {
		case <-c.Quit:
			close(c.Accept)

			if err := ln.Close(); err != nil {
				panic(err)
			}

			c.Quit <- true
		}
	}()

	go func() {
		for {
			conn, err := ln.Accept()

			// An error means that the listener was closed, or another event
			// happened where we can't continue listening for connections.
			if err != nil {
				return
			}

			c.Accept <- conn
		}
	}()

	return c
}

