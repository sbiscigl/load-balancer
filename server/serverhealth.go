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
	intervalCheck = time.Second * 10
)

/*HealthMap type for dealing with server use and health*/
type HealthMap struct {
	serverMap map[string]*Server
}

/*NewServerHealthMap constructor for a HealthMap*/
func NewServerHealthMap(p *params.Params) *HealthMap {
	serverMap := make(map[string]*Server, 0)
	for _, host := range p.GetServerAddresses() {
		serverMap[host] = NewServer(host, true, 0)
	}
	health := &HealthMap{
		serverMap: serverMap,
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

/*IncrimentUseCount incriments the useage of one of the instances*/
func (hm *HealthMap) IncrimentUseCount(host string) {
	hm.serverMap[host].IsUsed++
}

/*DecrimentUseCount decirments the usages of one of the instances*/
func (hm *HealthMap) DecrimentUseCount(host string) {
	hm.serverMap[host].IsUsed--
}

func (hm *HealthMap) checkHealthOnServers() {
	client := &http.Client{}
	for {
		for k, v := range hm.serverMap {
			resp, err := client.Get("http://" + k + "/_health")
			if err != nil {
				log.Println("error in requesting health")
			}
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
		hm.PrintMap()
		/*sleep on check interval*/
		time.Sleep(intervalCheck)
	}
}
