见版本V2.0
# 问题引入
在4.1节中,我们已经实现了使用宿主机/root/busybox目录作为文件的根目录，但在容器内对文件的操作仍然会直接影响到宿主机的/root/busybox目录
![alt text](/WYFdocker/docs/res/image-1.png)

# 目标实现
**本节要进一步进行容器和镜像隔离，实现在容器中进行的操作不会对镜像产生任何影响的功能。**

UFS（Union File System）和 AUFS（Another Union File System）都是联合文件系统，它们允许将多个目录或文件系统合并到单个目录中。这种合并是虚拟的，意味着原始目录和文件系统的内容并没有实际合并，而是通过联合文件系统的机制进行管理和访问。

>举个生活中的例子：
>假设你有一本书，其中有几个章节是你自己写的，另外一些章节是你的朋友写的。使用 UFS 或 AUFS 就像是你在阅读这本书
>时，同时看到了你的章节和你朋友的章节，但这并不意味着你的章节和你朋友的章节被合并成一篇新的章节，而是它们仍然保持>各自独立的存在。

在 Docker 中，利用 AUFS 可以实现写时复制（copy-on-write）功能，从而使得容器中的操作不会对基础镜像产生影响。

具体来说，Docker 使用 AUFS 实现如下功能：

1.容器的文件系统与镜像的分离： 当你在 Docker 中启动一个容器时，Docker 使用 AUFS 将容器的文件系统挂载到基础镜像的文件系统之上。这样，容器中的文件系统就成为了一个独立的层（layer），与基础镜像分离开来。

2.写时复制： 当容器中的进程试图进行写操作（例如创建、修改或删除文件）时，AUFS 会拦截这些操作，并在容器的文件系统中创建一个新的文件副本，而不是直接修改基础镜像的文件系统。这种写时复制的机制使得容器中的操作不会影响到基础镜像，同时保持了容器与基础镜像之间的分离。

![alt text](/WYFdocker/docs/res/images.jpg/)

>可以说容器是镜像的一个实例化运行。每个容器都是基于特定的镜像创建的，并且可以根据需要配置不同的参数。然而，容器
>身并不是镜像，而是镜像的一个运行时实例。


# 代码实现

NewWorkSpace函数是用来创建容器文件系统的，它包括CreateReadOnlyLayer、CreateWriteLayer和CreateMountPoint。1.CreateReadOnlyLayer函数新建busybox文件夹，将busybox.tar解压到busybox目录下，作为容器的只读层。
2.CreateWriteLayer函数创建了一个名为writeLayer的文件夹，作为容器唯一的可写层。
3.CreateMountPoint函数中，首先创建了mnt文件夹，作为挂载点，然后把writeLayer目录和busybox目录mount到mnt目录下。
4.最后，在NewParentProcess函数中将容器使用的宿主机目录/root/busybox替换成/root/mnt。

**注意：由于本人在centos系统上（不支持AUFS文件系统）进行编写，故以overlay文件系统进行替代，最新版本的docker也支持overlay文件系统，实现方法大同小异。**


>OverlayFS（overlay）和AUFS（Advanced Multi-Layered Unification Filesystem）都是 Linux 内核中的联合文件系统，用于在容器中实现文件层叠和联合挂载的功能。它们之间有一些相似之处，也有一些不同之处。
**相似之处：**
1.联合挂载： 无论是 OverlayFS 还是 AUFS，它们都允许多个文件系统层被联合挂载到一个目录中，形成一个统一的文件系统视图。
2.写时复制（Copy-on-Write）： 当一个文件系统层被修改时，OverlayFS 和 AUFS 都采用写时复制的策略，即只有在需要时才会复制被修改的部分，而不是整个文件。这可以提高性能和节省空间。
3.层叠文件系统： 它们都支持在容器中使用多个文件系统层，每个层都可以包含文件和目录，并且在统一的文件系统视图中合并。
**不同之处：**
1.实现机制： OverlayFS 是 Linux 内核的一部分，从 Linux 内核版本 3.18 开始就被引入了，而 AUFS 则是一个独立的文件系统，需要额外的内核模块来支持。
2.性能特性： 一般来说，OverlayFS 在性能方面比 AUFS 更好，特别是在大规模容器部署中。OverlayFS 的性能优化在内核中得到了持续改进，而 AUFS 由于不是内核的一部分，可能在一些情况下性能表现不如 OverlayFS。
3.支持程度： 由于 OverlayFS 是 Linux 内核的一部分，因此在大多数现代 Linux 发行版中都有良好的支持。相比之下，AUFS 需要额外的内核模块，并且可能需要额外的配置和维护。
稳定性和成熟度： 由于 OverlayFS 是 Linux 内核的一部分，因此在稳定性和成熟度方面可能更可靠一些。而 AUFS 虽然在一些场景下性能可能较好，但其稳定性和成熟度可能不如 OverlayFS。
**总的来说，OverlayFS 是当前 Linux 容器环境下的首选文件系统，它性能良好且稳定可靠，而 AUFS 则是一个备选方案，适用于一些特定的场景和需求。**

