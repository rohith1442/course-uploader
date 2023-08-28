package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"workspace/goproject/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type VideoHandler struct {
	S3        Config
	Service   service.VideoService
	Extension string
	VideoId   string
	VideoDir  string
	Logger    *logrus.Logger
}

type Config struct {
	BucketName string
	Host       string
}

func InitVideoHandler(bucketName, host string, service *service.VideoService) *VideoHandler {
	return &VideoHandler{
		S3: Config{
			BucketName: bucketName,
			Host:       host,
		},
		Service:   *service,
		Logger:    service.Logger,
		Extension: "",
		VideoId:   "",
		VideoDir:  "",
	}
}

func (vh *VideoHandler) HandleTranscode() gin.HandlerFunc {
	return func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["file"]

		file := files[0]
		vh.Extension = strings.ToLower(filepath.Ext(file.Filename))[1:]
		vh.VideoId = uuid.NewString()
		vh.VideoDir = filepath.Join("uploads", vh.VideoId)

		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No input file received. Please send a video file in multipart/form-data format."})
		}
		videoFilePath, err := vh.Service.GenerateFilePath(vh.VideoDir, vh.Extension, vh.VideoId, vh.VideoDir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate file path"})
		}
		err = c.SaveUploadedFile(file, videoFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file."})
		}
		vh.Service.TranscodeVideo(vh.VideoId, vh.Extension, vh.S3.Host)
		c.String(http.StatusOK, "success")
	}
}

func (vh *VideoHandler) HandleBaremetal() gin.HandlerFunc {
	return func(c *gin.Context) {
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
				vh.Service.UploadVideo(&wg, fileName, video, vh.S3.BucketName)
			}(file.Name())
		}
		wg.Wait()
		c.String(http.StatusOK, "Uploaded Video. Video is ready to stream.")
	}
}
