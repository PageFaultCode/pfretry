package retry

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

const (
	simpleRetryCount        = 5
	simpleRetryName         = "test"
	simpleMinimumTimePeriod = 2
)

type RetryTestSuite struct {
	suite.Suite
}

func (suite *RetryTestSuite) TestSimpleRetry() {
	retry := AfmRetry{retryCount: simpleRetryCount, retryName: simpleRetryName}

	count := 0
	_ = retry.Retry(func() error {
		count++
		return fmt.Errorf("%d", count)
	}, nil, nil)

	suite.Assert().Equal(simpleRetryCount, count)
}

func (suite *RetryTestSuite) TestOnErrorCallback() {
	retry := AfmRetry{retryCount: simpleRetryCount, retryName: simpleRetryName}

	targetName := ""
	_ = retry.Retry(func() error {
		return fmt.Errorf("fail")
	}, nil, func(name string, err error) {
		targetName = name
	})

	suite.Assert().Equal(simpleRetryName, targetName)
}

func (suite *RetryTestSuite) TestOnSuccessCallback() {
	retry := AfmRetry{retryCount: simpleRetryCount, retryName: simpleRetryName}

	targetName := ""
	_ = retry.Retry(func() error {
		return nil
	}, func(name string) {
		targetName = name
	}, nil)

	suite.Assert().Equal(simpleRetryName, targetName)
}

func (suite *RetryTestSuite) TestStopFlag() {
	stopFlag := make(chan struct{})
	retry := AfmRetry{retryName: simpleRetryName, stopFlag: stopFlag}

	count := 0
	go func() {
		if count > simpleRetryCount {
			stopFlag <- struct{}{}
		}
	}()

	_ = retry.Retry(func() error {
		count++
		return fmt.Errorf("%d", count)
	}, nil, nil)

	suite.Assert().NotZero(count)
}

func (suite *RetryTestSuite) TestMinimumTime() {
	retry := AfmRetry{retryCount: simpleRetryCount, retryName: simpleRetryName, minimumRetryTime: time.Second * simpleMinimumTimePeriod}

	minimumTime := time.Now()
	minimumTime = minimumTime.Add(time.Second * simpleMinimumTimePeriod)
	count := 0
	_ = retry.Retry(func() error {
		count++
		return fmt.Errorf("%d", count)
	}, nil, nil)

	endTime := time.Now()

	suite.Assert().Equal(simpleRetryCount, count)
	suite.Assert().True(endTime.After(minimumTime))
}

func TestRetryTestSuite(t *testing.T) {
	suite.Run(t, new(RetryTestSuite))
}
