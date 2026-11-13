package server_http

import (
	"fmt"
	"strconv"

	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"

	"github.com/pavlo67/common/common/logger"
)

type Endpoints []Endpoint

func (eps Endpoints) Join(joinerOp joiner.Operator) error {
	for i, ep := range eps {

		if err := joinerOp.Join(&eps[i], ep.InternalKey); err != nil {
			return errors.CommonError(err, fmt.Sprintf("can't join %#v as Endpoint with key '%s'", ep, ep.InternalKey))
		}
	}

	return nil
}

const onHandleEndpoints = "on server_http.HandleEndpoints()"

func (c *Config) HandleEndpoints(srvOp Operator, l logger.Operator) error {
	if c == nil {
		return nil
	}

	if srvOp == nil {
		return errors.New(onHandleEndpoints + ": srvOp == nil")
	}

	for key, ep := range c.EndpointsSettled {
		if ep.Endpoint == nil {
			return fmt.Errorf(onHandleEndpoints+": endpoint %s to handle is nil (%#v)", key, ep)
		}
		if err := srvOp.HandleEndpoint(key, c.Prefix+ep.Path, *ep.Endpoint); err != nil {
			return fmt.Errorf(onHandleEndpoints+": handling %s, %s, %#v got %s", key, ep.Path, ep, err)
		}
	}

	return nil
}

// joining endpoints -----------------------------------------------------

func (c *Config) CompleteWithJoiner(joinerOp joiner.Operator, host string, port int) error {
	if c == nil {
		return errors.New("no server_http.Config to be completed")
	}

	var portStr string
	if port > 0 {
		portStr = ":" + strconv.Itoa(port)
	}
	c.Host, c.Port = host, portStr

	for key, ep := range c.EndpointsSettled {
		//if c.EndpointsSettled[key].Endpoint != nil {
		//	continue
		//}

		if endpoint, ok := joinerOp.Interface(key).(Endpoint); ok {
			ep.Endpoint = &endpoint
			c.EndpointsSettled[key] = ep
		} else if endpointPtr, _ := joinerOp.Interface(key).(*Endpoint); endpointPtr != nil {
			ep.Endpoint = endpointPtr
			c.EndpointsSettled[key] = ep
		} else {
			return fmt.Errorf("no server_http.Endpoint joined with key %s", key)
		}
	}

	return nil
}

// TODO: be careful, it's method for http-client only
// TODO: be careful, it shouldn't be used on server side because it uses non-initiated (without starter.Operator.Run()) endpoints
func (c *Config) CompleteDirectly(endpoints Endpoints, host string, port int, ignoreAbsent bool) error {
	if c == nil {
		return errors.New("no server_http.Config to be completed")
	}

	var portStr string
	if port > 0 {
		portStr = ":" + strconv.Itoa(port)
	}
	c.Host, c.Port = host, portStr

EP_SETTLED:
	for key, epSettled := range c.EndpointsSettled {
		// TODO??? use epSettled.InternalKey to correct the main key value

		for _, ep := range endpoints {
			if ep.InternalKey == key {
				epSettled.Endpoint = &ep
				c.EndpointsSettled[key] = epSettled
				continue EP_SETTLED
			}
		}
		if !ignoreAbsent {
			return fmt.Errorf("no server_http.Endpoint with key %s", key)
		}
		fmt.Printf("no server_http.Endpoint with key %s\n", key)
	}

	return nil
}
