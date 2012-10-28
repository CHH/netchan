package netchan

import netchan "./"
import "testing"
import "fmt"
import "net"
import "bufio"
import "io"

func ExampleListen() {
	ln, err := netchan.Listen("tcp", ":9999")

	if err != nil {
		panic(err)
	}

	for {
		select {
		case conn := <-ln.Accept:
			io.Copy(conn, conn)
			conn.Close()
		}
	}
}

func TestChan(t *testing.T) {
	ln, err := netchan.Listen("tcp", ":9999")

	defer func() {
		ln.Quit <- true
	}()

	if err != nil {
		t.Errorf("Error while listening: %q", err)
	}

	go func() {
		select {
		case c := <- ln.Accept:
			if _, err := fmt.Fprintln(c, "Hello World"); err != nil {
				t.Errorf("Failed writing to connection: %q", err)
			}
			c.Close()
		}
	}()

	conn, err := net.Dial("tcp", ":9999")

	if err != nil {
		t.Errorf("Error while connecting to server: %q", err)
	}

	r := bufio.NewReader(conn)

	msg, _ := r.ReadString('\n')

	if msg != "Hello World\n" {
		t.Errorf("Expected %q, Got %q", "Hello World\n", msg)
	}
}
