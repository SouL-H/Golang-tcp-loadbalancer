package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

var (
	counter    int
	listenAddr = "localhost:8080"
	//Server list
	server = []string{
		"localhost:5001",
		"localhost:5002",
		"localhost:5003",
	}
)

func main() {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal("Failed to listen: %s", err)
	}

	defer listener.Close() //Close listener
	for {
		conn, err := listener.Accept() //Accept connection
		if err != nil {
			log.Printf("Failed to accept: %s", err)
		}

		backend := cooseBackend() //Choose backend
		fmt.Printf("Counter:%d bacnend:%s\n", counter, backend)
		go func() {
			err := proxy(backend, conn) //Proxy connection
			if err != nil {
				log.Printf("Failed to proxy: %v", err)
			}
		}()
	}

}

func proxy(backend string, c net.Conn) error {
	bc, err := net.Dial("tcp", backend) //Dial backend
	if err != nil {
		return fmt.Errorf("Failed to connect to backend: %s: %v", backend, err)
	}
	//c->bc
	go io.Copy(bc, c)
	//bc->c
	go io.Copy(c, bc)
	return nil
}
func cooseBackend() string {
	//Server selection in order
	s := server[counter%len(server)]
	counter++
	return s
}
