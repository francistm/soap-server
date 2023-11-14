package soap

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/beevik/etree"
	"github.com/francistm/soap-server/internal"
	"github.com/francistm/soap-server/internal/serde"
	"github.com/pkg/errors"
)

type Service struct {
	name      string
	namespace string
	ports     map[string]*Port
	actions   map[string]*Action

	envelopeSpace string
	envelopeNS    map[string]string
}

func NewService(name string, namespace string, opts ...serviceOption) *Service {
	serv := &Service{
		name:      name,
		namespace: namespace,
		ports:     make(map[string]*Port, 20),
		actions:   make(map[string]*Action, 20),
	}

	if !strings.HasSuffix(serv.namespace, "/") {
		serv.namespace += "/"
	}

	for _, opt := range opts {
		opt(serv)
	}

	if serv.envelopeNS == nil {
		serv.envelopeNS = map[string]string{
			"xmlns:soapenv": internal.NsSoap,
		}
	}

	if serv.envelopeSpace == "" {
		serv.envelopeSpace = internal.XmlSoapNs
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
	soapOut, err := soapAction.Run(r.Context(), soapIn)

	if err != nil {
		s.handleSoapOutError(w, r, err)
		return
	}

	s.handleSoapOut(
		w,
		r,
		s.namespace,
		soapInElem.Tag+internal.ElemOutSuffix,
		soapAction.out,
		soapOut,
	)
}

func (s *Service) handleSoapOut(w http.ResponseWriter, r *http.Request, ns, soapOutName string, tSoapOut, soapOut interface{}) {
	if tSoapOut == nil {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	if !serde.IsActionReturnValid(tSoapOut, soapOut) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	doc := etree.NewDocument()
	doc.CreateProcInst("xml", internal.XmlProcInst)

	requestBodyElem := serde.BuildResponseBodyChild(ns, soapOutName, tSoapOut, soapOut)
	envelopeElem := serde.BuildEnvelope(s.envelopeSpace, s.envelopeNS, requestBodyElem)

	doc.AddChild(envelopeElem)

	doc.Indent(2)
	w.WriteHeader(http.StatusOK)
	doc.WriteTo(w)
}

func (s *Service) handleSoapOutError(w http.ResponseWriter, r *http.Request, err error) {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", internal.XmlProcInst)

	faultBodyElem := serde.BuildFaultBodyChild(s.envelopeSpace, err)
	envelopeElem := serde.BuildEnvelope(s.envelopeSpace, s.envelopeNS, faultBodyElem)

	doc.AddChild(envelopeElem)

	doc.Indent(2)
	w.WriteHeader(http.StatusInternalServerError)
	doc.WriteTo(w)
}
