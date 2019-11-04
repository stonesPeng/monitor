package main

import (
	"encoding/json"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/docker"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func Service() {
	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		tk := req.Header.Get("TOKEN")
		if tk != Conf.Server.Token {
			writer.WriteHeader(401)
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(200)
		writer.Write(tempMonitor())
	})
	http.HandleFunc("/docker", func(writer http.ResponseWriter, req *http.Request) {
		tk := req.Header.Get("TOKEN")
		if tk != Conf.Server.Token {
			writer.WriteHeader(401)
			return
		}
		if !tryInitClient() {
			writer.WriteHeader(500)
			return
		}
		switch req.Method {
		case "POST":
			id := req.URL.Query().Get("id")
			name := req.URL.Query().Get("name")
			if len(id) == 0 && len(name) == 0 {
				writer.WriteHeader(400)
				return
			}
			if len(id) != 0 {
				log.Printf("will start docker container of id %v", id)
				if startContainer(id) {
					writer.WriteHeader(200)
				} else {
					writer.WriteHeader(500)
				}
			} else {
				log.Printf("will start docker container of name %v", name)
				if startContainerByName(name) {
					writer.WriteHeader(200)
				} else {
					writer.WriteHeader(500)
				}
			}
			return
		case "PUT":
			id := req.URL.Query().Get("id")
			name := req.URL.Query().Get("name")
			if len(id) == 0 && len(name) == 0 {
				writer.WriteHeader(400)
				return
			}
			if len(id) != 0 {
				log.Printf("will start docker container of id %v", id)
				if stopContainer(id) {
					writer.WriteHeader(200)
				} else {
					writer.WriteHeader(500)
				}
			} else {
				log.Printf("will start docker container of name %v", name)

				if stopContainerByName(name) {
					writer.WriteHeader(200)
				} else {
					writer.WriteHeader(500)
				}
			}
			return

		default:
			all := req.URL.Query().Get("all")
			d, _ := json.Marshal(containers(len(all) != 0 && strings.ToLower(all) == "true"))
			writer.WriteHeader(200)
			writer.Write(d)
		}
	})
	log.Printf("start server on %v \n", Conf.Server.Addr)
	if e := http.ListenAndServe(Conf.Server.Addr, nil); e != nil {
		log.Fatal("start server error", e)
	}
}

type Status struct {
	Host        host.InfoStat             `json:"host"`
	Temperature []host.TemperatureStat    `json:"temperature"`
	User        []host.UserStat           `json:"user"`
	Cpu         float64                   `json:"cpu"`
	Memory      mem.VirtualMemoryStat     `json:"memory"`
	Swap        mem.SwapMemoryStat        `json:"swap"`
	Disk        []disk.UsageStat          `json:"disk"`
	Process     []ProcessInfo             `json:"process"`
	Docker      []docker.CgroupDockerStat `json:"docker"`
	Timestamp   int64                     `json:"timestamp"`
}
type ProcessInfo struct {
	Pid    int32                   `json:"pid"`
	Name   string                  `json:"name"`
	Memory *process.MemoryInfoStat `json:"memory"`
	CPU    float64                 `json:"cpu"`
}

func tempMonitor() []byte {
	r := new(Status)
	i, _ := host.Info()
	r.Host = *i
	t, _ := host.SensorsTemperatures()
	r.Temperature = t
	u, _ := host.Users()
	r.User = u

	path, _ := disk.Partitions(true)
	if r.Disk == nil {
		r.Disk = make([]disk.UsageStat, 0, len(path))
	}
	for _, x := range path {
		s, _ := disk.Usage(x.Mountpoint)
		if s != nil {
			r.Disk = append(r.Disk, *s)
		}
	}
	v, _ := mem.VirtualMemory()
	r.Memory = *v
	s, _ := mem.SwapMemory()
	r.Swap = *s
	r.Cpu = cpuTime(Conf.Cpu.Duration)
	p, _ := process.Processes()
	if r.Process == nil {
		r.Process = make([]ProcessInfo, 0, len(p))
	}
	for _, ps := range p {
		pm, _ := ps.MemoryInfo()
		pc, _ := ps.CPUPercent()
		pn, _ := ps.Name()
		r.Process = append(r.Process, ProcessInfo{
			Pid:    ps.Pid,
			Name:   pn,
			Memory: pm,
			CPU:    pc,
		})
	}

	_, err := exec.LookPath("docker")
	if err == nil {
		dc, _ := docker.GetDockerStat()
		r.Docker = dc
	}
	r.Timestamp = time.Now().Unix()
	data, _ := json.Marshal(r)
	return data
}

/*func tempMonitor(dur time.Duration, close chan bool) <-chan string {
	out := make(chan string)
	path, _ := disk.Partitions(true)
	go func(d time.Duration, c <-chan bool, o chan<- string) {
	l:
		for {
			select {
			case <-c:
				break l
			default:
				buf := new(bytes.Buffer)
				diskStat := make([]string, 0, len(path))
				for _, x := range path {
					s, _ := disk.Usage(x.Mountpoint)
					diskStat = append(diskStat, fmt.Sprintf(`{"path":"%v","total":%v,"used":%v,"used_percent":%.2f}`, s.Path, s.Total, s.Used, s.UsedPercent))
				}
				dc, _ := docker.GetDockerStat()
				dockerStat := make([]string, 0, len(dc))
				for _, t := range dc {
					dockerStat = append(dockerStat, fmt.Sprintf(`{"conainter":"%v","containerId":"%v","status":%v}`, t.Name, t.ContainerID, t.Running))
				}
				time.Sleep(d)
			}
		}
	}(dur, close, out)
	return out
}*/
