package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Denisowiec/FoleyBookkeeper/internal/db"
)

func commandCreateSession(cfg *config, args []string) error {
	// Takes project title, episode number, duration of session, part worked on, activity done and a list of usernames as input
	if len(args) < 6 {
		return fmt.Errorf("invalid number of arguments")
	}

	projectName := args[0]
	episodeNumber, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}
	partWorkedOn := args[3]
	activityDone := args[4]

	// This converts the duration into something that should match the postgresql preferences
	durationTime, err := time.ParseDuration(args[2])
	if err != nil {
		return err
	}

	duration := durationTime.Microseconds()

	// We convert the usernames in the arguments into a list of IDs
	users := []string{}

	for _, username := range args[5:] {
		userID, err := getUserID(cfg, username)
		if err != nil {
			return err
		}

		users = append(users, userID)
	}

	// Now we need the project id
	reqPrjBody := struct {
		ProjectTitle string `json:"title"`
	}{
		ProjectTitle: projectName,
	}

	prj, err := getThing(cfg, "/api/projects", reqPrjBody, db.Project{})
	if err != nil {
		return err
	}

	// Now we need the episode ID
	reqEpBody := struct {
		ProjectID     string `json:"project_id"`
		EpisodeNumber int    `json:"episode_number"`
	}{
		ProjectID:     prj.ID.String(),
		EpisodeNumber: episodeNumber,
	}

	ep, err := getThing(cfg, "/api/episodes", reqEpBody, db.Episode{})
	if err != nil {
		return err
	}

	// Now we have everything we need to record a session

	type createSesType struct {
		Duration     int64  `json:"duration"`
		EpisodeID    string `json:"episode_id"`
		PartWorkedOn string `json:"part_worked_on"`
		ActivityDone string `json:"activity_done"`
	}
	createSesReq := createSesType{
		Duration:     duration,
		EpisodeID:    ep.ID.String(),
		PartWorkedOn: partWorkedOn,
		ActivityDone: activityDone,
	}

	url := fmt.Sprintf("%s/api/sessions", cfg.serverAddress)
	resp, err := sendRequest(createSesReq, "POST", url, cfg.jwt)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusAccepted {
		return processErrorResponse(resp)
	}

	ses := db.Session{}

	err = processResponse(resp, &ses)
	if err != nil {
		return err
	}

	// Now we add the users to the session
	reqUsSesBody := struct {
		Users []string `json:"user_ids"`
	}{
		Users: users,
	}

	url = fmt.Sprintf("%s/api/sessions/%s", cfg.serverAddress, ses.ID.String())

	resp2, err := sendRequest(reqUsSesBody, "POST", url, cfg.jwt)
	if err != nil {
		return err
	}
	if resp2.StatusCode != http.StatusAccepted {
		return processErrorResponse(resp2)
	}

	// Commented out for the time being, because there's no need to extract the json data from the response
	/*resp2Body := struct {
		Users []string `json:"user_ids"`
	}{}

	err = processResponse(resp2, &resp2Body)
	if err != nil {
		return err
	}*/

	fmt.Printf("Session for episode %d of %s created successfully\n", episodeNumber, projectName)

	return nil
}

func commandGetSessions(cfg *config, args []string) error {
	// This command Gets a list of sessions for a given project/episode
	// Takes project title, episode number and a limit of sessions to return
	if len(args) < 2 {
		return fmt.Errorf("invalid number of arguments")
	}

	limit := 0
	if len(args) == 3 {
		var err error
		limit, err = strconv.Atoi(args[2])
		if err != nil {
			return err
		}
	}

	projectName := args[0]
	episodeNumber, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}

	// Now we need the project id
	reqPrjBody := struct {
		ProjectTitle string `json:"title"`
	}{
		ProjectTitle: projectName,
	}

	prj, err := getThing(cfg, "/api/projects", reqPrjBody, db.Project{})
	if err != nil {
		return err
	}

	// Now we need the episode ID
	reqEpBody := struct {
		ProjectID     string `json:"project_id"`
		EpisodeNumber int    `json:"episode_number"`
	}{
		ProjectID:     prj.ID.String(),
		EpisodeNumber: episodeNumber,
	}

	ep, err := getThing(cfg, "/api/episodes", reqEpBody, db.Episode{})
	if err != nil {
		return err
	}

	episodeID := ep.ID

	// TODO: The actual query

	return nil
}
