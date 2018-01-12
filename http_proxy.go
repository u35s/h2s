package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

func httpProxy(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	debug.Printf("starting local http proxy server at %v ...\n", addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			debug.Printf("listen accept error %v\n", err)
			break
		}
		go handleHttpProxyConn(conn)
	}
}

func readHttpHead(conn net.Conn, stop int) *bytes.Buffer {
	cur := &bytes.Buffer{}
	bts := make([]byte, 1)
	var last byte
	for {
		_, err := conn.Read(bts)
		if err != nil {
			break
		} else if bts[0] == '\r' {
			continue
		} else if bts[0] == '\n' {
			if stop == 1 {
				break
			} else if last == '\n' && stop == 2 {
				break
			}
		}
		last = bts[0]
		cur.Write(bts)
	}
	return cur
}

func getHostPortType(line []byte) (host, port, tp string, err error) {
	if n := len(line); n > 0 && line[n-1] == '\r' {
		line = line[:n-1]
	}
	slc := strings.Split(string(line), " ")
	if len(slc) < 2 {
		err = fmt.Errorf("first line err %v", string(line))
		return
	}
	switch slc[0] {
	case "CONNECT":
		hp := strings.Split(slc[1], ":")
		if len(hp) != 2 {
			err = fmt.Errorf("connect extract host,port err %v", slc[1])
			return
		}
		host, port, tp = hp[0], hp[1], "https"
	default:
		thp := strings.Split(slc[1], "/")
		if len(thp) < 3 {
			err = fmt.Errorf("%v host err %v", slc[0], slc[1])
			return
		}
		hp := strings.Split(thp[2], ":")
		if len(hp) == 1 {
			host, port, tp = hp[0], "80", "http"
		} else if len(hp) == 2 {
			host, port, tp = hp[0], hp[1], "http"
		} else {
			err = fmt.Errorf("%v extract host,port err %v", slc[0], slc[1])
			return
		}
	}
	return
}

func handleHttpProxyConn(conn net.Conn) {
	defer conn.Close()
	first := readHttpHead(conn, 1)
	later := readHttpHead(conn, 2)

	hosts, ports, tps, err := getHostPortType(first.Bytes())
	if err != nil {
		debug.Printf("get host,port,type error %v", err)
		return
	}

	addr := bytes.Buffer{}
	addr.WriteByte(0x03)
	addr.WriteByte(byte(uint8(len(hosts))))
	addr.WriteString(hosts)
	port, err := strconv.ParseInt(ports, 10, 0)
	if err != nil {
		debug.Printf("https host port conv err,%v\n", ports)
		return
	}
	//大端序
	port16 := uint16(port)
	addr.WriteByte(byte(uint8(port16 >> 8)))
	addr.WriteByte(byte(uint8(port16 & 0x00ff)))
	debug.Printf("connecting to %v:%v over %v\n", hosts, ports, tps)
	debug.Printf("\n%v\n%v", first.String(), later.String())

	remote, err := createServerConn(addr.Bytes(), hosts+":"+ports)
	if err != nil {
		debug.Printf("%v\n", err)
		if len(servers.srvCipher) > 1 {
			debug.Println("Failed connect to all avaiable shadowsocks server")
		}
		return
	}
	defer remote.Close()
	if tps == "https" {
		conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	} else {
		remote.Write(first.Bytes())
		remote.Write([]byte{'\n'})
		remote.Write(later.Bytes())
		remote.Write([]byte{'\n'})
	}
	go ss.PipeThenClose(conn, remote)
	ss.PipeThenClose(remote, conn)
	debug.Printf("closed connection to %v:%v", hosts, ports)
}
