package handlers

import (
	"context"
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/benbjohnson/hashfs"
	"github.com/elielamora/kratos-selfservice-ui-go/apiclient"
)

// RegistrationParams configure the Registration http handler
type RegistrationParams struct {
	// FS provides access to static files
	FS *hashfs.FS

	// FlowRedirectURL is the kratos URL to redirect the browser to,
	// when the user wishes to register, and the 'flow' query param is missing
	FlowRedirectURL string
}

// Registration directs the user to a page where they can sign-up or
// register to use the site
func (rp RegistrationParams) Registration(w http.ResponseWriter, r *http.Request) {

	// Start the login flow with Kratos if required
	flow := r.URL.Query().Get("flow")
	if flow == "" {
		log.Printf("No flow ID found in URL, initializing login flow, redirect to %s", rp.FlowRedirectURL)
		http.Redirect(w, r, rp.FlowRedirectURL, http.StatusSeeOther)
		return
	}

	log.Print("Calling Kratos API to get self service registration")
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
	defer cancel()
	req := apiclient.PublicClient().FrontendAPI.GetRegistrationFlow(ctx)
	res, _, err := req.Id(flow).
		Cookie(r.Header.Get("Cookie")).
		Execute()
	if err != nil {
		log.Printf("Error getting self service registration flow %v, redirecting to /", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	ui := res.GetUi()
	data, _ := json.Marshal(ui)
	log.Println(string(data))
	dataMap := map[string]interface{}{
		"method":      ui.Method,
		"action":      ui.Action,
		"nodes":       ui.Nodes,
		"messages":    ui.Messages,
		"flow":        flow,
		"fs":          rp.FS,
		"pageHeading": "Registration",
	}

	if err = GetTemplate(registrationPage).Render("layout", w, r, dataMap); err != nil {
		ErrorHandler(w, r, err)
	}
}
