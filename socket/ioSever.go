package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sort"
	"strings"
)

func main() {
	fmt.Println("::server 실행::")

	//open socket
	l, err := net.Listen("tcp", "127.0.0.1:10000")
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
			//return
		} else {
			//log.Println("connection client successful :", conn.RemoteAddr())

		}
		//프로세스 종료시 연결도 종료
		//defer conn.Close()

		//handler에 연결 전달 : 연결에 대한 처리를 여러개 하기위해 go 루틴 사용
		//go connHandler(conn)

		go bufioHandler(conn)
		//go fileIO(conn)
	}
}

func bufioHandler(conn net.Conn) {
	from := conn.RemoteAddr().String()
	log.Println("client connect from : ", from)

	//req 한번에 다 스캔
	reader := bufio.NewReader(conn)
	l := 0
	reqHeader := make(map[string]string)
	for {
		l++

		//line, err := reader.ReadString('\n') //\r\n이 불가하니 scanner func을 사용해야할듯...
		line, isPrefix, err := reader.ReadLine()

		if isPrefix || err != nil {
			log.Println("Error Read : ", err)
		}

		if len(line) == 0 {
			break
		}

		// byte to string
		strLine := string(line)
		log.Println("strLine : ", strLine)
		if l == 1 {
			_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
			conn.Write([]byte("content-type: text/html; charset=UTF-8\r\n"))
			conn.Write([]byte("\r\n"))

			firstData := strings.Split(strLine, " ")
			method := firstData[0]
			path := firstData[1]
			version := firstData[2]

			conn.Write([]byte("method is " + method + "<br>\r\n"))
			conn.Write([]byte("path is " + path + "<br>\r\n"))
			conn.Write([]byte("version is " + version + "<br>\r\n"))
			conn.Write([]byte("#########start header data #########<br>\r\n"))
		} else {
			tHeader := strings.Split(strLine, ": ")
			reqHeader[tHeader[0]] = tHeader[1]
		}
	}

	writeHeader(conn, reqHeader)
	conn.Close()

}

func getSortKeys(object map[string]string) []string {
	///////////////////////////////
	//map은 순서를 보장하지 않아서 sort할 key들의 배열이 필요
	sortKeys := make([]string, 0, len(object))

	//key 기준으로 sort 할꺼니깐 배열에 키만 담음
	for k := range object {
		sortKeys = append(sortKeys, k)
	}

	//string 배열을 오름차순으로 정렬
	sort.Strings(sortKeys)

	return sortKeys
	///////////////////////////
}

func writeHeader(conn net.Conn, reqHeader map[string]string) {
	//header sort
	sortKeys := getSortKeys(reqHeader)
	//write : key 기준으로
	for _, k := range sortKeys {
		conn.Write([]byte(k + " : " + reqHeader[k] + "<br>\r\n"))
	}

}
func connHandler(conn net.Conn) {
	fmt.Println("::connHandler 실행::")

	data := make([]byte, 0)
	// byte  buf 생성
	recvBuf := make([]byte, 4096)

	//반복하여 읽음
	//for {
	//연결이 client에서 온걸 읽음 : client가 값을 줄때까지 blocking되어 대기하다가 값을 주면 읽어들인다
	n, err := conn.Read(recvBuf)

	//log.Println("connect read :: ", n)
	//에러 처리
	if err != nil {
		//입력이 종료되면 종료
		if io.EOF == err {
			log.Println("eof error :", err)
			return
		}
		log.Println("fail to read : ", err)
		return
	}
	data = append(data, recvBuf[:n]...)
	fmt.Printf("READ  %d bytes\n", n)

	if 0 < n {
		log.Println("client start ennene")
		//buf 를 data에 할당
		//client 에서 받아온 값을 data에 할당 : 받아온 길이 만큼 슬라이스를 잘라서 출력
		data := recvBuf[:n]
		//log.Println("client send message :: ", string(data))
		method, path, version, mHeader := classifyData(data)

		///////////////////////////////
		//map은 순서를 보장하지 않아서 sort할 key들의 배열이 필요
		sortKeys := make([]string, 0, len(mHeader))

		//key 기준으로 sort 할꺼니깐 배열에 키만 담음
		for k := range mHeader {
			sortKeys = append(sortKeys, k)
		}

		//string 배열을 오름차순으로 정렬
		sort.Strings(sortKeys)
		///////////////////////////

		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
		conn.Write([]byte("content-type: text/html; charset=UTF-8\r\n"))
		conn.Write([]byte("\r\n"))
		//conn.Write([]byte("hello"))
		conn.Write([]byte("method is " + method + "<br>\r\n"))
		conn.Write([]byte("path is " + path + "<br>\r\n"))
		conn.Write([]byte("version is " + version + "<br>\r\n"))
		conn.Write([]byte("#########start header data #########<br>\r\n"))
		//sork 기준으로 값을 불러옴
		for _, k := range sortKeys {
			//fmt.Println(k, mHeader[k])
			conn.Write([]byte(k + " : " + mHeader[k] + "<br>\r\n"))
		}
		conn.Write([]byte("#########finish header data #########<br>\r\n"))

		//client data 파일로 떨굼
		//createNetInfoFile(data)

		//response:: client의 값을 받아서 그대로 다시 client에 전송
		//_, err = conn.Write(data[:n])

		//에러 처리
		if err != nil {
			log.Println("response err :: ", err)
			return
		}
	}
	conn.Close()
	//}
}

