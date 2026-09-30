package api
import("testing";"time")
func TestRateLimiter(t *testing.T){r:=NewRateLimiter(2,time.Minute);if !r.Allow("x")||!r.Allow("x"){t.Fatal("expected first two requests allowed")};if r.Allow("x"){t.Fatal("expected third request blocked")};r2:=NewRateLimiter(1,time.Millisecond);if !r2.Allow("x"){t.Fatal("first request blocked")};time.Sleep(2*time.Millisecond);if !r2.Allow("x"){t.Fatal("window did not reset")}}
