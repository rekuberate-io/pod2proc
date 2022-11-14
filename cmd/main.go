package main

import (
	"flag"
	"fmt"
	"github.com/rekuberate-io/pod2proc"

	"k8s.io/klog/v2"
)

var (
	podUID      = flag.String("pUID", "", "Pod unique ID")
	containerID = flag.String("cID", "", "Container ID")
)

func main() {
	if *podUID == "" || *containerID == "" {
		klog.Fatalln(fmt.Errorf("invalid arguments"))
	}

	procID, err := pod2proc.GetProcessIdFromPod(*podUID, *containerID, false)
	if err != nil {
		klog.Fatalln(err)
	}

	fmt.Println(procID)
}

func init() {
	klog.InitFlags(nil)
	flag.Parse()
}

func exit() {
	exitCode := 10
	klog.FlushAndExit(klog.ExitFlushTimeout, exitCode)
}
