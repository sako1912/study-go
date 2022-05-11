package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	fmt.Println(":start server:")
	//카페 오픈 해서 주문 받을 준비
	l, err := net.Listen("tcp", "127.0.0.1:10000")
	if err != nil {
		log.Println(err)
	}

	//main process 종류 직전 카페 문닫음
	defer l.Close()

	//여러 손님을 주문을 받기 위해 무한루틴
	for {
		//카페에 손님이 옴
		conn, err := l.Accept()
		if err != nil {
			//에러 날 경우 다음 손님 기다려 : 손님이 들어왓다 바로 나간다거나...
			log.Println(err)
			continue
			//return
		} else {
			log.Println("connection client successful :", conn.RemoteAddr())
		}

		//고 루틴 : 쓰레드, 손님 한명당 직원 한명 배정(페어) => 손님이 여러 명인걸 대비하기 위해 직원을 대기
		go connHandler(conn)

	}
}

func connHandler(conn net.Conn) {
	// byte  buf 생성 : 제한 주문 (4096만큼만 주문을 받을 수 있음)
	recvBuf := make([]byte, 4096)

	n, err := conn.Read(recvBuf)
	if err != nil {
		if io.EOF == err {
			log.Println("connection is closed from client :", err)
			return
		}
		log.Println("fail to read : ", err)
		return
	}

	if 0 < n {
		data := recvBuf[:n]

		//스트링으로 변환
		log.Println("client send message :: ", string(data))

		//http 형식 응답
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n"))               // 요청(의도)
		conn.Write([]byte("content-type: text/html; charset=UTF-8\r\n")) //header
		conn.Write([]byte("\r\n"))                                       //header end
		conn.Write([]byte("hello"))                                      //body

		if err != nil {
			log.Println("response err :: ", err)
			return
		}
	}
	//사용 후 종료
	conn.Close()
}
