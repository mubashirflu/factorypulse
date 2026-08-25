package maintenance

import "time"

type MaintenanceJob struct {
	ID          int       `json:"id"`
	MachineID   int       `json:"machine_id"`
	AlertID     *int      `json:"alert_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	AssignedTo  *int      `json:"assigned_to"`
	CreatedBy   *int      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateJobInput struct {
	MachineID   int    `json:"machine_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type UpdateStatusInput struct {
	Status string `json:"status" binding:"required"`
}
