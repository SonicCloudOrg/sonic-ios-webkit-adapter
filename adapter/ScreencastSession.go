/*
 *  Copyright (C) [SonicCloudOrg] Sonic Project
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */
package adapters

import (
	"github.com/tidwall/gjson"
	"log"
	"strconv"
	"strings"
	"time"
)

type screencastSession struct {
	adapter         *Adapter
	frameId         int
	framesAcked     map[int]bool
	frameInterval   time.Duration // default 250, 60 fps is 16ms
	format          string
	quality         int
	maxWidth        int
	maxHeight       int
	timerCookie     interface{}
	deviceWidth     int
	deviceHeight    int
	offsetTop       int
	pageScaleFactor int
	scrollOffsetX   int
	scrollOffsetY   int
	closeFlag       chan bool
}

func newScreencastSession(adapter *Adapter, optFuncs ...ScreencastOptFunc) *screencastSession {
	screencast := &screencastSession{
		adapter:   adapter,
		quality:   100,
		format:    "jpg",
		maxHeight: 1024,
		maxWidth:  1024,
	}
	screencast.frameInterval = 250

	screencast.closeFlag = make(chan bool)

	if len(optFuncs) != 0 {
		for _, optFunc := range optFuncs {
			optFunc(screencast)
		}
	}
	return screencast
}

func (s *screencastSession) start() {
	s.framesAcked = make(map[int]bool)
	s.frameId = 1
	var err error
	params := map[string]interface{}{
		"expression": `(window.innerWidth > 0 ? window.innerWidth : screen.width) + "," + (window.innerHeight > 0 ? window.innerHeight : screen.height) + "," + window.devicePixelRatio`,
	}
	s.adapter.CallTarget("Runtime.evaluate", params, func(message []byte) {
		parts := strings.Split(gjson.Get(string(message), "result.value").String(), ",")
		var deviceWidth int
		var deviceHeight int
		var pageScaleFactor int
		deviceWidth, err = strconv.Atoi(parts[0])
		if err != nil {
			log.Println(err)
		}
		deviceHeight, err = strconv.Atoi(parts[1])
		if err != nil {
			log.Println(err)
		}
		pageScaleFactor, err = strconv.Atoi(parts[2])
		if err != nil {
			log.Println(err)
		}
		s.deviceWidth = deviceWidth
		s.deviceHeight = deviceHeight
		s.pageScaleFactor = pageScaleFactor

		ticker := time.NewTicker(s.frameInterval * time.Millisecond)
		go func() {
			for {
				select {
				case <-ticker.C:
					s.recordingLoop()
				case <-s.closeFlag:
					return
				}
			}
		}()
	})
}

func (s *screencastSession) stop() {
	if s.closeFlag == nil {
		return
	}
	s.closeFlag <- true
}

func (s *screencastSession) ackFrame(frameNumber int) {
	s.framesAcked[frameNumber] = true
}

func (s *screencastSession) recordingLoop() {
	currentFrame := s.frameId
	frameAckFlag, ok := s.framesAcked[currentFrame-1]
	if currentFrame > 1 && (ok || !frameAckFlag) {
		return
	}
	s.frameId++
	params := map[string]interface{}{
		"expression": `window.document.body.offsetTop + "," + window.pageXOffset + "," + window.pageYOffset`,
	}
	s.adapter.CallTarget("Runtime.evaluate", params, func(message []byte) {
		if !gjson.Get(string(message), "wasThrown").Exists() || gjson.Get(string(message), "wasThrown").Bool() {
			return
		}
		parts := strings.Split(gjson.Get(string(message), "result.value").String(), ",")
		var offsetTop int
		var scrollOffsetX int
		var scrollOffsetY int
		var err error
		offsetTop, err = strconv.Atoi(parts[0])
		if err != nil {
			log.Println(err)
		}
		scrollOffsetX, err = strconv.Atoi(parts[1])
		if err != nil {
			log.Println(err)
		}
		scrollOffsetY, err = strconv.Atoi(parts[2])
		if err != nil {
			log.Println(err)
		}
		s.offsetTop = offsetTop
		s.scrollOffsetY = scrollOffsetY
		s.scrollOffsetX = scrollOffsetX

		snapshotRectParams := map[string]interface{}{
			"x":                0,
			"y":                0,
			"width":            s.deviceWidth,
			"height":           s.deviceHeight,
			"coordinateSystem": "Viewport",
		}
		s.adapter.CallTarget("Page.snapshotRect", snapshotRectParams, func(msg []byte) {
			dataURL := gjson.Get(string(msg), "dataURL").String()
			index := strings.Index(dataURL, "base64")

			frame := map[string]interface{}{
				"data": dataURL[index+7:],
				"metadata": map[string]interface{}{
					"pageScaleFactor": s.pageScaleFactor,
					"offsetTop":       s.offsetTop,
					"deviceWidth":     s.deviceWidth,
					"deviceHeight":    s.deviceHeight,
					"scrollOffsetX":   s.scrollOffsetX,
					"scrollOffsetY":   s.scrollOffsetY,
					"timestamp":       time.Now().UnixNano(),
				},
				"sessionId": currentFrame,
			}

			s.adapter.FireEventToTools("Page.screencastFrame", frame)
		})

	})
}

type ScreencastOptFunc func(screencast *screencastSession)

func WithFormat(format string) ScreencastOptFunc {
	return func(screencast *screencastSession) {
		screencast.format = format
	}
}

func WithQuality(quality int) ScreencastOptFunc {
	return func(screencast *screencastSession) {
		screencast.quality = quality
	}
}

func WithMaxHeight(maxHeight int) ScreencastOptFunc {
	return func(screencast *screencastSession) {
		screencast.maxHeight = maxHeight
	}
}

func WithMaxWidth(maxWidth int) ScreencastOptFunc {
	return func(screencast *screencastSession) {
		screencast.maxWidth = maxWidth
	}
}
