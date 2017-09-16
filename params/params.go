package params

import "os"

const (
	serverflag = "-i"
	portflag   = "-p"
)

/*Params type for housing the commandline parameters*/
type Params struct {
	paramMap map[string][]string
}

/*New contructor for parameter type*/
func New() *Params {
	cmdParams := os.Args[1:]
	/*cmdParams map should be even*/
	if len(cmdParams)%2 != 0 {
		panic("malformed cmd line params")
	}
	paramMap := make(map[string][]string, 0)
	for i := 0; i < len(cmdParams); i += 2 {
		if val, ok := paramMap[cmdParams[i]]; ok {
			paramMap[cmdParams[i]] = append(val, cmdParams[i+1])
		} else {
			paramMap[cmdParams[i]] = []string{cmdParams[i+1]}
		}
	}
	return &Params{
		paramMap: paramMap,
	}
}

/*GetServerAddresses returns all addresses of the servers that are available*/
func (p *Params) GetServerAddresses() []string {
	return p.paramMap[serverflag]
}

/*GetPort returns port param for load balancer to run on*/
func (p *Params) GetPort() string {
	port := p.paramMap[portflag]
	if len(port) != 1 {
		panic("you need/can only have one port")
	}
	return port[0]
}
