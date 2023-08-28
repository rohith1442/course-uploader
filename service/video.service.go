package service

import (
	"errors"
	"fmt"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"workspace/goproject/constants"
	"workspace/goproject/helpers"
	"workspace/goproject/utils"

	"github.com/sirupsen/logrus"
)

type VideoService struct {
	Logger *logrus.Logger
}

func (vs *VideoService) GenerateFilePath(filePath, extension, videoId, videoDir string) (string, error) {

	if !helpers.IsValidVideoExtension(extension, constants.VIDEO_EXTENSIONS) {
		vs.Logger.Error("Video format is not supported.")
		return "", errors.New("invalid Video Format")
	}

	err := os.Mkdir(videoDir, os.ModePerm)
	if err != nil {
		vs.Logger.Error("Unable to make Directory")
		return "", errors.New("mkdir error")
	}

	videoFilePath := filepath.Join(videoDir, fmt.Sprintf("%s.%s", videoId, extension))
	return videoFilePath, nil

}

func (vs *VideoService) TranscodeVideo(videoId, extension, s3Host string) {
	go vs.transcodeVideoToHLS(constants.V240P, videoId, extension, s3Host)
	go vs.transcodeVideoToHLS(constants.V360P, videoId, extension, s3Host)
	go vs.transcodeVideoToHLS(constants.V480P, videoId, extension, s3Host)
	go vs.transcodeVideoToHLS(constants.V720P, videoId, extension, s3Host)
	go vs.transcodeVideoToHLS(constants.V1080P, videoId, extension, s3Host)
}

func (vs *VideoService) transcodeVideoToHLS(renditions, videoId, entension, s3Host string) {
	cmd := exec.Command("bash", "create-hls-vod.sh", videoId, entension, s3Host, renditions)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		vs.Logger.Error("Failed to execute transcoding script:", err)
		// return err
	}
	// return nil
}

func (vs *VideoService) UploadVideo(wg *sync.WaitGroup, fileName, video, bucketName string) {
	defer wg.Done()
	AccessKeyID := utils.Getenv("AWS_ACCESS_KEY_ID")
	SecretAccessKey := utils.Getenv("AWS_SECRET_ACCESS_KEY")
	MyRegion := utils.Getenv("AWS_REGION")

	extension := strings.ToLower(filepath.Ext(fileName))[1:]
	mimeType := mime.TypeByExtension("." + extension)

	filePath := filepath.Join("uploads", video, fileName)

	sess, err := utils.CreateSession(AccessKeyID, SecretAccessKey, MyRegion)
	if err != nil {
		fmt.Println("Error creating AWS session:", err)
		return
	}

	if err := utils.UploadObject(sess, bucketName, video, fileName, filePath, mimeType); err != nil {
		fmt.Println("Error uploading object to S3:", err)
		return
	}
}