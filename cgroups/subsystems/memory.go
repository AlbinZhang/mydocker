package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

// memory subsystem的实现
type MemorySubSystem struct {
}

// 设置 cgroupPath 对应的cgroup的内存资源限制
func (s *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	// GetCgroupPath 的作用是获取当前Subsystem在虚拟文件系统中的路径，GetCgroupPath这个函数会在下面介绍
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.MemoryLimit != "" {
			// 设置这个cgroup的内存限制，即 将限制写入到cgroup对应目录的memory.limit_in_bytes文件中
			if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0664); err != nil {
				return fmt.Errorf("set cgroup memory fail %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}

// 删除 cgroupPath对应的cgroup
func (s *MemorySubSystem) Remove(cgroupPath string) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		//删除 cgroup 便是删除对应的 cgroupPath 的目录
		return os.Remove(subsysCgroupPath)
	} else {
		return err
	}
}

// 将一个进程加入到 cgroupPath 对应的 cgroup 中
func (s *MemorySubSystem) Apply(cgroupPath string, pid int) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		// 把进程的PID写到cgroup的虚拟文件及系统对应目录下的task文件中
		if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("set cgroup fail %v", err)
		} else {
			return err
		}
	} else {
		return err
	}
}

func (s *MemorySubSystem) Name() string {
	return "memory"
}
