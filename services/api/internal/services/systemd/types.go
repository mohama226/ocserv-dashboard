package systemd

import "time"

type OcservSystemdStatus struct {
	ID            string `json:"id"`
	Description   string `json:"description"`
	ActiveState   string `json:"active_state"`
	SubState      string `json:"sub_state"`
	UnitFileState string `json:"unit_file_state"`

	MainPID      int       `json:"main_pid"`
	StartTime    time.Time `json:"start_time"`
	Memory       int64     `json:"memory"`
	CPUUsageNSec int64     `json:"cpu_usage_nsec"`
	Tasks        int       `json:"tasks"`
}
