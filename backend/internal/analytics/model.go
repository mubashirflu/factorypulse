package analytics

type MachineAnalytics struct {
	MachineID     int     `json:"machine_id"`
	MachineName   string  `json:"machine_name"`
	TotalAlerts   int     `json:"total_alerts"`
	CriticalCount int     `json:"critical_count"`
	AvgRepairMins float64 `json:"avg_repair_minutes"`
}

type DowntimeSummary struct {
	MachineID    int    `json:"machine_id"`
	MachineName  string `json:"machine_name"`
	CriticalHits int    `json:"critical_hits"`
}
