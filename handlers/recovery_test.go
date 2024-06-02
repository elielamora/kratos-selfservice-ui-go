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

func TestRecoveryWithFlowStub(t *testing.T) {
	var flow *kratos.RecoveryFlow
	err := json.Unmarshal(getAsset(t, "recoveryflow.json"), &flow)
	ui := flow.Ui
	dataMap := map[string]interface{}{
		"method":      ui.Method,
		"action":      ui.Action,
		"nodes":       ui.Nodes,
		"messages":    ui.Messages,
		"flow":        flow,
		"fs":          hashfs.NewFS(fstest.MapFS{}),
		"pageHeading": "Recovery",
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/recovery", strings.NewReader(""))
	err = GetTemplate(recoveryPage).Render("layout", w, r, dataMap)
	require.NoError(t, err, "unexpected error getting template")
	actual, err := io.ReadAll(w.Body)
	require.NoError(t, err, "unexpected error reading http recorder body")
	golden := goldenValue(t, "TestRecoveryWithFlowStub.html", actual, true)
	assert.Equal(t, string(golden), string(actual))
}
