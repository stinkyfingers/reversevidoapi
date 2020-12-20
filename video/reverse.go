package video

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

const reversedDir = "reversed"

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
	defer os.Remove(tmp.Name())

	_, err = buf.WriteTo(tmp)
	if err != nil {
		return "", err
	}

	if _, err = os.Stat(reversedDir); os.IsNotExist(err) {
		err = os.Mkdir(reversedDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	id := fmt.Sprintf("%s.mov", uuid.New().String())
	err = exec.Command("ffmpeg", "-i", tmp.Name(), "-vf", "reverse", filepath.Join(reversedDir, id)).Run()
	if err != nil {
		return "", err
	}
	return id, nil
}

func GetVideo(key string) (io.ReadCloser, error) {
	f, err := os.Open(filepath.Join(reversedDir, key))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func Cleanup(key string) error {
	err := os.Remove(filepath.Join(reversedDir, key))
	if err != nil {
		return err
	}

	return cleanup()
}

func cleanup() error {
	infos, err := ioutil.ReadDir(reversedDir)
	if err != nil {
		return err
	}
	for _, info := range infos {
		if time.Now().Add(time.Hour).After(info.ModTime()) {
			err = os.Remove(filepath.Join(reversedDir, info.Name()))
			if err != nil {
				log.Print("Error cleaning up file: ", err)
			}
		}
	}
	return nil
}
