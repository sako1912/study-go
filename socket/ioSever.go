package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("::server 실행::")

	//open socket
	l, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Println(err)
	}

	//main process 종류 시 소켓도 종료
	defer l.Close()

	//연결을 무한이 받을수 있게 루프
	for {
		//소켓 연결 대기: 사용자 접속대기
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		} else {
			log.Println("connection client successful :", conn.RemoteAddr())

		}
		//프로세스 종료시 연결도 종료
		defer conn.Close()

		//handler에 연결 전달 : 연결에 대한 처리를 여러개 하기위해 go 루틴 사용
		go connHandler(conn)
		//go fileIO(conn)
	}
}

func connHandler(conn net.Conn) {
	// byte  buf 생성
	recvBuf := make([]byte, 1024)
	//반복하여 읽음
	for {
		//read request
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

func fileIO(conn net.Conn) {

	fileBuf := make([]byte, 1024)

	//무한 요청을 받으려고 for
	for {
		path := "C:\\Project\\server\\study-go\\sample\\"

		//////////////////////////
		//reqest start
		n, err := conn.Read(fileBuf)

		if err != nil {
			//입력이 종료되면 종료
			if io.EOF == err {
				log.Println("fisnish connect :", err)
				return
			}
			log.Println("connect fail : ", err)
		}

		if 0 < n {
			//buf 를 data에 할당
			//client 에서 받아온 값을 data에 할당 : 받아온 길이 만큼 슬라이스를 잘라서 출력
			data := fileBuf[:n]
			log.Println("client looking for :: ", string(data))
			path += string(data)
			log.Println("full path  :: ", path)
		}
		////////////reqest end////////////////

		//파일 존재 여부
		hasFile := fileExists(path)

		//파일 없을경우 message 전달 있을경우 파일을 읽어서 내용 전달
		if !hasFile {
			s := "No such fileName"
			log.Println(s)
			_, err = conn.Write([]byte(s))

		} else {

			//입력파일 열기
			fi, err := os.Open(path)
			if err != nil {
				log.Println("fail to open file", err)
				return
			}

			for {
				//buffer 만큼 슬라이스해서 읽음
				cnt, err := fi.Read(fileBuf)
				if err != nil {
					//read 완료시
					if io.EOF == err {
						log.Println("io.EOF : 파일 읽기 끝")
						//return
					} else {
						log.Println("fail to rad file", err)
						return
					}

				}
				//루프 종료
				if cnt == 0 {
					fi.Close()
					break
				}

				//응답 처리
				_, err = conn.Write(fileBuf[:cnt])
				if err != nil {
					log.Println("response err :: ", err)
					return
				}
			}
		}

	}

}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
