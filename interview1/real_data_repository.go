package interview1

import (
	"context"
	"math/rand"
	"time"
)

type RealDataRepository struct{}

func (it *RealDataRepository) Get(ctx context.Context, address, key string) (string, error) {
	sleepDuration := rand.Intn(100) + 100
	randomDuration := time.Duration(sleepDuration) * time.Microsecond
	time.Sleep(randomDuration)
	return key, nil
}
