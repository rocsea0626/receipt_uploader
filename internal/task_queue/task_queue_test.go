package task_queue

import (
	"testing"
	"time"

	"receipt_uploader/internal/constants"
	images_mock "receipt_uploader/internal/images/mock"

	"github.com/stretchr/testify/assert"
)

func TestEnqueue(t *testing.T) {
	queueCapacity := 3
	mockImagesService := &images_mock.ServiceMock{}
	queue := NewService(queueCapacity, mockImagesService)

	task := Task{
		Name: "mock_task_name",
		Func: func() error {
			return nil
		},
	}

	success := queue.Enqueue(task)
	assert.True(t, success)

	success1 := queue.Enqueue(task)
	assert.True(t, success1)

	queue.Enqueue(task)
	success = queue.Enqueue(task)
	assert.False(t, success)
}

func TestWithTimeout(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		task := Task{
			Name: "mock_task_name",
			Func: func() error {
				return nil
			},
		}

		timeout := constants.RESIZE_TIMEOUT
		err := withTimeout(task, timeout)
		assert.Nil(t, err)
	})

	t.Run("should fail, time out", func(t *testing.T) {
		task := Task{
			Name: "mock_task_name",
			Func: func() error {
				time.Sleep(constants.RESIZE_TIMEOUT)
				return nil
			},
		}
		timeout := constants.RESIZE_TIMEOUT - 1*time.Second
		err := withTimeout(task, timeout)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "resizeImages() timed out")
	})

}
