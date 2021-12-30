package soap

import (
	"net/http"

	"github.com/beevik/etree"
	"github.com/francistm/soap-server/internal"
	"github.com/francistm/soap-server/internal/model"
	"github.com/francistm/soap-server/internal/serde"
)

func (s *Service) printDefinition(w http.ResponseWriter, r *http.Request) {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", internal.XmlProcInst)

	actions := make(model.Actions)
	for portName, port := range s.ports {
		for actionName, action := range port.actions {
			actions[portName][actionName] = &model.Action{
				InType:  action.in,
				OutType: action.out,
			}
		}
	}

	doc.AddChild(serde.BuildDefinitions(s.name, actions, serde.WithNamespace(s.domain)))
	doc.Indent(2)

	w.WriteHeader(http.StatusOK)
	doc.WriteTo(w)
}
