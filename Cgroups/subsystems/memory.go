package subsystems

import(
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type MemorySubSystem struct {
}

func (MS *MemorySubSystem) Name() string {
	return "memory"
}

func (MS *MemorySubSystem)Set(CgroupPath string,res *ResourceConfig) error{

	 if VirtualFile_cgrouppath,err:=GetVFCgroupPath(MS.Name(),CgroupPath,true);err==nil{
	     if res.MemoryLimit!=""{
	         if err:=ioutil.WriteFile(path.Join(VirtualFile_cgrouppath,CgroupPath,"memory.limit_in_bytes"),[]byte(res.MemoryLimit),0644);err!=nil{
				return fmt.Errorf("set memory limit filed %v",err)
			 }
	     }
		 return nil
	 }else{
      return err
	 }
}

func (MS *MemorySubSystem)Apply(CgroupPath string,pid int) error{
	if VirtualFile_cgrouppath,err:=GetVFCgroupPath(MS.Name(),CgroupPath,false);err==nil{
			if err:=ioutil.WriteFile(path.Join(VirtualFile_cgrouppath,"tasks"),[]byte(strconv.Itoa(pid)),0644);err!=nil{
			   return fmt.Errorf("set cgroup proc filed %v",err)
		}
		return nil
	}else{
	 return fmt.Errorf("get cgroup %v err %v",CgroupPath,err)
    }
}

func (MS *MemorySubSystem)Remove(CgroupPath string) error{
    if VirtualFile_cgrouppath,err:=GetVFCgroupPath(MS.Name(),CgroupPath,false);err==nil{
		return os.Remove(VirtualFile_cgrouppath)
	}else{
		return err
	}
}
