package api

import("testing";"time")

func TestRateLimiter(t *testing.T){
 r:=NewRateLimiter(2,time.Minute)
 if !r.Allow("x")||!r.Allow("x"){t.Fatal("expected first two requests allowed")}
 if r.Allow("x"){t.Fatal("expected third request blocked")}
 r2:=NewRateLimiter(1,time.Millisecond)
 if !r2.Allow("x"){t.Fatal("first request blocked")}
 time.Sleep(2*time.Millisecond)
 if !r2.Allow("x"){t.Fatal("window did not reset")}
}

func TestRateLimiterPrunesExpiredEntries(t *testing.T){
 r:=NewRateLimiter(1,time.Millisecond)
 now:=time.Now().Add(-time.Second)
 r.entries["expired"]=&rateEntry{Started:now,Count:1}
 r.entries["live"]=&rateEntry{Started:time.Now(),Count:1}
 r.pruneExpired(time.Now())
 if _,ok:=r.entries["expired"];ok{t.Fatal("expired entry was not pruned")}
 if _,ok:=r.entries["live"];!ok{t.Fatal("live entry was pruned")}
}

func TestRateLimiterNormalizesInvalidConfig(t *testing.T){
 r:=NewRateLimiter(0,0)
 if r.Limit!=1||r.Window<=0{t.Fatal("invalid configuration was not normalized")}
}
