package main

import (
	"fmt"
	"strconv"
	"time"
)

func commandCreateSession(cfg *config, args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("invalid number of arguments")
	}

	url := fmt.Sprintf("%s/api/sessions", cfg.serverAddress)

	projectName := args[0]
	episodeNumber, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}
	durationTime, err := time.ParseDuration(args[2])
	var duration int64
	duration = durationTime.Microseconds()

	users := []string

	for _, username := range args[3:] {
		userID, err := getUserID(cfg, username)
		if err != nil {
			return err
		}

		users = append(users, userID)
	}

	type getEpisodeID

	type createSesType struct {
		Duration     int64  `json:"duration"`
		EpisodeID    string `json:"episode_id"`
		PartWorkedOn string `json:"part_worked_on"`
		ActivityDone string `json:"activity_done"`
	}
	createSesReq := createSesType{
		Duration: duration,
		EpisodeID: ,
	}

	return nil
}
