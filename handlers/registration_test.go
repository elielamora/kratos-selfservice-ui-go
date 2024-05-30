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

func TestRegistrationWithFlowStub(t *testing.T) {
	var flow *kratos.RegistrationFlow
	err := json.Unmarshal(getAsset(t, "registrationflow.json"), &flow)
	ui := flow.Ui
	dataMap := map[string]interface{}{
		"method":      ui.Method,
		"action":      ui.Action,
		"nodes":       ui.Nodes,
		"messages":    ui.Messages,
		"flow":        flow,
		"fs":          hashfs.NewFS(fstest.MapFS{}),
		"pageHeading": "Registration",
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/register", strings.NewReader(""))
	err = GetTemplate(registrationPage).Render("layout", w, r, dataMap)
	require.NoError(t, err, "unexpected error getting template")
	actual, err := io.ReadAll(w.Body)
	require.NoError(t, err, "unexpected error reading http recorder body")
	golden := goldenValue(t, "TestRegistrationWithFlowStub.html", actual, *update)
	assert.Equal(t, string(golden), string(actual))
}
