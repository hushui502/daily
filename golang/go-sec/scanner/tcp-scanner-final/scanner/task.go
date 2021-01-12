package scanner

import (
	"fmt"
	"go-sec/scanner/tcp-scanner-final/vars"
	"net"
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

func RunTask(tasks []map[string]int) {
	wg := &sync.WaitGroup{}

	taskChan := make(chan map[string]int, vars.ThreadNum*2)

	for i := 0; i < vars.ThreadNum; i++ {
		go Scan(taskChan, wg)
	}

	for _, task := range tasks {
		wg.Add(1)
		taskChan <- task
	}
}

func Scan(taskChan chan map[string]int, wg *sync.WaitGroup) {
	for task := range taskChan {
		for ip, port := range task {
			if strings.ToLower(vars.Mode) == "syn" {
				err := SaveResult(SynScan(ip, port))
				_ = err
			} else {
				err := SaveResult(Connect(ip, port))
				_ = err
			}
			wg.Done()
		}
	}
}

func SaveResult(ip string, port int, err error) error {
	if err != nil {
		return err
	}

	if port > 0 {
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
	}

	return err
}

func PrintResult() {
	vars.Result.Range(func(key, value interface{}) bool {
		fmt.Printf("ip:%v\n", key)
		fmt.Printf("ports: %v\n", value)
		fmt.Println(strings.Repeat("-", 100))
		return true
	})
}
