package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/benbjohnson/hashfs"
	"github.com/elielamora/kratos-selfservice-ui-go/apiclient"
)

// LoginParams configure the Login http handler
type LoginParams struct {
	// FS provides access to static files
	FS *hashfs.FS

	// FlowRedirectURL is the kratos URL to redirect the browser to,
	// when the user wishes to login, and the 'flow' query param is missing
	FlowRedirectURL string
}

// Login handler displays the login screen
func (lp LoginParams) Login(w http.ResponseWriter, r *http.Request) {

	// Start the login flow with Kratos if required
	flow := r.URL.Query().Get("flow")
	log.Printf("flow from query parameters '%v'", flow)
	if flow == "" {
		log.Printf("No flow ID found in URL, initializing login flow, redirect to %s", lp.FlowRedirectURL)
		http.Redirect(w, r, lp.FlowRedirectURL, http.StatusSeeOther)
		return
	}

	log.Print("Calling Kratos API to get self service login")
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
	defer cancel()
	loginFlowRequest := apiclient.PublicClient().FrontendAPI.GetLoginFlow(ctx)
	log.Printf("public client server: %v", apiclient.PublicClient().GetConfig().Servers[0].URL)
	res, _, err := loginFlowRequest.Id(flow).
		Cookie(r.Header.Get("Cookie")).
		Execute()
	if err != nil {
		log.Printf("Error getting self service login flow: %v, redirecting to /", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	ui := res.GetUi()
	data, _ := json.Marshal(ui)
	log.Println(string(data))
	dataMap := map[string]any{
		"flow":        flow,
		"method":      ui.Method,
		"action":      ui.Action,
		"nodes":       ui.Nodes,
		"messages":    ui.Messages,
		"fs":          lp.FS,
		"pageHeading": "Login",
	}
	if err = GetTemplate(loginPage).Render("layout", w, r, dataMap); err != nil {
		ErrorHandler(w, r, err)
	}
}
