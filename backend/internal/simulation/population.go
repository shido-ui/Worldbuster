package simulation

import ("math/rand";"sync")

type PopulationProfile struct {
 Names []string
 JobWeights map[string]int
 PersonalityBias Personality
}

type PopulationGenerator struct {
 mu sync.Mutex
 next int
 rng *rand.Rand
 profile PopulationProfile
}

func NewPopulationGenerator(seed int64, profile PopulationProfile)*PopulationGenerator{return &PopulationGenerator{rng:rand.New(rand.NewSource(seed)),profile:profile}}

func(g *PopulationGenerator) Generate(count int)[]SimCharacter{
 g.mu.Lock();defer g.mu.Unlock()
 out:=make([]SimCharacter,0,count)
 for i:=0;i<count;i++{
  g.next++
  name:="Citizen"
  if len(g.profile.Names)>0{name=g.profile.Names[g.rng.Intn(len(g.profile.Names))]}
  out=append(out,SimCharacter{ID:formatID(g.next),Name:name+" "+formatID(g.next),Controller:ControllerSimulated,Personality:g.profile.PersonalityBias,Goals:[]Goal{{ID:"work",Kind:"WORK",Priority:50}} ,Active:false})
 }
 return out
}

func formatID(n int)string{digits:="0123456789";s:="";if n==0{return "0"};for n>0{s=string(digits[n%10])+s;n/=10};return s}
