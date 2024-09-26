package producer

import (
	"math/rand"
	"testing"
	"time"

	"github.com/mchatzis/go/producer/pkg/base"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTask(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := 42

	assert.Panics(t, func() {
		generateTask(r, 0)
	}, "Expected panic when ID is 0")

	assert.Panics(t, func() {
		generateTask(r, -1)
	}, "Expected panic when ID is negative")

	task := generateTask(r, id)

	assert.Equal(t, int32(id), task.ID, "Task ID should match the input ID")

	assert.GreaterOrEqual(t, task.Type, int32(0), "Task Type should be >= 0")
	assert.Less(t, task.Type, int32(10), "Task Type should be < 10")

	assert.GreaterOrEqual(t, task.Value, int32(0), "Task Value should be >= 0")
	assert.Less(t, task.Value, int32(100), "Task Value should be < 100")

	assert.Equal(t, base.TaskStatePending, task.State, "Task State should be pending")

	assert.NotZero(t, task.Creationtime, "Task Creationtime should not be zero")
	assert.NotZero(t, task.Lastupdatetime, "Task Lastupdatetime should not be zero")

	assert.LessOrEqual(t, task.Creationtime, float64(time.Now().Unix()), "Task creation time should not be in the future")
	assert.LessOrEqual(t, task.Lastupdatetime, float64(time.Now().Unix()), "Task creation time should not be in the future")
}
