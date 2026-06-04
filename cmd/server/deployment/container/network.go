package container

import "fmt"

func ProjectNetworkName(projectId string) string {
	return fmt.Sprintf("pushnpray-%s", projectId)
}
