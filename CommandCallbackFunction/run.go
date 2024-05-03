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

func Run(tty bool, cmdarry []string,res *subsystems.ResourceConfig,volume string) {
    
	
	pwd,err:=os.Getwd()
	if err != nil {
	    log.Errorf("get current work dir error %v", err)
		return
	}

    RootUrl:=pwd+"/"
	log.Errorf("RunUrl is %s",RootUrl)
	MntUrl:=pwd+"/mnt/"
	log.Errorf("MntUrl is %s",MntUrl)

	parent,writePipe := container.NewParentProcess(tty,RootUrl,MntUrl,volume)

	if parent == nil {
	    log.Errorf("New parent process error")
		return 
	}

	if err := parent.Start(); err != nil {
		deleteWorkSpace(RootUrl, MntUrl, volume)
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
    
	deleteWorkSpace(RootUrl,MntUrl,volume)

	os.Exit(-1) 
}

func deleteWorkSpace(RunUrl string,MntUrl string,volume string) {
	deleteMountPoint(MntUrl)
	deleteWirteLayer(RunUrl)
	unmountVolume(MntUrl, volume)
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

func unmountVolume(mntUrl string, volume string) {
	if volume == "" {
		return
	}
	volumeUrls := strings.Split(volume, ":")
	if len(volumeUrls) != 2 || volumeUrls[0] == "" || volumeUrls[1] == "" {
		return
	}

	containerUrl := mntUrl + volumeUrls[1]
	cmd := exec.Command("umount", containerUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("ummount volume failed: %v", err)
	}
}