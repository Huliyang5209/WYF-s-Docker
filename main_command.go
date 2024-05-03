package main

import (
	"WYFdocker/CommandCallbackFunction"
	"WYFdocker/container"

	"github.com/urfave/cli"
	log "github.com/sirupsen/logrus"

	"WYFdocker/Cgroups/subsystems"

	"fmt"

	
)

var runCommand = cli.Command{
	Name:  "run",
	Usage: "create a container with namespace and Cgroups limit WYFdocker run -ti [command]",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "mem", 
			Usage: "memory limit,e.g.: -mem 100m",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
		    Name: "v",
			Usage: "volume",
		},
	},
	Action: func(context *cli.Context) error {

		if len(context.Args()) < 1 {
			return fmt.Errorf("can't run without command")
		}
		
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		
		resConf := &subsystems.ResourceConfig{
			MemoryLimit: context.String("mem"),
			CpuShare:    context.String("cpuShare"),
			CpuSet:      context.String("cpuSet"),
		}
		log.Info("resConf:", resConf)
		
		tty := context.Bool("ti")
		log.Infof("createTty %v", tty)

		volume := context.String("v")
		
		CommandCallbackFunction.Run(tty, cmdArray,resConf,volume)
		return nil
	},
}

var initCommand=cli.Command{
	Name:"init",
	Usage:"Init container process run user's process in container. Do not call it outside",
	Action:func(context *cli.Context)error{
    log.Infof("init come on")

	cmd:=context.Args().Get(0)
	log.Infof("command:%s",cmd)
    
	err:=container.Runcontainerinit()
	return err
	},
}

var CommitCommand=cli.Command{
    Name:"commit",
	Usage:"commit a container into image",
	Action:func(context *cli.Context)error{
	    if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container name")
		}
		imageName := context.Args().Get(0)
		return container.CommitContainer(imageName)
	},
}