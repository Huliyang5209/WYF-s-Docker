���汾V2.0
# ��������
��4.1����,�����Ѿ�ʵ����ʹ��������/root/busyboxĿ¼��Ϊ�ļ��ĸ�Ŀ¼�����������ڶ��ļ��Ĳ�����Ȼ��ֱ��Ӱ�쵽��������/root/busyboxĿ¼
![alt text](/WYFdocker/docs/res/image-1.png)

# Ŀ��ʵ��
**����Ҫ��һ�����������;�����룬ʵ���������н��еĲ�������Ծ�������κ�Ӱ��Ĺ��ܡ�**

UFS��Union File System���� AUFS��Another Union File System�����������ļ�ϵͳ�������������Ŀ¼���ļ�ϵͳ�ϲ�������Ŀ¼�С����ֺϲ�������ģ���ζ��ԭʼĿ¼���ļ�ϵͳ�����ݲ�û��ʵ�ʺϲ�������ͨ�������ļ�ϵͳ�Ļ��ƽ��й���ͷ��ʡ�

>�ٸ������е����ӣ�
>��������һ���飬�����м����½������Լ�д�ģ�����һЩ�½����������д�ġ�ʹ�� UFS �� AUFS �����������Ķ��Ȿ��
>ʱ��ͬʱ����������½ں������ѵ��½ڣ����Ⲣ����ζ������½ں������ѵ��½ڱ��ϲ���һƪ�µ��½ڣ�����������Ȼ����>���Զ����Ĵ��ڡ�

�� Docker �У����� AUFS ����ʵ��дʱ���ƣ�copy-on-write�����ܣ��Ӷ�ʹ�������еĲ�������Ի����������Ӱ�졣

������˵��Docker ʹ�� AUFS ʵ�����¹��ܣ�

1.�������ļ�ϵͳ�뾵��ķ��룺 ������ Docker ������һ������ʱ��Docker ʹ�� AUFS ���������ļ�ϵͳ���ص�����������ļ�ϵͳ֮�ϡ������������е��ļ�ϵͳ�ͳ�Ϊ��һ�������Ĳ㣨layer���������������뿪����

2.дʱ���ƣ� �������еĽ�����ͼ����д���������紴�����޸Ļ�ɾ���ļ���ʱ��AUFS ��������Щ�����������������ļ�ϵͳ�д���һ���µ��ļ�������������ֱ���޸Ļ���������ļ�ϵͳ������дʱ���ƵĻ���ʹ�������еĲ�������Ӱ�쵽��������ͬʱ�������������������֮��ķ��롣

![alt text](/WYFdocker/docs/res/images.jpg/)

>����˵�����Ǿ����һ��ʵ�������С�ÿ���������ǻ����ض��ľ��񴴽��ģ����ҿ��Ը�����Ҫ���ò�ͬ�Ĳ�����Ȼ��������
>�����Ǿ��񣬶��Ǿ����һ������ʱʵ����


# ����ʵ��

NewWorkSpace�������������������ļ�ϵͳ�ģ�������CreateReadOnlyLayer��CreateWriteLayer��CreateMountPoint��1.CreateReadOnlyLayer�����½�busybox�ļ��У���busybox.tar��ѹ��busyboxĿ¼�£���Ϊ������ֻ���㡣
2.CreateWriteLayer����������һ����ΪwriteLayer���ļ��У���Ϊ����Ψһ�Ŀ�д�㡣
3.CreateMountPoint�����У����ȴ�����mnt�ļ��У���Ϊ���ص㣬Ȼ���writeLayerĿ¼��busyboxĿ¼mount��mntĿ¼�¡�
4.�����NewParentProcess�����н�����ʹ�õ�������Ŀ¼/root/busybox�滻��/root/mnt��

**ע�⣺���ڱ�����centosϵͳ�ϣ���֧��AUFS�ļ�ϵͳ�����б�д������overlay�ļ�ϵͳ������������°汾��dockerҲ֧��overlay�ļ�ϵͳ��ʵ�ַ�����ͬС�졣**


