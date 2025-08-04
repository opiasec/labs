package k8s

import "time"

type ContainerStatusInfo struct {
	Name       string
	Ready      bool
	State      string
	Reason     string
	RestartCnt int32
	ExitCode   int32
	StartedAt  time.Time
}
