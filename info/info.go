package info

import (
	"fmt"
)

var (
	Version = ""
	Commit  = ""
	Date    = ""
	BuiltBy = ""
)

func BuildVersion(version, commit, date, builtBy string) string {
	coalesce := func(input, fallback string) string {
		if input != "" {
			return input
		}
		return fallback
	}
	result := coalesce(version, Version)
	result = fmt.Sprintf("%s\ncommit: %s", result, coalesce(commit, Commit))
	result = fmt.Sprintf("%s\nbuilt at: %s", result, coalesce(date, Date))
	result = fmt.Sprintf("%s\nbuilt by: %s", result, coalesce(builtBy, BuiltBy))
	return result
}

func PrintLogo() {
	fmt.Println("                                                                             ")
	fmt.Println(" ,-----.                       ,--.      ,--------.,--.  ,--.        ,--.    ")
	fmt.Println("'  .-.  ' ,---.  ,---. ,--,--, |  | ,---.'--.  .--'|  '--'  |,--.,--.|  |-.  ")
	fmt.Println("|  | |  || .-. || .-. :|      \\|  || .-. |  |  |   |  .--.  ||  ||  || .-. ' ")
	fmt.Println("'  '-'  '| '-' '\\   --.|  ||  ||  |' '-' '  |  |   |  |  |  |'  ''  '| `-' | ")
	fmt.Println(" `-----' |  |-'  `----'`--''--'`--' `---'   `--'   `--'  `--' `----'  `---'   ")
	fmt.Println("         `--'               form https://github.com/OpenIoTHub/gateway-go/v2     ")
}
