package gotest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeEqual(t *testing.T) {
	t.Log(time.RFC3339)
	assert := assert.New(t)
	t1, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	assert.Nil(err)
	t2, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	assert.Nil(err)
	assert.Equal(t1, t2, "t1 == t2")
}
