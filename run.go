package main

import (
	"os"

	"github.com/AlbinZhang/mydocker/container"
	log "github.com/Sirupsen/logrus"
)

/*
这里的Start方法是真正开始前面创建好的Command的调用，他首先会clone出来一个namespace隔离的进程，
然后在子进程中，调用/proc/self/exe,也就是到用自己，发送init参数，调用我们写的init方法，去初始化容器的一些资源
*/
func Run(tty bool, command string) {

	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	parent.Wait()
	os.Exit(-1)
}
