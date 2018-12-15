//Package pinger package pings the Hostname or Ipaddress and calculates the Statistics such as
//Number of Packets sent and received , Min Max and Average Round Trip time and The Packet loss statistics.
//todo: Support  IpV6
package pinger

import (
	"errors"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/net/ipv4"

	"golang.org/x/net/icmp"
)

const (
	timeSliceLength   = 8
	protocolICMP      = 1 //IP protocol number [https://en.wikipedia.org/wiki/List_of_IP_protocol_numbers]
	timeIntervalRange = time.Second / 1000
	timeoutSecond     = time.Second * 5
	packetsToBeSent   = 25
)

var (
	ipv4Proto = map[string]string{"ip": "ip4:icmp", "udp": "udp4"}
)

//Pinger represents ICMP packet send/receiver
type Pinger struct {
	//Interval is the wait time between each packet send.  default is 1s
	Interval time.Duration
	//Timeout is the duration before ping exits
	Timeout time.Duration
	//Count is the number of packets to be sent by the pinger
	Count int
	//PacketsSent is the number of Packets sent
	PacketsSent int
	//PacketsRecv is the number of Packets received
	PacketsRecv int
	//rtts is the list of Round Trip times
	rtts []time.Duration
	//OnRecv is called when the pinger receives and process a packet
	OnRecv func(*Packet)
	//OnFinish is called when Pinger exits
	OnFinish func(*Statistics)
	//stop chan bool
	done chan bool
	//IPAddr resolved from the given address string
	IPAddr *net.IPAddr
	//Addr is the Address host
	addr     string
	ipv4     bool
	source   string
	size     int
	sequence int
	network  string
	conn     *icmp.PacketConn
}
type packet struct {
	bytes  []byte
	nbytes int
}

//Packet represents the received and processed ICMP packet
type Packet struct {
	//Round trip time
	Rtt time.Duration
	//IpAddress to ping
	IPAddr *net.IPAddr
	//Nbytes is the number of bytes in the message
	Nbytes int
	//Seq is the ICMP sequence number
	Seq int
}

// Statistics represents the stats of currently running or finished pinger operation
type Statistics struct {
	//PacketsRecv is the number of packets received
	PacketsRecv int
	//PacketsSent is the number of packets sent
	PacketsSent int
	//PacketsLoss is the loss of Packets in Percentage
	PacketsLoss float64
	//IPAddr is the Ip being pinged
	IPAddr *net.IPAddr
	//Addr is the address of the host being pinged
	Addr string
	//Rtts is the 	all of the Round trip times
	Rtts []time.Duration
	//MinRtt is the Minimum Round trip time
	MinRtt time.Duration
	//MaxRtt is the maximum round trip time
	MaxRtt time.Duration
	//AvgRtt is the average Round trip time
	AvgRtt time.Duration
}

//NewPinger returns a ping.Pinger pointer
func NewPinger(addr string) (*Pinger, error) {
	ipaddr, err := net.ResolveIPAddr("ip", addr)
	if err != nil {
		return nil, err
	}
	var ipv4 bool
	if isIpv4(ipaddr.IP) {
		ipv4 = true
	} else {
		return nil, errors.New("IP is not in ipv4 address space")
	}
	return &Pinger{
		IPAddr:   ipaddr,
		addr:     addr,
		Interval: timeIntervalRange, //If the Interval denominator is increased,  then the program runs faster.
		//Now it is configured to send ICMp packets for 1 milli second
		Timeout: timeoutSecond,   //Time to wait for packets from the Host
		Count:   packetsToBeSent, //total Number of Packets to be sent
		network: "ip",            //Set as ip sends a raw privileged ICMP string
		ipv4:    ipv4,
		size:    timeSliceLength,
		done:    make(chan bool),
	}, nil
}

func isIpv4(ip net.IP) bool {
	return len(ip.To4()) == net.IPv4len
}

// Run runs the Pinger. This is a blocking function unitl it exits
func (p *Pinger) Run() error {
	// fmt.Println("Inside Run function")
	// fmt.Println(p)
	return p.run()
}

//run sends the ICMP packets to the targetted IP
func (p *Pinger) run() error {
	conn := new(icmp.PacketConn)
	if conn = p.listen(ipv4Proto[p.network], p.source); conn == nil {
		return errors.New("unable to ping ip address")
	}
	defer conn.Close()
	defer p.finish()
	var wg sync.WaitGroup
	recv := make(chan *packet, 5)
	wg.Add(1)
	p.conn = conn
	go p.recvICMP(conn, recv, &wg)
	err := p.sendICMP(conn)
	if err != nil {
		log.Printf("Error in sending ICMP packets to %s: %v", p.addr, err)
	}
	timeout := time.NewTicker(p.Timeout)
	interval := time.NewTicker(p.Interval)
	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, os.Interrupt)
	signal.Notify(osChan, syscall.SIGTERM)

	for {
		select {
		case <-osChan:
			close(p.done)
		case <-p.done:
			wg.Wait()
			return nil
		case <-timeout.C:
			close(p.done)
			wg.Wait()
			return nil
		case <-interval.C:
			err := p.sendICMP(conn)
			if err != nil {
				log.Printf("Fatal Error sending ICMP packets %v", err)
			}
		case r := <-recv:
			err := p.processPacket(r)
			if err != nil {
				log.Printf("Fatal Error Processing ICMP packets %v", err)
			}
		}
	}
}

//finish runs when the ping configured packet count finishes.
func (p *Pinger) finish() {
	handler := p.OnFinish
	if handler != nil {
		s := p.Statistics()
		handler(s)
	}
}

