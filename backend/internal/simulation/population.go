package simulation

import (
 crand "crypto/rand"
 "fmt"
 mrand "math/rand"
 "sync"
)

type PopulationProfile struct { Names []string; JobWeights map[string]int; PersonalityBias Personality }
type PopulationGenerator struct { mu sync.Mutex; next int; rng *mrand.Rand; profile PopulationProfile }
func NewPopulationGenerator(seed int64, profile PopulationProfile)*PopulationGenerator{return &PopulationGenerator{rng:mrand.New(mrand.NewSource(seed)),profile:profile}}
func(g *PopulationGenerator) Generate(count int)[]SimCharacter{g.mu.Lock();defer g.mu.Unlock();out:=make([]SimCharacter,0,count);for i:=0;i<count;i++{g.next++;name:="Citizen";if len(g.profile.Names)>0{name=g.profile.Names[g.rng.Intn(len(g.profile.Names))]};out=append(out,SimCharacter{ID:newUUID(),Name:name+" "+formatID(g.next),Controller:ControllerSimulated,Personality:g.profile.PersonalityBias,Goals:[]Goal{{ID:"work",Kind:"WORK",Priority:50}},Active:false})};return out}
func newUUID()string{b:=make([]byte,16);if _,err:=crand.Read(b);err!=nil{return fmt.Sprintf("00000000-0000-4000-8000-%012d",0)};b[6]=(b[6]&15)|64;b[8]=(b[8]&63)|128;return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",b[0:4],b[4:6],b[6:8],b[8:10],b[10:16])}
func formatID(n int)string{digits:="0123456789";s:="";if n==0{return "0"};for n>0{s=string(digits[n%10])+s;n/=10};return s}
