package soap

type serviceOption func(s *Service)

func WithEnvelopeSpace(s string) serviceOption {
	return func(serv *Service) {
		serv.envelopeSpace = s
	}
}

func WithEnvelopeNS(ns map[string]string) serviceOption {
	return func(serv *Service) {
		serv.envelopeNS = ns
	}
}
