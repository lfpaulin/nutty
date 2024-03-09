package main

import (
	"sniffles2_helper_go/config"
	"sniffles2_helper_go/utils"
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
	}
}
