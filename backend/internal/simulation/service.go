package simulation
import ("errors";"sync";"time")
var ErrInvalidCharacter=errors.New("invalid simulated character")
type Service struct{mu sync.RWMutex;characters map[string]*SimCharacter}
func NewService()*Service{return &Service{characters:map[string]*SimCharacter{}}}
func(s *Service)Register(c SimCharacter)error{if c.ID==""||c.Name==""||c.Controller!=ControllerSimulated{return ErrInvalidCharacter};s.mu.Lock();defer s.mu.Unlock();if _,ok:=s.characters[c.ID];ok{return ErrInvalidCharacter};c.Active=true;s.characters[c.ID]=&c;return nil}
func(s *Service)Get(id string)(SimCharacter,bool){s.mu.RLock();defer s.mu.RUnlock();c,ok:=s.characters[id];if !ok{return SimCharacter{},false};out:=*c;out.Goals=append([]Goal(nil),c.Goals...);return out,true}
func(s *Service)List(limit int)[]SimCharacter{s.mu.RLock();defer s.mu.RUnlock();if limit<=0||limit>10000{limit=100};out:=make([]SimCharacter,0,len(s.characters));for _,c:=range s.characters{out=append(out,*c);if len(out)>=limit{break}};return out}
func(s *Service)TickTier(now time.Time,tier Tier)int{s.mu.Lock();defer s.mu.Unlock();count:=0;for _,c:=range s.characters{if !c.Active{continue};age:=now.Unix()-c.LastTick;if tier==TierActive||age<60{c.LastTick=now.Unix();count++}};return count}
func(s *Service)Tick(now time.Time)int{return s.TickTier(now,TierActive)}

func(s *Service)Touch(id string,now time.Time){s.mu.Lock();defer s.mu.Unlock();if ch,ok:=s.characters[id];ok{ch.LastTick=now.Unix()}}
