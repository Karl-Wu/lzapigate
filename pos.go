package main

import (
	"log"
	"net"
	"time"
)

var (
	positionConn *net.UDPConn
)

func listenPositionRequest() (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", posReqAddr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}

func servePositionReq() error {
	positionConn, _ = listenPositionRequest()
	for appStop == false {
		buffer := make([]byte, maxDatagramSize)
		numBytes, src, err := positionConn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("(pos req)ReadFromUDP failed:", err)
			time.Sleep(1 * time.Second)
			positionConn, _ = listenPositionRequest()
			continue
		}
		posReqHandler(src, buffer, numBytes)
	}

	return nil
}

func posReqAck() {

}

func posReqHandler(src *net.UDPAddr, buf []byte, len int) {
	DEBUG.Println(POS, "receive position request", src, len, buf)
}
