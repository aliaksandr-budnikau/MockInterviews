package interview1

import (
	"context"
	"log"
)

type DistDataRepositoryDecorator struct {
	dataRepository DataRepository
}

func NewDistDataRepositoryDecorator(dataRepository DataRepository) *DistDataRepositoryDecorator {
	return &DistDataRepositoryDecorator{dataRepository: dataRepository}
}

func (it *DistDataRepositoryDecorator) Get(ctx context.Context, addresses []string, key string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	resultChannel := make(chan string)
	defer close(resultChannel)

	for i, address := range addresses {
		go func() {
			log.Printf("Started gorouting #%d", i)
			defer log.Printf("Finished gorouting #%d", i)

			value, err := it.dataRepository.Get(ctx, address, key)
			if err != nil {
				log.Printf(err.Error()+" gorouting #%d", i)
				return
			}
			defer func() {
				err := recover()
				if err != nil {
					log.Printf("Failed to send to closed channel #%d", i)
				}
			}()

			resultChannel <- value
		}()
	}

	log.Printf("Waiting result in main...")
	select {
	case result := <-resultChannel:
		log.Printf("Result was gotten: %s", result)
		return result, nil
	case <-ctx.Done():
		log.Printf("Operation interrupted")
		return "", ctx.Err()
	}
}
