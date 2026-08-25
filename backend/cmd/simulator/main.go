// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"math/rand"
// 	"net/http"
// 	"time"
// )

// const (
// 	baseURL    = "http://localhost:8080"
// 	loginEmail = "mubashir@gmail.com"
// 	loginPass  = "mubashir"
// 	machineID  = 1
// )

// type LoginResponse struct {
// 	Token string `json:"token"`
// }

// type Reading struct {
// 	MachineID   int     `json:"machine_id"`
// 	Temperature float64 `json:"temperature"`
// 	Vibration   float64 `json:"vibration"`
// 	Pressure    float64 `json:"pressure"`
// }

// func main() {
// 	fmt.Println("🔧 Starting FactoryPulse simulator...")

// 	token, err := login()
// 	if err != nil {
// 		fmt.Println("❌ Login failed:", err)
// 		return
// 	}
// 	fmt.Println("✅ Logged in, token acquired")

// 	// Baseline values — inke aas paas readings ghoomengi
// 	baseTemp := 70.0
// 	baseVibration := 3.0
// 	basePressure := 6.0

// 	tick := 0
// 	for {
// 		tick++

// 		// Normal random fluctuation
// 		temp := baseTemp + (rand.Float64()*4 - 2)
// 		vibration := baseVibration + (rand.Float64()*1 - 0.5)
// 		pressure := basePressure + (rand.Float64()*0.6 - 0.3)

// 		// Har ~15th reading pe ek "spike" simulate karo (jaise machine problem)
// 		if tick%15 == 0 {
// 			vibration += 4.5
// 			temp += 8
// 			fmt.Println("⚠️  Simulating anomaly spike...")
// 		}

// 		reading := Reading{
// 			MachineID:   machineID,
// 			Temperature: round2(temp),
// 			Vibration:   round2(vibration),
// 			Pressure:    round2(pressure),
// 		}

// 		err := sendReading(token, reading)
// 		if err != nil {
// 			fmt.Println("❌ Failed to send reading:", err)
// 		} else {
// 			fmt.Printf("📊 Sent: temp=%.2f vibration=%.2f pressure=%.2f\n",
// 				reading.Temperature, reading.Vibration, reading.Pressure)
// 		}

// 		time.Sleep(3 * time.Second)
// 	}
// }

// func login() (string, error) {
// 	payload := map[string]string{
// 		"email":    loginEmail,
// 		"password": loginPass,
// 	}
// 	body, _ := json.Marshal(payload)

// 	resp, err := http.Post(baseURL+"/api/auth/login", "application/json", bytes.NewBuffer(body))
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	var loginResp LoginResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
// 		return "", err
// 	}

// 	if loginResp.Token == "" {
// 		return "", fmt.Errorf("empty token, check credentials")
// 	}

// 	return loginResp.Token, nil
// }

// func sendReading(token string, reading Reading) error {
// 	body, _ := json.Marshal(reading)

// 	req, err := http.NewRequest("POST", baseURL+"/api/readings", bytes.NewBuffer(body))
// 	if err != nil {
// 		return err
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+token)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusCreated {
// 		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
// 	}

// 	return nil
// }

// func round2(val float64) float64 {
// 	return float64(int(val*100)) / 100
// }

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const (
	baseURL    = "http://localhost:8080"
	loginEmail = "mubashir@gmail.com"
	loginPass  = "mubashir"
)

type LoginResponse struct {
	Token string `json:"token"`
}

type Machine struct {
	ID int `json:"id"`
}

type Reading struct {
	MachineID   int     `json:"machine_id"`
	Temperature float64 `json:"temperature"`
	Vibration   float64 `json:"vibration"`
	Pressure    float64 `json:"pressure"`
}

func main() {
	fmt.Println("🔧 Starting FactoryPulse simulator...")

	token, err := login()
	if err != nil {
		fmt.Println("❌ Login failed:", err)
		return
	}
	fmt.Println("✅ Logged in, token acquired")

	tick := 0
	for {
		tick++

		machineIDs, err := fetchMachineIDs(token)
		if err != nil {
			fmt.Println("❌ Failed to fetch machines:", err)
			time.Sleep(3 * time.Second)
			continue
		}

		for _, id := range machineIDs {
			baseTemp := 70.0
			baseVibration := 3.0
			basePressure := 6.0

			temp := baseTemp + (rand.Float64()*4 - 2)
			vibration := baseVibration + (rand.Float64()*1 - 0.5)
			pressure := basePressure + (rand.Float64()*0.6 - 0.3)

			if tick%15 == 0 {
				vibration += 4.5
				temp += 8
				fmt.Printf("⚠️  Simulating anomaly spike on machine %d...\n", id)
			}

			reading := Reading{
				MachineID:   id,
				Temperature: round2(temp),
				Vibration:   round2(vibration),
				Pressure:    round2(pressure),
			}

			err := sendReading(token, reading)
			if err != nil {
				fmt.Printf("❌ Failed to send reading for machine %d: %v\n", id, err)
			} else {
				fmt.Printf("📊 Machine %d: temp=%.2f vibration=%.2f pressure=%.2f\n",
					id, reading.Temperature, reading.Vibration, reading.Pressure)
			}
		}

		time.Sleep(3 * time.Second)
	}
}

func login() (string, error) {
	payload := map[string]string{
		"email":    loginEmail,
		"password": loginPass,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(baseURL+"/api/auth/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", err
	}

	if loginResp.Token == "" {
		return "", fmt.Errorf("empty token, check credentials")
	}

	return loginResp.Token, nil
}

func fetchMachineIDs(token string) ([]int, error) {
	req, err := http.NewRequest("GET", baseURL+"/api/machines", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var machines []Machine
	if err := json.NewDecoder(resp.Body).Decode(&machines); err != nil {
		return nil, err
	}

	ids := make([]int, len(machines))
	for i, m := range machines {
		ids[i] = m.ID
	}
	return ids, nil
}

func sendReading(token string, reading Reading) error {
	body, _ := json.Marshal(reading)

	req, err := http.NewRequest("POST", baseURL+"/api/readings", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}

func round2(val float64) float64 {
	return float64(int(val*100)) / 100
}
