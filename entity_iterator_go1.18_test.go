//go:build go1.18
// +build go1.18

package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testIterator struct {
	mock.Mock
}

func (ti *testIterator) Close() error {
	return ti.Called().Error(0)
}

func (ti *testIterator) Next(entity any) error {
	return ti.Called(entity).Error(0)
}

func TestEntityRepository_Close(t *testing.T) {
	var (
		iterator       = &testIterator{}
		entityIterator = newEntityIterator[User](iterator)
	)

	iterator.On("Close").Return(nil)
	assert.Nil(t, entityIterator.Close())
	iterator.AssertExpectations(t)
}

func TestEntityRepository_Next(t *testing.T) {
	var (
		iterator       = &testIterator{}
		entityIterator = newEntityIterator[User](iterator)
	)

	iterator.On("Next", mock.Anything).Return(nil)
	_, err := entityIterator.Next()
	assert.Nil(t, err)
	iterator.AssertExpectations(t)
}
