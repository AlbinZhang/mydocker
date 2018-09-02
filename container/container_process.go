package container

import (
	"os"
	"os/exec"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

/*
这里是父进程，也就是当前进程执行的内容，根据上一章介绍的内容，应该比较容易明白.
1. 这里的/proc/self/exe 调用中， /proc/self/指的是当前运行进程自己的环境，exec其实就是自己调用了自己，使用这种方式对创建出来的进程进行初始
2. 后面的args是参数，其中init是传递给本进程的第一个参数，在本例中，其实就是会去调用initCommand去初始化进程的一些环境和资源
3. 下面的clone参数就是去fork出来一个新进程，并且使用了namespace隔离新创建的进程和外部环境
4. 如果用户指定了-ti参数，就需要吧当前进程的输入输出导出到标准输入输出上
*/

func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil, nil
	}

	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	cmd.Dir = "/root/busybox"

	cmd.ExtraFiles = []*os.File{readPipe}
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}
