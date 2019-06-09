package object

import "github.com/adambaumeister/panfw-util/panos/api/deviceconfig"

/*
API Generics
*/
type ApiObject interface {
	Add(string, string, []string) deviceconfig.MsgJobResponse
	GetType() string
}
