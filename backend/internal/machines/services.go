package machines

import (
	"context"

	"factorypulse/backend/internal/database"
)

func CreateMachine(input CreateMachineInput) (int, error) {
	var id int
	query := `INSERT INTO machines (name, type, location, status) 
	          VALUES ($1, $2, $3, 'STOPPED') RETURNING id`

	err := database.Pool.QueryRow(context.Background(), query, input.Name, input.Type, input.Location).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetAllMachines() ([]Machine, error) {
	query := `SELECT id, name, type, status, location, created_at FROM machines ORDER BY id`

	rows, err := database.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var machinesList []Machine
	for rows.Next() {
		var m Machine
		err := rows.Scan(&m.ID, &m.Name, &m.Type, &m.Status, &m.Location, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		machinesList = append(machinesList, m)
	}

	return machinesList, nil
}

func GetMachineByID(id int) (Machine, error) {
	var m Machine
	query := `SELECT id, name, type, status, location, created_at FROM machines WHERE id = $1`

	err := database.Pool.QueryRow(context.Background(), query, id).Scan(
		&m.ID, &m.Name, &m.Type, &m.Status, &m.Location, &m.CreatedAt,
	)
	return m, err
}
