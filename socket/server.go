package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	fmt.Println("start sercer go")
	//소켓 대기중
	l, err := net.Listen("tcp", "127.0.0.1:10000")
	if err != nil {
		log.Println(err)
	}

	//main process 종류 시 소켓도 종료
	defer l.Close()

	//연결을 무한이 받을수 있게 루프
	for {
		//클라이언트에서 연결 받음
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		//main 프로세스 종료시 연결도 종료
		defer conn.Close()

		//handler에 연결 전달 : 연결에 대한 처리를 go 루틴 사용 :점원이 여명 쓰레드멀티
		go connHandler(conn)
	}
}

func connHandler(conn net.Conn) {
	// byte  buf 생성
	recvBuf := make([]byte, 4096)
	//반복하여 읽음test

	for {
		//연결이 client에서 온걸 읽음 : client가 값을 줄때까지 blocking되어 대기하다가 값을 주면 읽어들인다
		n, err := conn.Read(recvBuf) //length 몇

		log.Println("coon Read :: ", n)
		//에러 처리
		if err != nil {
			//입력이 종료되면 종료
			if io.EOF == err {
				log.Println("connection is closed from client :", err)
				return
			}
			log.Println("fail to read : ", err)
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

func asdasd() {
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
		//프로세스 종료 전 직원들 퇴근
		//defer conn.Close()

		//고 루틴 : 쓰레드, 손님 한명당 직원 한명 배정(페어) => 손님이 여러 명인걸 대비하기 위해 직원을 대기
		go connHandler(conn)
		//go fileIO(conn)
	}
}
