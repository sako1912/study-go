package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	log.Printf("##client 실행 ")

	//client 서버와 연결
	conn, err := net.Dial("tcp", ":10000")
	if err != nil {
		log.Println("fail connect server :", err)
	}

	//고루틴을 생성해서 서버가 값을 던질때까지 기다렸다가 던지면 값을 출력한다
	go func() {
		data := make([]byte, 1000)
		for {
			n, err := conn.Read(data)
			if err != nil {
				log.Println("fail send :", err)
				return
			}

			log.Println("sever send : " + string(data[:n]))
			time.Sleep(time.Duration(3) * time.Second)
		}
	}()

	//사용자의 입력이 들어올때까지 blocking했다가 입력을 마치면 서버로 전송한다.
	for {
		var s string
		fmt.Scanln(&s)        //사용자 입력값
		conn.Write([]byte(s)) //서버로 전송
		time.Sleep(time.Duration(3) * time.Second)
	}

	// if nil != error {
	// 	log.Printf("접속 실패: %v", error)
	// } else {
	// 	_, error := connection.Write([]byte("붐"))
	// 	if nil == error {
	// 		log.Printf("전송 성공")
	// 	}
	// 	connection.Close()
	// }
}
