package main

import (
	"bufio"
	"io"
	"os"
)

var gfwList = make(map[string]struct{})
var fgfw bool

func readGfwlist(file string) {
	fileHandle, err := os.Open(file)
	if err != nil {
		debug.Printf("open gfwlist %v,error: %v", file, err)
		return
	}
	defer fileHandle.Close()
	read := bufio.NewReader(fileHandle)
	for {
		bts, _, err := read.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			debug.Printf("read gfwlist %v,error: %v", file, err)
			break
		}
		if domain := readDomain(string(bts)); domain != "" {
			gfwList[domain] = struct{}{}
		}
	}
	debug.Printf("gfwlist load %v", len(gfwList))
}

func readDomain(bts string) string {
	var start int
	length := len(bts)
	for i := 0; i < length; i++ {
		if bts[i] == '/' {
			if start > 0 {
				return string(bts[start:i])
			} else {
				start = i + 1
			}
		}
	}
	return ""
}
