package main
//
//import (
//	"bytes"
//	"context"
//	"fmt"
//	"github.com/prometheus/common/log"
//	"os"
//	"os/exec"
//	"syscall"
//	"time"
//)
//
//func Run(instance string, env map[string]string) bool {
//	var (
//		cmd         *exec.Cmd
//		proc        *Process
//		sysProcAttr *syscall.SysProcAttr
//	)
//
//	t := time.Now()
//	sysProcAttr = &syscall.SysProcAttr{
//		Setpgid: true, // 使子进程拥有自己的 pgid，等同于子进程的 pid
//		Credential: &syscall.Credential{
//			Uid: uint32(uid),
//			Gid: uint32(gid),
//		},
//	}
//
//	// 超时控制
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(j.Timeout)*time.Second)
//	defer cancel()
//
//	if j.ShellMode {
//		cmd = exec.Command("/bin/bash", "-c", j.Command)
//	} else {
//		cmd = exec.Command(j.cmd[0], j.cmd[1:]...)
//	}
//
//	cmd.SysProcAttr = sysProcAttr
//	var b bytes.Buffer
//	cmd.Stdout = &b
//	cmd.Stderr = &b
//
//	if err := cmd.Start(); err != nil {
//		j.Fail(t, instance, fmt.Sprintf("%s\n%s", b.String(), err.Error()), env)
//		return false
//	}
//
//	waitChan := make(chan struct{}, 1)
//	defer close(waitChan)
//
//	// 超时杀掉进程组 或正常退出
//	go func() {
//		select {
//		case <-ctx.Done():
//			log.Warnf("timeout kill job %s-%s %s ppid:%d", j.Group, j.ID, j.Name, cmd.Process.Pid)
//			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
//		case <-waitChan:
//		}
//	}()
//
//	if err := cmd.Wait(); err != nil {
//		j.Fail(t, instance, fmt.Sprintf("%s\n%s", b.String(), err.Error()), env)
//		return false
//	}
//	return true
//}
