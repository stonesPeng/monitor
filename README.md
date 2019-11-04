# GMonitor

server monitor of go

## configuration
```yaml
server:
  enable: true     # http server for command
  addr: ""         # default 0.0.0.0:80
  token: ""        # access token
client:            # webclient to post warning to some http server
  enable: false    # enable webclient
  url: ""          # url to put data
  method: ""       # default POST (could be POST,PUT,PATCH..)
memory:
  enable: false     # watch server memeory
  limit: 42         # percent
  frequcey: 1000    # loop time in ms
disk:               # watch server disk
  enable: false
  frequcey: 1000
  limit: 40          # optional
  paths:             # optional
    - path: "c:"     # mount point to watch
      limit: 40      # max used percent
cpu:                 # watch server cpu (
  enable: true
  duration: 100      # measure duration in ms,default 100 ms
  limit: 64.0
  frequcey: 1000
docker:          # watch docker container is running
  enable: false
  frequcey: 0
  containers:
    - id: 123213213sadsf # container id
       name: container   # container name
```
## server api

### fetch server status

method: `GET`

url: `/`

header: `TOKEN:${token}`

```json
{
  "host": {
    "hostname": "centos764",
    "uptime": 3172,
    "bootTime": 1561344011,
    "procs": 270,
    "os": "linux",
    "platform": "centos",
    "platformFamily": "rhel",
    "platformVersion": "7.6.1810",
    "kernelVersion": "3.10.0-957.1.3.el7.x86_64",
    "virtualizationSystem": "",
    "virtualizationRole": "",
    "hostid": "1e3c59fc-0b1f-4e3f-a4e0-75edccb8a20e"
  },
  "temperature": null,
  "user": [
    {
      "user": "peter",
      "terminal": "pts/0",
      "host": "192.168.8.173",
      "started": 1561344042
    }
  ],
  "cpu": 0,
  "memory": {
    "total": 3161473024,
    "available": 1532825600,
    "used": 1386000384,
    "usedPercent": 43.840335611859395,
    "free": 1113698304,
    "active": 1249697792,
    "inactive": 502370304,
    "wired": 0,
    "laundry": 0,
    "buffers": 2166784,
    "cached": 659607552,
    "writeback": 0,
    "dirty": 4096,
    "writebacktmp": 0,
    "shared": 15216640,
    "slab": 145891328,
    "sreclaimable": 48402432,
    "pagetables": 25657344,
    "swapcached": 0,
    "commitlimit": 4936175616,
    "committedas": 4927365120,
    "hightotal": 0,
    "highfree": 0,
    "lowtotal": 0,
    "lowfree": 0,
    "swaptotal": 3355439104,
    "swapfree": 3355439104,
    "mapped": 154009600,
    "vmalloctotal": 35184372087808,
    "vmallocused": 188772352,
    "vmallocchunk": 35183933779968,
    "hugepagestotal": 0,
    "hugepagesfree": 0,
    "hugepagesize": 2097152
  },
  "swap": {
    "total": 3355439104,
    "used": 0,
    "free": 3355439104,
    "usedPercent": 0,
    "sin": 0,
    "sout": 0
  },
  "disk": [
    {
      "path": "/boot",
      "fstype": "xfs",
      "total": 1063256064,
      "free": 820068352,
      "used": 243187712,
      "usedPercent": 22.871979783037474,
      "inodesTotal": 524288,
      "inodesUsed": 348,
      "inodesFree": 523940,
      "inodesUsedPercent": 0.066375732421875
    }
  ],
  "process": [
    {
      "pid": 1,
      "name": "systemd",
      "memory": {
        "rss": 4366336,
        "vms": 128716800,
        "data": 0,
        "stack": 0,
        "locked": 0,
        "swap": 0
      },
      "cpu": 0.09802304231770462
    }
  ],
  "docker": [],
  "timestamp": 1561347183
}

```

## fetch docker status
method: `GET`

url: `/docker?all={false|true}`

header: `TOKEN:${token}`

