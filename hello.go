package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	fmt.Println("hellow go")
	//open socket
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println(err)
	}

	//main process 종류 시 소켓도 종료
	defer l.Close()

	for {
		//소켓에 연결
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		//프로세스 종료시 연결도 종료
		defer conn.Close()

		//handler에 연결 전달
		go connHandler(conn)
	}
}

func connHandler(conn net.Conn) {
	//4096 byte  buf 생성
	recvBuf := make([]byte, 4096)
	for {
		//연결이 but 읽음 : client가 값을 줄때까지 blocking되어 대기하다가 값을 주면 읽어들인다??
		n, err := conn.Read(recvBuf)

		//에러 처리
		if err != nil {
			if io.EOF == err {
				log.Println(err)
				return
			}
			log.Println(err)
			return
		}

		if 0 < n {
			//buf 를 data에 할당
			data := recvBuf[:n]
			log.Println(string(data))
			//client가 던진 값을 다시 client에게 던진다.?
			_, err = conn.Write(data[:n])

			//에러 처리
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
