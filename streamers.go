package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"syscall"
	"time"

	"github.com/subosito/gozaru"
)

type Streamer struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Image         string `json:"image"`
	FollowerCount int    `json:"followerCount"`
	IsLive        bool   `json:"isLive"`
	IsDownloading bool   `json:"isDownloading"`
	AutoDownload  bool   `json:"autoDownload"`
}

type DownloadProcess struct {
	Command *exec.Cmd
	Err     error
	StdErr  chan string
}

type Streamers struct {
	Infos     map[string]*Streamer
	Processes map[string]*DownloadProcess
}

func (s *Streamers) UpdateFile() error {
	idFile, err := os.OpenFile(os.Getenv("id_file"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer idFile.Close()

	for k := range s.Infos {
		_, err = idFile.Write([]byte(k + "\n"))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Streamers) Load(idFilePath string) error {
	if s.Infos == nil {
		s.Infos = map[string]*Streamer{}
	}
	if s.Processes == nil {
		s.Processes = map[string]*DownloadProcess{}
	}

	idFile, err := os.OpenFile(idFilePath, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer idFile.Close()

	idList := []string{}
	scanner := bufio.NewScanner(idFile)
	for scanner.Scan() {
		idList = append(idList, scanner.Text())
	}

	for _, id := range idList {
		err := s.UpdateStreamer(id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Streamers) Watch() {
	ticker := time.NewTicker(30 * time.Second)

	go func() {
		for ; true; <-ticker.C {
			log.Println("fetching infos")
			for id := range s.Infos {
				err := s.UpdateStreamer(id)
				if err != nil {
					log.Printf("update streamer error: %s\n", err)
					continue
				}
				if s.Infos[id].IsLive && s.Infos[id].AutoDownload && s.Processes[id] == nil {
					s.StartDownload(id)
				}
			}
		}
	}()
}

func (s *Streamers) UpdateStreamer(id string) error {
	info, err := GetChannelData(id)
	if err != nil {
		return err
	}
	if s.Infos[id] == nil {
		s.Infos[id] = &Streamer{
			Id:            info.Id,
			Name:          info.Name,
			Image:         info.Image,
			FollowerCount: info.FollowerCount,
			IsLive:        info.IsLive,
			IsDownloading: s.Processes[id] != nil,
			AutoDownload:  true,
		}
	} else {
		s.Infos[id].Name = info.Name
		s.Infos[id].Image = info.Image
		s.Infos[id].FollowerCount = info.FollowerCount
		s.Infos[id].IsLive = info.IsLive
	}
	return nil
}

func (s *Streamers) AddStreamer(id string) error {
	err := s.UpdateStreamer(id)
	if err != nil {
		return err
	}

	err = s.UpdateFile()
	if err != nil {
		return err
	}

	return nil
}

func (s *Streamers) DeleteStreamer(id string) error {
	if s.Infos[id] == nil {
		return errors.New("streamer with id doesnt exist")
	}

	if s.Processes[id] != nil {
		err := s.StopDownload(id)
		if err != nil {
			return err
		}
	}

	delete(s.Infos, id)
	err := s.UpdateFile()
	if err != nil {
		return err
	}

	return nil
}

func (s *Streamers) StartDownload(id string) error {
	if s.Infos[id] == nil {
		return errors.New("streamers info not found")
	}
	if s.Processes[id] != nil {
		return errors.New("download process already exists")
	}
	log.Printf("start downloading %s\n", id)
	playlistUrl, err := GetPlaylistUrl(id)
	if err != nil {
		return err
	}

	outDir := os.Getenv("dir")
	filename := gozaru.Sanitize(fmt.Sprintf("%s-%s.mp4", s.Infos[id].Name, FormatDate()))
	outputPath := path.Join(outDir, filename)
	args := []string{
		"-headers", "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
		"-loglevel", "error",
		"-i", playlistUrl,
		"-c", "copy",
		outputPath,
	}
	log.Printf("ffmpeg %+v\n", args)
	cmd := exec.Command("ffmpeg", args...)
	go (func() {
		s.Processes[id] = &DownloadProcess{
			Command: cmd,
			Err:     nil,
		}
		s.Infos[id].IsDownloading = true
		CmdLogger(s.Infos[id].Name, cmd)

		err = cmd.Run()
		s.Infos[id].IsDownloading = false
		log.Printf("stopped downloading %s\n", id)
		delete(s.Processes, id)
		if err != nil {
			log.Printf("download %s process error: %s", id, err)
		}
	})()
	return nil
}

func (s *Streamers) StopDownload(id string) error {
	if s.Infos[id] == nil {
		return errors.New("streamers info not found")
	}
	if s.Processes[id] == nil {
		return errors.New("download process doesn't exists")
	}

	s.Infos[id].IsDownloading = false
	err := s.Processes[id].Command.Process.Signal(syscall.SIGINT)
	if err != nil {
		s.Infos[id].IsDownloading = true
		return err
	}
	return nil
}
