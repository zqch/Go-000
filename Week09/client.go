package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	r := bufio.NewReader(conn)

	var message string
	for {
		fmt.Print("send message: ")
		fmt.Scanf("%s", &message)

		if message == "exit" {
			return
		}
		fmt.Fprintf(conn, "%s\n", message)
		line, _, err := r.ReadLine()
		if err != nil {
			fmt.Printf("read err: %v\n", err)
			return
		} else {
			fmt.Printf("receive: %s\n", line)
		}
	}
}
