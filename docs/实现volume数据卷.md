# 实现volume数据卷
## 问题引入
上一小节介绍了如何使用AUFS包装busybox，从而实现容器和镜像的分离。但是一旦容器退出，容器可写层的所有内容都会被删除。那么，如果用户需要持久化容器里的部分数据该怎么办呢？

## 目标实现
**本节要实现将宿主机的目录作为数据卷挂载到容器中，并且在容器退出后，数据卷中的内容仍然能够保存在宿主机上**
### 什么是volume数据卷？
![alt text](/WYFdocker/docs/res/volume.png)
>在 Docker 中，数据卷（Volume）是一个可持久化的数据存储机制，用于在容器之间共享和持久化数据。数据卷可以绕过容器文件系统，并且可以由一个或多个容器共享和访问。
使用数据卷的主要优点包括：
1.持久化存储： 数据卷中的数据在容器销毁后仍然存在，因此可以用于持久化存储应用程序数据。
2.容器之间共享数据： 多个容器可以共享同一个数据卷，以便共享数据或状态。
3.与主机解耦： 数据卷可以在主机文件系统中的特定位置，也可以由 Docker 管理，从而解耦了容器的数据存储与主机的关系。

## 代码思路
- **首先，在main_command.go文件的runCommand命令中添加-v标签。**
```go
var RunCommand = cli.Command{
	Name:  "run",
	Usage: `Create a container with namespace and cgroups limit mydocker run -ti [command]`,
	Flags: []cli.Flag{
		......
		// 添加-v标签
		cli.StringFlag{
			Name:  "v",
			Usage: "volume",
		},
	},
	/*
		这里是run命令执行的真正函数
		1.判断参数是否包含command
		2.获取用户指定的command
		3.调用Run function 去准备启动容器
	*/
	Action: func(context *cli.Context) error {
		......
		volume := context.String("v")
		run.Run(tty, cmdArray, resConfig, volume)
		return nil
	},
}
```
- **在Run函数中，把volume传给创建容器的NewParentProcess函数和删除容器文件系统的DeleteWorkSpace函数。**
```go
func Run(tty bool, cmdarry []string,res *subsystems.ResourceConfig,volume string) {
    
	...........
    // 传入初始化进程中
	parent,writePipe := container.NewParentProcess(tty,RunUrl,MntUrl,volume)
    //如果fork进程出现异常，但有相关的文件已经进行了挂载，需要进行清理，避免后面运行报错时，需要手工清理
    if err := parent.Start(); err != nil {
		deleteWorkSpace(rootUrl, mntUrl, volume)
		log.Errorf("Run parent.Start err:%v", err)
		return
	}
    ......
    // 传入退出时的清理函数中
	deleteWorkSpace(RunUrl,MntUrl,volume)
    ...........
}
```
- 在NewWorkSpace函数中，继续把volume值传给创建容器文件系统的NewWorkSpace方法。
  - 创建只读层
  - 创建只读层
  - 创建只读层
```go
func newWorkSpace(rootUrl string,MntUrl string,volume string) error{
    ............
	if err := mountExtractVolume(MntUrl, volume); err != nil {
		return err
	}
	return nil
}
```
  - **接下来，首先判断volume是否为空，如果是，就表示用户并没有使用挂载标签，结束创建过程。如果不为空，则使用volumeUrlExtract函数解析volume字符串。**

```go
func volumeUrlExtract(volume string)([]string){
	var volumeUrl []string
	volumeUrl:=strings.Split(volume,"/")
	return volumeUrl
}
```
- volumeUrlExtract函数返回的字符数组长度为2，并且数据元素均不为空的时候，则执行MountVolume函数来挂载数据卷。否则，提示用户创建数据卷输入值不对。

```go
func mountExtractVolume(MntUrl string, volume string) error {
    if volume != "" {
        volumeUrl:=volumeUrlExtract(volume)
		length:=len(volumeUrl)
		if length==2 && volumeUrl[0]!=""&&volumeUrl[1]!=""{
			mountvolume(volumeUrl)
			log.Infof("%q",volumeUrl)
		}else{
            log.Infof("volume parameter input is not correct")
			return fmt.Errorf("volume parameter input is not correct")
		}    
	}
	return nil

}
```
--------------------------------------------------
- **挂载数据卷的过程如下：**
   - 1.首先，读取宿主机文件目录URL，创建宿主机文件目录（/root/${parentUrl}）。
   - 2.然后，读取容器挂载点URL，在容器文件系统里创建挂载点（/root/mnt/${containerUrl}）。
   - 3.最后，把宿主机文件目录挂载到容器挂载点。这样启动容器的过程，对数据卷的处理也就完成了。

