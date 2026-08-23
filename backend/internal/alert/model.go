package alerts

import "time"

type Alert struct {
	ID        int       `json:"id"`
	MachineID int       `json:"machine_id"`
	Severity  string    `json:"severity"`
	Message   string    `json:"message"`
	Value     float64   `json:"value"`
	Threshold float64   `json:"threshold"`
	Resolved  bool      `json:"resolved"`
	CreatedAt time.Time `json:"created_at"`
}
