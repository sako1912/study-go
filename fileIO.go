package main

import (
	"io"
	"log"
	"os"
)

func main() {
	path := "C:\\Project\\server\\study-go\\sample\\test.txt"
	outputPath := "C:\\Temp\\2.txt"

	//입력 파일
	fi, err := os.Open(path)
	if err != nil {
		log.Println("fail to open file", err)
		panic(err)
	}

	//메인 함수가 끝날 때 파일 닫음
	defer fi.Close()

	// 출력파일 생성
	fo, err := os.Create(outputPath)
	if err != nil {
		log.Println("fail to create output file", err)
		panic(err)
	}
	defer fo.Close()

	buff := make([]byte, 1024)
	for {
		//buffer 만큼 슬라이스해서 읽음
		cnt, err := fi.Read(buff)
		if err != nil {
			//read 완료시
			if io.EOF == err {
				log.Println("파일 읽기 종료")
				return
			}
			log.Println("fail to rad file", err)
			return

		}
		log.Println(cnt)
		//끝이면 루프 종료
		if cnt == 0 {
			break
		}

		// 쓰기
		_, err = fo.Write(buff[:cnt])
		if err != nil {
			//현재 함수를 즉시 멈추고 현재 함수에 defer 함수들을 모두 실행한 후 즉시 리턴한다
			//...? 전체 종료??
			panic(err)
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