```json
[
  {
    "Id": "ff842328eac6c52869fe1467d28b931dba8bcf33ce382fc7f19bd4e1af9223ff",
    "Names": [
      "/mysql"
    ],
    "Image": "mysql:5.7",
    "ImageID": "sha256:a1aa4f76fab910095dfcf4011f32fbe7acdb84c46bb685a8cf0a75e7d0da8f6b",
    "Command": "docker-entrypoint.sh mysqld",
    "Created": 1561348363,
    "Ports": [

    ],
    "Labels": {

    },
    "State": "exited",
    "Status": "Exited (1) 6 minutes ago",
    "HostConfig": {
      "NetworkMode": "default"
    },
    "NetworkSettings": {
      "Networks": {
        "bridge": {
          "IPAMConfig": null,
          "Links": null,
          "Aliases": null,
          "NetworkID": "2722f74c9d78e33a64b824c3ad630925ea1c3faf305a4adf0e941b591854a354",
          "EndpointID": "",
          "Gateway": "",
          "IPAddress": "",
          "IPPrefixLen": 0,
          "IPv6Gateway": "",
          "GlobalIPv6Address": "",
          "GlobalIPv6PrefixLen": 0,
          "MacAddress": ""
        }
      }
    },
    "Mounts": [
      {
        "Type": "volume",
        "Name": "28058a39e143b818a107797e1f231de566926667b3a847f68d7e40b42598443c",
        "Source": "",
        "Destination": "/var/lib/mysql",
        "Driver": "local",
        "Mode": "",
        "RW": true,
        "Propagation": ""
      }
    ]
  },
  {
    "Id": "e3c5e5dce58264d5244f68eb22794f0c8eed23583212b81f38abbe6cfde49e4d",
    "Names": [
      "/dreamy_shaw"
    ],
    "Image": "hello-world",
    "ImageID": "sha256:fce289e99eb9bca977dae136fbe2a82b6b7d4c372474c9235adc1741675f587e",
    "Command": "/hello",
    "Created": 1555560728,
    "Ports": [

    ],
    "Labels": {

    },
    "State": "exited",
    "Status": "Exited (0) 2 months ago",
    "HostConfig": {
      "NetworkMode": "default"
    },
    "NetworkSettings": {
      "Networks": {
        "bridge": {
          "IPAMConfig": null,
          "Links": null,
          "Aliases": null,
          "NetworkID": "b60e858ea6e2c9f2462feae70c33e6f3e8067d30f4ff51db383bfec526924e5d",
          "EndpointID": "",
          "Gateway": "",
          "IPAddress": "",
          "IPPrefixLen": 0,
          "IPv6Gateway": "",
          "GlobalIPv6Address": "",
          "GlobalIPv6PrefixLen": 0,
          "MacAddress": ""
        }
      }
    },
    "Mounts": [

    ]
  }
]
```
## start docker container

method: `POST`

url: `/docker?id={container_id}`|| `/docker?name=/{container_name}`

header: `TOKEN:${token}`
## stop docker container

method: `PUT`

url: `/docker?id={container_id}`|| `/docker?name=/{container_name}`

header: `TOKEN:${token}`

## client api

```json
{
  "server": "iz8vbiwllfxxr77qmmlicrz",
  "message": {
    "cpu": 100.0000
  },
  "timestamp": 1561353342
}
{
  "server": "iz8vbiwllfxxr77qmmlicrz",
  "message": {
    "memory": 90.0000
  },
  "timestamp": 1561353342
}
{
  "server": "iz8vbiwllfxxr77qmmlicrz",
  "message": {
    "disk": "/"
    "used": 78.00
  },
  "timestamp": 1561353342
}
{
  "server": "iz8vbiwllfxxr77qmmlicrz",
  "message": {
    "container": "/test_container",
    "container_id":"12312312312312",
    "status": "Exited (0) 2 months ago",
    "running": false
  },
  "timestamp": 1561353342
}
```
## docker image useage
```yml
version: '2'
services:
  monitor:
    image: monitor
    expose:
      - 8080:80
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - /usr/local/bin/docker:/usr/local/bin/docker:ro
      - /config.yml:/app/config.yml
```
