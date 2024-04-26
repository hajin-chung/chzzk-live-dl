package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/subosito/gozaru"
)

type Streamer struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Image         string `json:"image"`
	IsLive        bool   `json:"isLive"`
	IsDownloading bool   `json:"isDownloading"`
	AutoDownload  bool   `json:"autoDownload"`
}

type DownloadProcess struct {
	Command *exec.Cmd
	Err     error
}

type Streamers struct {
	Infos     map[string]*Streamer
	Processes map[string]*DownloadProcess
}

func (s *Streamers) UpdateFile() error {
	idFile, err := os.OpenFile("ids.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
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

func (s *Streamers) Load() error {
	if s.Infos == nil {
		s.Infos = map[string]*Streamer{}
	}
	if s.Processes == nil {
		s.Processes = map[string]*DownloadProcess{}
	}

	idFile, err := os.OpenFile("ids.txt", os.O_RDONLY|os.O_CREATE, 0600)
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
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for ; true; <-ticker.C {
			log.Println("fetching infos")
			for id := range s.Infos {
				err := s.UpdateStreamer(id)
				if err != nil {
					log.Printf("update streamer error: %s\n", err)
					continue
				}
				if s.Infos[id].AutoDownload == true && s.Infos[id].IsDownloading == false && s.Processes[id] == nil {
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
	s.Infos[id] = &Streamer{
		Id:            info.Id,
		Name:          info.Name,
		Image:         info.Image,
		IsLive:        info.IsLive,
		IsDownloading: s.Processes[id] != nil,
		AutoDownload:  true,
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
	delete(s.Infos, id)
	err := s.UpdateFile()
	if err != nil {
		return err
	}

	err = s.StopDownload(id)
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

	outputPath := gozaru.Sanitize(fmt.Sprintf("%s-%s.mp4", s.Infos[id].Name, FormatDate()))
	args := []string{"-loglevel", "error", "-i", playlistUrl, "-c", "copy", outputPath}
	log.Printf("ffmpeg %+v\n", args)
	cmd := exec.Command("ffmpeg", args...)
	go (func() {
		s.Processes[id] = &DownloadProcess{
			Command: cmd,
			Err:     nil,
		}
		s.Infos[id].IsDownloading = true

		err := cmd.Run()
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

	err := s.Processes[id].Command.Process.Signal(syscall.SIGINT)
	if err != nil {
		return err
	}
	return nil
}
