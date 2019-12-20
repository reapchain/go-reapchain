package main

import "fmt"

func sum(a, b int ) <-chan int{  // 함수의 리턴값은  int형 받기 전용 채널
	out := make(chan int)
	go func(){
		out <- a + b
		fmt.Printf("out=%x, %v\n", out, out )  //out 채널의 포인터만 출력 ?
	}()
	return out  // 채널 변수 자체 리턴
}

func main() {
	c := sum(1,2) // 채널을 리턴값을 받아서 c에 대입
	//c는 <-chan int : int형 받기 전용 채널
	fmt.Println( <- c )  // 3 : 채널에서 값을 꺼냄

	ch := make(chan string, 1) //size 1
	sendChan(ch)
	receiveChan(ch)
}
func sendChan(ch chan<-string){  //채널에 값을 주기
	ch <- "Data"  // 채널의 buf에 [0]string 에 저장됨.
	// x:= <-ch //에러

}
func receiveChan(ch <-chan string){  //<-채널,, 채널로 부터 값을 가져오기
	data := <-ch
	fmt.Println(data)
}