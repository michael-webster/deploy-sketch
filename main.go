// call the circleci api to trigger a pipeline
// fill in your pipeline definition id from project settings
// put an API token in an env variable called CIRCLE_TOKEN
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"

	"net/http"
)

const pipelineDefinitionId = "<put your pipeline definition here>"

type CircleCICLient struct {
	client http.Client
	token  string
}

func (c *CircleCICLient) TriggerPipeline(ctx context.Context, pipelineDefinition string) (map[string]any, error) {
	payload := map[string]any{
		"definition_id": pipelineDefinition,
		"config": map[string]any{
			"branch": "main",
		},
		"checkout": map[string]any{
			"branch": "main",
		},
	}
	marshalled, _ := json.Marshal(payload)
	bodyReader := bytes.NewReader(marshalled)
	resp := make(map[string]any)
	req, _ := http.NewRequest("POST", "https://circleci.com/api/v2/project/github/michael-webster/deploy-sketch/pipeline/run", bodyReader)
	req.Header.Set("Circle-Token", c.token)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	json.NewDecoder(res.Body).Decode(&resp)
	return resp, nil
}

func getPipeline(pipeline string) string {
	commandMap := map[string]string{
		"deploy": pipelineDefinitionId,
	}
	return commandMap[pipeline]
}

func main() {
	println("Let's get situated....")
	ctx := context.Background()
	client := http.Client{}

	circle := CircleCICLient{
		client: client,
		token:  os.Getenv("CIRCLE_TOKEN"),
	}

	pipelineDefinition := getPipeline("deploy")

	resp, err := circle.TriggerPipeline(ctx, pipelineDefinition)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(resp)
	}
}
