// package sensors

// import (
// 	"context"

// 	"factorypulse/backend/internal/database"
// )

// func CreateReading(input CreateReadingInput) error {
// 	query := `INSERT INTO sensor_readings (machine_id, temperature, vibration, pressure)
// 	          VALUES ($1, $2, $3, $4)`

// 	_, err := database.Pool.Exec(context.Background(), query,
// 		input.MachineID, input.Temperature, input.Vibration, input.Pressure)
// 	return err
// }

// func GetLatestReading(machineID int) (SensorReading, error) {
// 	var r SensorReading
// 	query := `SELECT id, machine_id, temperature, vibration, pressure, recorded_at
// 	          FROM sensor_readings
// 	          WHERE machine_id = $1
// 	          ORDER BY recorded_at DESC
// 	          LIMIT 1`

// 	err := database.Pool.QueryRow(context.Background(), query, machineID).Scan(
// 		&r.ID, &r.MachineID, &r.Temperature, &r.Vibration, &r.Pressure, &r.RecordedAt,
// 	)
// 	return r, err
// }

// func GetReadingHistory(machineID int, limit int) ([]SensorReading, error) {
// 	query := `SELECT id, machine_id, temperature, vibration, pressure, recorded_at
// 	          FROM sensor_readings
// 	          WHERE machine_id = $1
// 	          ORDER BY recorded_at DESC
// 	          LIMIT $2`

// 	rows, err := database.Pool.Query(context.Background(), query, machineID, limit)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var readings []SensorReading
// 	for rows.Next() {
// 		var r SensorReading
// 		err := rows.Scan(&r.ID, &r.MachineID, &r.Temperature, &r.Vibration, &r.Pressure, &r.RecordedAt)
// 		if err != nil {
// 			return nil, err
// 		}
// 		readings = append(readings, r)
// 	}
// 	return readings, nil
// }

package sensors

import (
	"context"
	"fmt"

	alerts "factorypulse/backend/internal/alert"
	"factorypulse/backend/internal/database"
)

const (
	vibrationWarning  = 5.0
	vibrationCritical = 7.0
)

func CreateReading(input CreateReadingInput) error {
	query := `INSERT INTO sensor_readings (machine_id, temperature, vibration, pressure) 
	          VALUES ($1, $2, $3, $4)`

	_, err := database.Pool.Exec(context.Background(), query,
		input.MachineID, input.Temperature, input.Vibration, input.Pressure)
	if err != nil {
		return err
	}

	return evaluateThresholds(input)
}

func evaluateThresholds(input CreateReadingInput) error {
	switch {
	case input.Vibration >= vibrationCritical:
		alerts.UpdateMachineStatus(input.MachineID, "CRITICAL")
		return alerts.CreateAlert(
			input.MachineID,
			"CRITICAL",
			fmt.Sprintf("Vibration critically high: %.2f mm/s", input.Vibration),
			input.Vibration,
			vibrationCritical,
		)

	case input.Vibration >= vibrationWarning:
		alerts.UpdateMachineStatus(input.MachineID, "WARNING")
		return alerts.CreateAlert(
			input.MachineID,
			"WARNING",
			fmt.Sprintf("Vibration above normal: %.2f mm/s", input.Vibration),
			input.Vibration,
			vibrationWarning,
		)

	default:
		alerts.UpdateMachineStatus(input.MachineID, "RUNNING")
	}

	return nil
}

func GetLatestReading(machineID int) (SensorReading, error) {
	var r SensorReading
	query := `SELECT id, machine_id, temperature, vibration, pressure, recorded_at 
	          FROM sensor_readings 
	          WHERE machine_id = $1 
	          ORDER BY recorded_at DESC 
	          LIMIT 1`

	err := database.Pool.QueryRow(context.Background(), query, machineID).Scan(
		&r.ID, &r.MachineID, &r.Temperature, &r.Vibration, &r.Pressure, &r.RecordedAt,
	)
	return r, err
}

func GetReadingHistory(machineID int, limit int) ([]SensorReading, error) {
	query := `SELECT id, machine_id, temperature, vibration, pressure, recorded_at 
	          FROM sensor_readings 
	          WHERE machine_id = $1 
	          ORDER BY recorded_at DESC 
	          LIMIT $2`

	rows, err := database.Pool.Query(context.Background(), query, machineID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readings []SensorReading
	for rows.Next() {
		var r SensorReading
		err := rows.Scan(&r.ID, &r.MachineID, &r.Temperature, &r.Vibration, &r.Pressure, &r.RecordedAt)
		if err != nil {
			return nil, err
		}
		readings = append(readings, r)
	}
	return readings, nil
}
