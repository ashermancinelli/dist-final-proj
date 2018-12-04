

package main

import (
	"log"
	"net"
	"os"
	"io"
	"flag"
	"fmt"
)

type IncrementalMessage struct {
	bytes uint64
}

func CoordinateStreams(con net.Conn) {
	c := make(chan IncrementalMessage)

	handleGameInput := func(r io.ReadCloser, w io.WriteCloser) {
		defer r.Close()
		defer w.Close()

		n, _ := io.Copy(w, r)




	}

	handleUserInput := func(r io.ReadCloser, w io.WriteCloser) {
		defer r.Close()
		defer w.Close()

		n, _ := io.Copy(w, r)

		input := fmt.Sprint(w)
		log.Println(input)
		c <- IncrementalMessage{bytes: uint64(n)}

	}

	go handleUserInput(con, os.Stdout)
	go handleUserInput(os.Stdin, con)

	p := <-c
	log.Printf("[%s]: Connection closed by remote. %d bytes received.\n", con.RemoteAddr(), p.bytes)
	p = <-c
	log.Printf("[%s]: Connection closed locally. %d bytes sent.\n", con.RemoteAddr(), p.bytes)

}

func StartClient(host string, port string) {
	con, err := net.Dial("tcp", host + port)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Connected to ", host + port)
	CoordinateStreams(con)
}


func main() {

	var host, port string
	flag.StringVar(&host, "host", "", "Remote host to connect to")
	flag.StringVar(&port, "port", ":8080", "Port of remote host")
	flag.Parse()

	if host == "" {
		flag.Usage()
	}

	StartClient(host, port)

}
