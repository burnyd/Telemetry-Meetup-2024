//Add curl commands for logic.
//go run main.go -gnmicapi http://clab-cluster-gnmic1:7890 -gettargets
//go run main.go -gnmicapi http://clab-cluster-gnmic1:7890 -addtarget -target 172.20.20.200:6030 -username admin -password admin -insecure=true
//go run main.go -gnmicapi http://clab-cluster-gnmic1:7890 -delete -target 172.20.20.200:6030

package main

import (
	"flag"
	"fmt"

	"github.com/burnyd/Telemetry-Meetup-2024/cli/models"
	"github.com/burnyd/Telemetry-Meetup-2024/cli/pkg/rest"
)

func main() {
	Rest := rest.RestClient{}
	target := flag.String("target", "", "Target address and port example 1.2.3.4:6030")
	username := flag.String("username", "admin", "username for targets")
	password := flag.String("password", "admin", "password for targets")
	insecure := flag.Bool("insecure", true, "Insecure Skip TLS for connections")
	findleader := flag.Bool("findleader", false, "Initial leader finder to figure out which is the leader for gnmic")
	gnmicapi := flag.String("gnmicapi", "http://clab-cluster-gnmic1:7890", "api for gnmic")
	gettargets := flag.Bool("gettargets", false, "Prints out a list of targets")
	getsubs := flag.Bool("getsubs", false, "Prints out a list of subscriptions")
	addtarget := flag.Bool("addtarget", false, "Adds a target")
	delete := flag.Bool("delete", false, "Deletes a target")
	flag.Parse()

	Rest.Server = *gnmicapi
	if *gettargets {
		RestTargets := Rest.GetTargets()
		fmt.Print(RestTargets)
	}
	if *getsubs {
		RestSubs := Rest.GetSubs()
		fmt.Print(RestSubs)
	}
	if *findleader {
		RestLeader, _ := Rest.GetLeader()
		fmt.Println(RestLeader)
	}
	if *addtarget {
		NewDev := models.NewTarget{
			Name:         *target,
			Address:      *target,
			Username:     *username,
			Password:     *password,
			Insecure:     *insecure,
			Skipverify:   true,
			Buffersize:   100,
			RetryTimer:   10000000000,
			Logtlssecret: false,
			Gzip:         false,
			Timeout:      10000000000,
			Token:        "",
		}
		Rest.PostTarget(NewDev)
	}
	if *delete {
		NewDev := models.NewTarget{
			Name: *target,
		}
		Rest.DeleteTarget(NewDev)
	}
}
