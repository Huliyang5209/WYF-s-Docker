package CommandCallbackFunction

import (
	
	"os"
	"WYFdocker/container"
	"WYFdocker/Cgroups"
	"WYFdocker/Cgroups/subsystems"
	log "github.com/sirupsen/logrus"
	"strings"

	
)

func Run(tty bool, cmdarry []string,res *subsystems.ResourceConfig) {

	parent,writePipe := container.NewParentProcess(tty)

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
		
	sendInitCommand(cmdarry,writePipe)
	parent.Wait()
	os.Exit(-1) 
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