package main

import (
	"log"
	"net"
	"time"
)

var (
	pwrConn *net.UDPConn
)

//test command:  echo "dddf" | nc -u 228.8.8.8 11000

func listenPowerOn() (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", powerOnAddr)
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

func servePowerOn() error {
	pwrConn, _ = listenPowerOn()
	for appStop == false {
		buffer := make([]byte, maxDatagramSize)
		numBytes, src, err := pwrConn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("(power on)ReadFromUDP failed:", err)
			time.Sleep(1 * time.Second)
			pwrConn, _ = listenPowerOn()
			continue
		}
		powerOnCtrlHandler(src, buffer, numBytes)
	}

	return nil
}

func powerOnCtrlAck() {
	dstAddr := &net.UDPAddr{
		Port: int(myID),
		IP:   net.ParseIP(commanderIP),
	}
	n, err := pwrConn.WriteToUDP([]byte("power on ack"), dstAddr)
	DEBUG.Println(PWR, "send request ack:", n, err)
}

func powerOnCtrlHandler(src *net.UDPAddr, buf []byte, len int) {
	DEBUG.Println(PWR, "receive power on request", src, len, buf[0:len])
	powerOnCtrlAck()
}
