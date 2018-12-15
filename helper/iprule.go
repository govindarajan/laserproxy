package helper

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// CreateOrUpdateRoute used to create a route table for the given IP.
// Assumption: Table = 100+id
func CreateOrUpdateRoute(id int, ip string) error {
	// ip rule add from <ip> table <rtable>
	// ip rule add to <ip> table <rtable>
	table := 100 + id
	if e := clearOldRule(table); e != nil {
		return e
	}
	cmd := exec.Command("/sbin/ip", "rule", "add", "from", ip, "table", strconv.Itoa(table))
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	_, err = exec.Command("ip", "rule", "add", "to", ip, "table", strconv.Itoa(table)).Output()
	if err != nil {
		return err
	}

	return nil
}

func clearOldRule(tableId int) error {
	out, err := exec.Command("ip", "rule", "show").Output()
	if err != nil {
		return err
	}
	lines := strings.Split(string(out), "\n")
	rules := 0
	for _, line := range lines {
		regex := regexp.MustCompile(":.*" + strconv.Itoa(tableId))
		if regex.MatchString(line) {
			fmt.Println("Match", line)
			rules++
		}
	}

	for i := 0; i < rules; i++ {
		_, err := exec.Command("ip", "rule", "del", "table", strconv.Itoa(tableId)).Output()
		if err != nil {
			return err
		}
	}
	return nil
}
