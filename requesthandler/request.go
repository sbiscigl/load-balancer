package requesthandler

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/sbiscigl/load-balancer/entities"
	"github.com/sbiscigl/load-balancer/server"
)

/*RequestHandler a type for routing requests to different instances*/
/*of the serever                                                   */
type RequestHandler struct {
	servers *server.HealthMap
	client  http.Client
}

/*New consrtuctor for RequestHandler*/
func New(s *server.HealthMap) *RequestHandler {
	return &RequestHandler{
		servers: s,
		client:  http.Client{},
	}
}

/*HandleRequest handles a request and sends it to a server instance*/
func (rh *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s, found := rh.servers.FindHealthy()
	if found == false {
		/*return 503 and response with information*/
		w.WriteHeader(503)
		w.Write(entities.NewResponse(503, "there are no healthy servers").ToJSON())
	} else {
		req, err := http.NewRequest(r.Method, "http://"+s.Host, r.Body)
		if err != nil {
			log.Println("could not create request")
		}
		/*Loop through headers*/
		for name, headers := range r.Header {
			name = strings.ToLower(name)
			for _, h := range headers {
				req.Header.Add(name, h)
			}
		}

		/*record in use*/
		rh.servers.IncrimentUseCount(s.Host)
		/*execute request*/
		resp, err := rh.client.Do(req)
		/*record out of use*/
		rh.servers.DecrimentUseCount(s.Host)

		if err != nil {
			log.Println("response error")
		}

		if resp.StatusCode >= 500 && resp.StatusCode < 600 {
			/*set instance to unhealthy*/
			rh.servers.SetHealth(s.Host, false)
			/*retry*/
			rh.ServeHTTP(w, r)
		} else {
			defer resp.Body.Close()
			read, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("error reading body")
			}

			w.WriteHeader(resp.StatusCode)
			w.Write(read)
		}
	}
}