//Statistics returns the statistics of the pinger
func (p *Pinger) Statistics() *Statistics {
	loss := float64(p.PacketsSent-p.PacketsRecv) / float64(p.PacketsSent) * 100
	var min, max, total time.Duration
	if len(p.rtts) > 0 {
		min = p.rtts[0]
		max = p.rtts[0]
	}
	for _, rtt := range p.rtts {
		if rtt < min {
			min = rtt
		}
		if rtt > max {
			max = rtt
		}
		total += rtt
	}
	stats := Statistics{
		PacketsSent: p.PacketsSent,
		PacketsRecv: p.PacketsRecv,
		PacketsLoss: loss,
		Rtts:        p.rtts,
		Addr:        p.addr,
		IPAddr:      p.IPAddr,
		MaxRtt:      max,
		MinRtt:      min,
	}

	rttLength := len(p.rtts)
	if rttLength == 0 {
		rttLength = 1
	}
	stats.AvgRtt = total / time.Duration(rttLength)
	//fmt.Println(stats)
	return &stats
}

//listen listens for the echo packets
func (p *Pinger) listen(netProto string, source string) *icmp.PacketConn {
	conn, err := icmp.ListenPacket(netProto, source)
	if err != nil {
		log.Printf("Error listening ICMP packets %s\n", err.Error())
		close(p.done)
		return nil
	}
	return conn
}

//sendICMP issues ICMP packets
func (p *Pinger) sendICMP(conn *icmp.PacketConn) error {
	if p.PacketsSent == p.Count || p.PacketsRecv == p.Count {
		close(p.done)
		return nil
	}
	var typ icmp.Type
	typ = ipv4.ICMPTypeEcho

	var dst net.Addr = p.IPAddr
	//convert current time to bytes and send as the ICMP data. This is one efficient way of sending datas.
	//So that the difference in time can be measured when the packet is echoed.
	t := timeToBytes(time.Now())
	bytes, err := (&icmp.Message{
		Type: typ, Code: 0,
		Body: &icmp.Echo{
			ID:   rand.Intn(65535),
			Seq:  p.sequence,
			Data: t,
		},
	}).Marshal(nil)
	if err != nil {
		return err
	}

	for {
		if _, err := conn.WriteTo(bytes, dst); err != nil {
			if neterr, ok := err.(*net.OpError); ok {
				if neterr.Err == syscall.ENOBUFS {
					continue
				}
			}
		}
		p.PacketsSent++
		p.sequence++
		break
	}
	return nil
}

//recvICMP reads the echoed packet and checks its genuinity
func (p *Pinger) recvICMP(conn *icmp.PacketConn, recv chan<- *packet, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-p.done:
			return
		default:
			bytes := make([]byte, 512)
			conn.SetDeadline(time.Now().Add(time.Millisecond * 100))
			n, addr, err := conn.ReadFrom(bytes)
			if err != nil {
				if neterr, ok := err.(*net.OpError); ok {
					if neterr.Timeout() {
						continue
					} else {
						close(p.done)
						return
					}
				}
			}
			if p.conn == conn && p.IPAddr.String() == addr.String() {

				recv <- &packet{bytes: bytes, nbytes: n}
				//fmt.Printf("recv %v : %v : %v\n", recv, addr, n)
			}
		}
	}

}

//Converts the current nano timestamp to bytes
func timeToBytes(t time.Time) []byte {
	nsec := t.UnixNano()
	b := make([]byte, 8)
	for i := uint8(0); i < 8; i++ {
		b[i] = byte((nsec >> ((7 - i) * 8)) & 0xff)
	}
	return b
}

//processPacket processes the received packet and calculates the time elapsed.
func (p *Pinger) processPacket(recv *packet) error {
	var bytes []byte
	var proto int
	bytes = ipv4Payload(recv.bytes)
	proto = protocolICMP

	var msg *icmp.Message
	var err error

	if msg, err = icmp.ParseMessage(proto, bytes[:recv.nbytes]); err != nil {
		return errors.New("unable to parse icmp message")
	}

	if msg.Type != ipv4.ICMPTypeEchoReply {
		//Not an echo reply
		return nil
	}

	outPkt := &Packet{
		Nbytes: recv.nbytes,
		IPAddr: p.IPAddr,
	}

	switch pkt := msg.Body.(type) {
	case *icmp.Echo:
		outPkt.Rtt = time.Since(bytesToTime(pkt.Data[:timeSliceLength]))
		outPkt.Seq = pkt.Seq
		p.PacketsRecv++
	default:
		return errors.New("invalid icmp packets")
	}

	p.rtts = append(p.rtts, outPkt.Rtt)
	handler := p.OnRecv
	if handler != nil {
		handler(outPkt)
	}
	return nil
}

//Converts the bytes to time . The ICMP packets data is nothing but the time converted into bytes.
// This is a reverse function to convert bytes to time so the elapsed time can be calculated
func bytesToTime(b []byte) time.Time {
	var nsec int64
	for i := uint8(0); i < 8; i++ {
		nsec += int64(b[i]) << ((7 - i) * 8)
	}
	return time.Unix(nsec/1000000000, nsec%1000000000)
}

//ipv4Payload checks if the ICMP packet as the IpV4 header.
func ipv4Payload(b []byte) []byte {
	if len(b) < ipv4.HeaderLen {
		return b
	}
	hdrlen := int(b[0]&0x0f) << 2
	return b[hdrlen:]
}
