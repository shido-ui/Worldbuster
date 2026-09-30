package world

import ("errors";"sync";"time")

var ErrInvalidLocation=errors.New("invalid location")
var ErrNoRoute=errors.New("no route")
var ErrAlreadyTraveling=errors.New("already traveling")

type TravelState struct{CharacterID string `json:"characterId"`;From string `json:"from"`;To string `json:"to"`;StartedAt time.Time `json:"startedAt"`;ArrivesAt time.Time `json:"arrivesAt"`}

type TravelService struct{mu sync.RWMutex; active map[string]TravelState}

func NewTravelService()*TravelService{return &TravelService{active:make(map[string]TravelState)}}

func(s *TravelService)Start(characterID,from,to string,now time.Time)(TravelState,error){
 if from==to||!validLocation(from)||!validLocation(to){return TravelState{},ErrInvalidLocation}
 seconds:=routeSeconds(from,to);if seconds<=0{return TravelState{},ErrNoRoute}
 s.mu.Lock();defer s.mu.Unlock()
 if _,ok:=s.active[characterID];ok{return TravelState{},ErrAlreadyTraveling}
 t:=TravelState{CharacterID:characterID,From:from,To:to,StartedAt:now,ArrivesAt:now.Add(time.Duration(seconds)*time.Second)}
 s.active[characterID]=t;return t,nil
}

func(s *TravelService)Get(characterID string,now time.Time)(TravelState,bool){
 s.mu.RLock();defer s.mu.RUnlock()
 t,ok:=s.active[characterID];if !ok||!now.Before(t.ArrivesAt){return t,ok}
 return t,true
}

func(s *TravelService)Complete(characterID string,now time.Time)(TravelState,bool){
 s.mu.Lock();defer s.mu.Unlock()
 t,ok:=s.active[characterID];if !ok||now.Before(t.ArrivesAt){return t,ok}
 delete(s.active,characterID);return t,true
}
func validLocation(id string)bool{for _,l:=range locations{if l.ID==id{return true}};return false}
func routeSeconds(from,to string)int{for _,r:=range routes{if r.From==from&&r.To==to{return r.TravelSeconds}};return 0}
