package main

import (
	"encoding/json"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	netUtil "github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"net"
	"strings"
	"time"
)

var ns = netSpeed{}

func setUpTimer(conn net.Conn, ipType int) {
	go poolSpeed()
	go startProbe()

	for {
		load1, load5, load15 := getLoad()
		memoryTotal, memoryUsed := getMemory()
		swapTotal, swapUsed := getSwap()
		diskTotal, diskUsed := getDisk()
		networkIn, networkOut := getTraffic()
		tcp, udp, ps, th := getVitalCounts()
		online4 := ipType == 4
		online6 := ipType == 6
		data := sysInfo{
			Uptime:      getUptime(),
			Load1:       load1,
			Load5:       load5,
			Load15:      load15,
			MemoryTotal: memoryTotal,
			MemoryUsed:  memoryUsed,
			SwapTotal:   swapTotal,
			SwapUsed:    swapUsed,
			HDDTotal:    diskTotal,
			HDDUsed:     diskUsed,
			Cpu:         getCpuPercent(),
			IpStatus:    true,
			NetworkRx:   ns.netrx,
			NetworkTx:   ns.netrx,
			NetworkIn:   networkIn,
			NetworkOut:  networkOut,
			Online4:     online4,
			Online6:     online6,
			Tcp:         tcp,
			Udp:         udp,
			Process:     ps,
			Thread:      th,
			Ping10010:   lostRate["10010"] * 100,
			Ping189:     lostRate["189"] * 100,
			Ping10086:   lostRate["10086"] * 100,
			Time10010:   pingTime["10010"],
			Time189:     pingTime["189"],
			Time10086:   pingTime["10086"],
		}

		// 转换为Json数据
		jsonStr, err := json.Marshal(data)
		if err != nil {
			logger.Fatalf("Can not convert to json: %s", err.Error())
		}

		logf("Submitting data: %s", jsonStr)

		// 发送实时数据
		_, err = conn.Write([]byte("update " + string(jsonStr) + "\n"))
		if err != nil {
			logger.Fatalf("Can not send data: %s", err.Error())
		}

		// 转换间隔为时间间隔
		inter := time.Second * time.Duration(config.Interval)

		// 间隔
		time.Sleep(inter)
	}
}

// 获取负载信息
func getLoad() (float64, float64, float64) {
	Load, err := load.Avg()
	if err != nil {
		return 0.0, 0.0, 0.0
	}
	return Load.Load1, Load.Load5, Load.Load15
}

// 获取CPU使用率
func getCpuPercent() float64 {
	Cpu, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0.0
	}
	return Cpu[0]
}

// 获取内存信息
func getMemory() (uint64, uint64) {
	Mem, err := mem.VirtualMemory()
	if nil != err {
		return 0, 0
	}
	return Mem.Total / 1024, Mem.Used / 1024
}

// 获取启动时间
func getUptime() uint64 {
	up, err := host.Uptime()
	if nil != err {
		return 0
	}
	return up
}

// 获取交换空间信息
func getSwap() (uint64, uint64) {
	Swap, err := mem.SwapMemory()
	if nil != err {
		return 0, 0
	}
	return Swap.Total / 1024, Swap.Used / 1024
}

// 获取硬盘信息
func getDisk() (uint64, uint64) {
	// 获取所有分区
	partitions, err := disk.Partitions(true)

	if nil != err {
		return 0, 0
	}

	// 总空间及使用空间变量
	var Total, Used uint64
	var validFS = []string{"ext4", "ext3", "ext2", "reiserfs", "jfs", "btrfs", "fuseblk", "zfs", "simfs", "ntfs", "fat32", "exfat", "xfs", "apfs"}
	var disks = map[string]string{}

	// 循环所有分区
	for _, d := range partitions {
		if _, ok := disks[d.Device]; !ok {
			for _, fsType := range validFS {
				if fsType == strings.ToLower(d.Fstype) {
					disks[d.Device] = d.Mountpoint
				}
			}
		}
	}

	for _, mountPoint := range disks {
		// 读取分区使用情况
		Disk, err := disk.Usage(mountPoint)
		if nil != err {
			continue
		}
		// 加上空间总量
		Total += Disk.Total
		// 加上使用总量
		Used += Disk.Used
	}

	return Total / 1024 / 1024, Used / 1024 / 1024
}

// 获取流量（入，出）
func getTraffic() (uint64, uint64) {
	stats, err := linuxproc.ReadNetworkStat("/proc/net/dev")

	if err != nil {
		return 0.0, 0.0
	}

	for _, stat := range stats {
		if strings.Contains(stat.Iface, "eth") {
			return stat.RxBytes, stat.TxBytes
		}
	}

	return 0.0, 0.0
}

func getVitalCounts() (int, int, int, int) {
	tcp, udp, ps, th := 0, 0, 0, 0
	processes, err := process.Processes()

	if err == nil {
		ps = len(processes)
	}

	tcps, err := netUtil.Connections("tcp")

	if err == nil {
		tcp = len(tcps)
	}

	udps, err := netUtil.Connections("udp")

	if err == nil {
		udp = len(udps)
	}

	return tcp, udp, ps, th
}

// 获取实时网速
func poolSpeed() {
	for {
		// 读取所有网卡网速
		Net, err := netUtil.IOCounters(true)
		if nil != err {
			break
		}
		// 定义网速保存变量
		var rx, tx uint64
		// 循环网络信息
		for _, nv := range Net {
			// 去除多余信息
			if "lo" == nv.Name || strings.Contains(nv.Name, "tun") {
				continue
			}
			// 加上网速信息
			rx += nv.BytesRecv
			tx += nv.BytesSent
		}

		// 暂停一秒
		time.Sleep(time.Second)

		// 重新读取网络信息
		Net, err = netUtil.IOCounters(true)
		if nil != err {
			break
		}
		// 网络信息保存变量
		var rx2, tx2 uint64
		// 循环网络信息
		for _, nv := range Net {
			// 去除多余信息
			if "lo" == nv.Name || strings.Contains(nv.Name, "tun") {
				continue
			}
			// 加上网速信息
			rx2 += nv.BytesRecv
			tx2 += nv.BytesSent
		}

		ns.netrx = rx2 - rx
		ns.nettx = tx2 - tx
	}
}
