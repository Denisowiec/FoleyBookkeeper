package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Denisowiec/FoleyBookkeeper/internal/db"
)

func commandCreateSession(cfg *config, args []string) error {
	// Takes project title, episode number, date of session,
	// duration of session, part worked on, activity done and a list of usernames as input
	if len(args) < 6 {
		return fmt.Errorf("invalid number of arguments")
	}

	projectName := args[0]
	episodeNumber, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}
	dateFromInput, err := time.Parse(time.DateOnly, args[2])
	if err != nil {
		return err
	}
	date := dateFromInput.Format(time.DateOnly)
	partWorkedOn := args[4]
	activityDone := args[5]

	// This converts the duration into something that should match the postgresql preferences
	durationTime, err := time.ParseDuration(args[3])
	if err != nil {
		return err
	}

	duration := int32(durationTime.Minutes())

	// We convert the usernames in the arguments into a list of IDs
	users := []string{}

	for _, username := range args[6:] {
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
		Duration     int32  `json:"duration"`
		SessionDate  string `json:"session_date"`
		EpisodeID    string `json:"episode_id"`
		PartWorkedOn string `json:"part_worked_on"`
		ActivityDone string `json:"activity_done"`
	}
	createSesReq := createSesType{
		Duration:     duration,
		SessionDate:  date,
		EpisodeID:    ep.ID.String(),
		PartWorkedOn: partWorkedOn,
		ActivityDone: activityDone,
	}

	url := fmt.Sprintf("%s/api/sessions", cfg.serverAddress)
	resp, err := sendRequest(createSesReq, "POST", url, cfg.jwt)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return processErrorResponse(resp)
	}

	ses := db.Session{}

	err = processResponse(resp, &ses)
	if err != nil {
		return err
	}

	// Now we add the users to the session
	reqUsSesBody := struct {
		UserIDs []string `json:"user_ids"`
	}{
		UserIDs: users,
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
	// Takes number of items, project title, episode number as arguments
	if len(args) < 2 {
		return fmt.Errorf("invalid number of arguments")
	}

	limit := 20
	if len(args) == 3 {
		var err error
		limit, err = strconv.Atoi(args[0])
		if err != nil {
			return err
		}
	}

	projectName := args[1]

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
	projectID := prj.ID.String()

	type listItem struct {
		db.Session `json:"session"`
		Users      []db.GetUsersForSessionRow `json:"users"`
	}
	list := []listItem{}

	// The request will be different depending on the arguments given
	var episodeNumber int
	if len(args) >= 3 {
		episodeNumber, err = strconv.Atoi(args[2])
		if err != nil {
			return err
		}

		reqBody := struct {
			ProjectID     string `json:"project_id"`
			EpisodeNumber int    `json:"episode_number"`
		}{
			ProjectID:     projectID,
			EpisodeNumber: episodeNumber,
		}

		ep, err := getThing(cfg, "/api/episodes", reqBody, db.Episode{})
		if err != nil {
			return err
		}

		episodeID := ep.ID.String()

		req2Body := struct {
			Limit     int    `json:"limit"`
			EpisodeID string `json:"episode_id"`
		}{
			Limit:     limit,
			EpisodeID: episodeID,
		}

		list, err = getThing(cfg, "/api/sessions", req2Body, list)
		if err != nil {
			return err
		}

		fmt.Printf("Sessions for project %s, episode %d:\n", prj.Title, episodeNumber)
	} else {
		reqBody := struct {
			Limit     int    `json:"limit"`
			ProjectID string `json:"project_id"`
		}{
			Limit:     limit,
			ProjectID: projectID,
		}
		list, err = getThing(cfg, "/api/sessions", reqBody, list)
		if err != nil {
			return err
		}

		fmt.Printf("Sessions for project %s:\n", prj.Title)
	}

	for _, item := range list {
		fmt.Printf("%s, %s, %s: ", item.SessionDate.Format(time.DateOnly), item.ActivityDone, item.PartWorkedOn)
		for i, u := range item.Users {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%s", u.Username)
		}
		fmt.Printf("\n")
	}

	return nil
}
