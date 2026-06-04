package container

import "fmt"

func ProjectNetworkName(projectID string) string {
	return fmt.Sprintf("pushnpray-%s", projectID)
}
