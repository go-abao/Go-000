package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {

	// 1、测试 Errgroup
	testErrgroup()

	// 2、优雅退出go守护进程
	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	mc := make(chan int)
	cg := CreateGeneral(mc)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("退出", s)
				mc <- 1
				ExitFunc(cg)
			case syscall.SIGUSR1:
				fmt.Println("usr1", s)
			case syscall.SIGUSR2:
				fmt.Println("usr2", s)
			default:
				fmt.Println("other", s)
			}
		}
	}()
	fmt.Println("进程启动...")
	Worker(cg)

}

// CreateGeneral CreateGeneral
func CreateGeneral(ch chan int) chan int {
	sum := 0
	c := make(chan int, 100)
	stop := 0
	go func() {
		for {
			select {
			case n := <-ch:
				fmt.Println(n, "要关闭了...")
				if n == 1 {
					stop = n
				}
			default:
				time.Sleep(1 * time.Second)
				sum++
				c <- sum
			}
			if stop == 1 {
				close(c)
				break
			}
		}
	}()
	return c
}

// Worker Worker
func Worker(ch chan int) {
	for {
		time.Sleep(2 * time.Second)
		fmt.Println(<-ch, len(ch))
	}
}

// ExitFunc exit func
func ExitFunc(ch chan int) {
	fmt.Println("开始退出...")
	fmt.Println("执行清理...")
	for {
		if len(ch) == 0 {
			fmt.Println("结束退出...")
			time.Sleep(1 * time.Second)
			break
		}
	}
	os.Exit(0)
}

// testErrgroup test Errgroup
func testErrgroup() {
	group, _ := errgroup.WithContext(context.Background())
	for i := 0; i < 9; i++ {
		index := i
		group.Go(func() error {
			fmt.Printf("start to execute the %d gorouting\n", index)
			time.Sleep(time.Duration(index) * time.Second)
			if index%2 == 0 {
				return fmt.Errorf("something has failed on grouting:%d", index)
			}
			fmt.Printf("gorouting:%d end\n", index)
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		fmt.Println(err)
	}
}
