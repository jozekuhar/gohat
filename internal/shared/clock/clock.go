package clock

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type Clock interface {
	NowUTC() time.Time
}

type SystemClock struct{}

func NewSystemClock() *SystemClock {
	return &SystemClock{}
}

func (c *SystemClock) NowUTC() time.Time {
	return time.Now().UTC()
}

type MockClock struct {
	mock.Mock
}

func NewMockClock() *MockClock {
	return &MockClock{}
}

func (c *MockClock) NowUTC() time.Time {
	args := c.Called()
	return args.Get(0).(time.Time)
}
