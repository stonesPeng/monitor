package main // import "github.com/ZenLiu/GMonitor"

func main() {
	if Conf.Server.Enable {
		go Service()
	}
	DoMonitor()
}
