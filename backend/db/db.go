package db

import (
	"log"

	"github.com/twinone/iot/backend/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DBName                = "iot-backend"
	CookieCollection      = "cookies"
	UsersCollection       = "users"
	AccesTokensCollection = "accesstokens"
	FunctionsCollection   = "functions"
)

var defaultSession *mgo.Session

func GetCookieCollection() *mgo.Collection {
	return defaultSession.Copy().DB(DBName).C(CookieCollection)
}

// Gets the requested user or nil they it doesn't exist
func FindUserByEmail(email string) *model.User {
	s := defaultSession.Copy()
	defer s.Close()

	u := &model.User{}
	c := s.DB(DBName).C(UsersCollection)
	if err := c.Find(bson.M{"email": email}).One(u); err != nil {
		return nil
	}
	return u
}

func FindUserByAccessToken(token string) *model.User {
	s := defaultSession.Copy()
	defer s.Close()

	t := &model.AccessToken{}
	c := s.DB(DBName).C(AccesTokensCollection)
	if err := c.Find(bson.M{"token": token}).One(t); err != nil {
		return nil
	}
	return FindUserByEmail(t.Email)
}

func InsertAccessToken(t *model.AccessToken) {
	s := defaultSession.Copy()
	defer s.Close()

	c := s.DB(DBName).C(AccesTokensCollection)
	c.Insert(t)
}

func FindFunctionsByEmail(email string) []*model.Function {
	s := defaultSession.Copy()
	defer s.Close()

	var f []*model.Function
	c := s.DB(DBName).C(FunctionsCollection)
	if err := c.Find(bson.M{"owner": email}).All(&f); err != nil {
		log.Println(err)
		return nil
	}
	log.Println("Found functions:", f)
	return f
}

func FindFunctionById(id string) *model.Function {
	s := defaultSession.Copy()
	defer s.Close()

	f := &model.Function{}
	c := s.DB(DBName).C(FunctionsCollection)
	if err := c.Find(bson.M{"id": id}).One(f); err != nil {
		log.Println(err)
		return nil
	}
	log.Println("Found function:", f)
	return f
}

func RemoveFunction(id string, email string) {
	s := defaultSession.Copy()
	defer s.Close()
	defer func() {
		if err := recover(); err != nil {
			log.Println("Error removing function:", err)
		}
	}()
	c := s.DB(DBName).C(FunctionsCollection)
	if err := c.Remove(bson.M{"_id": bson.ObjectIdHex(id), "owner": email}); err != nil {
		log.Println("Error removing function:", err)
	}

}

func InsertFunction(f *model.Function) string {
	s := defaultSession.Copy()
	defer s.Close()

	c := s.DB(DBName).C(FunctionsCollection)
	i := bson.NewObjectId()
	f.Id = i
	c.Insert(f)
	return i.Hex()
}

func RemoveAccessToken(token string) {
	s := defaultSession.Copy()
	defer s.Close()

	c := s.DB(DBName).C(AccesTokensCollection)
	c.Remove(bson.M{"token": token})
}

func InsertUser(u *model.User) {
	s := defaultSession.Copy()
	defer s.Close()

	c := s.DB(DBName).C(UsersCollection)
	c.Insert(u)
}

func Init() *mgo.Session {
	var err error
	defaultSession, err = mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}
	return defaultSession
}
