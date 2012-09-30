package mgou

import (
	"flag"
	"fmt"
	"testing"

	. "github.com/araddon/gou"
	"labix.org/v2/mgo/bson"
)

// go test -host "192.168.1.25"
var (
	host   *string = flag.String("host", "localhost", "Host")
	mgo_db string  = "testdb"
)

func init() {
	flag.Parse()
	SetMongoInfo(*host)
	ModelsDelete(mgo_db, nil, NewWork())
}

// Work is a set of work to be done, a generic container that describes
// work with a unique name, and a service to complete the work
type Work struct {
	//ObjectId     string `json:"oid" bson:"oid,omitempty"`
	//SegmentId  string     `json:"sid" bson:"sid,omitempty"` // indexed
	Mid     bson.ObjectId `bson:"_id,omitempty" json:"-"`
	Mhex    string        `bson:"mhex,omitempty" json:"mhex"`
	Name    string        `json:"name"`
	Service string        `json:"service"`
	Config  JsonHelper    `json:"config"`
}

func NewWork() *Work {
	w := Work{}
	return &w
}
func (w *Work) Init() {

}
func (m *Work) OnLoad() {
	if len(m.Mid) > 0 {
		m.Mhex = bson.ObjectId(m.Mid).Hex()
	}
}
func (m *Work) OidGet() string {
	return m.Name
}
func (m *Work) OidSet(oid string) {
	//m.ObjectId = oid
}
func (m *Work) MidSet(mid bson.ObjectId) {
	m.Mid = mid
}
func (m *Work) MidGet() bson.ObjectId {
	return m.Mid
}
func (w *Work) Type() string {
	return "work"
}

func (w *Work) String() string {
	return fmt.Sprintf("<Work type:work svc:%s Name:%s Oid:%s Mid='%s' config='%v'>",
		w.Service, w.Name, w.OidGet(), w.MidGet(), w.Config)
}

func TestMgoStore(t *testing.T) {

	w := NewWork()
	w.Name = "test1"
	SaveModel(mgo_db, w, nil)
	wl := make([]*Work, 0)
	ModelsLoad(mgo_db, &wl, nil, NewWork())
	Assert(len(wl) == 1, t, "should only have 1")
	SaveModel(mgo_db, w, nil)
	wl2 := make([]*Work, 0)
	ModelsLoad(mgo_db, &wl2, nil, NewWork())
	Assert(len(wl2) == 1, t, "should only have 1")
}
