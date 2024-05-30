package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func InitS3Client() (s3Client *s3.S3) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return s3.New(sess)
}

func UploadFileToS3(fileName string, fileContent []byte) (string, error) {
	input := &s3.PutObjectInput{
		Bucket: aws.String("neekcodesolutions"),
		Key:    aws.String(fileName),
		Body:   aws.ReadSeekCloser(bytes.NewReader(fileContent)),
	}

	_, err := s3Client.PutObject(input)

	if err != nil {
		log.Print(err)
	}

	s3URL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", "neekcodesolutions", fileName)

	return s3URL, err
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		b[i] = charset[num.Int64()]
	}
	return string(b)
}
