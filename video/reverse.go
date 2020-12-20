package video

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

const tempFilename = "reversed.mp4"

func Session() (*session.Session, error) {
	options := session.Options{}
	// Heroku specifies port
	if os.Getenv("PORT") == "" {
		options.Profile = "jds"
	}

	sess, err := session.NewSessionWithOptions(options)
	if err != nil {
		return nil, err
	}

	sess.Config.WithRegion("us-west-1")
	// if os.Getenv("AWS_ACCESS_KEY_ID") != "" && os.Getenv("AWS_SECRET_ACCESS_KEY") != "" {
	// 	sess.Config.WithCredentials(credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), os.Getenv("AWS_SESSION_TOKEN")))
	// }
	return sess, nil
}

func Reverse(reader io.Reader) (string, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, reader)
	if err != nil {
		return "", err
	}

	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}

	_, err = buf.WriteTo(tmp)
	if err != nil {
		return "", err
	}

	err = exec.Command("ffmpeg", "-i", tmp.Name(), "-vf", "reverse", tempFilename).Run()
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFilename)

	f, err := os.Open(tempFilename)
	if err != nil {
		return "", err
	}

	sess, err := Session()
	if err != nil {
		return "", err
	}

	id := fmt.Sprintf("%s.mov", uuid.New().String())

	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Body:   f,
		Bucket: aws.String("reversevideo"),
		Key:    &id,
	})
	return id, err
}

func GetVideo(key string) (io.ReadCloser, error) {
	sess, err := Session()
	if err != nil {
		return nil, err
	}
	resp, err := s3.New(sess).GetObject(&s3.GetObjectInput{
		Bucket: aws.String("reversevideo"),
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
