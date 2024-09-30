package image_worker

import (
	"log"
	"os"
	"path/filepath"
	"receipt_uploader/internal/images"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/models/image_meta"
	"receipt_uploader/internal/test_utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// type MockImageService struct {
// 	mock.Mock
// }

// func (m *MockImageService) GenerateResizedImages(srcPath string, destDir string) error {
// 	args := m.Called(srcPath, destDir)
// 	return args.Error(0)
// }

func TestResizeImages(t *testing.T) {
	imageService := images.NewService(&configs.AllowedDimensions)
	service := &Service{
		ImageService: imageService,
	}

	baseDir := "image_worker"
	uploadsDir := filepath.Join(baseDir, "uploads")
	destDir := filepath.Join(baseDir, "resized")
	os.MkdirAll(uploadsDir, 0755)
	os.MkdirAll(destDir, 0755)
	defer os.RemoveAll(baseDir)

	username := "user_1"
	extension := "jpg"
	testFilename := username + "#receiptupload123.jpg"
	testFilePath := filepath.Join(uploadsDir, testFilename)
	imageFile := image_meta.FromFormData(username, extension, uploadsDir)
	test_utils.CreateTestImage(imageFile.Path, 1000, 1200)

	resizeErr := service.ResizeImages(uploadsDir, destDir)
	assert.Nil(t, resizeErr)

	smallImagePath := image_meta.GetResizedPath(imageFile, filepath.Join(destDir, username), "small")
	mediumImagePath := image_meta.GetResizedPath(imageFile, filepath.Join(destDir, username), "medium")
	largeImagePath := image_meta.GetResizedPath(imageFile, filepath.Join(destDir, username), "large")

	_, smallErr := os.Stat(smallImagePath)
	assert.Nil(t, smallErr)
	log.Println("smallErr: ", smallErr)
	_, mediumErr := os.Stat(mediumImagePath)
	assert.Nil(t, mediumErr)
	_, largeErr := os.Stat(largeImagePath)
	assert.Nil(t, largeErr)

	_, err := os.Stat(testFilePath)
	assert.True(t, os.IsNotExist(err))
}

// func TestResizeImages_FileError(t *testing.T) {
// 	t.Skip()
// 	mockImageService := new(MockImageService)
// 	service := &Service{
// 		ImageService: mockImageService,
// 	}

// 	// This directory will not exist
// 	srcDir := "non_existent_src"
// 	destDir := "test_dest"

// 	// Act
// 	err := service.ResizeImages(srcDir, destDir)

// 	// Assertions
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "os.ReadDir() failed")
// }

// func TestResizeImages_GenerateFailed(t *testing.T) {
// 	t.Skip()

// 	mockImageService := new(MockImageService)
// 	service := &Service{
// 		ImageService: mockImageService,
// 	}

// 	srcDir := "test_src"
// 	destDir := "test_dest"
// 	testFilename := "image#001.jpg"
// 	testFilePath := filepath.Join(srcDir, testFilename)

// 	// Create a temporary file for testing
// 	os.Mkdir(srcDir, 0755)
// 	defer os.RemoveAll(srcDir) // Cleanup
// 	defer os.RemoveAll(destDir)

// 	file, err := os.Create(testFilePath)
// 	assert.NoError(t, err)
// 	file.Close()

// 	// Expect GenerateResizedImages to return an error
// 	destPathPrefix := filepath.Join(destDir, "image")
// 	mockImageService.On("GenerateResizedImages", testFilePath, destPathPrefix).Return(errors.New("failed to resize"))

// 	// Act
// 	err = service.ResizeImages(srcDir, destDir)

// 	// Assertions
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "s.ImageService.GenerateResizedImages() failed")
// }
