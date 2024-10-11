package ops_monitor

import (
	"errors"
	"fmt"
	"net"
	"opsPilot/internal/pkg/log"
	"time"
)

// CheckThirdPort 检查三方端口连通性
func CheckThirdPort(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	timeout := 5 * time.Second
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		log.Logger.Errorf("Error connecting failed! host: %s err: %v", address, err)
		panic(errors.New("connect port failed"))
		return false
	}
	defer conn.Close()
	return true
}
