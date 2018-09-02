package cgroups

import (
	"github.com/AlbinZhang/mydocker/cgroups/subsystems"
	"github.com/Sirupsen/logrus"
)

type CgroupManager struct {
	// cgroup 在 hierarchy 中的路径，相对于创建的 cgroup 目录对于个 root cgroup 目录的路径
	Path string

	//资源配置
	Reasource *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

// 将进程PID加入到每个cgroup中
func (c *CgroupManager) Apply(pid int) error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		subSysIns.Apply(c.Path, pid)
	}
	return nil
}

// 设置cgroup 资源限制
func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		subSysIns.Set(c.Path, res)
	}
	return nil
}

// 设置各个subsystem 挂载中的 cgroup
func (c *CgroupManager) Destroy() error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		if err := subSysIns.Remove(c.Path); err != nil {
			logrus.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}
