package project

import "strings"

func GetEnvForLinkedApps(links []string, containerNames map[string]string) map[string]string {
	env := map[string]string{}

	for _, link := range links {
		name := "APP_" + strings.ToUpper(link) + "_HOST"
		env[name] = containerNames[link]
	}

	return env
}
