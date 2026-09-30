package progression

import "sync"

type PlayerProgressionStore struct {
 mu sync.RWMutex
 data map[string]Progression
}

func NewPlayerProgressionStore()*PlayerProgressionStore{
 return &PlayerProgressionStore{data:map[string]Progression{}}
}

func(s *PlayerProgressionStore) Get(playerID string) Progression{
 s.mu.RLock();defer s.mu.RUnlock()
 if p,ok:=s.data[playerID];ok{return p}
 return NewProgression()
}

func(s *PlayerProgressionStore) Set(playerID string,p Progression) bool{
 if playerID=="" || !p.Stats.Valid() || p.Level<1{return false}
 s.mu.Lock();defer s.mu.Unlock();s.data[playerID]=p;return true
}

func(s *PlayerProgressionStore) AddXP(playerID string,amount int64)(Progression,int){
 if playerID==""||amount<=0{return Progression{},0}
 s.mu.Lock();defer s.mu.Unlock()
 p,ok:=s.data[playerID];if !ok{p=NewProgression()}
 gained:=p.AddXP(amount);s.data[playerID]=p
 return p,gained
}
