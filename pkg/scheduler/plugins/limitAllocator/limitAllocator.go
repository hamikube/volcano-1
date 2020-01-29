

package limitAllocator


import (
	"github.com/hamikube/kube-batch/pkg/scheduler/api"
	"github.com/hamikube/kube-batch/pkg/scheduler/util"
	"math"
	"StrConv"

	)

import (NodeInfo "../../api")
import (job_info "../../api")

func letsSee()  {


	var Nondomsharelist[] float64
	var initalloccpulist[] float64
	var domsharelist[] float64
	var initallocmemorylist[] float64
	var finalallocatedMemory float64
	var finalallocatedCPU float64
	var availableCPU float64
	var reqCPU float64
	type NodeInformationFunctionAccess = NodeInfo.NodeInfo
	type JobInfoFunctionAccess = job_info.JobInfo

	nodeinfo:=NodeInfo.NodeInfo{}

	cpucapacity:=nodeinfo.Capability.MilliCPU

	memorycapacity:=nodeinfo.Capability.Memory
	Pods := NodeInformationFunctionAccess.Pods(nodeinfo)
	// get the number of pods, count the amount in array and convert to float 64
	var numberOfPod = float64(len(Pods));
	//var fairshareCpu float64
	fairshareCpu := cpucapacity/numberOfPod

	fairshareMemory:=memorycapacity/numberOfPod



	for _, pod :=range Pods {

		initResreq := JobInfoFunctionAccess.Clone(pod)

		RequestedNodeObject := initResreq.Allocated
		reqCPU := RequestedNodeObject.MilliCPU

		reqMemory := RequestedNodeObject.Memory

		dominantshare:= math.Max(reqCPU/float64(cpucapacity), reqMemory/float64(memorycapacity))

		domsharelist = append(domsharelist, dominantshare)

		Nondominantshare := math.Min(reqCPU/float64(cpucapacity), reqMemory/float64(memorycapacity))

		Nondomsharelist = append(Nondomsharelist, Nondominantshare)

	}
		
		for _, dmlist :=range domsharelist {
			// how will we access this reqCPU as it is a moving variable
			if dmlist <= fairshareCpu {

				allocatedCpu := fairshareCpu

				initalloccpulist = append(initalloccpulist, allocatedCpu)
			} else {

				allocatedCpu := reqCPU

				initalloccpulist = append(initalloccpulist, allocatedCpu)
			}
			// where else is allocated memory referenced
			if dmlist <= fairshareMemory {

				allocatedMemory := fairshareMemory

				initallocmemorylist := append(initallocmemorylist, allocatedMemory)
			} else {

				allocatedMemory := dmlist

				initallocmemorylist = append(initallocmemorylist,allocatedMemory)

			}
		}

	sumallocCPU := float64(0)

	for _, alloc_CPU := range  initalloccpulist{

		sumallocCPU += alloc_CPU

	}

	sumallocMemory := float64(0)

	for _, alloc_memory := range   initallocmemorylist{

		sumallocMemory += alloc_memory

	}
	availableCPU = cpucapacity - sumallocCPU

	availableMemory := memorycapacity - sumallocMemory


	for _, initialCPU :=range initalloccpulist {

		finalallocatedCPU = (initialCPU * availableCPU) / sumallocCPU
	}


	for _, initialMemory :=range initalloccpulist {

		finalallocatedMemory= initialMemory * availableMemory // sumallocMemory

	}

	updateAPI(finalallocatedMemory, finalallocatedCPU)


}
/**
 * this needs to update the resource in the node info
 */
func updateAPI(finalAllocatedMemory float64, finalallocatedCPU float64){
	StringOfAllocatedMemory := strconv.FormatFloat(finalAllocatedMemory, 'E', -1, 64)
	StringOfAllocatedCPU := strconv.FormatFloat(finalallocatedCPU, 'E', -1, 64)

	api.NewResource(util.BuildResourceList(StringOfAllocatedMemory, StringOfAllocatedCPU))
}