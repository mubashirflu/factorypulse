package maintenance

import (
	"context"

	"factorypulse/backend/internal/database"
)

func CreateJob(input CreateJobInput, createdBy int) (int, error) {
	var id int
	query := `INSERT INTO maintenance_jobs (machine_id, title, description, created_by, status) 
	          VALUES ($1, $2, $3, $4, 'OPEN') RETURNING id`

	err := database.Pool.QueryRow(context.Background(), query,
		input.MachineID, input.Title, input.Description, createdBy).Scan(&id)
	return id, err
}

func GetAllJobs() ([]MaintenanceJob, error) {
	query := `SELECT id, machine_id, alert_id, title, description, status, assigned_to, created_by, created_at, updated_at 
	          FROM maintenance_jobs ORDER BY created_at DESC`

	rows, err := database.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []MaintenanceJob
	for rows.Next() {
		var j MaintenanceJob
		err := rows.Scan(&j.ID, &j.MachineID, &j.AlertID, &j.Title, &j.Description,
			&j.Status, &j.AssignedTo, &j.CreatedBy, &j.CreatedAt, &j.UpdatedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func UpdateJobStatus(jobID int, status string) error {
	query := `UPDATE maintenance_jobs SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := database.Pool.Exec(context.Background(), query, status, jobID)
	return err
}

func AssignJob(jobID int, userID int) error {
	query := `UPDATE maintenance_jobs SET assigned_to = $1, status = 'ASSIGNED', updated_at = NOW() WHERE id = $2`
	_, err := database.Pool.Exec(context.Background(), query, userID, jobID)
	return err
}
