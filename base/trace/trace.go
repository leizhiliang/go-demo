package main

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"runtime/trace"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup

var httpClient = &http.Client{Timeout: 30 * time.Second}

// SleepSomeTime :生成一个介于0到9之间的微秒级别的时间段，用于模拟程序中的等待或延迟操作。
// 返回值: 返回一个 time.Duration 类型的值。
func SleepSomeTime() time.Duration {
	return time.Microsecond * time.Duration(rand.Int()%10)
}

func create(readChan chan int) {
	defer wg.Done()
	for i := 0; i < 50; i++ {
		readChan <- getBodySize()
		SleepSomeTime()
	}
	close(readChan)
}

func convert(readChan chan int, output chan string) {
	defer wg.Done()
	for readChan := range readChan {
		output <- strconv.Itoa(readChan)
		SleepSomeTime()
	}
	close(output)
}

func outputStr(output chan string) {
	defer wg.Done()
	for _ = range output {
		// do nothing
		SleepSomeTime()
	}
}

// 获取百度页面大小
func getBodySize() int {
	resp, _ := httpClient.Get("https://baidu.com")
	res, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return len(res)
}

func run() {
	readChan, output := make(chan int), make(chan string)
	wg.Add(3)
	go create(readChan)
	go convert(readChan, output)
	go outputStr(output)
}

func main() {
	// 将trace数据输出到trace.out
	f, _ := os.Create("trace.out")
	defer f.Close()

	// trace 的开启和停止
	_ = trace.Start(f)
	defer trace.Stop()
	run()
	wg.Wait()
}
