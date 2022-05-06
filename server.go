package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	fmt.Println("start sercer go")
	//open socket
	l, err := net.Listen("tcp", ":10000")
	if err != nil {
		log.Println(err)
	}

	//main process 종류 시 소켓도 종료
	defer l.Close()

	//연결을 무한이 받을수 있게 루프
	for {
		//소켓에 연결
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		//프로세스 종료시 연결도 종료
		defer conn.Close()

		//handler에 연결 전달 : 연결에 대한 처리를 go 루틴 사용
		go connHandler(conn)
	}
}

func connHandler(conn net.Conn) {
	// byte  buf 생성
	recvBuf := make([]byte, 1000)
	//반복하여 읽음
	for {
		//연결이 client에서 온걸 읽음 : client가 값을 줄때까지 blocking되어 대기하다가 값을 주면 읽어들인다
		n, err := conn.Read(recvBuf)

		//에러 처리
		if err != nil {
			//입력이 종료되면 종료
			if io.EOF == err {
				log.Println("fisnish connect :", err)
				return
			}
			log.Println("connect fail : ", err)
			return
		}

		if 0 < n {
			//buf 를 data에 할당
			//client 에서 받아온 값을 data에 할당 : 받아온 길이 만큼 슬라이스를 잘라서 출력
			data := recvBuf[:n]
			log.Println("client send message :: ", string(data))
			//response:: client의 값을 받아서 다시 client에 전송
			_, err = conn.Write(data[:n])

			//에러 처리
			if err != nil {
				log.Println("response err :: ", err)
				return
			}
		}
	}
}
