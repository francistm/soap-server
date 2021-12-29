package soap

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/francistm/soap-server/internal/model"
)

func (s *Service) printDefinition(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ContentType, ContentTypeXml)
	w.Write([]byte(xml.Header))

	defs := &model.Definitions{
		XMLNsS:    model.NsXml,
		XMLNsSOAP: model.NsSoap,
		XMLNsWSDL: model.NsWsdl,
		XMLNsTns:  s.domain,

		TNS: s.domain,

		Types: &model.Types{
			Schema: &model.Schema{
				TNS:      s.domain,
				Elements: s.buildDefElems(),
			},
		},
		Messages:  s.buildDefMessages(),
		PortTypes: s.buildPortTypes(),
		Bindings:  s.buildBindings(),
		Service:   s.buildServices(),
	}

	xml.NewEncoder(w).Encode(defs)
	w.WriteHeader(http.StatusOK)
}

func (s *Service) buildDefElems() []*model.SchemaElement {
	var out []*model.SchemaElement

	for portName, port := range s.ports {
		for actionName, action := range port.actions {
			inElem := &model.SchemaElement{
				Name: fmt.Sprintf("%s%s", portName, actionName),
				Type: &model.SchemaElementComplexType{
					Sequence: &model.SchemaElementComplexTypeSequence{
						Elements: model.StructToSeqElems(action.in),
					},
				},
			}

			out = append(out, inElem)

			if !action.IsOneWay() {
				outElem := &model.SchemaElement{
					Name: fmt.Sprintf("%s%sResponse", portName, actionName),
					Type: &model.SchemaElementComplexType{
						Sequence: &model.SchemaElementComplexTypeSequence{
							Elements: model.StructToSeqElems(action.out),
						},
					},
				}

				out = append(out, outElem)
			}
		}
	}

	return out
}

func (s *Service) buildDefMessages() []*model.Message {
	var messages []*model.Message

	for portName, port := range s.ports {
		for actionName, action := range port.actions {
			in := &model.Message{
				Name: fmt.Sprintf("%s%sSoapIn", portName, actionName),
				Part: &model.MessagePart{
					Name:     "parameters",
					ElemName: fmt.Sprintf("tns:%s%s", portName, actionName),
				},
			}

			messages = append(messages, in)

			if !action.IsOneWay() {
				out := &model.Message{
					Name: fmt.Sprintf("%s%sSoapOut", portName, actionName),
					Part: &model.MessagePart{
						Name:     "parameters",
						ElemName: fmt.Sprintf("tns:%s%sResponse", portName, actionName),
					},
				}

				messages = append(messages, out)
			}
		}
	}

	return messages
}

func (s *Service) buildPortTypes() []*model.PortType {
	var elems []*model.PortType
	for portName, port := range s.ports {
		portType := &model.PortType{
			Name:       fmt.Sprintf("%sSoap", portName),
			Operations: make([]*model.PortTypeOperation, 0, len(port.actions)),
		}

		for actionName, action := range port.actions {
			op := &model.PortTypeOperation{
				Name: actionName,
				Doc: &model.PortTypeOperationDoc{
					NsWsdl: model.NsWsdl,
					Body:   action.documentation,
				},
				Input: &model.PortTypeOperationItem{
					Message: fmt.Sprintf("tns:%s%sSoapIn", portName, actionName),
				},
			}

			if !action.IsOneWay() {
				op.Output = &model.PortTypeOperationItem{
					Message: fmt.Sprintf("tns:%s%sSoapOut", portName, actionName),
				}
			}

			portType.Operations = append(portType.Operations, op)
		}

		elems = append(elems, portType)
	}

	return elems
}

func (s *Service) buildBindings() []*model.Binding {
	var items []*model.Binding

	for portName, port := range s.ports {
		binding := &model.Binding{
			Name: fmt.Sprintf("%sSoap", portName),
			Type: fmt.Sprintf("tns:%sSoap", portName),
			Binding: &model.BindingSoap{
				Transport: model.TransportHttp,
			},
			Operations: make([]*model.BindingWsdlOperation, 0, len(port.actions)),
		}

		for actionName, action := range port.actions {
			op := &model.BindingWsdlOperation{
				Name: actionName,
				Operation: &model.BindingSoapOperation{
					Action: fmt.Sprintf("%s%s", s.domain, actionName),
					Style:  "document",
				},
				Input: &model.BindingWsdlOperationItem{
					Body: &model.BindingWsdlOperationItemBody{
						Use: "literal",
					},
				},
			}

			if !action.IsOneWay() {
				op.Output = &model.BindingWsdlOperationItem{
					Body: &model.BindingWsdlOperationItemBody{
						Use: "literal",
					},
				}
			}

			binding.Operations = append(binding.Operations, op)
		}

		items = append(items, binding)
	}

	return items
}

func (s *Service) buildServices() *model.Service {
	serv := &model.Service{
		Name:  s.name,
		Ports: make([]*model.ServicePort, 0, len(s.ports)),
	}

	for portName := range s.ports {
		p := &model.ServicePort{
			Name:    portName,
			Binding: fmt.Sprintf("tns:%s", portName),
			Location: &model.ServicePortAddr{
				Location: "",
			},
		}

		serv.Ports = append(serv.Ports, p)
	}

	return serv
}
