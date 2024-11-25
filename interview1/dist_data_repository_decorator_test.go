package interview1

import (
	"context"
	"errors"
	"mock_interview/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const ip1 = "1.1.1.1"
const ip2 = "1.1.1.2"
const ip3 = "1.1.1.3"
const dummyKey = "dummyKey"
const expectedValue = "ValueWeGot"
const workerSleepTime = time.Duration(5) * time.Second
const expectedWorkerTimeout = time.Duration(500) * time.Millisecond
const workerTimeout = time.Duration(300) * time.Millisecond

func TestGetFromServersNormalCase(t *testing.T) {
	ctx := context.Background()

	mockedRepository := new(mocks.DataRepository)
	mockedRepository.On("Get", mock.Anything, ip1, dummyKey).Return(expectedValue, nil).Maybe()
	mockedRepository.On("Get", mock.Anything, ip2, dummyKey).Return(expectedValue, nil).Maybe()
	mockedRepository.On("Get", mock.Anything, ip3, dummyKey).Return(expectedValue, nil).Maybe()

	decorator := NewDistDataRepositoryDecorator(mockedRepository)

	result, error := decorator.Get(ctx, []string{ip1, ip2, ip3}, dummyKey)
	assert.NoError(t, error)
	assert.Equal(t, expectedValue, result)

	mockedRepository.AssertExpectations(t)
}

func TestGetFromServersNormalCaseWithNoMocks(t *testing.T) {
	ctx := context.Background()

	decorator := NewDistDataRepositoryDecorator(&RealDataRepository{})

	result, error := decorator.Get(ctx, []string{ip1, ip2, ip3}, dummyKey)
	assert.NoError(t, error)
	assert.Equal(t, dummyKey, result)
}

func TestGetFromServersCaseWithAFewErrors(t *testing.T) {
	ctx := context.Background()

	mockedRepository := new(mocks.DataRepository)
	mockedRepository.On("Get", mock.Anything, ip1, dummyKey).Return("", errors.New("Server error")).Maybe()
	mockedRepository.On("Get", mock.Anything, ip2, dummyKey).Return("", errors.New("Server error")).Maybe()
	mockedRepository.On("Get", mock.Anything, ip3, dummyKey).Return(expectedValue, nil)

	decorator := NewDistDataRepositoryDecorator(mockedRepository)

	result, error := decorator.Get(ctx, []string{ip1, ip2, ip3}, dummyKey)
	assert.NoError(t, error)
	assert.Equal(t, expectedValue, result)

	mockedRepository.AssertExpectations(t)
}

func TestGetFromServersCaseWithDelays(t *testing.T) {
	ctx := context.Background()

	mockedRepository := new(mocks.DataRepository)
	mockedRepository.On("Get", mock.Anything, ip1, dummyKey).Return("", nil).Run(func(args mock.Arguments) { time.Sleep(workerSleepTime) })
	mockedRepository.On("Get", mock.Anything, ip2, dummyKey).Return("", nil).Run(func(args mock.Arguments) { time.Sleep(workerSleepTime) })
	mockedRepository.On("Get", mock.Anything, ip3, dummyKey).Return(expectedValue, nil)

	decorator := NewDistDataRepositoryDecorator(mockedRepository)

	startTime := time.Now()
	result, error := decorator.Get(ctx, []string{ip1, ip2, ip3}, dummyKey)
	elapsedTime := time.Since(startTime)

	assert.NoError(t, error)
	assert.Equal(t, expectedValue, result)
	assert.True(t, elapsedTime < expectedWorkerTimeout, "Test took %v", elapsedTime)

	mockedRepository.AssertExpectations(t)
}

func TestGetFromServersCaseWithTimeout(t *testing.T) {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), workerTimeout)
	defer cancel()

	mockedRepository := new(mocks.DataRepository)
	mockedRepository.On("Get", mock.Anything, ip1, dummyKey).Return(expectedValue, nil).Run(func(args mock.Arguments) { time.Sleep(workerSleepTime) })
	mockedRepository.On("Get", mock.Anything, ip2, dummyKey).Return(expectedValue, nil).Run(func(args mock.Arguments) { time.Sleep(workerSleepTime) })
	mockedRepository.On("Get", mock.Anything, ip3, dummyKey).Return(expectedValue, nil).Run(func(args mock.Arguments) { time.Sleep(workerSleepTime) })

	decorator := NewDistDataRepositoryDecorator(mockedRepository)

	startTime := time.Now()
	result, error := decorator.Get(ctxWithTimeout, []string{ip1, ip2, ip3}, dummyKey)
	elapsedTime := time.Since(startTime)

	assert.Empty(t, result)
	assert.Error(t, error)
	assert.True(t, elapsedTime < expectedWorkerTimeout, "Test took %v", elapsedTime)

	mockedRepository.AssertExpectations(t)
}
