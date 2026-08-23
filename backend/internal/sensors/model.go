package sensors

import "time"

type SensorReading struct {
	ID          int       `json:"id"`
	MachineID   int       `json:"machine_id"`
	Temperature float64   `json:"temperature"`
	Vibration   float64   `json:"vibration"`
	Pressure    float64   `json:"pressure"`
	RecordedAt  time.Time `json:"recorded_at"`
}

type CreateReadingInput struct {
	MachineID   int     `json:"machine_id" binding:"required"`
	Temperature float64 `json:"temperature"`
	Vibration   float64 `json:"vibration"`
	Pressure    float64 `json:"pressure"`
}
