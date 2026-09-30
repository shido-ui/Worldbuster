package simulation

import (
 "sort"
 "sync"
 "time"
)

type Tier string
const(
 TierActive Tier="ACTIVE"
 TierRecent Tier="RECENT"
 TierBackground Tier="BACKGROUND"
 TierDormant Tier="DORMANT"
)

type LifecyclePolicy struct{ActivePercent int;RecentPercent int;BackgroundPercent int}
type TierScheduler struct{Policy LifecyclePolicy}
func(s TierScheduler)Classify(chars []SimCharacter,now time.Time)map[string]Tier{
 out:=map[string]Tier{};sorted:=append([]SimCharacter(nil),chars...)
 sort.Slice(sorted,func(i,j int)bool{return sorted[i].LastTick>sorted[j].LastTick})
 n:=len(sorted);if n==0{return out}
 active:=n*s.Policy.ActivePercent/100;recent:=n*(s.Policy.ActivePercent+s.Policy.RecentPercent)/100;background:=n*(s.Policy.ActivePercent+s.Policy.RecentPercent+s.Policy.BackgroundPercent)/100
 for i,c:=range sorted{t:=TierDormant;if i<active{t=TierActive}else if i<recent{t=TierRecent}else if i<background{t=TierBackground};out[c.ID]=t}
 _=now;return out
}

type PopulationScheduler struct{Population *Service;Policy TierScheduler;mu sync.Mutex;LastRun time.Time}
func(s *PopulationScheduler)Run(now time.Time)(map[Tier]int,error){
 s.mu.Lock();defer s.mu.Unlock();chars:=s.Population.List(1000);tiers:=s.Policy.Classify(chars,now);counts:=map[Tier]int{}
 for _,c:=range chars{t:=tiers[c.ID];counts[t]++}
 s.LastRun=now;return counts,nil
}
