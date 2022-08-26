package main

import (
	"log"
	"net"
)

func main() {
	s := newServer()
	go s.run()
	//msg(">Hos Geldiniz")
	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("Server başlatılamadı: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("Server başlatıldı :5050")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Giris yapilamadi", err.Error())
			continue
		}

		c := s.newClient(conn)
		go c.readInput()
	}
}
