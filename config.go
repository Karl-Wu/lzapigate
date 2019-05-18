package main

import (
	"bytes"
	"encoding/binary"
	"net"
)

var (
	authURL    = "http://localhost:8080/login"
	graphqlURL = "http://localhost:8080/rf/graphql"

	username = "demo"
	password = "secret"
)

var (
	myID   = uint16(13001)
	myIP   = "172.19.9.135"
	myPort = uint16(10000)

	commanderIP   = "127.0.0.1" //"172.19.9.109"
	commanderPort = myID
	powerOnAddr   = "228.8.8.8:11000"
	timeSyncAddr  = "228.8.8.8:11001"
	posReqAddr    = "228.8.8.8:11002"
)

const (
	maxDatagramSize = 8192
)

func ip2Long(ip string) uint32 {
	var long uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}
