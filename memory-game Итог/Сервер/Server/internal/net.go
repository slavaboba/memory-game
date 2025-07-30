package internal

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func HandleConnection(conn net.Conn, scores []Score) []Score {
	conn.SetWriteDeadline(time.Time{})
	messag, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Printf("Получено: %s", messag)
	if messag == "Read\n" {
		fmt.Println("Start of sending")
		//conn.Write([]byte(fmt.Sprintf("%v\n", len(scores))))
		for i := 0; i < len(scores); i++ {
			//fmt.Print(scores[i].Name, " ", scores[i].Text)
			_, err := conn.Write([]byte(fmt.Sprintf("%v:%v|", scores[i].Name, scores[i].Text)))
			fmt.Println(fmt.Sprintf("%v:%v", scores[i].Name, scores[i].Text))
			if err != nil {
				fmt.Println(err)
			}
		}
		conn.Write([]byte("&"))
		for u := range AllUsersNames {
			_, err := conn.Write([]byte(fmt.Sprintf("%v|", u)))
			fmt.Println(fmt.Sprintf("%v", u))
			if err != nil {
				fmt.Println(err)
			}
		}
		conn.Write([]byte("\n"))
		fmt.Println("End of sending")
	} else {
		parts := strings.SplitN(string(messag), ":", 2)
		fmt.Println(parts)
		scores = append(scores, Score{parts[0], parts[1]})
	}
	for i := 0; i < len(scores); i++ {
		fmt.Println(scores[i].Name, " ", scores[i].Text)
	}
	conn.Close()
	return scores
}