>OverlayFS��overlay����AUFS��Advanced Multi-Layered Unification Filesystem������ Linux �ں��е������ļ�ϵͳ��������������ʵ���ļ���������Ϲ��صĹ��ܡ�����֮����һЩ����֮����Ҳ��һЩ��֮ͬ����
**����֮����**
1.���Ϲ��أ� ������ OverlayFS ���� AUFS�����Ƕ��������ļ�ϵͳ�㱻���Ϲ��ص�һ��Ŀ¼�У��γ�һ��ͳһ���ļ�ϵͳ��ͼ��
2.дʱ���ƣ�Copy-on-Write���� ��һ���ļ�ϵͳ�㱻�޸�ʱ��OverlayFS �� AUFS ������дʱ���ƵĲ��ԣ���ֻ������Ҫʱ�ŻḴ�Ʊ��޸ĵĲ��֣������������ļ��������������ܺͽ�ʡ�ռ䡣
3.����ļ�ϵͳ�� ���Ƕ�֧����������ʹ�ö���ļ�ϵͳ�㣬ÿ���㶼���԰����ļ���Ŀ¼��������ͳһ���ļ�ϵͳ��ͼ�кϲ���
**��֮ͬ����**
1.ʵ�ֻ��ƣ� OverlayFS �� Linux �ں˵�һ���֣��� Linux �ں˰汾 3.18 ��ʼ�ͱ������ˣ��� AUFS ����һ���������ļ�ϵͳ����Ҫ������ں�ģ����֧�֡�
2.�������ԣ� һ����˵��OverlayFS �����ܷ���� AUFS ���ã��ر����ڴ��ģ���������С�OverlayFS �������Ż����ں��еõ��˳����Ľ����� AUFS ���ڲ����ں˵�һ���֣�������һЩ��������ܱ��ֲ��� OverlayFS��
3.֧�̶ֳȣ� ���� OverlayFS �� Linux �ں˵�һ���֣�����ڴ�����ִ� Linux ���а��ж������õ�֧�֡����֮�£�AUFS ��Ҫ������ں�ģ�飬���ҿ�����Ҫ��������ú�ά����
�ȶ��Ժͳ���ȣ� ���� OverlayFS �� Linux �ں˵�һ���֣�������ȶ��Ժͳ���ȷ�����ܸ��ɿ�һЩ���� AUFS ��Ȼ��һЩ���������ܿ��ܽϺã������ȶ��Ժͳ���ȿ��ܲ��� OverlayFS��
**�ܵ���˵��OverlayFS �ǵ�ǰ Linux ���������µ���ѡ�ļ�ϵͳ���������������ȶ��ɿ����� AUFS ����һ����ѡ������������һЩ�ض��ĳ���������**

>**tips**
1.aufs��������ֿ�д���ֻ��ģ�
���� AUFS �У��ļ�ϵͳ�Ĳ�νṹ�ǰ���˳��ѵ��ģ���ײ��ǻ�������㣬������ǿ�д�㡣��ˣ�AUFS ����ͨ���鿴�ļ�ϵͳ�Ĳ�νṹ�������ĸ��ļ��ǿ�д�㣬�ĸ��ǻ�������㡣ͨ����AUFS �Ὣ�����Ϊ��д�㣬���ڽ���д����ʱ����д���������������д���ϣ�����Ӱ��ײ�Ļ�������㡣(�ļ�ϵͳ�Ĳ�νṹ�ǰ��չ��ص�ǰ��˳��ȷ���ġ��� AUFS �����������ļ�ϵͳ�У�����˳������˲�νṹ��ÿ�����λ��)
�������� AUFS �ļ�ϵͳ�а���˳�� a, b, c, d ���ص� /mnt����ô��
��㣨�������ص��Ǹ��㣩�ǿ�д�㡣����������У�d ���ǿ�д�㡣
����Ĳ㶼��ֻ���㡣��ˣ�a, b, c ����������ֻ���㡣
�������������У�����д�������������ڿ�д�� d �ϣ�����ȡ�����ᰴ�մ���㵽��ײ��˳����С����һ���ļ����ϲ���ڣ������ڵ����²��ͬ���ļ������ֻ��ƽ���"����" (overlay)�������ļ�ϵͳ�ܹ���Ч����϶��ֻ���㣬ͬʱ������ṩ��д����������ԡ�

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
	// ȥ��ĩβ�Ļ��з�
	upperDirPath := string(upperDir[:len(upperDir)-1])
	workDirPath := string(workDir[:len(workDir)-1])
	// ������������
	cmd := exec.Command("sudo", "mount", "-t", "overlay",
		"overlay",
		"-o", fmt.Sprintf("lowerdir=%s:%s,upperdir=%s,workdir=%s", aPath, bPath, upperDirPath, workDirPath),
		cPath)

    //aufsʵ��
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

## ʵ��Ч��
����������������һ���ļ����۲������;�����ļ�ϵͳ�Ƿ����仯��
�����۲죺

![alt text](/WYFdocker/docs/res/image2.png)

�������۲죺
![alt text](/WYFdocker/docs/res/image.png)

