package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	log.Printf(":: client 실행 ::")

	//client 서버와 연결
	conn, err := net.Dial("tcp", "127.0.0.1:10000")
	if err != nil {
		log.Println("fail connect server :", err)
	}

	//고루틴을 생성해서 서버가 값을 던질때까지 기다렸다가 던지면 값을 출력한다
	//직원이 커피 주기를 기다림
	go func() {
		data := make([]byte, 4096)
		for {
			//제조된 음료를 받음
			res, err := conn.Read(data)
			//받기 실패!!
			if err != nil {
				log.Println("fail read :", err)
				return
			}

			log.Println("sever send : " + string(data[:res]))
			time.Sleep(time.Duration(3) * time.Second)
		}
	}()

	//직원한테 커피를 주문
	for {
		var s string
		fmt.Scanln(&s)        //사용자 입력값 s변수에 담기
		conn.Write([]byte(s)) //서버로 s 전송
		time.Sleep(time.Duration(3) * time.Second)
	}

	//서버에 데이타 전송 후 connection 닫기
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
