package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func handleGameString(str string) []byte {
	commands := strings.Split(str, ";")
	var finalValue string
	switch commands[0] {
	case "attack":
		finalValue = fmt.Sprintf("I am attacking %s", commands[1])
		finalValue = fmt.Sprint(finalValue + " doing " + commands[2] + "damage.\n")
	default:
		finalValue = "This is not an attack\n"
	}
	return []byte(finalValue)
}

func handleInputString(str string) []byte {
	// rVal := fmt.Sprintf("attack;%s\n", str)
	return []byte(str)
}

func copyToString(r io.Reader) (res string, err error, n int64) {
	var sb strings.Builder
	if n, err = io.Copy(&sb, r); err == nil {
		res = sb.String()
	}
	return
}

func streamCpy(src io.Reader, dst io.Writer, isOutgoing bool) <-chan int {
	buf := make([]byte, 1024)
	sync := make(chan int)

	go func() {
		defer func() {
			if con, ok := dst.(net.Conn); ok {
				con.Close()
				log.Printf("Con from %v closed\n", con.RemoteAddr())
			}
			sync <- 0
		}()

		for {

			nBytes, err := src.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("Read error: %s\n", err)
				}
				break
			}

			if !isOutgoing {
				_, err = dst.Write(buf[0:nBytes])
			} else {
				str := string(buf[0:nBytes])
				data := handleGameString(str)
				_, err = dst.Write(data)
			}
			if err != nil {
				log.Fatalf("Write error: %s\n", err)
			}
		}
	}()

	return sync

}

func HandleCons(con net.Conn) {

	stdoutChan := streamCpy(con, os.Stdout, false)
	remoteChan := streamCpy(os.Stdin, con, true)

	select {
	case <-stdoutChan:
		log.Println("Remote connection broken.")
	case <-remoteChan:
		log.Println("Local connection broken.")
	}
}

func StartClient(host string, port string) {
	con, err := net.Dial("tcp", host+port)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Connected to ", host+port)
	HandleCons(con)
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
