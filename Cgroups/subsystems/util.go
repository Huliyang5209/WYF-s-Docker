package subsystems
import (
    "bufio"
    "fmt"
    "os"
    "path"
    "strings"
)


func FindCgroupMountPoint(subsys string) string{
	f,err:=os.Open("/proc/self/mountinfo")
	if err!=nil{
		return ""
	}
	defer f.Close()
	scanner:=bufio.NewScanner(f)

	for scanner.Scan(){
    text:=scanner.Text()
	fileds:=strings.Split(text,"")
	for _,opt:=range strings.Split(fileds[len(fileds)-1],","){
		if opt==subsys{
		    return fileds[4]
		}
	}
	}
    if err:=scanner.Err();err!=nil{
	    return ""
	}
return ""
}

func GetVFCgroupPath(subsys string,cgroupPath string,autoCreate bool)(string,error){
   cgrouproot:=FindCgroupMountPoint(subsys)
   if _,err:=os.Stat(path.Join(cgrouproot,cgroupPath));err==nil||(autoCreate&&os.IsNotExist(err)){
       if os.IsNotExist(err){
		if err:=os.Mkdir(path.Join(cgrouproot,cgroupPath),0755);err==nil{	    
		}else{
		return "", fmt.Errorf("fail create cgroup %v",err)
	   }
	}
	return path.Join(cgrouproot,cgroupPath),nil
   }else
   {
	 return "", fmt.Errorf("cgrouppath err %v",err)
   }
}
