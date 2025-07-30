package internal

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

var IP string

func Writer(name string, text string) {
	conn, err := net.Dial("tcp", "192.168.1.148:8080")
	if err != nil {
		return
	}
	defer conn.Close()
	message := fmt.Sprintf("%s:%s", name, text)
	_, err = conn.Write([]byte(message))
	if err != nil {
		return
	}
}

func Writing() ([]Pair, error, map[string]struct{}) {
	conn, err := net.Dial("tcp", "192.168.1.148:8080")
	if err != nil {
		return nil, nil, nil
	}
	defer conn.Close()
	_, err = conn.Write([]byte("Read\n"))
	if err != nil {
	}
	messages := make([]Pair, 0)
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return nil, err, nil
	}
	text := strings.SplitN(string(response), "&", 2)
	parts := strings.SplitN(string(text[0]), "|", 4096)
	for i := 0; i < len(parts); i++ {
		if parts[i] == "\n" {
			break
		}
		part := strings.SplitN(string(parts[i]), ":", 2)
		if len(part) != 2 {
			continue
		}
		if part[0] == "" || part[1] == "" {
			continue
		}
		messages = append(messages, Pair{Name: part[0], Text: part[1]})
	}
	allU := strings.SplitN(string(text[1]), "|", 4096)
	res := make(map[string]struct{})
	for i := 0; i < len(allU); i++ {
		if strings.Trim(allU[i], "\n") == "" {
			break
		}
		res[strings.Trim(allU[i], "\n")] = struct{}{}
	}
	return messages, nil, res
}
