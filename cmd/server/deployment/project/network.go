package project

import "fmt"

func NetworkName(projectID string) string {
	return fmt.Sprintf("pushnpray-%s", projectID)
}
