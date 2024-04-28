package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
)

type LiveDetail struct {
	Content LiveDetailContent `json:"content"`
}

type LiveDetailContent struct {
	LivePlaybackJson string `json:"livePlaybackJson"`
}

type PlaybackData struct {
	Media []MediaData `json:"media"`
}

type MediaData struct {
	MediaId  string `json:"mediaId"`
	Path     string `json:"path"`
	Protocol string `json:"protocol"`
}

func GetPlaylistUrl(channelId string) (string, error) {
	url := fmt.Sprintf("https://api.chzzk.naver.com/service/v2/channels/%s/live-detail", channelId)
	res, err := client.Get(url)
	if err != nil {
		return "", err
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	liveDetail := &LiveDetail{}
	err = json.Unmarshal(bytes, liveDetail)
	if err != nil {
		return "", err
	}

	playbackStr := liveDetail.Content.LivePlaybackJson
	playbackData := &PlaybackData{}
	err = json.Unmarshal([]byte(playbackStr), playbackData)
	if err != nil {
		return "", nil
	}

	if len(playbackData.Media) == 0 {
		return "", errors.New("no media data")
	}

	var targetMediaData *MediaData = nil
	for _, mediaData := range playbackData.Media {
		if mediaData.MediaId == "LLHLS" {
			targetMediaData = &mediaData
		}
	}
	if targetMediaData == nil {
		targetMediaData = &playbackData.Media[0]
	}

	return targetMediaData.Path, nil
}

type ChannelResponse struct {
	Content ChannelData `json:"content"`
}

type ChannelData struct {
	Id     string `json:"channelId"`
	Name   string `json:"channelName"`
	Image  string `json:"channelImageUrl"`
	IsLive bool   `json:"openLive"`
	FollowerCount int `json:"followerCount"`
}

func GetChannelData(id string) (*ChannelData, error) {
	res, err := client.Get(fmt.Sprintf("https://api.chzzk.naver.com/service/v1/channels/%s", id))
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	data := &ChannelResponse{}
	err = json.Unmarshal(bytes, data)
	if err != nil {
		return nil, err
	}

	return &data.Content, nil
}

type SearchResponse struct {
	Content SearchContent `json:"content"`
}

type SearchContent struct {
	Data []SearchData `json:"data"`
}

type SearchData struct {
	Channel ChannelData `json:"channel"`
}

func SearchChannel(query string) ([]*ChannelData, error) {
	escapedQuery := url.QueryEscape(query)
	url := fmt.Sprintf("https://api.chzzk.naver.com/service/v1/search/channels?keyword=%s&offset=0", escapedQuery)
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	searchResponse := &SearchResponse{}
	err = json.Unmarshal(bytes, searchResponse)
	if err != nil {
		return nil, err
	}

	channels := []*ChannelData{}
	for _, searchData := range searchResponse.Content.Data {
		channels = append(channels, &searchData.Channel)
	}

	return channels, nil
}
