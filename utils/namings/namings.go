package namings

import "fmt"

// If something occurs unusual we do the mappings here
var appNameToServiceNameMap = map[string]string{}

var appNameToRepoNameMap = map[string]string{}

var appNameToDataDirNameMap = map[string]string{}

func GetServiceNameFromAppName(appName string) string {
	if value, found := appNameToServiceNameMap[appName]; found {
		return value
	}
	return fmt.Sprintf("nubeio-%s.service", appName)
}

func GetAppNameFromRepoName(repoName string) string {
	for k := range appNameToRepoNameMap {
		if appNameToRepoNameMap[k] == repoName {
			return k
		}
	}
	return repoName
}

func GetRepoNameFromAppName(appName string) string {
	if value, found := appNameToRepoNameMap[appName]; found {
		return value
	}
	return appName
}

func GetDataDirNameFromAppName(appName string) string {
	if value, found := appNameToDataDirNameMap[appName]; found {
		return value
	}
	return appName
}
