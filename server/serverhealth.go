package server

import (
	"fmt"
	"math"

	"github.com/sbiscigl/load-balancer/params"
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
	/*TODO: start a go func that polls forever*/
	return &HealthMap{
		serverMap: serverMap,
	}
}

/*FindHealthy To balence traffic we use a counter that shows how many requests*/
/*are currently being routed through a endpoint, thus we can spread out 		  */
/*requests to the least utilized server instance														  */
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
	fmt.Println("{")
	for k, v := range hm.serverMap {
		fmt.Println("\t[" + k + ":" + v.ToString())
	}
	fmt.Println("}")
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
