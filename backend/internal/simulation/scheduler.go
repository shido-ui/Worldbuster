package simulation

import "time"

type TierScheduler struct {
 Policy LifecyclePolicy
}

func (s TierScheduler) ShouldTick(index,total int, now time.Time) bool {
 tier:=ClassifyPopulation(index,total,s.Policy)
 switch tier {
 case TierActive:
  return true
 case TierRecent:
  return now.Unix()%2==0
 case TierBackground:
  return now.Unix()%10==0
 default:
  return false
 }
}
