package util

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"go-sec/scanner/password_crack/logger"
	"go-sec/scanner/password_crack/models"
	"go-sec/scanner/password_crack/plugins"
	"go-sec/scanner/password_crack/util/hash"
	"go-sec/scanner/password_crack/vars"
	"gopkg.in/cheggaaa/pb.v2"
	"runtime"
	"strings"
	"sync"
	"time"
)

func GenerateTask(ipList []models.IpAddr, users []string, passwords []string) (tasks []models.Service, taskNum int) {
	tasks = make([]models.Service, 0)

	for _, user := range users {
		for _, password := range passwords {
			for _, addr := range ipList {
				service := models.Service{Ip: addr.Ip, Port: addr.Port, Protocol: addr.Protocol, Username: user, Password: password}
				tasks = append(tasks, service)
			}
		}
	}

	return tasks, len(tasks)
}

func Scan(ctx *cli.Context) (err error) {
	if ctx.IsSet("debug") {
		vars.DebugMode = ctx.Bool("debug")
	}
	if vars.DebugMode {
		logger.Log.Level = logrus.DebugLevel
	}

	if ctx.IsSet("timeout") {
		vars.Timeout = time.Duration(ctx.Int("timeout")) * time.Second
	}
	if ctx.IsSet("scan_num") {
		vars.ScanNum = ctx.Int("scan_num")
	}

	if ctx.IsSet("ip_list") {
		vars.IpList = ctx.String("ip_list")
	}

	if ctx.IsSet("user_dict") {
		vars.UserDict = ctx.String("user_dict")
	}

	if ctx.IsSet("pass_dict") {
		vars.PassDict = ctx.String("pass_dict")
	}

	if ctx.IsSet("outfile") {
		vars.ResultFile = ctx.String("outfile")
	}

	vars.StartTime = time.Now()

	userDict, uErr := ReadUserDict(vars.UserDict)
	if uErr != nil {
		return err
	}
	passDict, pErr := ReadPasswordDict(vars.PassDict)
	if pErr != nil {
		return err
	}

	ipList := ReadIpList(vars.IpList)

	aliveIpList := CheckAlive(ipList)

	tasks, _ := GenerateTask(aliveIpList, userDict, passDict)
	RunTask(tasks)

	return nil
}

func RunTask(tasks []models.Service) {
	tasksCount := len(tasks)
	vars.ProgressBar = pb.StartNew(tasksCount)
	vars.ProgressBar.SetTemplate(`{{ rndcolor "Scanning progress: " }} {{  percent . "[%.02f%%]" "[?]"| rndcolor}} {{ counters . "[%s/%s]" "[%s/?]" | rndcolor}} {{ bar . "「" "-" (rnd "ᗧ" "◔" "◕" "◷" ) "•" "」" | rndcolor }} {{rtime . | rndcolor}} `)

	wg := &sync.WaitGroup{}

	// 创建一个线程数为缓冲区大小的channel
	taskChan := make(chan models.Service, vars.ScanNum)

	for i := 0; i < vars.ScanNum; i++ {
		go checkPassword(taskChan, wg)
	}

	for _, task := range tasks {
		wg.Add(1)
		taskChan <- task
	}

	defer close(taskChan)
	waitTimeout(wg, vars.Timeout)

	{
		_ = models.SaveResultToFile()
		models.ResultTotal()
		_ = models.DumpToFile(vars.ResultFile)
	}
}

func checkPassword(taskChan chan models.Service, wg *sync.WaitGroup) {
	for task := range taskChan {
		vars.ProgressBar.Increment()
		if vars.DebugMode {
			logger.Log.Debugf("checking: Ip: %v, Port: %v, Protocol: %v, Username: %v, Password: %v, goroutineNum: %v", task.Ip,
				task.Port, task.Protocol, task.Username, task.Password, runtime.NumGoroutine())
		}

		var k string
		protocol := strings.ToUpper(task.Protocol)

		if protocol == "REDIS" {
			k = fmt.Sprintf("%v-%v-%v", task.Ip, task.Port, task.Protocol)
		} else {
			k = fmt.Sprintf("%v-%v-%v", task.Ip, task.Port, task.Username)
		}

		h := hash.MakeTaskHash(k)
		// 如果已经保存过了就跳过
		if hash.CheckTaskHash(h) {
			wg.Done()
			continue
		}

		fn := plugins.ScanFuncMap[protocol]
		models.SaveResult(fn(task))
		wg.Done()
	}
}

// returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		// when close <-c will send a '0'
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		return false
	case <-time.After(timeout):
		return true
	}
}
