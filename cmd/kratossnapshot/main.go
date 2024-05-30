package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/elielamora/kratos-selfservice-ui-go/apiclient"
	"log"
	"net/url"
	"os"
)

var (
	adminURL string
	flowType string
)

func main() {
	log.SetOutput(os.Stderr)
	flag.StringVar(
		&adminURL,
		"kratos-admin-url",
		"",
		"the url of the admin api",
	)
	flag.StringVar(
		&flowType,
		"kratos-flow-type",
		"",
		"the type of flow e.g. login",
	)
	flag.Parse()
	u, err := url.Parse(adminURL)
	if err != nil {
		log.Fatal("could not parse kratos-admin-url", err)
	}
	apiclient.InitAdminClient(*u)
	client := apiclient.AdminClient()
	switch flowType {
	case "login":
		req := client.FrontendAPI.CreateBrowserLoginFlow(context.Background())
		flow, _, err := req.Execute()
		if err != nil {
			log.Fatal("could not execute create browser login flow", err)
		}
		flowAsJSON, err := json.Marshal(flow)
		if err != nil {
			log.Fatal("could not marsha flow to JSON", err)
		}
		fmt.Print(string(flowAsJSON))
	case "registration":
		req := client.FrontendAPI.CreateBrowserRegistrationFlow(context.Background())
		flow, _, err := req.Execute()
		if err != nil {
			log.Fatal("could not execute create browser registration flow", err)
		}
		flowAsJSON, err := json.Marshal(flow)
		if err != nil {
			log.Fatal("could not marsha flow to JSON", err)
		}
		fmt.Print(string(flowAsJSON))
	case "recovery":
		req := client.FrontendAPI.CreateBrowserRecoveryFlow(context.Background())
		req.ReturnTo("/recovered")
		flow, _, err := req.Execute()
		if err != nil {
			log.Fatal("could not execute create browser recovery flow", err)
		}
		flowAsJSON, err := json.Marshal(flow)
		if err != nil {
			log.Fatal("could not marsha flow to JSON", err)
		}
		fmt.Print(string(flowAsJSON))
	default:
		log.Fatal("flow type must be login or registration")
	}
}
