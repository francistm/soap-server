package model

type (
	Actions map[string]map[string]*Action
)

type Action struct {
	InType  interface{}
	OutType interface{}
}
