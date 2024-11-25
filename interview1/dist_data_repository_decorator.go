package interview1

import (
	"context"
	"log"
	"sync"
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

	resultChannel := make(chan string, len(addresses))
	var wg sync.WaitGroup

	for i, address := range addresses {
		wg.Add(1)
		go func() {
			defer wg.Done()

			log.Printf("Started gorouting #%d", i)
			defer log.Printf("Finished gorouting #%d", i)

			value, err := it.dataRepository.Get(ctx, address, key)
			if err != nil {
				log.Printf("Error in goroutine #%d: %v", i, err)
				return
			}

			select {
			case resultChannel <- value:
				log.Printf("Successfully sent to channel by goroutine #%d", i)
			case <-ctx.Done():
				log.Printf("Context canceled, goroutine #%d", i)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	log.Printf("Waiting for result in main...")
	select {
	case result, ok := <-resultChannel:
		if ok {
			log.Printf("Result received: %s", result)
			return result, nil
		}
		log.Printf("Channel closed without results")
		return "", nil
	case <-ctx.Done():
		log.Printf("Operation interrupted")
		return "", ctx.Err()
	}
}
