package main

type sysInfo struct {
	Uptime      uint64  `json:"uptime"`
	Load1       float64 `json:"load_1"`
	Load5       float64 `json:"load_5"`
	Load15      float64 `json:"load_15"`
	MemoryTotal uint64  `json:"memory_total"`
	MemoryUsed  uint64  `json:"memory_used"`
	SwapTotal   uint64  `json:"swap_total"`
	SwapUsed    uint64  `json:"swap_used"`
	HDDTotal    uint64  `json:"hdd_total"`
	HDDUsed     uint64  `json:"hdd_used"`
	Cpu         float64 `json:"cpu"`
	NetworkRx   uint64  `json:"network_rx"`
	NetworkTx   uint64  `json:"network_tx"`
	NetworkIn   uint64  `json:"network_in"`
	NetworkOut  uint64  `json:"network_out"`
	IpStatus    bool    `json:"ip_status"`
	Ping10010   float64 `json:"ping_10010"`
	Ping189     float64 `json:"ping_189"`
	Ping10086   float64 `json:"ping_10086"`
	Time10010   int     `json:"time_10010"`
	Time189     int     `json:"time_189"`
	Time10086   int     `json:"time_10086"`
	Online4     bool    `json:"online4"`
	Online6     bool    `json:"online6"`
	Tcp         int     `json:"tcp"`
	Udp         int     `json:"udp"`
	Process     int     `json:"process"`
	Thread      int     `json:"thread"`
}

type netSpeed struct {
	netrx uint64
	nettx uint64
	avgrx uint64
	avgtx uint64
}
