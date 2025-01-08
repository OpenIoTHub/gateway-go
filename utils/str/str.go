package str

import (
	"fmt"
	version2 "github.com/OpenIoTHub/gateway-go/info"
)

func BuildVersion(version, commit, date, builtBy string) string {
	var result = ""
	//TODO ..
	if version != "" {
		result = version
	} else {
		result = version2.Version
	}
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	} else {
		result = fmt.Sprintf("%s\ncommit: %s", result, version2.Commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	} else {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, version2.Date)
	}
	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	} else {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, version2.BuiltBy)
	}
	return result
}
