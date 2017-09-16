package server

import "strconv"

/*Server a type that represents a server that we are balencing*/
type Server struct {
	Host      string
	IsHealthy bool
	IsUsed    int
}

/*NewServer constructor fort server instance*/
func NewServer(host string, isHealthy bool, isUsed int) *Server {
	return &Server{
		Host:      host,
		IsHealthy: isHealthy,
		IsUsed:    isUsed,
	}
}

/*ToString a function for converting this to a string for debugging*/
func (s *Server) ToString() string {
	health := strconv.FormatBool(s.IsHealthy)
	usage := strconv.Itoa(s.IsUsed)
	return "[" + s.Host + ":" + health + ":" + usage + "]"
}
