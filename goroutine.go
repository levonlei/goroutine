package main

import (
	"fmt"
	"sync"
	"time"
)

// LEN 循环次数控制
const LEN = 2

type myLock struct {
	sync.Mutex
}

var (
	lock myLock
	//二维切片
	status [LEN][2]bool
)

// 状态设计
var (
	firstSuccess  = [2]bool{true, false}
	secondSuccess = [2]bool{true, true}
)

func first(i int) {
	status[i] = firstSuccess
}

func second(i int) {
	for {
		lock.Lock()
		curStatus := status[i]
		lock.Unlock()
		if curStatus == firstSuccess {
			status[i] = secondSuccess
			break
		}
	}
}

func three(i int) {
	for {
		lock.Lock()
		curStatus := status[i]
		lock.Unlock()
		//不能在for循环中加入defer，这样是有可能造成资源泄漏的，锁根本得不到释放。在这里也不能放在外面，原因是，第一遍还没有达到状态的时候，
		//上了锁了，但是锁没有得到及时的释放，造成了死锁。
		//这里的lock不能放到if中去，假如刚开始的状态达不到if的判断条件
		if curStatus == secondSuccess {
			fmt.Println("one > two > three")
			break
		}
	}
}

var control chan int
var wg sync.WaitGroup

func PrintfHelloWorld(i int) {
	defer wg.Done()
	fmt.Printf("goroutine:%d,input hello world!\n", i)
	time.Sleep(1 * time.Second)
	<-control
}

func test4() {
	control = make(chan int, 2)
	for i := 0; i < 10; i++ {
		control <- i
		wg.Add(1)
		go PrintfHelloWorld(i)
	}
	//等待线程结束，不存在数据竞争，那么可以直接使用这种方式来进行，其他的不行
	wg.Wait()
	fmt.Println("关闭信道")
	close(control)

	exit := make(chan bool)
	//wg.Add(1)
	go func() {

		for {
			select {
			case <-exit:
				fmt.Println("收到退出信号")
				close(exit)
				//wg.Done()
				return
			default:
				fmt.Println("监控中")
			}
		}
	}()

	//传递false也是可以的，这里就代表说只要收到信道的传值就是直接可以的。
	//time.Sleep(1*time.Millisecond)
	exit <- false
	//wg.Wait()
	//下面这种方式还是不能够实现正确关闭，还是会丢掉边界的执行，并不安全。
	//
	//for {
	//	if _,ok:=<-control;ok {
	//		fmt.Println("函数执行完毕")
	//		break
	//	}
	//}

	//这种使用可以先满足题目要求
	//time.Sleep(1*time.Second)
}

func main() {
	for i := 0; i < LEN; i++ {
		go first(i)
		go second(i)
		go three(i)
	}
	time.Sleep(1 * time.Second)
}
