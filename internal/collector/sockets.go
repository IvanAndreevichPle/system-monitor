package collector

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	pb "github.com/IvanAndreevichPle/system-monitor/api/proto"
)

// TCP state constants from /proc/net/tcp
const (
	tcpStateLISTEN      = "0A"
	tcpStateESTABLISHED = "01"
	tcpStateSYN_SENT    = "02"
	tcpStateSYN_RECV    = "03"
	tcpStateFIN_WAIT1   = "04"
	tcpStateFIN_WAIT2   = "05"
	tcpStateTIME_WAIT   = "06"
	tcpStateCLOSE       = "07"
	tcpStateCLOSE_WAIT  = "08"
	tcpStateLAST_ACK    = "09"
	tcpStateCLOSING     = "0B"
)

// parseTCPAddress parses a TCP address from /proc/net/tcp format
// Format: "00000000:0016" where first part is IP (hex) and second is port (hex)
func parseTCPAddress(addr string) (int32, error) {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid address format")
	}

	port, err := strconv.ParseInt(parts[1], 16, 32)
	if err != nil {
		return 0, err
	}

	return int32(port), nil
}

// CollectSocketStats reads /proc/net/tcp, /proc/net/udp, and /proc/net/sockstat
// Returns socket statistics with listening sockets and TCP connection states
func (c *Collector) CollectSocketStats() (*pb.SocketStats, error) {
	var listeningSockets []*pb.ListeningSocket
	tcpStates := &pb.TCPConnectionStates{}

	// Parse TCP sockets
	tcpFile, err := os.Open("/proc/net/tcp")
	if err == nil {
		defer func() {
			_ = tcpFile.Close()
		}()

		scanner := bufio.NewScanner(tcpFile)
		// Skip header
		_ = scanner.Scan()

		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			if len(fields) < 4 {
				continue
			}

			// Parse local address and state
			localAddr := fields[1]
			state := fields[3]

			// Parse port from local address
			port, err := parseTCPAddress(localAddr)
			if err != nil {
				continue
			}

			// Check if it's a listening socket (state 0A = LISTEN)
			if state == tcpStateLISTEN {
				listeningSockets = append(listeningSockets, &pb.ListeningSocket{
					Command:  "unknown", // Would need /proc/<pid>/ to get actual command
					Pid:      0,         // Would need /proc/<pid>/ to get actual PID
					User:     "unknown", // Would need /proc/<pid>/ to get actual user
					Protocol: "TCP",
					Port:     port,
				})
			}

			// Count TCP states
			switch state {
			case tcpStateESTABLISHED:
				tcpStates.Established++
			case tcpStateFIN_WAIT1, tcpStateFIN_WAIT2:
				tcpStates.FinWait++
			case tcpStateSYN_RECV:
				tcpStates.SynReceived++
			case tcpStateTIME_WAIT:
				tcpStates.TimeWait++
			case tcpStateCLOSE_WAIT:
				tcpStates.CloseWait++
			default:
				tcpStates.Other++
			}
		}
	}

	// Parse UDP sockets for listening sockets
	udpFile, err := os.Open("/proc/net/udp")
	if err == nil {
		defer func() {
			_ = udpFile.Close()
		}()

		scanner := bufio.NewScanner(udpFile)
		// Skip header
		_ = scanner.Scan()

		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			if len(fields) < 4 {
				continue
			}

			// Parse local address
			localAddr := fields[1]

			// Parse port from local address
			port, err := parseTCPAddress(localAddr)
			if err != nil {
				continue
			}

			// UDP sockets in /proc/net/udp are typically listening
			listeningSockets = append(listeningSockets, &pb.ListeningSocket{
				Command:  "unknown",
				Pid:      0,
				User:     "unknown",
				Protocol: "UDP",
				Port:     port,
			})
		}
	}

	return &pb.SocketStats{
		ListeningSockets: listeningSockets,
		TcpStates:       tcpStates,
	}, nil
}

