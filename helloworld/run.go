package helloworld

import (
	"context"
	"log"

	"github.com/pkg/errors"

	"go.temporal.io/sdk/client"
)

// Run encapsultes workflow code used in starter and helloworld_test
func Run(c client.Client, taskQ string) (string, error) {
	workflowOptions := client.StartWorkflowOptions{
		ID:        "hello_world_workflowID",
		TaskQueue: taskQ,
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, Workflow, "Temporal")
	if err != nil {
		return "", errors.Wrap(err, "Unable to execute workflow")
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	// Synchronously wait for the workflow completion.
	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		return "", errors.Wrap(err, "Unable get workflow result")
	}

	return result, nil
}
