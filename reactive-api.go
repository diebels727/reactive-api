package main

import (
  "net/http"
  "github.com/gorilla/mux"
  "fmt"
  "time"
  "strconv"
  "encoding/json"
  "flag"
  "gopkg.in/mgo.v2"
  // "gopkg.in/mgo.v2/bson"
)

var port string
var mongo string
func init() {
  flag.StringVar(&port,"port","8080","HTTP Server port")
  flag.StringVar(&mongo,"mongo","localhost","Mongo address")
}

type Event struct {
  Timestamp int64
  Message string
  Bot string
  Source string
  Command string
}

type Datastore struct {
  Session *mgo.Session
  Collection *mgo.Collection
}

func slug(str string) string {
  str = strings.ToLower(str)
  return strings.Replace(str,".","-",-1)
}

func NewDatastore(host string,server string,session *mgo.Session) (*Datastore) {
  local := session.Copy()
  collection := local.DB(slug(server)).C("events")
  datastore := Datastore{local,collection}
  return &datastore
}

func (d *Datastore) Events() {

}

func ServerHandler(response http.ResponseWriter,request *http.Request) {
  response.Header().Set("Content-Type", "application/json")
  params := mux.Vars(request)



  // fmt.Fprint(response,events)
}

var session *mgo.Session

func main() {
  session,err := mgo.Dial(mongo)
  if err != nil {
    panic(err)
  }
  defer session.Close()

  router := mux.NewRouter()
  router.HandleFunc("/{server}",ServerHandler)
  router.HandleFunc("/{server}/minutes/{minutes}",ServerHandler)
  http.Handle("/",router)
  http.ListenAndServe(port,nil)
  fmt.Println("done!")
}