package subsystems


type ResourceConfig struct{
MemoryLimit string
CpuShare    string
CpuSet      string
}

type Subsystem interface{

	
    Name() string

	
	Set(CgroupPath string,res *ResourceConfig) error

	
	Apply(CgroupPath string,pid int) error

	
	Remove(CgroupPath string) error


}


var SubsystemsIns=[]Subsystem{ 
   &CpuShareSubsystem{},
   &CpuSetSubsystem{},
   &MemorySubSystem{},
}