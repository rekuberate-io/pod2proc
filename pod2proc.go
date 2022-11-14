package pod2proc

import (
	"fmt"
	"github.com/fntlnz/mountinfo"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"unicode"
)

const (
	procMountInfoPath      = "mountinfo"
	procFsPathEnvKey       = "NODE_PROC"
	procFsPathDefaultValue = "/proc"
)

var ProcFsPath string

func init() {
	procFsPath, found := os.LookupEnv(procFsPathEnvKey)
	if !found {
		ProcFsPath = procFsPathDefaultValue
		return
	}

	ProcFsPath = procFsPath
}

//GetProcessIdFromPod retrieves the process ID for a container living in a kubernetes pod. The "mountinfo" of a kubernetes pod should look like
// /kubepods/burstable/pod{podUID}/{containerID} or kubepods/besteffort/pod{podUID}/{containerID} e.g: /kubepods/burstable/pod286d1910-a652-4d2f-929b-e300d9d9ed83/58184e99d9d04874131829262cb1cc1658cded5376b00a26ee9d8c524858ee67
func GetProcessIdFromPod(podUID, containerID string, getFullPath bool) (string, error) {
	files, err := ioutil.ReadDir(ProcFsPath)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		directoryName := file.Name()
		if !unicode.IsDigit(rune(directoryName[0])) {
			continue
		}

		mountInfoCollection, err := getMountInfo(path.Join(ProcFsPath, directoryName, procMountInfoPath))
		if err != nil {
			continue
		}

		for _, mountInfo := range mountInfoCollection {
			root := mountInfo.Root
			podContainerCombo := path.Join(podUID, containerID)

			if strings.Contains(root, podContainerCombo) {
				if getFullPath {
					return path.Join(ProcFsPath, directoryName), nil
				}

				return directoryName, nil
			}
		}
	}

	return "", fmt.Errorf("no process found for pod: %s and container: %s", podUID, containerID)
}

func getMountInfo(fd string) ([]mountinfo.Mountinfo, error) {
	file, err := os.Open(fd)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return mountinfo.ParseMountInfo(file)
}
