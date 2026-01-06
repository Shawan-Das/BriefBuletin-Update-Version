package util

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var _awsLogger = logrus.New()

func NewS3Session() (*session.Session, error) {
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(viper.GetViper().GetStringMapString("aws")["s3_region"]),
		Credentials: credentials.NewStaticCredentials(
			viper.GetViper().GetStringMapString("aws")["s3_access_key_id"],     // id
			viper.GetViper().GetStringMapString("aws")["s3_secret_access_key"], // secret
			""), // token can be left blank for now
	})
	if err != nil {
		_awsLogger.Errorf("Error while connecting aws %v", err)
	}
	return s, err
}

// UploadFileToS3 saves a file to aws bucket and returns the url to the file and an error if there's any
func UploadFileToS3(s *session.Session, file multipart.File, fileHeader *multipart.FileHeader, fileName string) (string, error) {
	// get the file size and read
	// the file content into a buffer
	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)

	// create a unique file name for the file
	tempFileName := fileName + filepath.Ext(fileHeader.Filename)
	_awsLogger.Info("UploadFileToS3 ", tempFileName)
	// config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file
	// you're uploading
	//input *s3.PutObjectInput

	_, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(viper.GetViper().GetStringMapString("aws")["s3_bucket"]),
		Key:                  aws.String(tempFileName),
		ACL:                  aws.String("public-read"), // could be private if you want it to be access by only authorized users
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("inline"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})
	if err != nil {
		return "", err
	}

	return tempFileName, err
}

func GetfromS3(s *session.Session, fileHeader *multipart.FileHeader, path string) (*s3.GetObjectOutput, error) {

	// create a unique file name for the file
	tempFileName := viper.GetViper().GetStringMapString("aws")["s3_folder"] + "/" + path + filepath.Ext(fileHeader.Filename)
	_awsLogger.Info("UploadFileToS3 ", tempFileName)

	getObjectOutput, err := s3.New(s).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(viper.GetViper().GetStringMapString("aws")["s3_bucket"]),
		Key:    aws.String(tempFileName),
	})
	if err != nil {
		return nil, err
	}

	return getObjectOutput, err
}

// GetfromS3 get files from aws bucket and returns the url to the file and an error if there's any
func GetListOfFileFromS3(s *session.Session, path string) ([]string, error) {

	// create a unique file name for the file
	ctx := context.Background()
	bucket := viper.GetStringMapString("aws")["s3_bucket"]
	s3Keys := make([]string, 0)
	//tempFileName := viper.GetViper().GetStringMapString("aws")["s3_folder"] + "/" +
	// list files under `blog` directory in `work-with-s3` bucket.
	if err := s3.New(s).ListObjectsPagesWithContext(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(path + "/"), // list files in the directory.
	}, func(o *s3.ListObjectsOutput, b bool) bool { // callback func to enable paging.
		for _, o := range o.Contents {
			s3Keys = append(s3Keys, *o.Key)
		}
		return true
	}); err != nil {
		_awsLogger.Fatalf("failed to list items in s3 directory: %v", err)
	}

	return s3Keys, nil
}

// DeletefromS3 delete files from aws bucket
func DeletefromS3(s *session.Session, path string) (string, error) {

	// create a unique file name for the file
	tempFileName := path
	_awsLogger.Info("UploadFileToS3 : ", tempFileName)

	_, err := s3.New(s).DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(viper.GetViper().GetStringMapString("aws")["s3_bucket"]),
		Key:    aws.String(tempFileName),
	})

	if err != nil {
		return "", err
	}

	return tempFileName, err
}

func DeleteFilesFromS3(s *session.Session, paths []string) (bool, []string, []string) {
	if len(paths) == 0 {
		return false, nil, []string{"No file paths provided"}
	}

	s3BaseURL := viper.GetViper().GetStringMapString("aws")["s3_url"]
	bucket := viper.GetViper().GetStringMapString("aws")["s3_bucket"]

	s3Client := s3.New(s)

	var deletedFiles []string
	var errors []string

	for _, path := range paths {
		if path == "" {
			errors = append(errors, "Empty path provided")
			continue
		}
		// Remove base URL if full link provided
		fileKey := strings.TrimPrefix(path, s3BaseURL+"/")
		if fileKey == "" {
			errors = append(errors, fmt.Sprintf("Invalid file path: %s", path))
			continue
		}

		// Try to delete files
		_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(fileKey),
		})
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to delete %s: %v", fileKey, err))
			continue
		}

		// confirm deletion
		err = s3Client.WaitUntilObjectNotExists(&s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(fileKey),
		})
		if err != nil {
			errors = append(errors, fmt.Sprintf("Deletion not confirmed for %s: %v", fileKey, err))
			continue
		}
		// file deleted successfully
		deletedFiles = append(deletedFiles, fileKey)
	}

	success := len(deletedFiles) > 0

	return success, deletedFiles, errors
}
