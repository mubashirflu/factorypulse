package alerts

import (
	"context"

	"factorypulse/backend/internal/database"
)

func CreateAlert(machineID int, severity, message string, value, threshold float64) error {
	query := `INSERT INTO alerts (machine_id, severity, message, value, threshold) 
	          VALUES ($1, $2, $3, $4, $5)`

	_, err := database.Pool.Exec(context.Background(), query, machineID, severity, message, value, threshold)
	return err
}

func GetActiveAlerts() ([]Alert, error) {
	query := `SELECT id, machine_id, severity, message, value, threshold, resolved, created_at 
	          FROM alerts 
	          WHERE resolved = FALSE 
	          ORDER BY created_at DESC`

	rows, err := database.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Alert
	for rows.Next() {
		var a Alert
		err := rows.Scan(&a.ID, &a.MachineID, &a.Severity, &a.Message, &a.Value, &a.Threshold, &a.Resolved, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}

func UpdateMachineStatus(machineID int, status string) error {
	query := `UPDATE machines SET status = $1 WHERE id = $2`
	_, err := database.Pool.Exec(context.Background(), query, status, machineID)
	return err
}
