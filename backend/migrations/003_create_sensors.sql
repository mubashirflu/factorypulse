CREATE TABLE sensor_readings (
    id SERIAL PRIMARY KEY,
    machine_id INT NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
    temperature DECIMAL(5,2),
    vibration DECIMAL(5,2),
    pressure DECIMAL(5,2),
    recorded_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_sensor_readings_machine_time ON sensor_readings (machine_id, recorded_at DESC);