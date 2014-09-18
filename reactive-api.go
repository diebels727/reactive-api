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
  "time"
  "strconv"
  "encoding/json"
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

func (d *Datastore) Events() (m []spyglass.Event,err error){
  m = make([]spyglass.Event,0)
  err = d.Collection.Find(bson.M{}).All(&m)
  return
}

func (d *Datastore) EventsWithMinutes(minutes string) (m []spyglass.Event,err error){
  m = make([]spyglass.Event,0)
  current_time := time.Now().Unix()
  minutes_int,err := strconv.Atoi(minutes)
  seconds := minutes_int * 60
  since := current_time - int64(seconds)

  err = d.Collection.Find(bson.M{"timestamp": bson.M{"$gte": since}}).All(&m)
  return
}

func Handler(response http.ResponseWriter,request *http.Request) {
  response.Header().Set("Content-Type", "application/json")
  params := mux.Vars(request)
  datastore := NewDatastore(params["server"])
  var events []spyglass.Event
  var err error
  if len(params["minutes"]) > 0 {
    events,err = datastore.EventsWithMinutes(params["minutes"])
  } else {
    events,err = datastore.Events()
  }

  if err != nil {
    http.Error(response,http.StatusText(500),500) //probably should not be a 500 -- this *shoulb* be a client error
  }

  bytes,err := json.Marshal(events)
  if err != nil {
    fmt.Println("Error marshalling events, aborting...")
    return
  }
  jsonEvents := string(bytes)

  fmt.Fprint(response,"{\"events\": "+jsonEvents+"}")


  // fmt.Fprint(response,events)
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
  fmt.Printf("Listening on :%s\n",port)
  http.ListenAndServe(":"+port,nil)
  fmt.Println("done!")
}