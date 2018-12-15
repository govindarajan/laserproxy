package helper_test

import (
	"fmt"
	"testing"

	"github.com/govindarajan/laserproxy/helper"
)

func Test_GetIPs(t *testing.T) {
	//ipRoutes := helper.GetLocalIPs()
}

func Test_GetHostIps(t *testing.T) {
	ips, err := helper.GetHostIPs("google.com")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ips)
}
