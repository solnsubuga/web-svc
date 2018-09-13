package service

import (
	"fmt"
	"log"
	"net/http"
	"time"

	consul "github.com/hashicorp/consul/api"
)

//Service is a service to register with consul
type Service struct {
	Name        string
	TTL         time.Duration
	ConsulAgent *consul.Agent
}

//New creates a new instance of a service
func New(ttl time.Duration) (*Service, error) {
	svc := new(Service)
	svc.TTL = ttl
	svc.Name = "testsvc"
	config := consul.Config{
		Address: "192.168.64.2:8500",
		Scheme:  "http",
	}

	//create a new consul client
	c, err := consul.NewClient(&config)
	if err != nil {
		return nil, err
	}
	svc.ConsulAgent = c.Agent()

	ok, err := svc.Check()
	if !ok {
		return nil, err
	}
	return svc, nil
}

//RegisterSvc registers the service with consul
func (s *Service) RegisterSvc() error {
	serviceDef := &consul.AgentServiceRegistration{
		Name: s.Name,
		Check: &consul.AgentServiceCheck{
			TTL: s.TTL.String(),
		},
	}
	return s.ConsulAgent.ServiceRegister(serviceDef)
}

//Check determines if the service is ok
func (s *Service) Check() (bool, error) {
	//simulate service is ok
	return true, nil
}

//ServeHttp ...
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	val := "testsvc\n"
	fmt.Fprint(w, val)
	log.Printf("url=\"%s\" remote=\"%s\" status=%d\n",
		r.URL, r.RemoteAddr, http.StatusOK)
}

//UpdateConsul sends a heart beat to the consul service discovery
func (s *Service) UpdateConsul(check func() (bool, error)) {
	ticker := time.NewTicker(s.TTL / 2)
	for range ticker.C {
		s.update(check)
	}
}

func (s *Service) update(check func() (bool, error)) {
	ok, err := check()
	if !ok {
		log.Printf("err=\"Check failed\" msg=\"%s\"", err.Error())
		if agentErr := s.ConsulAgent.FailTTL("service:"+s.Name, err.Error()); agentErr != nil {
			log.Print(agentErr)
		}
	} else {
		if agentErr := s.ConsulAgent.PassTTL("service:"+s.Name, ""); agentErr != nil {
			log.Print(agentErr)
		}
	}
}
