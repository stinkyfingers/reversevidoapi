package video

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const reversedDir = "reversed"
const errorLog = "error.json"

func Reverse(reader io.Reader, id string) error {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, reader)
	if err != nil {
		return err
	}

	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	_, err = buf.WriteTo(tmp)
	if err != nil {
		return err
	}

	if _, err = os.Stat(reversedDir); os.IsNotExist(err) {
		err = os.Mkdir(reversedDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	err = exec.Command("ffmpeg", "-i", tmp.Name(), "-vf", "reverse", filepath.Join(reversedDir, id)).Run()
	info, _ := os.Stat(filepath.Join(reversedDir, id))
	log.Print("SIZE SAVE", info.Name(), info.Size())
	return err
}

func GetVideo(key string) (io.ReadCloser, error) {
	f, err := os.Open(filepath.Join(reversedDir, key))
	if err != nil {
		return nil, err
	}
	info, _ := os.Stat(filepath.Join(reversedDir, key))
	log.Print("SIZE GET", info.Name(), info.Size())
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

func CheckStatus(key string) (bool, error) {
	_, err := os.Stat(filepath.Join(reversedDir, key))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func LogError(key string, errString string) error {
	if _, err := os.Stat(errorLog); os.IsNotExist(err) {
		_, err := os.Create(errorLog)
		if err != nil {
			return err
		}
	}
	f, err := os.Open(errorLog)
	if err != nil {
		return err
	}
	defer f.Close()
	var errors map[string]string
	err = json.NewDecoder(f).Decode(&errors)
	if err != nil {
		return err
	}
	errors[key] = errString
	return json.NewEncoder(f).Encode(&errors)
}

func CheckError(key string) (string, error) {
	if _, err := os.Stat(errorLog); os.IsNotExist(err) {
		return "", nil
	}
	f, err := os.Open(errorLog)
	if err != nil {
		return "", err
	}
	defer f.Close()
	var errors map[string]string
	err = json.NewDecoder(f).Decode(&errors)
	if err != nil {
		return "", err
	}
	return errors[key], nil
}
