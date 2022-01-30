package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/shirou/gopsutil/v3/process"
)

//csgo.exe, hl2.exe, hl.exe, Pirates!.exe,

//wpb_settings.ini: isAutoRestore | game process | targets to kill
//read .ini
//if file doesnt exist -> create with basics
//if line starts with '#' -> ignore (for comments)
//processes cuz of which we need to run the program
//Goroutine: check every 2s if processes exist (vrednie) + dont check some process which surely wont run again
//log statistics in wpb_stats.log (num times of killed apps, game session time -> create on exit in defer())
//count how many process were killed during the session
//MessageBox to show up after 3s when game is closed
//Next MessageBox shows up asking if user wants to restore (must be optional -> automatic in .ini)
//user can edit file any moment -> check hash file if there are changes -> check

//program is WinMain
//ini -> json
//csgo cpu usage statistics
//msbox new ui
//hourboostr minimizer
//ask if user wants to close program if duplicate is running
//if steam is running in -no-browser then in several seconds close hourboostr
//restore in background hourboostr if steam is closed or if idle for 20 minutes
//log data only if 5+ session
//Windows sounds when error, kill app or success

type Settings struct {
	AutoRestore         bool
	GameProcesses       []string
	ProcessesToKillOnce []string
	ProcessesToKill     []string
}

func main() {
	var err error

	defer func() {
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("\n\nPress any key to continue...")
		fmt.Scanln()
	}()

	/*var s Settings

	s.AutoRestore = true
	s.GameProcesses = []string{"csgo.exe", "hl2.exe", "hl.exe", "Pirates!.exe"}
	s.ProcessesToKill = []string{"explorer.exe", "Rainmeter.exe", "TrafficMonitor.exe", "ElevenClock.exe", "ModernFlyoutsHost.exe"}
	s.ProcessesToKillOnce = []string{"Widgets.exe", "msedgewebview2.exe"}

	b, _ := json.MarshalIndent(s, "", "\t")

	f, _ := os.Create("test.json")
	f.Write(b)
	f.Close()*/

	settingsFilePath, statFilePath, err := workFiles()
	if err != nil {
		return
	}

	_ = statFilePath

	b, err := ioutil.ReadFile(settingsFilePath)
	if err != nil {
		err = fmt.Errorf("failed to read %s, error: %v", settingsFilePath, err)
	}

	var settings Settings
	err = json.Unmarshal(b, &settings)
	if err != nil {
		fmt.Println(err)
	}

	inGame, err := isPlaying(settings.GameProcesses)
	if err != nil {
		return
	}

	_ = inGame

	fmt.Println(inGame)
}

func workFiles() (settings, statistic string, err error) {
	pwd, err := os.Getwd()
	if err != nil {
		err = fmt.Errorf("failed to get directory of this program running in | error: %v", err)
	}

	settings = pwd + "\\wpb_settings.json"
	statistic = pwd + "\\wpb_settings.log"

	return
}

func isPlaying(gameProcNames []string) (bool, error) {
	processes, err := process.Processes()
	if err != nil {
		return false, fmt.Errorf("failed to obtain Windows processes list | error: %v", err)
	}

	for _, proc := range processes {
		name, err := proc.Name()
		if err != nil {
			return false, fmt.Errorf("failed to obtain process name | error: %v", err)
		}

		for _, gpName := range gameProcNames {
			if name == gpName {
				return true, nil
			}
		}
	}

	return false, nil
}

func killProcs(procNames []string) error {
	processes, err := process.Processes()
	if err != nil {
		return fmt.Errorf("failed to obtain Windows processes list | error: %v", err)
	}

	for _, proc := range processes {
		n, err := proc.Name()
		if err != nil {
			return fmt.Errorf("failed to obtain process name | error: %v", err)
		}

		for _, name := range procNames {
			if n == name {
				if err := proc.Kill(); err != nil {
					return fmt.Errorf("failed to kill process %s | error: %v", name, err)
				}
			}
		}
	}

	return nil
	//return fmt.Errorf("process not found")
}
