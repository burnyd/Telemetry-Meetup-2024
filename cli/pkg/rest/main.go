package rest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/burnyd/Telemetry-Meetup-2024/cli/models"
)

type RestClient struct {
	Server string
	Url    string
}

func (c RestClient) MakeRestCall(method string, requestBody []byte) []byte {
	var req *http.Request
	var err error

	if method == "GET" {
		req, err = http.NewRequest("GET", c.Server+c.Url, nil)

	} else if method == "POST" {
		req, err = http.NewRequest("POST", c.Server+c.Url, bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
	} else if method == "DELETE" {
		req, err = http.NewRequest("DELETE", c.Server+c.Url, nil)
	} else {
		fmt.Println("HTTP METHOD NOT SUPPORTED")
	}

	if err != nil {
		log.Println("response error", err)
	}
	req.Header.Add("Accept", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("response error", err)
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return responseData
}

func (c RestClient) GetTargets() []string {
	c.Url = "/api/v1/config/targets"
	GetTargets := c.MakeRestCall("GET", nil)
	var Targetslice []string
	Targetresponse := make(map[string]interface{})

	_ = json.Unmarshal(GetTargets, &Targetresponse)

	for k, _ := range Targetresponse {
		Targetslice = append(Targetslice, k)
	}
	return Targetslice
}

func (c RestClient) GetTargetData() []models.Target {
	TargetSlice := c.GetTargets()
	TargetData := models.Target{}
	var Targets []models.Target
	for _, target := range TargetSlice {
		c.Url = c.Url + "/" + target
		_ = json.Unmarshal(c.MakeRestCall("GET", nil), &TargetData)
		//fmt.Println(TargetData)
		Targets = append(Targets, TargetData)
		c.Url = strings.ReplaceAll(c.Url, "/"+target, "")
	}
	return Targets
}

func (c RestClient) GetSubs() []string {
	c.Url = "/api/v1/config"
	GetPaths := c.MakeRestCall("GET", nil)
	PathData := models.Subscriptions{}
	var Paths []string
	_ = json.Unmarshal(GetPaths, &PathData)

	for _, p := range PathData.Subs.Sub1.Paths {
		Paths = append(Paths, p)
	}
	return Paths
}

func (c RestClient) GetLeader() (string, error) {
	c.Url = "/api/v1/cluster"
	GetLeader := c.MakeRestCall("GET", nil)
	LeaderModel := models.Leader{}
	if err := json.Unmarshal(GetLeader, &LeaderModel); err != nil {
		return "", err
	}

	return LeaderModel.Leader, nil
}

func (c RestClient) PostTarget(Dev models.NewTarget) {
	c.Url = "/api/v1/config/targets"
	jsonPost, err := json.Marshal(Dev)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print("Adding device " + Dev.Name + "\n")
	Post := c.MakeRestCall("POST", jsonPost)
	if len(Post) != 0 {
		fmt.Println("Issue adding device " + Dev.Name)
	}
}

func (c RestClient) DeleteTarget(Dev models.NewTarget) {
	c.Url = "/api/v1/config/targets/" + Dev.Name
	Delete := c.MakeRestCall("DELETE", nil)
	if len(Delete) != 0 {
		fmt.Println("Issue removing device Device " + Dev.Name)
		fmt.Println(string(Delete))
	}
}
