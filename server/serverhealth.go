package server

import (
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/sbiscigl/load-balancer/entities"
	"github.com/sbiscigl/load-balancer/params"
)

const (
	intervalCheck = time.Second * 2
	debug         = false
)

/*HealthMap type for dealing with server use and health*/
type HealthMap struct {
	serverMap   map[string]*Server
	usageChanel chan IncrimentMessage
}

/*IncrimentMessage message that will be used to change the usage number*/
/*on a host*/
type IncrimentMessage struct {
	Host string
	Num  int
}

/*NewServerHealthMap constructor for a HealthMap*/
func NewServerHealthMap(p *params.Params) *HealthMap {
	serverMap := make(map[string]*Server, 0)
	for _, host := range p.GetServerAddresses() {
		serverMap[host] = NewServer(host, true, 0)
	}
	health := &HealthMap{
		serverMap:   serverMap,
		usageChanel: make(chan IncrimentMessage),
	}
	/*start health check pool*/
	go health.checkHealthOnServers()
	return health
}

/*FindHealthy To balence traffic we use a counter that shows how many requests*/
/*are currently being routed through a endpoint, thus we can spread out*/
/*requests to the least utilized server instance*/
func (hm *HealthMap) FindHealthy() (*Server, bool) {
	var s *Server
	found := false
	requestMin := math.MaxInt32
	for _, server := range hm.serverMap {
		if server.IsHealthy && server.IsUsed < requestMin {
			s = server
			found = true
		}
	}
	return s, found
}

/*PrintMap a debugging statement*/
func (hm *HealthMap) PrintMap() {
	log.Println("{")
	for k, v := range hm.serverMap {
		log.Println("\t[" + k + ":" + v.ToString())
	}
	log.Println("}")
}

/*SetHealth sets health of one of the servers*/
func (hm *HealthMap) SetHealth(host string, health bool) {
	hm.serverMap[host].IsHealthy = health
}

/*PublishUsage publishes the usage of one of the instances*/
func (hm *HealthMap) PublishUsage(host string, num int,
	incrimentChan chan<- IncrimentMessage) {
	incrimentChan <- IncrimentMessage{Host: host, Num: num}
}

/*ConsumeUsage consumes from the channel and actually edits the instances*/
func (hm *HealthMap) ConsumeUsage(incrimentChan <-chan IncrimentMessage) {
	edit := <-incrimentChan
	hm.serverMap[edit.Host].IsUsed += edit.Num
}

/*EditUsage edits a hosts usage without exposing underlying logic*/
func (hm *HealthMap) EditUsage(host string, num int) {
	go hm.PublishUsage(host, num, hm.usageChanel)
	go hm.ConsumeUsage(hm.usageChanel)
}

/*RemoveServerAddresses removes server address from parameters*/
func (hm *HealthMap) RemoveServerAddresses(host string) {
	key := ""
	found := false
	for k := range hm.serverMap {
		if k == host {
			found = true
			key = k
		}
	}
	if found {
		delete(hm.serverMap, key)
	}
}

func (hm *HealthMap) checkHealthOnServers() {
	client := &http.Client{}
	for {
		for k, v := range hm.serverMap {
			resp, err := client.Get("http://" + k + "/_health")
			if err != nil {
				log.Println("error in requesting health -- server most likely offline" +
					" removing " + k + " from available servers")
				/*Remove server from server list*/
				hm.RemoveServerAddresses(k)
			} else {
				if resp.StatusCode != 200 {
					log.Println("health check failing on host: " + k + " with status: " +
						resp.Status)
				}
				defer resp.Body.Close()
				read, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Println("error reading body")
				}

				health := entities.NewHealthCheckResponseFromBytes(read).IsHealthy()
				hm.serverMap[k] = NewServer(v.Host, health, v.IsUsed)
			}
		}
		/*debugging for looking at the health status of serevers, is very verbose*/
		if debug {
			hm.PrintMap()
		}
		/*sleep on check interval*/
		time.Sleep(intervalCheck)
	}
}
