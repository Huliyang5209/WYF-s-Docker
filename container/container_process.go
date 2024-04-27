package container

import (
	"os/exec"
	"syscall"
	"os"
	log "github.com/sirupsen/logrus"
)


func NewParentProcess(tty bool)(*exec.Cmd,*os.File) {

	readPipe,writePipe,err:=os.Pipe()
	if(err!=nil){
    	log.Errorf("new Pipe error %v",err)
		return nil,nil
	}

    //os.Args[0]
    //path, err := os.Executable("/proc/self/exe")
    //path, err := os.Readlink("/proc/self/exe")
	cmd := exec.Command("/proc/self/exe","init")

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if tty{
		cmd.Stdin=os.Stdin
		cmd.Stdout=os.Stdout
		cmd.Stderr=os.Stderr
	}

	cmd.ExtraFiles=[]*os.File{readPipe}

	return cmd,writePipe
}

