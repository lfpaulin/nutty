package main

import (
	"fmt"
	"nutty/config"
	"nutty/papers"
	"nutty/utils"
	"os"
)

func main() {
	userParsedParams := config.GetParams()
	switch userParsedParams.SubCMD {
	case "sv":
		utils.ParseSV(&userParsedParams)
	case "pop":
		utils.ParsePop(&userParsedParams)
	case "cancer":
		utils.ParseCancer(&userParsedParams)
	case "paper":
		papers.PaperAnalysis(&userParsedParams)
	default:
		fmt.Printf("Unknown subcommand: %s\n", userParsedParams.SubCMD)
		os.Exit(1)
	}
}
