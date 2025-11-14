package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Denisowiec/FoleyBookkeeper/internal/db"
)

func commandCreateEpisode(cfg *config, args []string) error {
	// This command creates a new episodes
	// Takes project title, episode number and title as arguments
	if len(args) < 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	prj, err := getProjectByName(cfg, args[0])
	if err != nil {
		return nil
	}

	projectID := prj.ID.String()
	var title string
	var epNumber int

	if len(args) >= 2 {
		epNumber, err = strconv.Atoi(args[1])
		if err != nil {
			return err
		}
	}
	if len(args) >= 3 {
		title = args[2]
	}

	url := fmt.Sprintf("%s/api/episodes", cfg.serverAddress)

	type createEpisodeReqType struct {
		Title         string `json:"title"`
		ProjectID     string `json:"project_id"`
		EpisodeNumber int    `json:"episode_number"`
	}
	createEpisodeReq := createEpisodeReqType{
		Title:         title,
		EpisodeNumber: epNumber,
		ProjectID:     projectID,
	}

	resp, err := sendRequest(createEpisodeReq, "POST", url, cfg.jwt)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return processErrorResponse(resp)
	}

	ep := db.Episode{}
	err = processResponse(resp, &ep)
	if err != nil {
		return err
	}

	if epNumber > 0 {
		fmt.Printf("Episode for project %s created successfully\n", prj.Title)
	} else {
		fmt.Printf("Episode %d for project %s created successfully\n", ep.EpisodeNumber, prj.Title)
	}
	return nil
}

func commandGetEpisodesForProjects(cfg *config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	prj, err := getProjectByName(cfg, args[0])
	if err != nil {
		return err
	}

	projectID := prj.ID.String()

	url := fmt.Sprintf("%s/api/episodes", cfg.serverAddress)

	type getEpReqType struct {
		ProjectID string `json:"project_id"`
	}
	getEpReq := getEpReqType{
		ProjectID: projectID,
	}

	resp, err := sendRequest(getEpReq, "GET", url, cfg.jwt)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var eps []db.Episode
	err = processResponse(resp, &eps)
	if err != nil {
		return err
	}
	fmt.Println("Episodes of project", prj.Title)
	for _, e := range eps {
		fmt.Printf("Title: %s, Number: %d\n", e.Title.String, e.EpisodeNumber)
	}
	return nil
}
