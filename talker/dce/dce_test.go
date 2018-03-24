package dce

import (
	"encoding/json"
	"testing"
)

var (
	testJSON = `
  {
    "repo":"daocloud/api",
    "image":"daocloud.io/daocloud/api:master-init",
    "name":"api",
    "build": {
      "build_flow_id":"8d7622ea-9323-4489-8c8e-fc4bed448961",
      "stages": [
        {
          "name": "test",
          "status": "Success"
        },
        {
          "name": "build",
          "status": "Success"
        },
        {
          "name": "deploy",
          "status": "Success"
        }
      ],
      "status":"Success",
      "duration_seconds":180,
      "author":"DaoCloud",
      "triggered_by":"tag",
      "sha":"a7c35d9dc7e93788ce81befbadeb0108de495e5e",
      "ref": "v1.0",
      "ref_is_branch": false,
      "ref_is_tag": true,
      "tag":"v1.0",
      "branch":null,
      "pull_request":"",
      "message":"init build ",
      "started_at":"2015-01-01T08:20:00+00:00",
      "build_type":"ci_build"
    }
  }
  `
)

func TestMarshalDCEStr(t *testing.T) {
	d := new(DCE)
	if err := json.Unmarshal([]byte(testJSON), d); err != nil {
		t.Errorf("unMarshal err:%s", err.Error())
	}
	if d.Name != "api" {
		t.Errorf("not expected %s", d.Name)
	}
}
