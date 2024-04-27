package Cgroups

import(
   "WYFdocker/Cgroups/subsystems"
   "github.com/sirupsen/logrus"
)

type CgroupsManager struct{
	Path string
	resoursre *subsystems.ResourceConfig
}

func NewCgroupsManager(path string) *CgroupsManager{
    return &CgroupsManager{
		Path:path,
	}
}

func (c *CgroupsManager)Apply(pid int)error{
	for _,SubIns:=range(subsystems.SubsystemsIns){
	    SubIns.Apply(c.Path,pid)
	}
	return nil
}

func (c *CgroupsManager)Set(res *subsystems.ResourceConfig)error{
    for _,SubIns:=range(subsystems.SubsystemsIns){
        err:=SubIns.Set(c.Path,res)
		if err!=nil{
		    logrus.Errorf("apply subsystem:%s err:%s", SubIns.Name(), err)
		}
    }
	return nil
}

func (c *CgroupsManager)Destory()error{
    for _,SubIns:=range(subsystems.SubsystemsIns){
       if err:=SubIns.Remove(c.Path);err!=nil{
           logrus.Warnf("Remove cgroup fail %v",err)
       }
    }
	return nil
}