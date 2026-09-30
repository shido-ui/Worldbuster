package simulation

import "time"

type Relationship struct {
 CharacterID string
 TargetID string
 Affinity int
 Trust int
 Familiarity int
 Interactions int
 LastInteraction *time.Time
}
func (r Relationship) SocialWeight() int { return r.Affinity/10 + r.Trust/20 + r.Familiarity/40 }
func applySocialInteraction(r Relationship, now time.Time, positive bool) Relationship {
 if positive { r.Affinity=clamp(r.Affinity+4,-1000,1000); r.Trust=clamp(r.Trust+2,-1000,1000) } else { r.Affinity=clamp(r.Affinity-3,-1000,1000); r.Trust=clamp(r.Trust-1,-1000,1000) }
 r.Familiarity=clamp(r.Familiarity+5,0,1000); r.Interactions++; r.LastInteraction=&now; return r
}
func clamp(v,min,max int) int { if v<min{return min}; if v>max{return max}; return v }
