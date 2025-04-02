package main

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

func main() {
	cycles := 3
	workDuration := 25 * time.Minute
	breakDuration := 5 * time.Minute

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

	fmt.Print("Use default durations (25 min work, 5 min break)? (Y/N): ")
	var response string
	fmt.Scanln(&response)
	if response != "Y" && response != "y" {
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

	Pomodoro(workDuration, breakDuration, cycles)
	os.Exit(0)
}

func openDunst() {
	cmd := exec.Command("/home/chigan/.config/scripts/pomodoro.sh")
	cmd.Run()
}

func playSound(soundFile string) {
	cmd := exec.Command("ffplay", "-nodisp", "-autoexit", html.EscapeString(soundFile))
	cmd.Run()
}

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

func Pomodoro(workDuration, breakDuration time.Duration, cycles int) {
	go openDunst()
	for i := 1; i <= cycles; i++ {
		if i == cycles {
			go playSound("/home/chigan/Music/Wow-Sound-Effect.mp3")
			countdown(workDuration, "Working", i)
		} else {
			go playSound("/home/chigan/Music/Wow-Sound-Effect.mp3")
			countdown(workDuration, "Working", i)

			go playSound("/home/chigan/Music/Commercial-break-Sound-Effect.mp3")
			countdown(breakDuration, "Break", i)
		}
	}
	countdown(3*time.Second, "Finished", cycles)

	updateJSON(PomodoroStatus{
		Status: "Done",
		Time:   "00:00",
		Cycle:  cycles,
	})

	playSound("/home/chigan/Music/Gudjob!.mp3")
}

func updateJSON(status PomodoroStatus) {
	file, _ := os.Create("/tmp/pomodoro.json")
	defer file.Close()

	jsonData, _ := json.Marshal(status)
	file.Write(jsonData)
}
