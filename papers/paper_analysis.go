package papers

import (
	"fmt"
	"nutty/config"
	"os"
)

const availablePapers string = "cancer_t2t => PMID: \n"
const analysisT2T string = "colo829, pog"

func PaperAnalysis(params *config.UserParam) {
	switch params.PaperID {
	case "cancer_t2t":
		switch params.PaperAnalysis {
		case "colo829":
			CancerT2T(params)
		case "pog":
			CancerT2T(params)
		default:
			fmt.Printf("Paper Analysis %s not known\n.Available analysis for %s are:\n%s", params.PaperAnalysis,
				params.PaperID, analysisT2T)
			os.Exit(1)
		}
	case "yeast":
		YeastSpace(params)
	default:
		fmt.Printf("Paper ID %s not known\n.Available paprs are:\n%s", params.PaperID, availablePapers)
	}
}
