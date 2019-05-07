package scheduler

import "testing"

func TestCheckStatus(t *testing.T) {
	var wantedStatus Status
	var currentStatusList []Status
	currentStatusList = []Status{
		SCHED_STATUS_INITIALIZING,
		SCHED_STATUS_STARTING,
		SCHED_STATUS_STOPPING,
	}
	wantedStatus = SCHED_STATUS_INITIALIZING
	for _, currentStatus := range currentStatusList {
		if err := checkStatus(currentStatus, wantedStatus, nil); err == nil {
			t.Fatalf("It still can check status with incorrect current status %q!",
				GetStatusDescription(currentStatus))
		}
	}
}
