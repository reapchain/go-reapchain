package main

import (
	"fmt"

)


type Duck struct {

}
type Person struct {

}

func ( d Duck) quack() {
	fmt.Println(" 꽥 ~!")
}

func( d Duck) feathers() {
	fmt.Println(" 오리는 흰색과 회색털을 가지고 있다. ")
}

func ( p Person) quack() {
	fmt.Println(" 꽥 ~!")
}

func( p Person) feathers() {
	fmt.Println(" 오리는 흰색과 회색털을 가지고 있다. ")
}


type Quacker interface { //인터페이스는 메서드들의 모음
	quack()
	feathers()
}

func inTheForest( q Quacker ){
	q.quack()  // Quacker 인터페이스로 quack 메서드 호출
	q.feathers() //

}



func main(){

	var donald Duck
	var john   Person

	inTheForest(donald)
	inTheForest(john)

	if v, ok := interface{} (donald).(Quacker); ok {  // 타입이 특정인터페이스를 구현하는지 검사하려면,
	                                                  // interface{}(인스턴스).(인터페이스)
		fmt.Println(v, ok )
	}

	var x interface{}  // C#, Java에서는 object, C/C++ 에서는 void*   같은것임.

	x =1

	x ="Tom"
	printemptyinterface(x)

	// Type Assertion

	// Interface type 의 x와 타입 T에 대하여 x.(T) 로 표현햇을때,


	var a interface{} = 1

	i := a  // i 와 a 는 다이나믹 타입, 값은 1

	j := a.(int) // j는 int 타입이면서, 값은 1

	println(i)  // 포인터 주소 출력
	println(j)  // 값 1 출력




















}

func printemptyinterface( v interface{}){    //Empty interface 는 메서드를 전혀 갖지 않는 빈 인터페이스
// 어떠한 타입도 담을수 있음, 다른 언어에서는   Dynamic Type 또는 Void * 같은 것임
	fmt.Println(v) //
}
