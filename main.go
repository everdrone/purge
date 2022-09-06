package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/docker/go-units"
)

type Memory struct {
	Total int
	Free  int
}

func dropPageCache() error {
	mu := sync.Mutex{}
	mu.Lock()

	if err := os.WriteFile("/proc/sys/vm/drop_caches", []byte("3"), 0644); err != nil {
		fmt.Println("failed to drop caches:", err)
		mu.Unlock()

		return err
	}

	mu.Unlock()
	return nil
}

func parseMemInfo() *Memory {
	file, err := os.OpenFile("/proc/meminfo", os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("unable to open meminfo:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	res := Memory{}

	for scanner.Scan() {
		k, v := parseLine(scanner.Text())
		switch k {
		case "MemTotal":
			res.Total = v
		case "MemFree":
			res.Free = v
		}
	}

	return &res
}

func parseLine(line string) (key string, val int) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return
	}

	key = strings.TrimSpace(parts[0])
	a := parts[1][:len(parts[1])-2]
	val, _ = strconv.Atoi(strings.TrimSpace(a))
	return
}

type Flags struct {
	Help bool
}

func parseArgs(args []string) (*Flags, error) {
	if len(args) == 0 {
		return &Flags{}, nil
	}

	var showHelp bool

	for _, arg := range args {
		switch arg {
		case "-h", "--help":
			showHelp = true
		default:
			return &Flags{Help: showHelp}, fmt.Errorf("unknown argument: %s", arg)
		}
	}

	return &Flags{Help: showHelp}, nil
}

func main() {
	flags, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if flags.Help {
		fmt.Println(`Drop page cache, dentries and inodes

Usage:
  sudo purge [-h|--help]`)
		os.Exit(0)
	}

	if os.Geteuid() != 0 {
		fmt.Println("This program must be run as root")
		os.Exit(1)
	}

	before := parseMemInfo()

	syscall.Sync()

	if err := dropPageCache(); err != nil {
		os.Exit(1)
	}

	after := parseMemInfo()

	fmt.Printf(
		"Before : %s\nAfter  : %s\nFree   : %.1f%%\n",
		units.HumanSizeWithPrecision(float64(after.Total)*units.KB, 2),
		units.HumanSizeWithPrecision(float64(after.Free)*units.KB, 2),
		float64(after.Free-before.Free)/float64(after.Total)*100,
	)
}
