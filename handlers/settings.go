package handlers

import (
	"context"
	_ "embed"
	"log"
	"net/http"
	"time"

	"github.com/benbjohnson/hashfs"
	"github.com/elielamora/kratos-selfservice-ui-go/apiclient"
)

// SettingsParams configure the Settings http handler
type SettingsParams struct {
	// FS provides access to static files
	FS *hashfs.FS

	// FlowRedirectURL is the kratos URL to redirect the browser to,
	// when the user wishes to edit their settings, and the 'flow' query param is missing
	FlowRedirectURL string
}

// Settings handler displays the Settings screen that allows the user to change their details
func (lp SettingsParams) Settings(w http.ResponseWriter, r *http.Request) {

	// Start the Settings flow with Kratos if required
	flow := r.URL.Query().Get("flow")
	log.Printf("flow from query parameters '%v'", flow)
	if flow == "" {
		log.Printf("No flow ID found in URL, initializing Settings flow, redirect to %s", lp.FlowRedirectURL)
		http.Redirect(w, r, lp.FlowRedirectURL, http.StatusSeeOther)
		return
	}

	log.Print("Calling Kratos API to get self service seettings")
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
	defer cancel()
	req := apiclient.PublicClient().FrontendAPI.GetSettingsFlow(ctx)
	res, _, err := req.Id(flow).
		Cookie(r.Header.Get("Cookie")).
		Execute()
	if err != nil {
		log.Printf("Error getting self service settings flow %v, redirecting to /", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	ui := res.GetUi()

	dataMap := map[string]interface{}{
		"flow":     flow,
		"method":   ui.Method,
		"action":   ui.Action,
		"nodes":    ui.GetNodes()[0],
		"messages": ui.Messages,
		//"password":    res.GetPayload().Methods["password"].Config,
		//"profile":     res.GetPayload().Methods["profile"].Config,
		"fs":          lp.FS,
		"pageHeading": "Update Profile",
	}
	if err = GetTemplate(settingsPage).Render("layout", w, r, dataMap); err != nil {
		ErrorHandler(w, r, err)
	}
}
