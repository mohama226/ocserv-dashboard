package repository

import (
	"context"
)

type SystemdRepository struct {
}

type SystemdRepositoryInterface interface {
	Status(ctx context.Context) (string, error)
}

func NewSystemdRepository() *SystemdRepository {
	return &SystemdRepository{}
}

func (s *SystemdRepository) Status(ctx context.Context) (string, error) {
	//cmd := exec.CommandContext(
	//	ctx,
	//	"sudo", "systemctl", "show", "ocserv",
	//	"-p", "Id",
	//	"-p", "Description",
	//	"-p", "ActiveState",
	//	"-p", "SubState",
	//	"-p", "UnitFileState",
	//	"-p", "MainPID",
	//	"-p", "ExecMainStartTimestamp",
	//	"-p", "MemoryCurrent",
	//	"-p", "CPUUsageNSec",
	//	"-p", "TasksCurrent",
	//	"--no-page",
	//)
	//
	//var out bytes.Buffer
	//var stderr bytes.Buffer
	//
	//cmd.Stdout = &out
	//cmd.Stderr = &stderr
	//
	//err := cmd.Run()
	//if err != nil {
	//	return "", fmt.Errorf("systemctl error: %v - %s", err, stderr.String())
	//}
	//
	//return out.String(), nil
	return "MainPID=1195\nExecMainStartTimestamp=Fri 2026-04-24 14:35:38 +0330\nMemoryCurrent=351567872\nCPUUsageNSec=105474810000\nTasksCurrent=100\nId=docker.service\nDescription=Docker Application Container Engine\nActiveState=active\nSubState=running\nUnitFileState=enabled\n", nil
}
