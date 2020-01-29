//Unbuffered Channel에서는 수신자가 데이타를 받을 때까지 송신자가 데이타를 보내는 채널에 묶여 있게 된다.

//Buffered Channel을 사용하면 비록 수신자가 받을 준비가 되어 있지 않을 지라도 지정된 버퍼만큼 데이타를 보내고 계속 다른 일을 수행할 수 있다.

package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan bool)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
		}
		done <- true
	}()

	// 위의 Go루틴이 끝날 때까지 대기
	<-done
    var done chan
	var oxtest bool
	done , oxtest = schedule(10)

	if (oxtest ) {

	fmt.Println("schedule is finished")
    }

}

// 일정 시간마다 스케쥴 반복하기
func schedule(delay time.Duration) chan bool {

	stop := make(chan bool) // 채널을 만들어 놓고
	go func() { // 고루틴을 돌렸는데
		for {
			select {
			case <-time.After(delay): //일정 시간이 지나야
			case <-stop: // 채널에 값을 보낸 후 리턴한다!
				return
			}
		}
	}()

	return stop
}