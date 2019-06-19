// Package dce defines dce API webhook package
package dce

import (
	"encoding/json"
	"fmt"
)

const (
	buildSuccess = "Success"
)

// DCE defines dce object
type DCE struct {
	Repo  string `json:"repo"`
	Image string `json:"image"`
	Name  string `json:"name"`
	Build struct {
		BuildFlowID string `json:"build_flow_id"`
		Stages      []struct {
			Name   string `json:"name"`
			Status string `json:"status"`
		} `json:"stages"`
		Status          string      `json:"status"`
		DurationSeconds int         `json:"duration_seconds"`
		Author          string      `json:"author"`
		TriggeredBy     string      `json:"triggered_by"`
		Sha             string      `json:"sha"`
		Ref             string      `json:"ref"`
		RefIsBranch     bool        `json:"ref_is_branch"`
		RefIsTag        bool        `json:"ref_is_tag"`
		Tag             string      `json:"tag"`
		Branch          interface{} `json:"branch"`
		PullRequest     string      `json:"pull_request"`
		Message         string      `json:"message"`
		StartedAt       string      `json:"started_at"`
		BuildType       string      `json:"build_type"`
	} `json:"build"`
}

// NewDCEObj init the dce callback obj
func NewDCEObj(str string) (*DCE, error) {
	d := new(DCE)
	if err := json.Unmarshal([]byte(str), d); err != nil {
		return nil, err
	}
	return d, nil
}

// GetCommitMsg fetch git commit message
func (d *DCE) GetCommitMsg() string {
	return d.Build.Message
}

// GetBuildDuration fetch build duration
func (d *DCE) GetBuildDuration() int64 {
	return int64(d.Build.DurationSeconds)
}

// GetBuildStatus gets build status
func (d *DCE) GetBuildStatus() bool {
	return d.Build.Status == buildSuccess
}

// GetSha gets commit sha
func (d *DCE) GetSha() string {
	return d.Build.Sha
}

// GetStageMap gets stages
func (d *DCE) GetStageMap() map[string]string {
	stagemap := make(map[string]string)
	for _, s := range d.Build.Stages {
		stagemap[s.Name] = s.Status
	}
	return stagemap
}

// GetRepoName fetch repo's name
func (d *DCE) GetRepoName() string {
	return d.Name
}

func (d *DCE) Fmt() string {

	text := fmt.Sprintf("[Commit Note]: `%s` ", d.GetCommitMsg())
	testStage, ok := d.GetStageMap()["test"]
	if ok {
		text = fmt.Sprintf("%s\n **Test Status**: `%v`", text, testStage)
	}
	buildStage, ok := d.GetStageMap()["build"]
	if ok {
		text = fmt.Sprintf("%s\n **Build Status**: `%v`", text, buildStage)
	}
	deployStage, ok := d.GetStageMap()["deploy"]
	if ok {
		text = fmt.Sprintf("%s\n **Deploy Status**: `%v`", text, deployStage)
	}
	text = fmt.Sprintf("%s\n **Duration**: `%v s`", text, d.GetBuildDuration())

	return text
}
