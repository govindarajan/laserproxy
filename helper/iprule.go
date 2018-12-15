package helper

import (
	"bytes"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/govindarajan/laserproxy/logger"
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

// ConfigureRoute used to add route for every IP address.
// So that, we can choose IP on the go while making request.
func ConfigureRoute() {
	// Get all the IP address
	// For each, Create Route
	ips, e := GetLocalIPs()
	if e != nil {
		logger.LogError("ConfigureRoute: " + e.Error())
		return
	}
	for i, ip := range ips {
		logger.LogDebug("Configuring route for " + ip)
		e := CreateOrUpdateRoute(i, ip)
		if e != nil {
			logger.LogError(e.Error())
		}
	}

}
