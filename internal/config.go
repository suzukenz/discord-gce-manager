package internal

var (
	projectID  string
	webhookURL string
)

// SetProjectID set GCP ProjectID to global variables.
func SetProjectID(pjID string) {
	projectID = pjID
}

// SetWebhookURL set Discord webhookURL to global variables.
func SetWebhookURL(url string) {
	webhookURL = url
}
