package controllers

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"workspace/goproject/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var VALID_VIDEO_EXTENSIONS = []string{"mp4", "avi", "mkv"}

const (
	S3_BUCKET = "vidhyatech-course"
	S3_URL    = "s3://vidhyatech-course/"
)

func HandleTranscode() gin.HandlerFunc {
	return func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["file"]
		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No input file received. Please send a video file in multipart/form-data format."})
			return
		}

		file := files[0]
		extension := strings.ToLower(filepath.Ext(file.Filename))[1:]

		if !isValidVideoExtension(extension) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Video format is not supported."})
			return
		}

		videoID := uuid.NewString()
		videoDir := filepath.Join("uploads", videoID)
		err := os.Mkdir(videoDir, os.ModePerm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory."})
			return
		}

		videoFilePath := filepath.Join(videoDir, fmt.Sprintf("%s.%s", videoID, extension))
		err = c.SaveUploadedFile(file, videoFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file."})
			return
		}

		go func() {
			cmd := exec.Command("bash", "create-hls-vod.sh", videoID, extension, S3_URL)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				fmt.Println("Failed to execute transcoding script:", err)
				// You can handle the error condition here if needed
			}
		}()

		c.String(http.StatusOK, "success")
	}

}

func HandleBaremetal() gin.HandlerFunc {
	return func(c *gin.Context) {
		AccessKeyID := utils.Getenv("AWS_ACCESS_KEY_ID")
		SecretAccessKey := utils.Getenv("AWS_SECRET_ACCESS_KEY")
		MyRegion := utils.Getenv("AWS_REGION")
		video := c.Param("video")
		if video == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required video id."})
			return
		}

		files, err := os.ReadDir(filepath.Join("uploads", video))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read directory."})
			return
		}

		var wg sync.WaitGroup
		for _, file := range files {
			if file.Name() == ".DS_Store" || strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())) == video {
				continue
			}

			wg.Add(1)
			go func(fileName string) {
				defer wg.Done()

				extension := strings.ToLower(filepath.Ext(fileName))[1:]
				mimeType := mime.TypeByExtension("." + extension)

				filePath := filepath.Join("uploads", video, fileName)

				// Create an AWS session.
				sess, err := utils.CreateSession(AccessKeyID, SecretAccessKey, MyRegion)
				if err != nil {
					fmt.Println("Error creating AWS session:", err)
					return
				}

				if err := utils.UploadObject(sess, S3_BUCKET, video, fileName, filePath, mimeType); err != nil {
					fmt.Println("Error uploading object to S3:", err)
					return
				}
			}(file.Name())
		}

		wg.Wait()

		c.String(http.StatusOK, "Uploaded Video. Video is ready to stream.")
	}
}

func isValidVideoExtension(extension string) bool {
	for _, validExt := range VALID_VIDEO_EXTENSIONS {
		if validExt == extension {
			return true
		}
	}
	return false
}

// func moveTranscodedFiles(videoID, extension string) {
// 	sourceDir := filepath.Join("uploads", videoID)
// 	destinationDir := "stream" // Change to the path of your local "stream" directory

// 	// List files in the source directory
// 	files, err := ioutil.ReadDir(sourceDir)
// 	if err != nil {
// 		fmt.Println("Failed to list files in source directory:", err)
// 		return
// 	}

// 	// Move each file to the destination directory
// 	for _, file := range files {
// 		sourceFilePath := filepath.Join(sourceDir, file.Name())
// 		destinationFilePath := filepath.Join(destinationDir, file.Name())

// 		// Move the file
// 		err := os.Rename(sourceFilePath, destinationFilePath)
// 		if err != nil {
// 			fmt.Println("Failed to move file:", err)
// 			// Handle error as needed
// 		}
// 	}

// 	// Remove the empty source directory
// 	err = os.Remove(sourceDir)
// 	if err != nil {
// 		fmt.Println("Failed to remove source directory:", err)
// 		// Handle error as needed
// 	}
// }
