package worker

import (
	"context"
	"fmt"
	"time"

	"factorypulse/backend/internal/database"
	"factorypulse/backend/internal/ws"
)

const (
	checkInterval  = 30 * time.Second
	staleThreshold = 60 * time.Second
)

func StartStaleMachineChecker() {
	go func() {
		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()

		for range ticker.C {
			checkStaleMachines()
		}
	}()

	fmt.Println("🔍 Stale machine checker started (runs every", checkInterval, ")")
}

func checkStaleMachines() {
	query := `
		UPDATE machines
		SET status = 'OFFLINE'
		WHERE status != 'OFFLINE'
		AND id IN (
			SELECT m.id FROM machines m
			LEFT JOIN sensor_readings sr ON sr.machine_id = m.id
			GROUP BY m.id
			HAVING MAX(sr.recorded_at) < NOW() - INTERVAL '60 seconds'
			OR MAX(sr.recorded_at) IS NULL
		)
		RETURNING id, name
	`

	rows, err := database.Pool.Query(context.Background(), query)
	if err != nil {
		fmt.Println("❌ Stale checker error:", err)
		return
	}
	defer rows.Close()

	var offlineMachines []string
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err == nil {
			offlineMachines = append(offlineMachines, name)

			ws.GlobalHub.Broadcast(map[string]interface{}{
				"type":       "status_change",
				"machine_id": id,
				"status":     "OFFLINE",
			})
		}
	}

	if len(offlineMachines) > 0 {
		fmt.Println("⚠️  Marked offline (no data received):", offlineMachines)
	}
}
