package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

type Data struct {
	b []byte
}

type HandleFunc func(data []byte) []byte

var replacer = strings.NewReplacer("吗", "", "？", "!", "?", "!")

func Reply(data []byte) []byte {
	reply := replacer.Replace(string(data))
	if string(reply[len(reply)-1]) != "!" {
		reply = reply + "!"
	}
	reply = reply + "\n"
	return []byte(reply)
}

func handleInput(r *bufio.Reader, f HandleFunc) chan Data {
	out := make(chan Data)
	go func() {
		defer close(out)
		for {
			line, _, err := r.ReadLine()
			if err != nil {
				return
			}
			out <- Data{f(line)}
		}
	}()
	return out
}

func HandleConn(conn net.Conn, f HandleFunc) {
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	dataChan := handleInput(r, f)
	go func() {
		for data := range dataChan {
			w.Write(data.b)
			w.Flush()
		}
	}()
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %v\n", err)
			continue
		}
		go HandleConn(conn, Reply)
	}
}
