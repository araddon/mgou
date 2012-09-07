package mgou

import (
//"encoding/json"
//"labix.org/v2/mgo/bson"
)

type DataModel interface {
	MidGet() string
	MidSet(string)
	OidGet() string
	OidSet(string)
	Type() string
	OnLoad()
}

/*
type ModelBase struct {
	Mid       string `bson:"_id,omitempty"`
	Oid       string `bson:"oid,omitempty"`
	ModelType string `bson:"-" json:"-"`
	//events    map[string]ServiceEvent
}

func NewModelBase() ModelBase {
	m := ModelBase{}
	//m.events = make(map[string]ServiceEvent)
	return m
}

func (m *ModelBase) OidSet(oid string) {
	m.Oid = oid
}
func (m *ModelBase) OidGet() string {
	return m.Oid
}
func (m *ModelBase) MidGet() string {
	return m.Mid
}
func (m *ModelBase) MidSet(mid string) {
	m.Mid = mid
}
func (m *ModelBase) Type() string {
	return m.ModelType
}

func (m *ModelBase) toMi(mi map[string]interface{}) {
	mi["mid"] = m.Mid
}
*/
