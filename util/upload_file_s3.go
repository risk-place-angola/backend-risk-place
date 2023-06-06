package util

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func UploadFileToS3(file *multipart.FileHeader) (*s3manager.UploadOutput, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(src)

	size := file.Size
	buffer := make([]byte, size)

	if _, err := src.Read(buffer); err != nil {
		return nil, err
	}
	BodyFile := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	path := fmt.Sprint("media/", generateFileName(file))

	sess, err := AwsSession()
	if err != nil {
		return nil, err
	}

	svc := s3manager.NewUploader(sess)

	putObjectOutput, err := svc.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(os.Getenv("AWS_BUCKET")),
		Key:         aws.String(path),
		Body:        BodyFile,
		ContentType: aws.String(fileType),
	})
	if err != nil {
		return nil, err
	}

	log.Println("Successfully uploaded")

	return putObjectOutput, nil
}

func generateFileName(file *multipart.FileHeader) string {
	return fmt.Sprint(time.Now().UnixNano(), "_", file.Filename)
}
