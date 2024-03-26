###  Unlocking Network Telemetry: A Deep Dive into gNMI/Telemetry/OpenConfig/YANG
Repo for the [Programmability and Automation Meetup](https://www.meetup.com/rtp-programmability-and-automation-meetup/events/299674458/)

### Summary

This is a demo example of gNMIC clustering.  This is meant to demo what a gnmi/telemetry collector would look like from a infrastructure perspective.

My thought here is to cluster 3 gNMIC devices together within containerlab and have them use their clustering with consul then using the service discovery aspect prometheus can then pickup on its targets through consul.  As well as gnmic can discover which targets to use.  Overall, you would end up within prometheus all of the metrics needed to create a robust telemetry collector.

The rest api allows for a user to make any sort of CRUD operations on both targets and subscriptions related to gNMIC which I find really amazing.  Where solutions like telegraf are typically either ran as a 1:1 telegraf to switch or 1 telegraf instance per many different devices which reads off of a telegraf.conf file that either needs to be reloaded or needs many changes to.  I prefer this method with gnmic clustering as it allows for a user to potentially lose one instance and still keep their observability working properly as well as have sharing of the load of telemetry.  Since these gNMIC pods/containers are taking in so many grpc connections it can be overwelming depending on the size of a infrastructure.

### Overall

![overall](/images/overall.jpg)

The binary or curl examples control the flow of devices being created, read, updated or deleted.  The gNMIC workers I will call them are responsible for connecting to the targets(switches gNMI ports).  Consul does the service discovery for each gNMIC device.  Prometheus will then scrape for each gNMIC device which is also a exporter.

### Requirements
- Containerlab
- cEOS 4.31.0F

### Start the lab

```
containerlab -t cluster.yaml deploy --reconfigure
```

```
+---+---------------------------+--------------+---------------------------------+-------+---------+----------------+----------------------+
| # |           Name            | Container ID |              Image              | Kind  |  State  |  IPv4 Address  |     IPv6 Address     |
+---+---------------------------+--------------+---------------------------------+-------+---------+----------------+----------------------+
| 1 | clab-cluster-ceos1        | 2913d9213741 | ceoslab:4.31.0F                 | ceos  | running | 172.20.20.2/24 | 2001:172:20:20::4/64 |
| 2 | clab-cluster-ceos2        | 115f613e268b | ceoslab:4.31.0F                 | ceos  | running | 172.20.20.3/24 | 2001:172:20:20::3/64 |
| 3 | clab-cluster-ceos3        | f4e60ef37ba1 | ceoslab:4.31.0F                 | ceos  | running | 172.20.20.4/24 | 2001:172:20:20::2/64 |
| 4 | clab-cluster-consul-agent | d56e23b8efb8 | consul:latest                   | linux | running | 172.20.20.9/24 | 2001:172:20:20::9/64 |
| 5 | clab-cluster-gnmic1       | 5752de2c94a2 | ghcr.io/openconfig/gnmic:latest | linux | running | 172.20.20.7/24 | 2001:172:20:20::7/64 |
| 6 | clab-cluster-gnmic2       | 97133f16a70c | ghcr.io/openconfig/gnmic:latest | linux | running | 172.20.20.5/24 | 2001:172:20:20::5/64 |
| 7 | clab-cluster-gnmic3       | de0355067ba8 | ghcr.io/openconfig/gnmic:latest | linux | running | 172.20.20.6/24 | 2001:172:20:20::6/64 |
| 8 | clab-cluster-prometheus   | 01fbe9f88c0e | prom/prometheus:latest          | linux | running | 172.20.20.8/24 | 2001:172:20:20::8/64 |
+---+---------------------------+--------------+---------------------------------+-------+---------+----------------+----------------------+
```

### Curl examples

Find all the containers names IPs if its not within DNS this is optional.
```
export GNMICAPI1=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' clab-cluster-gnmic1)
export GNMICAPI2=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' clab-cluster-gnmic2)
export GNMICAPI3=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' clab-cluster-gnmic3)
```

Find the leader this is related to clustering.  The leader has all the information due to the way gnmic uses sharding and consul.  So all API interaction needs to leverage the leader.  More about this can be found on the clustering portion of the gnmic docs [here](https://gnmic.openconfig.net/user_guide/HA/).

```
curl --request GET $GNMICAPI1:7890/api/v1/cluster | jq .leader
```

For example the response from this will be in my case any of the random gnmic docker containers.
```
"clab-cluster-gnmic1"
```

So all communication will be through clab-cluster-gnmic1 at this point.

Find all the targets.

curl --request GET $GNMICAPI1:7890/api/v1/config/targets | jq .
```
{
  "172.20.20.2:6030": {
    "name": "172.20.20.2:6030",
    "address": "172.20.20.2:6030",
    "username": "admin",
    "password": "admin",
    "timeout": 10000000000,
    "insecure": true,
    "skip-verify": true,
    "buffer-size": 100,
    "retry-timer": 10000000000,
    "log-tls-secret": false,
    "gzip": false,
    "token": ""
  },
  "172.20.20.3:6030": {
    "name": "172.20.20.3:6030",
    "address": "172.20.20.3:6030",
    "username": "admin",
    "password": "admin",
    "timeout": 10000000000,
    "insecure": true,
    "skip-verify": true,
    "buffer-size": 100,
    "retry-timer": 10000000000,
    "log-tls-secret": false,
    "gzip": false,
    "token": ""
  }
}
```

Find all the subscriptions
```
curl --request GET $GNMICAPI1:7890/api/v1/config/subscriptions | jq .
```

Add a target
```
curl -X POST -H "Content-Type: application/json" -d ' {"name":"172.20.20.4:6030","address":"172.20.20.4:6030","username":"admin","password":"admin","timeout": 10000000000,"insecure": true,"skip-verify": true,"buffer-size": 100,"retry-timer": 10000000000,"log-tls-secret": false,"gzip": false,"token": ""}' $GNMICAPI3:7890/api/v1/config/targets

```

Delete a target
```
curl -X DELETE -H "Content-Type: application/json" $GNMICAPI1:7890/api/v1/config/targets/172.20.20.4:6030
```

### Build the binary optional
### This is a small go based binary I put together to act like a CLI which does the same thing fwiw.

```
go build -o bin/cli main.go && cd cli/bin
```

Find the leader this is related to clustering.  The leader has all the information due to the way gnmic uses sharding and consul.  So all API interaction needs to leverage the leader.  More about this can be found on the clustering portion of the gnmic docs [here](https://gnmic.openconfig.net/user_guide/HA/).

```
./cli -findleader
http://clab-cluster-gnmic1:7890
```

Get the targets.

```
./cli -gnmicapi http://clab-cluster-gnmic1:7890 -gettargets
[172.20.20.2:6030 172.20.20.3:6030]
```

Get the subscriptions.

```
./cli -gnmicapi http://clab-cluster-gnmic1:7890 -getsubs
[/interfaces/interface/state/counters]
```

Add a gNMI target.

```
./cli -gnmicapi http://clab-cluster-gnmic1:7890 -addtarget -target 172.20.20.200:6030 -username admin -password admin -insecure=true
Adding device 172.20.20.200:6030
```

Delete a target.

```
./cli -gnmicapi http://clab-cluster-gnmic1:7890 -delete -target 172.20.20.200:6030
```

### Checking for promtheus metrics.
```
curl clab-cluster-gnmic1:9804/metrics
```

Truncating a bit here for only a few metrics.

```
interfaces_interface_state_counters_in_broadcast_pkts{interface_name="Ethernet2",source="172.20.20.3:6030",subscription_name="sub1"} 461
interfaces_interface_state_counters_in_multicast_pkts{interface_name="Ethernet1",source="172.20.20.3:6030",subscription_name="sub1"} 530
```
