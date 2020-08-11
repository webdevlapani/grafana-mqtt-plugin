//+build mage

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
)

var exname string

func getExecutableName(os string, arch string) (string, error) {
	if exname == "" {
		exename, err := getExecutableFromPluginJSON()
		if err != nil {
			return "", err
		}

		exname = exename
	}

	exeName := fmt.Sprintf("%s_%s_%s", exname, os, arch)
	if os == "windows" {
		exeName = fmt.Sprintf("%s.exe", exeName)
	}
	return exeName, nil
}

func getExecutableFromPluginJSON() (string, error) {
	byteValue, err := ioutil.ReadFile(path.Join("src", "plugin.json"))
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return "", err
	}
	executable := result["executable"]
	name, ok := executable.(string)
	if !ok || name == "" {
		return "", fmt.Errorf("plugin.json is missing an executable name")
	}
	return name, nil
}

func findRunningPIDs(exe string) []int {
	pids := []int{}
	out, err := sh.Output("pgrep", "-f", exe)
	if err != nil || out == "" {
		return pids
	}
	for _, txt := range strings.Fields(out) {
		pid, err := strconv.Atoi(txt)
		if err == nil {
			pids = append(pids, pid)
		} else {
			log.Printf("Unable to format %s (%s)", txt, err)
		}
	}
	return pids
}

func killAllPIDs(pids []int) error {
	for _, pid := range pids {
		log.Printf("Killing process: %d", pid)
		err := syscall.Kill(pid, 9)
		if err != nil {
			return err
		}
	}
	return nil
}

// checkLinuxPtraceScope verifies that ptrace is configured as required.
func checkLinuxPtraceScope() error {
	ptracePath := "/proc/sys/kernel/yama/ptrace_scope"
	byteValue, err := ioutil.ReadFile(ptracePath)
	if err != nil {
		return fmt.Errorf("unable to read ptrace_scope: %w", err)
	}
	val := strings.TrimSpace(string(byteValue))
	if val != "0" {
		log.Printf("WARNING:")
		fmt.Printf("ptrace_scope set to value other than 0 (currently: %s), this might prevent debugger from connecting\n", val)
		fmt.Printf("try writing \"0\" to %s\n", ptracePath)
		fmt.Printf("Set ptrace_scope to 0? y/N (default N)\n")

		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			if scanner.Text() == "y" || scanner.Text() == "Y" {
				// if err := sh.RunV("echo", "0", "|", "sudo", "tee", ptracePath); err != nil {
				//      return // Error?
				// }
				log.Printf("TODO, run: echo 0 | sudo tee /proc/sys/kernel/yama/ptrace_scope")
			} else {
				fmt.Printf("Did not write\n")
			}
		}
	}

	return nil
}

func Trace() error {
	// Debug build
	b := build.Build{}
	mg.Deps(b.Debug)

	// 1. kill any running instance
	exeName, err := getExecutableName(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return err
	}
	_ = killAllPIDs(findRunningPIDs(exeName))
	_ = sh.RunV("pkill", "strace")
	if runtime.GOOS == "linux" {
		if err := checkLinuxPtraceScope(); err != nil {
			return err
		}
	}

	// Wait for grafana to start plugin
	pid := -1
	for i := 0; i < 20; i++ {
		pids := findRunningPIDs(exeName)
		if len(pids) > 1 {
			return fmt.Errorf("multiple instances already running")
		}
		if len(pids) > 0 {
			pid = pids[0]
			log.Printf("Found plugin PID: %d", pid)
			break
		}

		log.Printf("Waiting for Grafana to start plugin: %q...", exeName)
		time.Sleep(250 * time.Millisecond)
	}
	if pid == -1 {
		return fmt.Errorf("could not find plugin process %q, perhaps Grafana is not running?", exeName)
	}

	pidStr := strconv.Itoa(pid)
	log.Printf("Attaching strace to plugin process %d", pid)
	if err := sh.RunV("strace",
		"-p", pidStr,
		"-s", "9999",
		"-f",
		"-e", "write"); err != nil {
		return err
	}

	return nil
}

// Default configures the default target.
var Default = build.BuildAll
