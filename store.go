package mgou

import (
	"errors"
	"launchpad.net/mgo/v2"
	"launchpad.net/mgo/v2/bson"
	"sync"

	. "github.com/araddon/gou"
)

var (
	mgo_db      string
	mgo_conn    string
	mgoMu       sync.Mutex
	mgoSessions = make(map[string]*MgoSession)
)

func SetMongoInfo(conn, db string) {
	mgo_db = db
	mgo_conn = conn
}

type MgoSession struct {
	S  *mgo.Session
	Ct int
}

func MgoConnCheckin(s *mgo.Session) {

}

// Manages creation of Mongo Connections w locking etc
func MgoConnGet(name string) (*mgo.Session, error) {
	var s *MgoSession
	var found bool
	mgoMu.Lock()
	if s, found = mgoSessions[name]; !found {
		s = new(MgoSession)
		session, err := mgo.Dial(mgo_conn)
		if err != nil {
			Log(ERROR, mgo_conn, " ", err)
		} else {
			s.S = session
		}
		mgoSessions[name] = s
	}
	mgoMu.Unlock()
	if s.S != nil {
		return s.S.Copy(), nil
	}
	return nil, errors.New("no session created")
}

// Save the DataModel to DataStore 
func SaveModel(m DataModel, conn *mgo.Session) (err error) {
	if conn == nil {
		conn, err = MgoConnGet(mgo_db)
		if err != nil {
			return
		}
		defer conn.Close()
	}
	if conn != nil {
		bsonMid := bson.ObjectId(m.MidGet())
		c := conn.DB(mgo_db).C(m.Type())
		Debug(mgo_db, " type=", m.Type(), " Mid=", bsonMid)
		if len(bsonMid) < 5 {
			m.MidSet(bson.NewObjectId().String())
			Debug("insert ", m)
			if err = c.Insert(m); err != nil {
				Log(ERROR, "ERROR on insert ", err, " TYPE=", m.Type(), " ", m.MidGet())
			} else {
				Log(WARN, "successfully inserted!!!!!!  ", m.MidGet())
			}
		} else {
			// YOU MUST NOT SEND Mid  "_id" to Mongo
			m.MidSet("") // omitempty means it doesn't get sent
			if err = c.Update(bson.M{"_id": bsonMid}, m); err != nil {
				Log(ERROR, "ERROR on update ", err, " ", bsonMid, " MID=?", m.MidGet())
			}
		}
	} else {
		Log(ERROR, "Nil connection")
		return errors.New("no db connection")
	}
	return
}

func ModelsDelete(qry interface{}, dm DataModel) error {
	if conn, c, ok := GetTableConn(dm); ok {

		info, err := c.RemoveAll(qry)
		Debug(info, " table=", dm.Type())
		if err != nil {
			Log(ERROR, "could not delete items? ", err)
			return err
		}
		conn.Close()
	} else {
		Log(ERROR, "Could not get conn for ", dm.Type())
	}
	return nil
}

// Load Models from Mongo
func ModelsLoad(list interface{}, qry interface{}, dm DataModel) {
	if conn, c, ok := GetTableConn(dm); ok {
		iter := c.Find(qry).Iter()
		err := iter.All(list)
		if err != nil {
			Log(ERROR, err)
		}
		//for iter.Next(&dm) {
		//	dm2 := dm
		//	list = append(list, dm2)
		//}
		conn.Close()
	} else {
		Log(ERROR, "Could not get conn for ", dm.Type())
	}
}

func GetTableConn(dm DataModel) (s *mgo.Session, c *mgo.Collection, ok bool) {
	conn, _ := MgoConnGet(mgo_db)
	if conn != nil {
		c := conn.DB(mgo_db).C(dm.Type())
		return conn, c, true
	}
	return nil, nil, false
}

func GetMgoCC(name string) (s *mgo.Session, c *mgo.Collection, ok bool) {
	conn, _ := MgoConnGet(mgo_db)
	if conn != nil {
		c := conn.DB(mgo_db).C(name)
		return conn, c, true
	}
	return nil, nil, false
}
