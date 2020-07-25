package webgin

import (
	"fmt"
	"testing"
)

var i =0
var num = make(chan int)
func TestContext_Bind(t *testing.T) {
	//go func() {<-num}()
	//for  {
	//	select {
	//	case num<-1:
	//		close(num)
	//		fmt.Println("will")
	//		num=make(chan int)
	//	case <-num:
	//		fmt.Println("close")
	//	}
	//}
	close(num)
	for  {
		fmt.Println(<-num)
	}

}
func test()  {
	go putnum()
	go readnum()
}
func putnum()  {
	for  {
		i++
		num<-i
	}
}
func readnum()  {
	for  {
		n:=<-num
		fmt.Println(n)
	}
}