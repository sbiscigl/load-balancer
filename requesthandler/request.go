package requesthandler

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

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
		/*
		   The load balancer’s job is to take incomming traffic and proxy it to one of
		   the available healthy instances. If none of the instances are healthy, the
		   load balancer should return an HTTP 503 response code.
		*/
		//TODO: custom return type
		w.WriteHeader(503)
	} else {
		req, err := http.NewRequest(r.Method, s.Host, r.Body)
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
		}

		defer resp.Body.Close()
		read, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("error reading body")
		}

		w.WriteHeader(resp.StatusCode)
		w.Write(read)
	}
}
