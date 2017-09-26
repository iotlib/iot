package db

import (
	"github.com/twinone/iot/backend/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DBName                = "iot-backend"
	CookieCollection      = "cookies"
	UsersCollection       = "users"
	AccesTokensCollection = "accesstokens"
)

var defaultSession *mgo.Session

func GetCookieCollection() *mgo.Collection {
	return defaultSession.Copy().DB(DBName).C(CookieCollection)
}

// Gets the requested user or nil they it doesn't exist
func GetUserByEmail(email string) *model.User {
	s := defaultSession.Copy()
	defer s.Close()

	u := &model.User{}
	c := s.DB(DBName).C(UsersCollection)
	if err := c.Find(bson.M{"email": email}).One(u); err != nil {
		return nil
	}
	return u
}

func GetUserByAccessToken(token string) *model.User {
	s := defaultSession.Copy()
	defer s.Close()

	t := &model.AccessToken{}
	c := s.DB(DBName).C(AccesTokensCollection)
	if err := c.Find(bson.M{"token": token}).One(t); err != nil {
		return nil
	}
	return GetUserByEmail(t.Email)
}

func InsertAccessToken(t *model.AccessToken) {
	s := defaultSession.Copy()
	defer s.Close()

	c := s.DB(DBName).C(AccesTokensCollection)
	c.Insert(t)
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
