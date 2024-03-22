// Find the locked target thing
// Show subs
// Show targets
// Add stuff
// Delete stuff
package main

import (
	"flag"
	"fmt"

	"github.com/burnyd/Telemetry-Meetup-2024/cli/models"
	"github.com/burnyd/Telemetry-Meetup-2024/cli/pkg/rest"
)

const (
	username = "admin"
	password = "admin"
	insecure = true
)

func main() {
	Rest := rest.RestClient{}
	target := flag.String("target", "", "Target address and port example 1.2.3.4:6030")
	username := flag.String("username", "admin", "username for targets")
	password := flag.String("password", "admin", "password for targets")
	insecure := flag.Bool("insecure", true, "Insecure Skip TLS for connections")
	findleader := flag.Bool("gnmicapi", false, "Initial leader finder to figure out which is the leader for gnmic")
	gnmicapi := flag.String("gnmicapi", "http://clab-cluster-gnmic1:7890", "api for gnmic")
	gettargets := flag.Bool("gettargets", false, "Prints out a list of targets")
	getsubs := flag.Bool("getsubs", false, "Prints out a list of subscriptions")
	addtarget := flag.Bool("addtarget", false, "Adds a target")
	delete := flag.Bool("addtarget", false, "Deletes a target")
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
			Name:     *target,
			Address:  *target,
			Username: *username,
			Password: *password,
			Insecure: *insecure,
		}
		Rest.PostTarget(NewDev)
	}
	if *delete {
		NewDev := models.NewTarget{
			Name:     *target,
			Address:  *target,
			Username: *username,
			Password: *password,
			Insecure: *insecure,
		}
		Rest.DeleteTarget(NewDev)
	}
}
