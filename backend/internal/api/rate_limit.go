package api

import("net";"net/http";"strings";"sync";"time")

type RateLimiter struct{mu sync.Mutex;entries map[string]*rateEntry;Limit int;Window time.Duration}
type rateEntry struct{Started time.Time;Count int}

const maxRateLimitEntries=10000

func NewRateLimiter(limit int,window time.Duration)*RateLimiter{
 if limit<1{limit=1}
 if window<=0{window=time.Minute}
 return &RateLimiter{entries:map[string]*rateEntry{},Limit:limit,Window:window}
}

func(r *RateLimiter)Allow(key string)bool{
 r.mu.Lock()
 defer r.mu.Unlock()
 now:=time.Now()
 if len(r.entries)>=maxRateLimitEntries{r.pruneExpired(now)}
 e:=r.entries[key]
 if e==nil||now.Sub(e.Started)>=r.Window{
  r.entries[key]=&rateEntry{Started:now,Count:1}
  return true
 }
 if e.Count>=r.Limit{return false}
 e.Count++
 return true
}

func(r *RateLimiter)pruneExpired(now time.Time){
 for key,e:=range r.entries{
  if now.Sub(e.Started)>=r.Window{delete(r.entries,key)}
 }
}

func RateLimitMiddleware(l *RateLimiter,next http.Handler)http.Handler{
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  host:=strings.TrimSpace(r.RemoteAddr)
  if h,_,err:=net.SplitHostPort(host);err==nil{host=h}
  key:=host+"|"+r.URL.Path
  if !l.Allow(key){w.Header().Set("Retry-After","1");http.Error(w,"rate limit exceeded",http.StatusTooManyRequests);return}
  next.ServeHTTP(w,r)
 })
}
