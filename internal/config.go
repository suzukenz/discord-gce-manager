package internal

var (
	projectID string
)

// SetProjectID set GCP ProjectID to global variables.
func SetProjectID(pjID string) {
	projectID = pjID
}
