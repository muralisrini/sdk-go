package helloworld

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"

	tsenv "github.com/temporalio/samples-go/helloworld/testserver"
)

func Test_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Mock activity implementation
	env.OnActivity(Activity, mock.Anything, "Temporal").Return("Hello Temporal!", nil)

	env.ExecuteWorkflow(Workflow, "Temporal")

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	var result string
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, "Hello Temporal!", result)
}

func Test_Activity(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()
	env.RegisterActivity(Activity)

	val, err := env.ExecuteActivity(Activity, "World")
	require.NoError(t, err)

	var res string
	require.NoError(t, val.Get(&res))
	require.Equal(t, "Hello World!", res)
}

func Test_Using_DevServer(t *testing.T) {
	testServerEnv := tsenv.New("", "hello-world")
	err := testServerEnv.Start()
	require.NoError(t, err)

	w, err := testServerEnv.Worker()
	require.NoError(t, err)
	require.NotNil(t, w)

	w.RegisterWorkflow(Workflow)
	w.RegisterActivity(Activity)

	c, err := testServerEnv.Client()
	require.NoError(t, err)
	require.NotNil(t, c)

	result, err := Run(c, testServerEnv.TaskQ)
	require.NoError(t, err)
	require.Equal(t, "Hello Temporal!", result)

	err = testServerEnv.Stop()
	require.NoError(t, err)
}
