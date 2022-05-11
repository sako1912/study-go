# study-go
Go lang

# 기타 설명
:n => 0:n =>  0이 생략가능, 0부터 n번재까지 
0: => 0부터 마지막 까지

# Server.go
카페를 예시로 설명

카페 = ip
카페 내 창고방인지 주문받는 곳 인지 위치 = port
카페 점원 = go 루틴 (쓰레드)

func main() {
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

func connHandler(conn net.Conn) {
	// byte  buf 생성 : 제한 주문 (4096만큼만 주문을 받을 수 있음)
	recvBuf := make([]byte, 4096)

    //같은 손님이 주문 시 커피 1잔 그리고 스콘 1개 그리고 초콜릿 2개....... 이어서 주문할경우를 대비하여 for를 사용 (connection 지향)
    //http 같은 경우 커피 1잔 까지만 주문가능 한번 요청으로 끝 (connection less)
	for {
		//직원이 손님의 주문을 받음 (buffer 사이즈 만큼만 받을 수 있음)
		n, err := conn.Read(recvBuf)

		if err != nil {
            //주문 도중에 중단하고 나가버림 혹은 주문 다하고 나가버림 
            if io.EOF == err {
                log.Println("connection is closed from client :", err)
                return
            }
            //손님의 주문을 직원이 못 알아들었을 경우
		    log.Println("fail to read : ", err)
		    return
        }
	
        //주문 내용이 1개이상 있을 주문 처리:음료제조 등
		if 0 < n {
			//주문이 담긴 buffer를 0부터 n번재 까지의 데이터를 가져옴
			data := recvBuf[:n]

            //스트링으로 변환 
			log.Println("client send message :: ", string(data))

			//제조한 음료를 손님한테 전달
			_, err = conn.Write(data[:n])

			//제조한 음료를 손님이 못받앗을 경우
			if err != nil {
				log.Println("response err :: ", err)
				return
			}
		}
	}
}
}