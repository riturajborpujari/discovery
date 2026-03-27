package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

const (
	DiscoveryPort = 9999
	ServerPort	  = 8080
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", conn.RemoteAddr(), err)
			}
			return
		}
		fmt.Printf("INFO: %s: %s\n", conn.RemoteAddr(), string(buf[0:n]))
	}
}

func broadcastDiscoveryMessage() {
	fmt.Printf("INFO: Discovery: Broadcasting presence\n")
	bcastClient, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP: net.IPv4bcast,
		Port: DiscoveryPort,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Discovery broadcast failed: %s\n", err)
		// TODO handle error a better way
		os.Exit(1)
	}
	_, err = bcastClient.Write([]byte("DISCOVERY"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Discovery broadcast failed: %s\n", err)
		// TODO handle error a better way
		os.Exit(1)
	}
	fmt.Printf("INFO: Discovery: Broadcast successful\n")
}

func connectPeer(peerIp net.IP) {
	fmt.Printf("INFO: Peer connecting: %s\n", peerIp)
	peer, err := net.Dial("tcp", fmt.Sprintf("%s:%d", peerIp, ServerPort))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Peer Connect failed: %s: %s\n", peerIp, err)
		return
	}
	defer peer.Close()

	fmt.Printf("INFO: Peer connected: %s\n", peerIp)
	_, err = peer.Write([]byte("Hello Peer. Found you via Discovery"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Peer Write failed: %s: %s\n", peerIp, err)
		return
	}
}

func startPeerDiscovery() {
	broadcastDiscoveryMessage()

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP: net.IPv4zero,
		Port: DiscoveryPort,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Discovery server failed: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("INFO: Discovery: Server started on %s:%d\n", net.IPv4zero, DiscoveryPort)

	buf := make([]byte, 512)
	for {
		n, peerAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Discovery: Peer read failed: %s\n", err)
		}
		fmt.Printf("INFO: Discovery: Peer read: %s: %s\n", peerAddr.IP, string(buf[0:n]))
		go connectPeer(peerAddr.IP)
	}
}

func main() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", ServerPort))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Application Server could not start: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("INFO: Application Server: started on :%d\n", ServerPort)

	go startPeerDiscovery()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Application Server connection accept failed: %s\n", err)
			continue
		}
		go handleConn(conn)
	}
}