>**tips**
1.aufs是如何区分可写层和只层的？
答：在 AUFS 中，文件系统的层次结构是按照顺序堆叠的，最底层是基础镜像层，而最顶层是可写层。因此，AUFS 可以通过查看文件系统的层次结构来区分哪个文件是可写层，哪个是基础镜像层。通常，AUFS 会将最顶层标记为可写层，并在进行写操作时，将写操作作用在这个可写层上，而不影响底层的基础镜像层。(文件系统的层次结构是按照挂载的前后顺序确定的。在 AUFS 或其他联合文件系统中，挂载顺序决定了层次结构中每个层的位置)
假如你在 AUFS 文件系统中按照顺序将 a, b, c, d 挂载到 /mnt，那么：
最顶层（即最后挂载的那个层）是可写层。在这个例子中，d 会是可写层。
其余的层都是只读层。因此，a, b, c 这三个会是只读层。
在这样的配置中，所有写操作都会作用在可写层 d 上，而读取操作会按照从最顶层到最底层的顺序进行。如果一个文件在上层存在，它会遮挡掉下层的同名文件，这种机制叫做"覆盖" (overlay)。这让文件系统能够有效地组合多个只读层，同时在最顶层提供可写操作的灵活性。

## NewWorkSpace
```go
func newWorkSpace(rootUrl string,MntUrl string) error{

    if err:=createReadOnlyLayer(rootUrl);err!=nil{
        return err
    }
	if err:=createWriteLayer(rootUrl);err!=nil{
	    return err
	}
	if err:=createMountPoint(rootUrl,MntUrl);err!=nil{
	    return err
	}
	return nil
}

```

## CreateReadOnlyLayer

```go
func createReadOnlyLayer(rootUrl string)error{
    readlayer:=rootUrl+"busybox/"
	readlayerpath,err:=pathExist(readlayer)
	if err!=nil{
	    return err
	}
	if !readlayerpath{
	    return fmt.Errorf("readlayer(busybox) dir don't exist: %s", readlayer)
	}
	return nil
}
```

## CreateWriteLayer

```go
func createWriteLayer(rootUrl string)error{
    writelayer:=rootUrl+"writelayer/"
	if err:=os.Mkdir(writelayer,0777);err!=nil{
	    return fmt.Errorf("create write layer failed: %v", err)
	}
	return nil
}
```

## CreateMountPoint
```go
func createMountPoint(rootUrl string,MntUrl string)error{

	if err := os.Mkdir(MntUrl, 0777); err != nil {
		log.Errorf("Mkdir mntUrl failed: %v", err)
		return fmt.Errorf("mkdir faild: %v", err)
	}
	
	aPath:=rootUrl+"writelayer/"
	bPath:=rootUrl+"busybox/"
	cPath:=MntUrl

	upperDir, _ := exec.Command("mktemp", "-d").Output()
	workDir, _ := exec.Command("mktemp", "-d").Output()
	// 去除末尾的换行符
	upperDirPath := string(upperDir[:len(upperDir)-1])
	workDirPath := string(workDir[:len(workDir)-1])
	// 构建挂载命令
	cmd := exec.Command("sudo", "mount", "-t", "overlay",
		"overlay",
		"-o", fmt.Sprintf("lowerdir=%s:%s,upperdir=%s,workdir=%s", aPath, bPath, upperDirPath, workDirPath),
		cPath)

    //aufs实现
    //dirs := "dirs=" + rootUrl + "writeLayer:" + rootUrl + "busybox"
	//cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err:=cmd.Run();err!=nil{
	     return fmt.Errorf("mmt dir err: %v", err)
	} 

	return nil
}

```
## NewParentProcess
```go
func NewParentProcess(tty bool,rootUrl string,MntUrl string)(*exec.Cmd,*os.File) {

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

	if err:=newWorkSpace(rootUrl,MntUrl);err!=nil{
		log.Errorf("newWorkSpace error %v",err)
		return nil,nil
	}
	cmd.Dir=MntUrl

	return cmd,writePipe
}
```
## deleteWorkSpace
```go
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
```

## 实现效果
运行容器，并创建一个文件，观察容器和镜像的文件系统是否发生变化。
容器观察：

![alt text](/WYFdocker/docs/res/image2.png)

宿主机观察：
![alt text](/WYFdocker/docs/res/image.png)

