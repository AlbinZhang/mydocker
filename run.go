package main

import (
	"os"
	"strings"

	"github.com/AlbinZhang/mydocker/cgroups"
	"github.com/AlbinZhang/mydocker/cgroups/subsystems"

	"github.com/AlbinZhang/mydocker/container"
	log "github.com/Sirupsen/logrus"
)

/*
这里的Start方法是真正开始前面创建好的Command的调用，他首先会clone出来一个namespace隔离的进程，
然后在子进程中，调用/proc/self/exe,也就是到用自己，发送init参数，调用我们写的init方法，去初始化容器的一些资源
*/
//func Run(tty bool, command string) {
func Run(tty bool, comArray []string, res *subsystems.ResourceConfig) {

	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	// use mydocker-cgroup as cgroup name
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)

	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(comArray, writePipe)

	parent.Wait()
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
