package pomodoro

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type PomodoroStatus struct {
	Status string
	Time   string
	Cycle  int
}

// StartPomodoro is the entry point to start the Pomodoro timer
func StartPomodoro(workDuration, breakDuration time.Duration, cycles int) {
	openDunst()
	for i := 1; i <= cycles; i++ {
		if i == cycles {
			// Last cycle (no break after work)
			go playSound("/home/chigan/Music/Wow-Sound-Effect.mp3")
			countdown(workDuration, "Working", i)
		} else {
			// Regular cycle (work followed by break)
			go playSound("/home/chigan/Music/Wow-Sound-Effect.mp3")
			countdown(workDuration, "Working", i)

			go playSound("/home/chigan/Music/Commercial-break-Sound-Effect.mp3")
			countdown(breakDuration, "Break", i)
		}
	}
	countdown(3*time.Second, "Finished", cycles)

	// Update JSON file with final status
	updateJSON(PomodoroStatus{
		Status: "Done",
		Time:   "00:00",
		Cycle:  cycles,
	})

	playSound("/home/chigan/Music/Gudjob!.mp3")
}

// openDunst triggers a notification via a shell script
func openDunst() {
	cmd := exec.Command("/home/chigan/.config/scripts/pomodoro.sh")
	cmd.Run()
}

// playSound plays a sound file
func playSound(soundFile string) {
	cmd := exec.Command("ffplay", "-nodisp", "-autoexit", html.EscapeString(soundFile))
	cmd.Run()
}

// countdown runs the countdown timer for work or break periods
func countdown(duration time.Duration, label string, cycle int) {
	for remaining := int(duration.Seconds()); remaining > 0; remaining-- {
		min := remaining / 60
		sec := remaining % 60
		timeStr := fmt.Sprintf("%02d:%02d", min, sec)

		updateJSON(PomodoroStatus{
			Status: label,
			Time:   timeStr,
			Cycle:  cycle,
		})

		time.Sleep(1 * time.Second)
	}
}

// updateJSON writes the current Pomodoro status to a JSON file
func updateJSON(status PomodoroStatus) {
	file, err := os.Create("/tmp/pomodoro.json")
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
		return
	}
	defer file.Close()

	// Marshal the status into JSON format
	jsonData, err := json.Marshal(status)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing to JSON file:", err)
	}
}

// ConfigurePomodoro takes user input and configures the number of cycles, work duration, and break duration
func ConfigurePomodoro() (int, time.Duration, time.Duration) {
	cycles := 3
	workDuration := 25 * time.Minute
	breakDuration := 5 * time.Minute

	// Ask for the number of cycles, use default if no input is provided
	fmt.Print("Enter number of Cycles (press Enter for default 3): ")
	var cyclesInput string
	fmt.Scanln(&cyclesInput)
	if cyclesInput != "" {
		var err error
		cycles, err = strconv.Atoi(cyclesInput)
		if err != nil {
			fmt.Println("Invalid input. Using default value 3 for cycles.")
			cycles = 3
		}
	}

	// Ask if the user wants to use default durations
	fmt.Print("Use default durations (25 min work, 5 min break)? (Y/N): ")
	var response string
	fmt.Scanln(&response)
	if response != "Y" && response != "y" {
		// Allow user to input custom durations
		fmt.Print("Enter work duration in minutes (default 25): ")
		var workInput string
		fmt.Scanln(&workInput)
		if workInput != "" {
			workDurationInt, err := strconv.Atoi(workInput)
			if err == nil {
				workDuration = time.Duration(workDurationInt) * time.Minute
			}
		}

		fmt.Print("Enter break duration in minutes (default 5): ")
		var breakInput string
		fmt.Scanln(&breakInput)
		if breakInput != "" {
			breakDurationInt, err := strconv.Atoi(breakInput)
			if err == nil {
				breakDuration = time.Duration(breakDurationInt) * time.Minute
			}
		}
	}

	return cycles, workDuration, breakDuration
}
