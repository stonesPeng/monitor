package main

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/docker"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Monitor struct {
	signal chan os.Signal
	out    chan<- string
	in     []chan string
	close  []chan bool
	host   string
}

func (m *Monitor) Watch(out chan string) {
	log.Println("prepare to watch")
	h, _ := host.Info()
	m.host = h.Hostname
	m.signal = make(chan os.Signal, 1)
	signal.Notify(m.signal, syscall.SIGINT, syscall.SIGTERM)
	m.out = out
	var wg sync.WaitGroup
	m.cpu()
	m.mem()
	m.disk()
	m.docker()
	wg.Add(len(m.in))
	for _, c := range m.in {
		go func(c <-chan string) {
			for v := range c {
				out <- fmt.Sprintf(`{"server":"%+v","message":%+v,"timestamp":%v}`, m.host, v, time.Now().Unix())
			}
		}(c)
	}
l:
	for {
		select {
		case <-m.signal:
			log.Println("closing monitors")
			for _, c := range m.close {
				c <- true
			}
			log.Println("closing channels")
			for _, c := range m.in {
				close(c)
			}
			close(m.out)
			log.Println("quite server monitor")
			break l
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}
func (m *Monitor) cpu() {
	if Conf.Cpu.Enable {
		log.Println("start to watch CPU")
		ch := make(chan string)
		w := make(chan bool)
		m.in = append(m.in, ch)
		m.close = append(m.close, w)
		go func(c chan string, s chan bool) {
		l:
			for {
				select {
				case <-w:
					log.Println("closing CPU monitor")
					break l
				default:
					t := cpuTime(Conf.Cpu.Duration)
					//log.Printf("%.4f\n",t)
					switch {
					case t > Conf.Cpu.Limit:
						c <- fmt.Sprintf(`{"cpu":%.4f}`, t)
					case t == -1:
						c <- fmt.Sprintf(`{"cpu":%.4f}`, t)
					}
					time.Sleep(time.Duration(uint(time.Millisecond) * Conf.Memory.Frequcey))
				}
			}
		}(ch, w)
		log.Println("start watching CPU")
	}
}
func (m *Monitor) mem() {
	if Conf.Memory.Enable {
		log.Println("start to watch memory")
		ch := make(chan string)
		w := make(chan bool)
		m.in = append(m.in, ch)
		m.close = append(m.close, w)
		go func(c chan string, s chan bool) {
		l:
			for {
				select {
				case <-w:
					log.Println("closing memeory monitor")
					break l
				default:
					v, _ := mem.VirtualMemory()
					if v.UsedPercent > Conf.Memory.Limit {
						c <- fmt.Sprintf(`{"memory":%.2f}`, v.UsedPercent)
					}
					time.Sleep(time.Duration(uint(time.Millisecond) * Conf.Memory.Frequcey))
				}
			}
		}(ch, w)
		log.Println("start watching memory")
	}
}
func (m *Monitor) disk() {
	if Conf.Disk.Enable {
		log.Println("start to watch disk")
		ch := make(chan string)
		w := make(chan bool)
		m.in = append(m.in, ch)
		m.close = append(m.close, w)
		paths := make([]DiskPath, 0, 3)
		if Conf.Disk.All {
			p, _ := disk.Partitions(true)
			for _, x := range p {
				paths = append(paths, DiskPath{
					Path:  x.Mountpoint,
					Limit: Conf.Disk.Limit,
				})
			}
		} else {
			paths = append(paths, Conf.Disk.Paths...)
		}
		go func(c chan string, s chan bool, path []DiskPath) {
		l:
			for {
				select {
				case <-s:
					log.Println("closing disk monitor")
					break l
				default:
					for _, x := range path {
						v1, _ := disk.Usage(x.Path)
						if v1.UsedPercent > x.Limit {
							c <- fmt.Sprintf(`{"disk":"%v","used":%.2f}`, v1.Path, v1.UsedPercent)
						}
					}
					time.Sleep(time.Duration(uint(time.Millisecond) * Conf.Disk.Frequcey))
				}
			}
		}(ch, w, paths)
		log.Println("start watching disk")
	}
}
func (m *Monitor) docker() {
	if Conf.Docker.Enable && len(Conf.Docker.Containers) > 0 {
		log.Println("start to watch docker")
		_, err := exec.LookPath("docker")
		if err != nil {
			log.Println("ERROR docker binary not found!")
			return
		}
		ch := make(chan string)
		w := make(chan bool)
		m.in = append(m.in, ch)
		m.close = append(m.close, w)

		go func(c chan string, s chan bool, path []DockerContainer) {
		l:
			for {
				select {
				case <-s:
					log.Println("closing docker monitor")
					break l
				default:
					d, _ := docker.GetDockerStat()
					for _, x := range path {
						t := findContainerIn(d, x.Name, x.Id)
						if !t.Running {
							c <- fmt.Sprintf(`{"conainter":"%v","container_id":"%v","status":%v,"running":%v}`, t.Name, t.ContainerID, t.Status, t.Running)
						}
					}
					time.Sleep(time.Duration(uint(time.Millisecond) * Conf.Docker.Frequcey))
				}
			}
		}(ch, w, Conf.Docker.Containers)
		log.Println("start watching docker")
	}
}
func findContainerIn(c []docker.CgroupDockerStat, name string, id string) *docker.CgroupDockerStat {
	for _, x := range c {
		if x.Name == name || x.ContainerID == id {
			return &x
		}
	}
	return nil
}
func cpuTime(duration int) float64 {
	pc, _ := cpu.Percent(time.Millisecond*time.Duration(duration), false)
	if len(pc) > 0 {
		return pc[0]
	} else {
		return -1
	}
}
func DoMonitor() {
	m := Monitor{}
	out := make(chan string, 3)
	if !Conf.Client.Enable || len(Conf.Client.Url) == 0 {
		go func(o chan string) {
			for x := range o {
				fmt.Println(x)
			}
		}(out)
	} else {
		go func(o chan string) {
			for x := range o {
				Do(Conf.Client.Method, Conf.Client.Url, x)
			}
		}(out)
	}
	m.Watch(out)
}
