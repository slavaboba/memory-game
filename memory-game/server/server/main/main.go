package main

import (
	"fmt"
	"memory-master/Server/internal"
	"net"
)

func main() {
	internal.AllUsersNames = make(map[string]struct{})
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()
	fmt.Println("Сервер запущен на :8080")

	internal.Scores = append(internal.Scores, internal.Score{Name: "", Text: ""})

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		internal.Scores = internal.HandleConnection(conn, internal.Scores)
	}
}
