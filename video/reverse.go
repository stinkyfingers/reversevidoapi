package video

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const reversedDir = "reversed"

type Log struct {
	Status bool   `json:"status"`
	Error  string `json:"error,omitempty"`
}

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
	if err != nil {
		log.Print("ffmpeg err", err)
		return err
	}

	info, err := os.Stat(filepath.Join(reversedDir, id))
	if err != nil {
		log.Print("STAT err", err)
		return err
	}

	log.Print("SIZE SAVE", info.Name(), " ", info.Size())

	return UpdateLog(id, true, "")
}

func GetVideo(key string) (io.ReadCloser, error) {
	f, err := os.Open(filepath.Join(reversedDir, key))
	if err != nil {
		return nil, err
	}
	info, _ := os.Stat(filepath.Join(reversedDir, key))
	log.Print("SIZE GET", info.Name(), " ", info.Size())
	return f, nil
}

func Cleanup(key string) error {
	err := os.Remove(filepath.Join(reversedDir, key))
	if err != nil {
		return err
	}
	err = os.Remove(filepath.Join(reversedDir, fmt.Sprintf("%s.json", key)))
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

func UpdateLog(key string, status bool, errString string) error {
	f, err := os.Open(filepath.Join(reversedDir, fmt.Sprintf("%s.json", key)))
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		f, err = os.Create(filepath.Join(reversedDir, fmt.Sprintf("%s.json", key)))
		if err != nil {
			return err
		}

	}
	defer f.Close()
	return json.NewEncoder(f).Encode(&Log{Status: status, Error: errString})
}

func CheckLog(key string) (*Log, error) {
	f, err := os.Open(filepath.Join(reversedDir, fmt.Sprintf("%s.json", key)))
	if err != nil {
		if os.IsNotExist(err) {
			return &Log{}, nil
		}
		return nil, err
	}
	defer f.Close()
	var l Log
	err = json.NewDecoder(f).Decode(&l)
	return &l, err
}
