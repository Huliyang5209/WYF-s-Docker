package CommandCallbackFunction

import (
	
	"os"
	"WYFdocker/container"
	"WYFdocker/Cgroups"
	"WYFdocker/Cgroups/subsystems"
	log "github.com/sirupsen/logrus"
	"strings"
	"os/exec"

	
)

func Run(tty bool, cmdarry []string,res *subsystems.ResourceConfig) {
    
	
	pwd,err:=os.Getwd()
	if err != nil {
	    log.Errorf("get current work dir error %v", err)
		return
	}

    RunUrl:=pwd+"/"
	log.Errorf("RunUrl is %s",RunUrl)
	MntUrl:=pwd+"/mnt/"
	log.Errorf("MntUrl is %s",MntUrl)

	parent,writePipe := container.NewParentProcess(tty,RunUrl,MntUrl)

	if parent == nil {
	    log.Errorf("New parent process error")
		return 
	}

	if err := parent.Start(); err != nil {
		log.Errorf("Run parent.Start err:%v", err)
		return
	}

	CgroupsManager:=Cgroups.NewCgroupsManager("mydocker-cgroups")
	defer CgroupsManager.Destory()

	CgroupsManager.Set(res)
	CgroupsManager.Apply(parent.Process.Pid)
	log.Infof("parent process run")
	sendInitCommand(cmdarry,writePipe)
	parent.Wait()
    
	deleteWorkSpace(RunUrl,MntUrl)

	os.Exit(-1) 
}

func deleteWorkSpace(RunUrl string,MntUrl string) {
	deleteMountPoint(MntUrl)
	deleteWirteLayer(RunUrl)
}

func deleteWirteLayer(runUrl string ){
    writeLayerUrl:=runUrl+"writeLayer/"
	if err:=os.RemoveAll(writeLayerUrl);err!=nil{
	    log.Errorf("deleteMountPoint remove %s err : %v", writeLayerUrl, err)
	}
}

func deleteMountPoint(mntUrl string){
    cmd := exec.Command("umount", mntUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("deleteMountPoint umount %s err : %v", mntUrl, err)
	}
	if err := os.RemoveAll(mntUrl); err != nil {
		log.Errorf("deleteMountPoint remove %s err : %v", mntUrl, err)
	}
}

func sendInitCommand(arry []string,writePipe *os.File) {
    command:=strings.Join(arry," ")
	log.Infof("all command is  %s",command)
	if _, err := writePipe.WriteString(command); err != nil {
		log.Errorf("write pipe write string err: %v", err)
		return
	}
	if err := writePipe.Close(); err != nil {
		log.Errorf("write pipe close err: %v", err)
	}
}