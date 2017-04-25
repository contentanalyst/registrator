package rancher

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	dockerapi "github.com/fsouza/go-dockerclient"
)

// BaseUrl - Rancher metadata service base url with API version
const BaseUrl = "http://rancher-metadata/2016-07-29"

// GetPortMappings - Get port mappings for a container from Rancher metadata service
func GetPortMappings(name string) map[dockerapi.Port][]dockerapi.PortBinding {
	portMappings := make(map[dockerapi.Port][]dockerapi.PortBinding)
	url := BaseUrl + "/containers/" + name + "/ports"
	httpClient := &http.Client{Timeout: time.Second * 2}
	req, err := http.NewRequest( "GET", url , nil)
	if err != nil {
		log.Println("Error building get ports request for Rancher metadata service: ", err)
		return portMappings
	}
	req.Header.Add("Accept", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("Error requesting ports from Rancher metadata service: ", err)
		return portMappings
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("Unexpected status code from requesting ports: ", resp.StatusCode)
		return portMappings
	}
	var ports []string
	derr := json.NewDecoder(resp.Body).Decode(&ports)
	if derr != nil {
		log.Println("Error unmarshalling ports json from Rancher metadata service: ", derr)
		return portMappings
	}
	for _, p := range ports {
		// expected string format "<host ip>:<host port>:<port (number/protocol)>"
		sections := strings.Split(string(p), ":")
		port := sections[2]
		published := []dockerapi.PortBinding{ {sections[0], sections[1]}, }
		portMappings[dockerapi.Port(port)] = published
	}
	return portMappings
}

// GetHostIp - Get host's IP from Rancher metadata service
func GetHostIp() (string, error) {
	var ip string
	url := BaseUrl + "/self/host/agent_ip"
	httpClient := &http.Client{Timeout: time.Second * 2}
	resp, err := httpClient.Get(url)
	if err != nil {
		return ip, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return ip, errors.New("Unexpected status code: " + string(resp.StatusCode))
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ip, err
	}
	ip = string(bodyText)
	return ip, nil
}
