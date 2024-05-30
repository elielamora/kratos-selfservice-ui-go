package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/benbjohnson/hashfs"
	"github.com/elielamora/kratos-selfservice-ui-go/apiclient"
)

// RecoveryParams configure the Recovery http handler
type RecoveryParams struct {
	// FS provides access to static files
	FS *hashfs.FS

	// FlowRedirectURL is the kratos URL to redirect the browser to,
	// when the user wishes to start recovery, and the 'flow' query param is missing
	FlowRedirectURL string
}

// Recovery handler displays the recovery screen, which allows the user to enter
// and email address, the email contains a link to authenticate the user
func (rp RecoveryParams) Recovery(w http.ResponseWriter, r *http.Request) {

	// Start the recovery flow with Kratos if required
	flow := r.URL.Query().Get("flow")
	if flow == "" {
		log.Printf("No flow ID found in URL, initializing login flow, redirect to %s", rp.FlowRedirectURL)
		http.Redirect(w, r, rp.FlowRedirectURL, http.StatusTemporaryRedirect)
		return
	}

	log.Printf("Calling Kratos API to get self service recovery")
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
	defer cancel()
	req := apiclient.PublicClient().FrontendAPI.GetRecoveryFlow(ctx)
	res, _, err := req.Id(flow).
		Cookie(r.Header.Get("Cookie")).
		Execute()
	if err != nil {
		log.Printf("Error getting self service recovery flow %v, redirecting to '/'", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	log.Printf("Recovery state: %v", res.GetState())
	ui := res.GetUi()
	dataMap := map[string]interface{}{
		"flow":        flow,
		"method":      ui.Method,
		"action":      ui.Action,
		"nodes":       ui.Nodes,
		"messages":    ui.Messages,
		"state":       res.State,
		"fs":          rp.FS,
		"pageHeading": "Recover your account",
	}
	if err = GetTemplate(recoveryPage).Render("layout", w, r, dataMap); err != nil {
		ErrorHandler(w, r, err)
	}
}
