package machines

import "time"

type Machine struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Location  *string   `json:"location"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateMachineInput struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Location string `json:"location"`
}
