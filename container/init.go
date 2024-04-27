package container
import(
	"os/exec"
	"syscall"
	"io/ioutil"
	"fmt"
	"strings"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)


func ReadUserCommand() []string{
	readPipe:=os.NewFile(uintptr(3),"pipe")
	msg,err:=ioutil.ReadAll(readPipe)
	if err!=nil{
	    log.Errorf("init readpipe failed %v",err)
		return nil
	}
	msgStr:=string(msg)
	return strings.Split(msgStr," ")
}



func Runcontainerinit() error{

	if err := SetUpMount(); err != nil {
		return err
	}

	cmdarry:=ReadUserCommand()
	if cmdarry==nil||len(cmdarry)==0{
	    return fmt.Errorf("Runcontainerinit can't get command from pipe")
	}

    path,err:=exec.LookPath(cmdarry[0])
	if err!=nil{
	    log.Errorf("can't find exec path: %s %v", cmdarry[0], err)
	    return err
	}
    log.Infof("find path %v",path)

    if err:=syscall.Exec(path,cmdarry,os.Environ());err!=nil{	
		log.Errorf("syscall exec err: %v", err.Error())
	}
	return nil
}

func SetUpMount() error{

	if err:=syscall.Mount("/","/","",syscall.MS_REC|syscall.MS_PRIVATE,"");err!=nil{
		return fmt.Errorf("setupmout mount proc err %v",err)
	}

	curpth,err:=os.Getwd()
	if err!=nil{
		return fmt.Errorf("get curpth err %v",err)
	}
	log.Infof("current location: %s", curpth)

	err=privotRoot(curpth)
	if err!=nil{	
		return err
	}

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		log.Errorf("mount proc failed: %v", err)
		return err
	}
	syscall.Mount("tmpfs", "/dev", "tempfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
	return nil
}

func privotRoot(root string) error{

   if err:=syscall.Mount(root,root,"bind",syscall.MS_BIND|syscall.MS_REC,"");err!=nil{
	   return fmt.Errorf("mount root to itselfe err %v",err)
   }

   privotDir:=filepath.Join(root,".pivot_root")

   if _,err:=os.Stat(privotDir);err==nil{
	  if err:=os.Remove(privotDir);err!=nil{
		return err
   }
   }

   if err:=os.Mkdir(privotDir,0700);err!=nil{
	return fmt.Errorf("mkdir pivot_root err %v",err)
   }

   if err:=syscall.PivotRoot(root,privotDir);err!=nil{
    return fmt.Errorf("pivot_root err %v",err)
   }

   if err:=os.Chdir("/");err!=nil{
	return fmt.Errorf("chdir root err: %v", err)
   }

    pivotDir:= filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir err: %v", err)
	}
	return os.Remove(pivotDir)
}