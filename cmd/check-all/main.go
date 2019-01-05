package main

import (
	"flag"

	"github.com/suzukenz/discord-gce-manager/internal"
)

// Variables used for command line parameters
var (
	Token     string
	ProjectID string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&ProjectID, "p", "", "GCP ProjectID")
	flag.Parse()

	internal.SetProjectID(ProjectID)
}

func main() {

}
