package analytics

import (
	"context"

	"factorypulse/backend/internal/database"
)

func GetMachineAnalytics() ([]MachineAnalytics, error) {
	query := `
		SELECT 
			m.id,
			m.name,
			COUNT(a.id) AS total_alerts,
			COUNT(CASE WHEN a.severity = 'CRITICAL' THEN 1 END) AS critical_count,
			COALESCE(AVG(
				CASE WHEN mj.status = 'COMPLETED' 
				THEN EXTRACT(EPOCH FROM (mj.updated_at - mj.created_at)) / 60 
				END
			), 0) AS avg_repair_minutes
		FROM machines m
		LEFT JOIN alerts a ON a.machine_id = m.id
		LEFT JOIN maintenance_jobs mj ON mj.machine_id = m.id
		GROUP BY m.id, m.name
		ORDER BY total_alerts DESC
	`

	rows, err := database.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []MachineAnalytics
	for rows.Next() {
		var r MachineAnalytics
		err := rows.Scan(&r.MachineID, &r.MachineName, &r.TotalAlerts, &r.CriticalCount, &r.AvgRepairMins)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}
