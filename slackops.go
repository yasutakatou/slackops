/*
 * help your operations like shell prompt when use slack.
 *
 * @author    yasutakatou
 * @copyright 2020 yasutakatou
 * @license   MIT License, Apache License 2.0, GNU General Public License v3.0
 */
package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"yasutakatou/chatops/winctl"

	"gopkg.in/ini.v1"

	prompt "github.com/c-bata/go-prompt"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

var (
	currentHwnd, targetHwnd uintptr
	debug                   bool
	wTitle, Sender, Enter   string
	ShortCuts               []csvData
	Converts                []csvData
	sleepWait               int
	singleLine              bool
)

type csvData struct {
	KEY string
	VAL string
}

func executor(in string) {
	if len(in) == 0 {
		return
	}
	if debug == true {
		fmt.Println("Your input: " + in)
	}

	for i := 0; i < len(Converts); i++ {
		if strings.Index(in, Converts[i].KEY) != -1 {
			in = strings.Replace(in, Converts[i].KEY, Converts[i].VAL, -1)
		}
	}
	ChangeTarget(targetHwnd)
	robotgo.TypeStr(in)
	robotgo.MilliSleep(sleepWait)
	if singleLine == true {
		robotgo.KeyTap("enter")
	}
	ChangeTarget(currentHwnd)

}

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func main() {
	currentHwnd = winctl.GetWindow("GetForegroundWindow", debug)

	_Debug := flag.Bool("debug", false, "[-debug=debug mode (true is enable)]")
	_Config := flag.String("config", ".slackops", "[-config=config file)]")

	flag.Parse()

	debug = bool(*_Debug)
	Config := string(*_Config)

	if Exists(Config) == true {
		loadConfig(Config)
	} else {
		fmt.Printf("Fail to read config file: %v\n", Config)
		os.Exit(1)
	}

	targetHwnd = winctl.FocusWindow(targetHwnd, currentHwnd, wTitle, debug)

	if targetHwnd == 0 {
		fmt.Printf("missing title: %s\n", wTitle)
		os.Exit(1)
	}

	go func() {
		ShortCutDo()
	}()

	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
	)
	p.Run()

	os.Exit(0)
}

func loadConfig(filename string) {
	loadOptions := ini.LoadOptions{}
	loadOptions.UnparseableSections = []string{"TITLE", "SHORTCUT", "CONVERT", "SEND", "ENTER", "WAIT", "SINGLELINE"}

	cfg, err := ini.LoadSources(loadOptions, filename)
	if err != nil {
		fmt.Printf("Fail to read config file: %v", err)
		os.Exit(1)
	}

	setSingleConfigStr(&wTitle, "TITLE", cfg.Section("TITLE").Body())
	setSingleConfigCSV(&ShortCuts, "SHORTCUT", cfg.Section("SHORTCUT").Body())
	setSingleConfigCSV(&Converts, "CONVERT", cfg.Section("CONVERT").Body())
	setSingleConfigStr(&Sender, "SEND", cfg.Section("SEND").Body())
	setSingleConfigStr(&Enter, "ENTER", cfg.Section("ENTER").Body())
	setSingleConfigInt(&sleepWait, "WAIT", cfg.Section("WAIT").Body())
	setSingleConfigBool(&singleLine, "SINGLELINE", cfg.Section("SINGLELINE").Body())
}

func setSingleConfigBool(config *bool, configType, datas string) {
	if debug == true {
		fmt.Println(" -- " + configType + " --")
	}
	for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(datas, -1) {
		if len(v) > 0 {
			if v == "Y" || v == "y" {
				*config = true
			} else {
				*config = false
			}
		}
		if debug == true {
			fmt.Println(v)
		}
	}
}

func setSingleConfigInt(config *int, configType, datas string) {
	if debug == true {
		fmt.Println(" -- " + configType + " --")
	}
	for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(datas, -1) {
		if len(v) > 0 {
			*config, _ = strconv.Atoi(v)
		}
		if debug == true {
			fmt.Println(v)
		}
	}
}

func setSingleConfigStr(config *string, configType, datas string) {
	if debug == true {
		fmt.Println(" -- " + configType + " --")
	}
	for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(datas, -1) {
		if len(v) > 0 {
			*config = v
		}
		if debug == true {
			fmt.Println(v)
		}
	}
}

func setSingleConfigCSV(config *[]csvData, configType, datas string) {
	if debug == true {
		fmt.Println(" -- " + configType + " --")
	}
	for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(datas, -1) {
		if len(v) > 0 {
			if strings.Index(v, ",") != -1 {
				strs := strings.Split(v, ",")
				*config = append(*config, csvData{KEY: strs[0], VAL: strs[1]})
			}
		}
		if debug == true {
			fmt.Println(v)
		}
	}
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func ShortCutDo() {
	EvChan := hook.Start()

	//shiftFlag := false
	ctrlFlag := false

	for ev := range EvChan {
		if ev.Kind == 4 { //KeyHold = 4, Shift = 160, Ctrl = 162
			// if int(ev.Rawcode) == 160 {
			// 	shiftFlag = true
			// }

			if int(ev.Rawcode) == 162 {
				ctrlFlag = true
			}
			//if shiftFlag == true && ctrlFlag == true {
			if ctrlFlag == true {
				keyTaps(string(ev.Rawcode), int(ev.Rawcode))
			}
		}

		if ev.Kind == 5 { //KeyUp = 5
			// if int(ev.Rawcode) == 160 {
			// 	shiftFlag = false
			// }
			if int(ev.Rawcode) == 162 {
				ctrlFlag = false
			}
		}
	}
}

func keyTaps(keystr string, keyint int) {
	for i := 0; i < len(ShortCuts); i++ {
		if keystr == ShortCuts[i].KEY {
			if debug == true {
				fmt.Println("Shortcut", ShortCuts[i].KEY, ShortCuts[i].VAL, keystr, keyint)
			}
			robotgo.TypeStr(ShortCuts[i].VAL)
			robotgo.MilliSleep(sleepWait)
		}
	}
	if singleLine == false {
		if keystr == Sender {
			if debug == true {
				fmt.Println("Sender:", keystr)
			}
			ChangeTarget(targetHwnd)
			robotgo.MilliSleep(sleepWait)
			robotgo.KeyTap("enter")
			robotgo.MilliSleep(sleepWait)
			ChangeTarget(currentHwnd)
		}
		if keystr == Enter {
			if debug == true {
				fmt.Println("enter", keystr)
			}
			ChangeTarget(targetHwnd)
			robotgo.MilliSleep(sleepWait)
			robotgo.KeyToggle("shift", "down")
			robotgo.KeyTap("enter")
			robotgo.KeyToggle("shift", "up")
			robotgo.MilliSleep(sleepWait)
			ChangeTarget(currentHwnd)
		}
	}
}

func ChangeTarget(setHwnd uintptr) bool {
	breakCounter := 10

	for {
		if setHwnd != winctl.GetWindow("GetForegroundWindow", debug) {
			winctl.SetActiveWindow(winctl.HWND(setHwnd), debug)
			time.Sleep(time.Duration(100) * time.Millisecond)
		} else {
			return true
		}
		breakCounter--
		if breakCounter == 0 {
			return false
		}
	}
}
