package coral

import (
	"errors"
	"io"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/superfishial/reef/server/config"
)

type FileMetadata struct {
	OwnerSub string
	Public   bool
}

func getS3Session(conf config.S3Config) *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:           aws.String(conf.Region),
			Endpoint:         aws.String(conf.Endpoint),
			Credentials:      credentials.NewStaticCredentials(conf.AccessKey, conf.SecretKey, ""),
			S3ForcePathStyle: aws.Bool(conf.ForcePathStyle),
		},
	}))
}

func DownloadFile(conf config.S3Config, filename string) (io.ReadCloser, error) {
	svc := s3.New(getS3Session(conf))
	out, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(conf.Bucket),
		Key:    aws.String(filename),
	})
	return out.Body, err
}

func UploadFile(conf config.S3Config, filename string, metadata FileMetadata, reader io.Reader) error {
	// The session the S3 Uploader will use
	sess := getS3Session(conf)

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Upload the file to S3
	public := strconv.FormatBool(metadata.Public)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(conf.Bucket),
		Key:    aws.String(filename),
		Body:   reader,
		Metadata: map[string]*string{
			"Owner-Sub": &metadata.OwnerSub,
			"Public":    &public,
		},
	})
	return err
}

func DeleteFile(conf config.S3Config, filename string) error {
	svc := s3.New(getS3Session(conf))
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(conf.Bucket),
		Key:    aws.String(filename),
	})
	return err
}

func GetFileMetadata(conf config.S3Config, filename string) (*FileMetadata, error) {
	svc := s3.New(getS3Session(conf))
	out, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(conf.Bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return &FileMetadata{}, err
	}

	OwnerSub, ok := out.Metadata["Owner-Sub"]
	if !ok {
		return nil, errors.New("could not find 'Owner-Sub' metadata in object")
	}
	Public, ok := out.Metadata["Public"]
	if !ok {
		return nil, errors.New("could not find 'Public' metadata in object")
	}

	return &FileMetadata{
		OwnerSub: *OwnerSub,
		Public:   *Public == "true",
	}, nil
}

func SetFileMetadata(conf config.S3Config, filename string, metadata *FileMetadata) error {
	svc := s3.New(getS3Session(conf))
	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(conf.Bucket),
		Key:    aws.String(filename),
		Metadata: map[string]*string{
			"ownerSub": aws.String(metadata.OwnerSub),
			"public":   aws.String(strconv.FormatBool(metadata.Public)),
		},
	})
	return err
}
