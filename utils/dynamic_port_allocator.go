package utils

import (
	"net"
	"strconv"
)

var (
	PORT string
)

func AssignDynamicPort() {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	// Get the actual port that was allocated
	port := listener.Addr().(*net.TCPAddr).Port

	PORT = strconv.Itoa(port)
}
