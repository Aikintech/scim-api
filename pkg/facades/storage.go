package facades

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type awsS3 struct {
	ctx      context.Context
	instance *s3.Client
	bucket   string
	url      *string
}

func NewS3() (*awsS3, error) {
	accessKeyId := Env().GetString("AWS_ACCESS_KEY_ID", "")
	accessKeySecret := Env().GetString("AWS_SECRET_ACCESS_KEY", "")
	region := Env().GetString("AWS_REGION", "")
	bucket := Env().GetString("AWS_BUCKET", "")
	url := Env().GetString("AWS_URL", "")
	if accessKeyId == "" || accessKeySecret == "" || region == "" || bucket == "" {
		return nil, fmt.Errorf("please set aws env configuration first")
	}

	client := s3.New(s3.Options{
		Region:      region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
	})

	return &awsS3{
		ctx:      context.TODO(),
		instance: client,
		bucket:   bucket,
		url:      &url,
	}, nil
}

func Storage() *awsS3 {
	s3, err := NewS3()
	if err != nil {
		panic(err.Error())
	}

	return s3
}

func (r *awsS3) Copy(originFile, targetFile string) error {
	_, err := r.instance.CopyObject(r.ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(r.bucket),
		CopySource: aws.String(r.bucket + "/" + originFile),
		Key:        aws.String(targetFile),
	})

	return err
}

func (r *awsS3) Delete(files ...string) error {
	var objectIdentifiers []types.ObjectIdentifier
	for _, file := range files {
		objectIdentifiers = append(objectIdentifiers, types.ObjectIdentifier{
			Key: aws.String(file),
		})
	}

	_, err := r.instance.DeleteObjects(r.ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(r.bucket),
		Delete: &types.Delete{
			Objects: objectIdentifiers,
			Quiet:   aws.Bool(true),
		},
	})

	return err
}

func (r *awsS3) DeleteDirectory(directory string) error {
	if !strings.HasSuffix(directory, "/") {
		directory += "/"
	}

	listObjectsV2Response, err := r.instance.ListObjectsV2(r.ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(r.bucket),
		Prefix: aws.String(directory),
	})
	if err != nil {
		return err
	}
	if len(listObjectsV2Response.Contents) == 0 {
		return nil
	}

	for {
		for _, item := range listObjectsV2Response.Contents {
			_, err = r.instance.DeleteObject(r.ctx, &s3.DeleteObjectInput{
				Bucket: aws.String(r.bucket),
				Key:    item.Key,
			})
			if err != nil {
				return err
			}
		}

		if *listObjectsV2Response.IsTruncated {
			listObjectsV2Response, err = r.instance.ListObjectsV2(r.ctx, &s3.ListObjectsV2Input{
				Bucket:            aws.String(r.bucket),
				ContinuationToken: listObjectsV2Response.ContinuationToken,
			})
			if err != nil {
				return err
			}
		} else {
			break
		}
	}

	return nil
}

func (r *awsS3) Directories(path string) ([]string, error) {
	var directories []string
	validPath := validPath(path)
	listObjsResponse, err := r.instance.ListObjectsV2(r.ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(r.bucket),
		Delimiter: aws.String("/"),
		Prefix:    aws.String(validPath),
	})
	if err != nil {
		return nil, err
	}
	for _, commonPrefix := range listObjsResponse.CommonPrefixes {
		directories = append(directories, strings.ReplaceAll(*commonPrefix.Prefix, validPath, ""))
	}

	return directories, nil
}

func (r *awsS3) Exists(file string) bool {
	_, err := r.instance.HeadObject(r.ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(file),
	})

	return err == nil
}

func (r *awsS3) Files(path string) ([]string, error) {
	var files []string
	validPath := validPath(path)
	listObjsResponse, err := r.instance.ListObjectsV2(r.ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(r.bucket),
		Delimiter: aws.String("/"),
		Prefix:    aws.String(validPath),
	})
	if err != nil {
		return nil, err
	}
	for _, object := range listObjsResponse.Contents {
		file := strings.ReplaceAll(*object.Key, validPath, "")
		if file == "" {
			continue
		}

		files = append(files, file)
	}

	return files, nil
}

func (r *awsS3) Get(file string) (string, error) {
	data, err := r.GetBytes(file)

	return string(data), err
}

func (r *awsS3) GetBytes(file string) ([]byte, error) {
	resp, err := r.instance.GetObject(r.ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(file),
	})
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *awsS3) LastModified(file string) (time.Time, error) {
	resp, err := r.instance.HeadObject(r.ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(file),
	})
	if err != nil {
		return time.Time{}, err
	}

	l, err := time.LoadLocation("UTC")
	if err != nil {
		return time.Time{}, err
	}

	return aws.ToTime(resp.LastModified).In(l), nil
}

func (r *awsS3) MimeType(file string) (string, error) {
	resp, err := r.instance.HeadObject(r.ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(file),
	})
	if err != nil {
		return "", err
	}

	return aws.ToString(resp.ContentType), nil
}

func (r *awsS3) Missing(file string) bool {
	return !r.Exists(file)
}

func (r *awsS3) Move(oldFile, newFile string) error {
	if err := r.Copy(oldFile, newFile); err != nil {
		return err
	}

	return r.Delete(oldFile)
}

func (r *awsS3) Path(file string) string {
	return file
}

func (r *awsS3) Put(file string, content *multipart.FileHeader) error {
	if content == nil {
		return fmt.Errorf("multipart.FileHeader is nil")
	}

	fileContents, err := content.Open()
	if err != nil {
		return err
	}
	defer fileContents.Close()

	_, err = r.instance.PutObject(r.ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(file),
		ContentType: aws.String(content.Header.Get("Content-Type")),
		Body:        fileContents,
	})

	return err
}

func (r *awsS3) Size(file string) (int64, error) {
	resp, err := r.instance.HeadObject(r.ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(file),
	})
	if err != nil {
		return 0, err
	}

	return *resp.ContentLength, nil
}

func (r *awsS3) TemporaryUrl(file string, t time.Time) (string, error) {
	presignClient := s3.NewPresignClient(r.instance)
	presignParams := &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(file),
	}
	presignDuration := func(po *s3.PresignOptions) {
		po.Expires = time.Until(t)
	}
	presignResult, err := presignClient.PresignGetObject(r.ctx, presignParams, presignDuration)
	if err != nil {
		return "", err
	}

	return presignResult.URL, nil
}

func (r *awsS3) Url(file string) string {
	return strings.TrimSuffix(aws.ToString(r.url), "/") + "/" + strings.TrimPrefix(file, "/")
}

func validPath(path string) string {
	realPath := strings.TrimPrefix(path, "./")
	realPath = strings.TrimPrefix(realPath, "/")
	realPath = strings.TrimPrefix(realPath, ".")
	if realPath != "" && !strings.HasSuffix(realPath, "/") {
		realPath += "/"
	}

	return realPath
}

func CloudfrontSignURL(key string) (string, error) {
	expiration := time.Hour * 24 * 7 // 1 week
	privateKey, err := os.ReadFile("pk-APKAW3W3GZGQSDOCLYU6.pem")
	if err != nil {
		return "", err
	}

	// Parse the private key
	privKey, err := parseECPrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	signer := sign.NewURLSigner(Env().GetString("CLOUDFRONT_KEY_ID"), privKey)

	signedURL, err := signer.Sign(Env().GetString("CLOUDFRONT_URL"), time.Now().Add(expiration))
	if err != nil {
		return "", err
	}

	return signedURL, nil
}

func parseECPrivateKey(privKeyBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privKeyBytes)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}
