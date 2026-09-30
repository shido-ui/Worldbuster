package simulation

type ActivityTier string
const(
 TierActive ActivityTier="ACTIVE"
 TierRecent ActivityTier="RECENT"
 TierBackground ActivityTier="BACKGROUND"
 TierDormant ActivityTier="DORMANT"
)
type PopulationRegion struct{Name string; Weight int}
type LifecyclePolicy struct{ActivePercent,RecentPercent,BackgroundPercent int; Regions []PopulationRegion}
func ClassifyPopulation(index,total int,p LifecyclePolicy)ActivityTier{
 if total<=0{return TierDormant}; pct:=index*100/total
 if pct<p.ActivePercent{return TierActive}; if pct<p.ActivePercent+p.RecentPercent{return TierRecent}
 if pct<p.ActivePercent+p.RecentPercent+p.BackgroundPercent{return TierBackground}; return TierDormant
}
func SelectRegion(seed uint64, regions []PopulationRegion)string{
 if len(regions)==0{return ""}; var total int
 for _,r:=range regions{if r.Weight>0{total+=r.Weight}}
 if total==0{return regions[0].Name}; n:=int(seed%uint64(total))
 for _,r:=range regions{if r.Weight<=0{continue};if n<r.Weight{return r.Name};n-=r.Weight}
 return regions[len(regions)-1].Name
}
