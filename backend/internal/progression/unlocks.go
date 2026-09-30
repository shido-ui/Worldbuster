package progression

import "sync"

type UnlockStore struct {
 mu sync.RWMutex
 owned map[string]map[string]ProgressionReward
}

func NewUnlockStore()*UnlockStore{return &UnlockStore{owned:map[string]map[string]ProgressionReward{}}}

func(s *UnlockStore) Grant(playerID string,reward ProgressionReward)bool{
 if s==nil||playerID==""||reward.ID==""{return false}
 s.mu.Lock();defer s.mu.Unlock()
 if s.owned[playerID]==nil{s.owned[playerID]=map[string]ProgressionReward{}}
 if _,exists:=s.owned[playerID][reward.ID];exists{return false}
 s.owned[playerID][reward.ID]=reward
 return true
}

func(s *UnlockStore) Has(playerID,rewardID string)bool{
 s.mu.RLock();defer s.mu.RUnlock()
 _,ok:=s.owned[playerID][rewardID];return ok
}

func(s *UnlockStore) List(playerID string)[]ProgressionReward{
 s.mu.RLock();defer s.mu.RUnlock()
 out:=make([]ProgressionReward,0)
 for _,r:=range s.owned[playerID]{out=append(out,r)}
 return out
}

func GrantEligible(store *UnlockStore,playerID string,state SkillState,rewards []ProgressionReward)int{
 if store==nil{return 0}
 count:=0
 for _,r:=range EligibleRewards(state,rewards){if store.Grant(playerID,r){count++}}
 return count
}
