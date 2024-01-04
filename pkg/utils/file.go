package utils

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/aikintech/scim-api/pkg/config"
	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// TODO: Refactor this to constructor pattern
func getS3Client() (*s3.Client, error) {
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.ClientLogMode = aws.LogSigning | aws.LogRequest | aws.LogResponseWithBody
	})

	return client, nil
}

func UploadFileS3(file *multipart.FileHeader, key string) (definitions.Map, error) {
	client, err := getS3Client()
	if err != nil {
		return nil, err
	}

	f, err := file.Open()
	if err != nil {
		return nil, err
	}

	uploader := manager.NewUploader(client)

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(key),
		Body:   f,
		ACL:    "public-read",
	})
	if err != nil {
		return nil, err
	}

	// Generate file URL
	location, err := GenerateS3FileURL(key)
	if err != nil {
		return nil, err
	}

	return definitions.Map{"key": result.Key, "url": location}, err
}

func GenerateS3FileURL(key string) (string, error) {
	if len(key) == 0 {
		return "", errors.New("key is required")
	}

	result, err := config.RedisStore.Get(key)
	if err != nil {
		return "", err
	}

	if len(result) > 0 {
		return string(result), nil
	}

	client, err := getS3Client()
	if err != nil {
		return "", err
	}

	expiration := time.Hour * 24 * 7 // 1 week
	preSignedClient := s3.NewPresignClient(client)
	preSignedURL, err := preSignedClient.PresignGetObject(context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(os.Getenv("AWS_BUCKET")),
			Key:    aws.String(key),
		},
		s3.WithPresignExpires(expiration))
	if err != nil {
		return "", err
	}

	err = config.RedisStore.Set(key, []byte(preSignedURL.URL), expiration)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return preSignedURL.URL, nil
}

func DeleteS3File(key string) error {
	client, err := getS3Client()
	if err != nil {
		return err
	}

	_, err = client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(key),
	})

	if err != nil {
		return err
	}

	// Delete from redis
	if err := config.RedisStore.Delete(key); err != nil {
		fmt.Printf("An error occurred while deleting %s. Error: %s", key, err.Error())
	}

	return nil
}

func UploadFileToYouTube() {}

func UploadFileToTikTok() {}
