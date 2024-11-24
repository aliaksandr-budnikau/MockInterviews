package interview1

import (
	"context"
)

type DataRepository interface {
	Get(ctx context.Context, address, key string) (string, error)
}