func classifyData(data []byte) (string, string, string, map[string]string) {
	mHeader := make(map[string]string)
	//mReq := make(map[string]string)

	var method, path, version string

	strData := string(data)
	spData := strings.Split(strData, "\r\n")

	for line, lData := range spData {
		//첫줄은 request 영역
		//두번째 줄부터 header \r\n
		//해더가 끝나는 기준이(\r\n\r\n) 다음 값이 공백일경우 = body 시작
		if len(lData) == 0 {
			//해더끝
			break
		}

		if line == 0 {
			firstData := strings.Split(lData, " ")
			method = firstData[0]
			path = firstData[1]
			version = firstData[2]
			continue
		}

		hData := strings.Split(lData, ": ")
		mHeader[hData[0]] = hData[1]
	}

	return method, path, version, mHeader
}

//client info 파일로 만듬
func createNetInfoFile(data []byte) {
	log.Println("::start create file::")
	outputPath := "C:\\Project\\server\\netInfo.txt"
	//outputPath := "C:\\Temp\\2.txt"
	fo, err := os.Create(outputPath)
	if err != nil {
		log.Println("fail to create output file", err)
		panic(err)
	}
	defer fo.Close()
	_, err = fo.Write(data)

	if err != nil {
		//현재 함수를 즉시 멈추고 현재 함수에 defer 함수들을 모두 실행한 후 즉시 리턴한다
		//...? 전체 종료??
		panic(err)
	}
}

func fileIO(conn net.Conn) {
	log.Println(" :: start FILE IO :: ")
	fileBuf := make([]byte, 4096)

	//무한 요청을 받으려고 for
	for {
		path := "C:\\Project\\server\\study-go\\sample\\"

		log.Println("path : ", path)

		//////////////////////////
		//reqest start
		n, err := conn.Read(fileBuf)
		log.Println("connect read :: ", n)

		if err != nil {
			//입력이 종료되면 종료
			if io.EOF == err {
				log.Println("fisnish connect :", err)
				return
			}
			log.Println("connect fail : ", err)
		}

		//if 0 < n {
		path += "test.txt"
		log.Println("full path : ", path)
		//}

		// if 0 < n {
		// 	//buf 를 data에 할당
		// 	//client 에서 받아온 값을 data에 할당 : 받아온 길이 만큼 슬라이스를 잘라서 출력
		// 	data := fileBuf[:n]
		// 	log.Println("client looking for :: ", string(data))
		// 	path += string(data)
		// 	log.Println("full path  :: ", path)
		// }
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
