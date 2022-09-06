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

func main() {
	if os.Geteuid() != 0 {
		fmt.Println("This program must be run as root")
		os.Exit(1)
	}

	before := getMemory()

	syscall.Sync()

	if err := dropPageCache(); err != nil {
		os.Exit(1)
	}

	after := getMemory()

	fmt.Printf(
		"Before : %s\nAfter  : %s\nFree   : %.1f%%\n",
		units.HumanSizeWithPrecision(float64(after.Total)*units.KB, 2),
		units.HumanSizeWithPrecision(float64(after.Free)*units.KB, 2),
		float64(after.Free-before.Free)/float64(after.Total)*100,
	)
}

func getMemory() *Memory {
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