```go
func mountvolume(mnturl string,volumeUrls []string)error{
	parentUrl := volumeUrls[0]
	exist, err := pathExist(parentUrl)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if !exist {
		// 使用mkdir all 递归创建文件夹
		if err := os.MkdirAll(parentUrl, 0777); err != nil {
			return fmt.Errorf("mkdir parent dir err: %v", err)
		}
	}

	// 在容器文件系统内创建挂载点
	containerUrl := mntUrl + volumeUrls[1]
	if err := os.Mkdir(containerUrl, 0777); err != nil {
		return fmt.Errorf("mkdir container volume err: %v", err)
	}


	aPath:=parentUrl
	cPath:=containerUrl

	upperDir, _ := exec.Command("mktemp", "-d").Output()
	workDir, _ := exec.Command("mktemp", "-d").Output()
	// 去除末尾的换行符
	upperDirPath := string(upperDir[:len(upperDir)-1])
	workDirPath := string(workDir[:len(workDir)-1])
	// 构建挂载命令
	cmd := exec.Command("sudo", "mount", "-t", "overlay",
		"overlay",
		"-o", fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", aPath,upperDirPath, workDirPath),
		cPath)

	// 把宿主机文件目录挂载到容器挂载点
	//dirs := "dirs=" + parentUrl
	//cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mount volume err: %v", err)
	}
	return nil
}
```
-------------------------------------------------
- **删除容器文件系统的过程如下：**
  - 1.只有在volume不为空，并且使用volumeUrlExtract函数解析volume字符串返回的字符数组长度为
  - 2，数据元素均不为空的时候，才执行DeleteMountPointWithVolume函数来处理。2.其余的情况仍然使用前面的DeleteMountPoint函数。


```go
func deleteWorkSpace(RunUrl string,MntUrl string,volume string) {
	deleteMountPoint(MntUrl)
	deleteWirteLayer(RunUrl)
	unmountVolume(MntUrl, volume)
}

```
```go
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
```

# 结果：

容器：
```shell
[root@localhost WYFdocker]# ./WYFdocker run -ti -v /root/volumn/test:/test sh
{"level":"info","msg":"resConf:\u0026{  }","time":"2024-05-03T16:51:53+08:00"}
{"level":"info","msg":"createTty true","time":"2024-05-03T16:51:53+08:00"}
{"level":"error","msg":"RunUrl is /home/gocodes/WYFdocker/","time":"2024-05-03T16:51:53+08:00"}
{"level":"error","msg":"MntUrl is /home/gocodes/WYFdocker/mnt/","time":"2024-05-03T16:51:53+08:00"}
{"level":"error","msg":"createMountPoint is running--------------------------","time":"2024-05-03T16:51:53+08:00"}
{"level":"error","msg":"createMountPoint is end---------------------","time":"2024-05-03T16:51:53+08:00"}
{"level":"info","msg":"[\"/root/volumn/test\" \"/test\"]","time":"2024-05-03T16:51:53+08:00"}
{"level":"info","msg":"parent process run","time":"2024-05-03T16:51:53+08:00"}
{"level":"info","msg":"all command is  sh","time":"2024-05-03T16:51:53+08:00"}
{"level":"info","msg":"init come on","time":"2024-05-03T16:51:53+08:00"}
{"level":"info","msg":"command:","time":"2024-05-03T16:51:53+08:00"}
{"level":"error","msg":"setupmout is running-------------------------------","time":"2024-05-03T16:51:53+08:00"}
{"level":"info","msg":"current location: /home/gocodes/WYFdocker/mnt","time":"2024-05-03T16:51:53+08:00"}
{"level":"error","msg":"privot root is running-------------------------------","time":"2024-05-03T16:51:53+08:00"}
{"level":"error","msg":"privot root is end-------------------------------","time":"2024-05-03T16:51:53+08:00"}
{"level":"error","msg":"setupmout is end-------------------------------","time":"2024-05-03T16:51:53+08:00"}
{"level":"info","msg":"find path /bin/sh","time":"2024-05-03T16:51:53+08:00"}
/ # ls
bin               hello.txt         proc              test              var
dev               home              root              tmp
etc               mydocker-cgroups  sys               usr
/ # touch /test/test.txt
/ # ls /test/
test.txt
/ # 

```

宿主机：
```shell
[root@localhost /]# ls /root/volumn/test
test.txt
```

退出后：
```shell
[root@localhost /]# ls /root/volumn/test
test.txt
```



# debug之旅：

{"level":"info","msg":"volume parameter input is not correct","time":"2024-05-03T16:17:10+08:00"}
{"level":"error","msg":"newWorkSpace error volume parameter input is not correct","time":"2024-05-03T16:17:10+08:00"}
{"level":"error","msg":"New parent process error","time":"2024-05-03T16:17:10+08:00"}

-------------代码volumeUrl=strings.Split(volume,":")


[admin@localhost WYFdocker]$ rm -fr mnt
rm: 无法删除"mnt": 设备或资源忙

-------------umount /home/gocodes/WYFdocker/mnt        [去掉挂载在上面的文件系统]

[root@localhost /]# ls /root/volumn/test
[root@localhost /]# 
宿主机下没有出现文件？

-----------，使用了 OverlayFS 文件系统来挂载卷到容器中。OverlayFS 是一种通过在不同层之间叠加文件系统的方式来提供文件系统的一致视图的技术。在这种情况下，当在容器内创建文件时，实际上是在 OverlayFS 的上层（upperdir）创建了文件。这个上层是一个临时目录，会被挂载到容器的指定目录上。
所以，当在容器内创建文件时，文件实际上被写入了 OverlayFS 的上层目录中，而不是真正的挂载目录。这就是为什么在宿主机相应目录下没有出现文件 test.txt 的原因。
为了解决这个问题，需要在挂载时将宿主机的目录（lowerdir）和容器内的挂载目录一起指定，这样文件就会被正确地写入到宿主机的目录中。

OverlayFS 文件系统:
https://zhuanlan.zhihu.com/p/392508816
https://zhuanlan.zhihu.com/p/436450556