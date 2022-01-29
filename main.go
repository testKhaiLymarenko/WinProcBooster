package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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

func main() {
	var err error

	defer func() {
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("\n\nPress any key to continue...")
		fmt.Scanln()
	}()

	settingsFile, statFile, err := workFiles()
	if err != nil {
		return
	}

	_ = statFile

	gameProcs, err := gameProcs(settingsFile)
	if err != nil {
		return
	}

	isPlaying, err := isInGame(gameProcs)
	if err != nil {
		return
	}

	_ = isPlaying
}

func workFiles() (ini, log string, err error) {
	pwd, err := os.Getwd()
	if err != nil {
		err = fmt.Errorf("failed to get directory of this program running in | error: %v", err)
	}

	ini = pwd + "\\wpb_settings.ini"
	log = pwd + "\\wpb_settings.log"

	return
}

func gameProcs(filePath string) ([]string, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s, error: %v", filePath, err)
	}

	//find "Game processes: "

	procsBuff := strings.Split(string(b), ",")
	var procs []string

	for _, pb := range procsBuff {
		if pb != "" && !strings.Contains(pb, ".exe") {
			procs = append(procs, pb)
		}
	}

	if procs == nil {
		return nil, fmt.Errorf("inappropriate game process names. Be sure it ends with '.exe'")
	}

	return procs, nil
}

func isInGame(names []string) (bool, error) {
	processes, err := process.Processes()
	if err != nil {
		return false, fmt.Errorf("failed to obtain Windows processes list | error: %v", err)
	}

	for _, p := range processes {
		n, err := p.Name()
		if err != nil {
			return false, fmt.Errorf("failed to obtain process name | error: %v", err)
		}

		for _, name := range names {
			if n == name {
				return true, nil
			}
		}
	}

	return false, nil
}

func killProcs(names []string) error {
	processes, err := process.Processes()
	if err != nil {
		return fmt.Errorf("failed to obtain Windows processes list | error: %v", err)
	}

	for _, p := range processes {
		n, err := p.Name()
		if err != nil {
			return fmt.Errorf("failed to obtain process name | error: %v", err)
		}

		for _, name := range names {
			if n == name {
				return p.Kill()
			}
		}
	}
	return fmt.Errorf("process not found")
}
