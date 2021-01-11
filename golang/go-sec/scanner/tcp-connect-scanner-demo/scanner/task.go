package scanner

import (
	"fmt"
	"go-sec/scanner/tcp-connect-scanner-demo/vars"
	"net"
	"runtime"
	"strings"
	"sync"
)

func GenerateTask(ipList []net.IP, ports []int) ([]map[string]int, int) {
	tasks := make([]map[string]int, 0)

	for _, ip := range ipList {
		for _, port := range ports {
			ipPort := map[string]int{ip.String(): port}
			tasks = append(tasks, ipPort)
		}
	}

	return tasks, len(tasks)
}

func AssigningTasks(tasks []map[string]int) {
	scanBatch := len(tasks) / vars.ThreadNum

	for i := 0; i < scanBatch; i++ {
		curTasks := tasks[vars.ThreadNum*i:vars.ThreadNum*(i+1)]
		RunTask(curTasks)
	}

	if len(tasks)%vars.ThreadNum > 0 {
		lastTasks := tasks[vars.ThreadNum*scanBatch:]
		RunTask(lastTasks)
	}
}

func RunTask(tasks []map[string]int) {
	wg := &sync.WaitGroup{}

	// create a buffer channel
	taskChan := make(chan map[string]int, vars.ThreadNum*2)

	for i := 0; i < vars.ThreadNum; i++ {
		go Scan(taskChan, wg)
	}

	// 生产者
	for _, task := range tasks {
		wg.Add(1)
		taskChan <- task
	}

	wg.Wait()
	close(taskChan)
}

func Scan(taskChan chan map[string]int, wg *sync.WaitGroup) {
	// 不断从channel中读取task
	for task := range taskChan {
		for ip, port := range task {
			err := SaveResult(Connect(ip, port))
			_ = err
			wg.Done()
		}
	}
}

func SaveResult(ip string, port int, err error) error {
	fmt.Printf("ip: %v, port: %v,err: %v, goruntineNum: %v\n", ip, port, err, runtime.NumGoroutine())
	if err != nil {
		return err
	}

	v, ok := vars.Result.Load(ip)
	if ok {
		ports, ok1 := v.([]int)
		if ok1 {
			ports = append(ports, port)
			vars.Result.Store(ip, ports)
		}
	} else {
		ports := make([]int, 0)
		ports = append(ports, port)
		vars.Result.Store(ip, ports)
	}

	return err
}

func PrintResult() {
	vars.Result.Range(func(key, value interface{}) bool {
		fmt.Printf("ip:%v\n", key)
		fmt.Printf("port:%v\n", value)
		fmt.Println(strings.Repeat("-", 100))
		return true
	})
}