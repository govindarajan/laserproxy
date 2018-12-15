package monitor

//Results is the output of monitor package
//Results contains  the ordered local interfaces and ordered targets ips
//Ordered means ordered by min packet loss and rtt
type Results struct {
	Interfaces []string
	TargetIPs  []string
}
