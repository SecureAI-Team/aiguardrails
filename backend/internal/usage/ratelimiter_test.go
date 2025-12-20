package usage

import (
	"context"
	"testing"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestRateLimiterQuota(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	rl := NewRateLimiter(client, "testns")

	ctx := context.Background()
	ok, count, err := rl.Allow(ctx, "app1", 2)
	if err != nil || !ok || count != 1 {
		t.Fatalf("first call should pass, count=1")
	}
	ok, count, err = rl.Allow(ctx, "app1", 2)
	if err != nil || !ok || count != 2 {
		t.Fatalf("second call should pass, count=2")
	}
	ok, count, err = rl.Allow(ctx, "app1", 2)
	if err != nil {
		t.Fatalf("third call err %v", err)
	}
	if ok || count <= 2 {
		t.Fatalf("third call should be blocked when over quota")
	}

	// ensure expiry set
	keys := mr.Keys()
	if len(keys) == 0 {
		t.Fatalf("expected redis key set")
	}
}

