package soap

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/beevik/etree"
	"github.com/francistm/soap-server/internal/serde"
	"github.com/pkg/errors"
)

type Service struct {
	name    string
	domain  string
	ports   map[string]*Port
	actions map[string]*Action
}

func NewService(name string, domain string) *Service {
	serv := &Service{
		name:    name,
		domain:  domain,
		ports:   make(map[string]*Port, 20),
		actions: make(map[string]*Action, 20),
	}

	if !strings.HasSuffix(serv.domain, "/") {
		serv.domain += "/"
	}

	return serv
}

func (s *Service) AddPort(name string, p *Port) {
	s.ports[name] = p
}

func (s *Service) cacheActions() {
	for portName, port := range s.ports {
		for actionName, action := range port.actions {
			requestName := fmt.Sprintf("%s%s", portName, actionName)
			s.actions[requestName] = action
		}
	}
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.cacheActions()

	switch r.Method {
	case http.MethodPost:
		s.executeAction(w, r)

	case http.MethodGet:
		s.printDefinition(w, r)

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (s *Service) executeAction(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	doc := etree.NewDocument()

	if _, err := doc.ReadFrom(r.Body); err != nil {
		s.handleSoapOutError(w, r, err)
		return
	}

	soapInElem, err := serde.ExtractSoapInElem(doc)

	if err != nil {
		s.handleSoapOutError(w, r, err)
		return
	}

	action, ok := s.actions[soapInElem.Tag]

	if !ok {
		s.handleSoapOutError(w, r, errors.Errorf("unknown request body element %s", soapInElem.Tag))
		return
	}

	soapIn, err := serde.ParseSoapIn(soapInElem, action.in)

	if err != nil {
		s.handleSoapOutError(w, r, err)
		return
	}

	soapAction := s.actions[soapInElem.Tag]
	soapOut, err := soapAction.Run(soapIn)

	if err != nil {
		s.handleSoapOutError(w, r, err)
		return
	}

	s.handleSoapOut(w, r, soapOut, soapAction.out)
}

func (s *Service) handleSoapOut(w http.ResponseWriter, r *http.Request, soapOut, soapOutType interface{}) {
	if soapOutType == nil {
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

func (s *Service) handleSoapOutError(w http.ResponseWriter, r *http.Request, err error) {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)

	bodyElem := serde.BuildEnvelope(doc)
	serde.BuildFaultBody(bodyElem, err)

	doc.Indent(2)
	doc.WriteTo(w)
	w.WriteHeader(http.StatusInternalServerError)

	return
}
