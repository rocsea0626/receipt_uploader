package resize_queue_test

import (
	"testing"
	"time"

	"receipt_uploader/internal/constants"
	images_mock "receipt_uploader/internal/images/mock"
	"receipt_uploader/internal/models/image_meta"
	"receipt_uploader/internal/models/tasks"
	"receipt_uploader/internal/resize_queue"

	"github.com/stretchr/testify/assert"
)

func TestEnqueue(t *testing.T) {
	queueCapacity := 3
	mockImagesService := &images_mock.ServiceMock{}
	queue := resize_queue.NewService(queueCapacity, mockImagesService)

	task := tasks.ResizeTask{
		ImageMeta: image_meta.ImageMeta{Path: "test/path"},
		DestDir:   "test/dest",
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
		mockImagesService := &images_mock.ServiceMock{}
		queue := resize_queue.NewService(2, mockImagesService)

		task := tasks.ResizeTask{ImageMeta: image_meta.ImageMeta{Path: "test/path"}, DestDir: "test/destDir"}
		timeout := constants.IMAGE_WORKER_TIMEOUT
		err := queue.WithTimeout(task, timeout)
		assert.Nil(t, err)
	})

	t.Run("should fail, WithTimeout()time out", func(t *testing.T) {
		mockImagesService := &images_mock.ServiceMock{}
		queue := resize_queue.NewService(2, mockImagesService)

		task := tasks.ResizeTask{ImageMeta: image_meta.ImageMeta{Path: "test/path"}, DestDir: "mock_generate_images_timeout"}
		timeout := constants.IMAGE_WORKER_TIMEOUT - 1*time.Second
		err := queue.WithTimeout(task, timeout)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "resizeImages() timed out")
	})

}
