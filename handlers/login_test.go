package handlers

import (
	"encoding/json"
	"github.com/benbjohnson/hashfs"
	kratos "github.com/ory/kratos-client-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func TestLoginWithFakeLoginFlow(t *testing.T) {
	for _, tc := range []struct {
		sourceFilePath string
		goldenFilePath string
	}{
		{
			sourceFilePath: "loginflow.json",
			goldenFilePath: "TestLoginWithFakeLoginFlow.html",
		},
		{
			sourceFilePath: "loginflowwithmessage.json",
			goldenFilePath: "TestLoginWithFakeLoginFlowWithMessage.html",
		},
	} {
		var flow *kratos.LoginFlow
		err := json.Unmarshal(getAsset(t, tc.sourceFilePath), &flow)
		ui := flow.Ui
		dataMap := map[string]any{
			"flow":        flow,
			"method":      ui.Method,
			"action":      ui.Action,
			"nodes":       ui.Nodes,
			"messages":    ui.Messages,
			"fs":          hashfs.NewFS(fstest.MapFS{}),
			"pageHeading": "Login",
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/login", strings.NewReader(""))
		err = GetTemplate(loginPage).Render("layout", w, r, dataMap)
		require.NoError(t, err, "unexpected error getting template")
		actual, err := io.ReadAll(w.Body)
		require.NoError(t, err, "unexpected error reading http recorder body")
		golden := goldenValue(t, tc.goldenFilePath, actual, *update)
		assert.Equal(t, string(golden), string(actual))
	}
}
