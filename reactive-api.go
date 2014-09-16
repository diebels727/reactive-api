package main

import (
  "net/http"
  "github.com/gorilla/mux"
  "fmt"
  "strings"
  "flag"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "github.com/diebels727/spyglass"
)

var port string
var mongo string
func init() {
  flag.StringVar(&port,"port","8080","HTTP Server port")
  flag.StringVar(&mongo,"mongo","localhost","Mongo address")
}

type Datastore struct {
  Collection *mgo.Collection
}

func slug(str string) string {
  str = strings.ToLower(str)
  return strings.Replace(str,".","-",-1)
}

func NewDatastore(server string) (*Datastore) {
  local := copySession()
  collection := local.DB(slug(server)).C("events")
  datastore := Datastore{collection}
  return &datastore
}

func (d *Datastore) Events() (m []spyglass.Event){
  m = make([]spyglass.Event,0)
  err := d.Collection.Find(bson.M{}).All(&m)
  if err != nil {
    panic(err)
  }
  return
}

func Handler(response http.ResponseWriter,request *http.Request) {
  response.Header().Set("Content-Type", "application/json")
  params := mux.Vars(request)
  datastore := NewDatastore(params["server"])
  events := datastore.Events()
  fmt.Fprint(response,events)
}

var session *mgo.Session

func copySession() (*mgo.Session) {
  return session.Copy()
}

func initSession(mongo string) (*mgo.Session){
  var err error
  session,err = mgo.Dial(mongo)
  if err != nil {
    panic(err)
  }
  return session
}


func main() {
  s := initSession(mongo)
  defer s.Close()
  router := mux.NewRouter()
  router.HandleFunc("/{server}",Handler)
  router.HandleFunc("/{server}/minutes/{minutes}",Handler)
  http.Handle("/",router)
  http.ListenAndServe(":"+port,nil)
  fmt.Println("done!")
}