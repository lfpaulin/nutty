package papers

import (
	"fmt"
	"nutty/config"
	"nutty/vcf"
	"strconv"
	"strings"
)

var (
	infoVCF     []string
	formatVCF   []string
	sampleNames []string
	info        map[string]string
	sampleSV    map[string]map[string]string
)
var contigsVCF = make(map[string]int)

const minReadSomatic int = 10
const minVAFCOLO829 float64 = 0.1 // 10%
const minVAFPOG float64 = 0.222   // 10%

func CancerT2T(params *config.UserParam) {
	VCFReader := vcf.VCFReaderMaker(params.VCF)
	if params.VCF != "-" && params.VCF != "stdin"{
        defer func(VCFReader *vcf.FileScanner) {
            err := VCFReader.Close()
            if err != nil {
                panic(err)
            }
        }(VCFReader)
    }
	// header metadata needed
	for VCFReader.Scan() {
		line := strings.TrimSpace(VCFReader.Text())
		switch {
		case strings.Contains(line, "##") && strings.Contains(line, "contig"):
			contigMatch := vcf.HeaderRegex(line, "contig")
			contigName := contigMatch[1]
			contigSize, err := strconv.Atoi(contigMatch[2])
			if err != nil {
				fmt.Println("[FAILED] strconv.Atoi(contigMatch[2]")
				panic(err)
			}
			if contigSize > params.MinContigLen {
				contigsVCF[contigName] = contigSize
			}
		case strings.Contains(line, "##") && strings.Contains(line, "INFO"):
			infoMatch := vcf.HeaderRegex(line, "info")
			infoVCF = append(infoVCF, infoMatch[1])
		case strings.Contains(line, "##") && strings.Contains(line, "FORMAT"):
			formatMatch := vcf.HeaderRegex(line, "format")
			formatVCF = append(formatVCF, formatMatch[1])
		case strings.Contains(line, "#CHROM"):
			lineSplit := strings.Split(line, "\t")
			for _, sample := range lineSplit[9:] {
				sampleNames = append(sampleNames, sample)
			}
			sampleNamesHeader := strings.Join(sampleNames, ",")
			// Here goes the parser header
			fmt.Println("##Sample names: ", sampleNamesHeader)
			fmt.Println("#CHROM\tSTART\tEND\tSVTYPE\tSVLEN\tID\tSUPPORT")
		case strings.Contains(line, "#"):
			//
		default:
			// each entry
			switch params.PaperAnalysis {
			case "colo829":
				ReadVCFCOLO829(line, &contigsVCF, &sampleNames)
			case "pog":
				ReadVCFPOG(line, &contigsVCF, &sampleNames)
			default:
				fmt.Printf("Paper Analysis %s not known.. exiting", params.PaperAnalysis)
			}
		}
	}
}

func ReadVCFCOLO829(VCFLineRaw string, contigs *map[string]int, sampleNames *[]string) {
	var (
		dr  int
		dv  int
		vaf float64
	)
	lineSplit := strings.Split(VCFLineRaw, "\t")
	VCFLineFormatted := new(vcf.VCF)
	VCFLineFormatted.Contig = lineSplit[0]
	if _, ok := (*contigs)[VCFLineFormatted.Contig]; ok {
		VCFPosInt, err := strconv.Atoi(lineSplit[1])
		if err != nil {
			fmt.Println("[FAILED] strconv.Atoi(lineSplit[1]")
			panic(err)
		}
		VCFLineFormatted.Pos = VCFPosInt
		VCFLineFormatted.Start = VCFPosInt
		VCFLineFormatted.ID = lineSplit[2]
		VCFLineFormatted.Ref = ""
		VCFLineFormatted.Alt = ""
		VCFLineFormatted.Quality = lineSplit[5]
		VCFLineFormatted.Filter = lineSplit[6]
		// split each key=value pair or flag
		info = make(map[string]string)
		for _, infoElem := range strings.Split(lineSplit[7], ";") {
			if strings.Contains(infoElem, "=") {
				infoKeyVal := strings.Split(infoElem, "=")
				info[infoKeyVal[0]] = infoKeyVal[1]
			} else {
				info[infoElem] = "flag"
			}
		}
		if info["SVTYPE"] == "BND" {
			VCFLineFormatted.End = VCFLineFormatted.Start + 1
			VCFLineFormatted.EndStr = lineSplit[4] // Alt
			info["SVLEN"] = "1"
		} else {
			end, err := strconv.Atoi(info["END"])
			if err != nil {
				fmt.Println("[FAILED] strconv.Atoi(info[\"END\"])")
				panic(err)
			}
			VCFLineFormatted.End = end
			VCFLineFormatted.EndStr = info["END"]
		}
		VCFLineFormatted.Info = info
		// we expect eight samples: 4 cancer, 4 controls
		formatSplit := strings.Split(lineSplit[8], ":")
		sampleSV = make(map[string]map[string]string)
		for sidx, sample := range *sampleNames {
			sampleSVSplit := strings.Split(lineSplit[9+sidx], ":")
			sampleSV[sample] = make(map[string]string)
			for idx := range formatSplit {
				sampleSV[sample][formatSplit[idx]] = sampleSVSplit[idx]
			}
		}
		VCFLineFormatted.Samples = sampleSV
		/* Filters from the paper
		1. support vector (SUPP_VEC) 11110000 cancer x4 -- control x4 (sample wide)
		2. AF >= 10% (per sample)
		3. DR+DV >= 10 (per sample)
		*/
		cancerSamples := 4
		filtersPassed := 0
		printCancer := make([]string, 0)
		printOut := ""
		if VCFLineFormatted.Info["SUPP_VEC"] == "11110000" {
			for _, sampleName := range (*sampleNames)[:4] {
				dr, err = strconv.Atoi(VCFLineFormatted.Samples[sampleName]["DR"])
				dv, err = strconv.Atoi(VCFLineFormatted.Samples[sampleName]["DV"])
				vaf = float64(dv) / float64(dr+dv)
				if dr+dv >= minReadSomatic && vaf >= minVAFCOLO829 {
					printOut = fmt.Sprintf("%s|%s|%0.3f:%d:%d", sampleName,
						VCFLineFormatted.Samples[sampleName]["GT"], vaf, dr, dv)
					printCancer = append(printCancer, printOut)
					filtersPassed += 1
				}
			}
			if filtersPassed == cancerSamples {
				samplePrint := strings.Join(printCancer, ", ")
				fmt.Printf("%s\t%d\t%s\t%s\t%s\t%s\t%s\n", VCFLineFormatted.Contig,
					VCFLineFormatted.Start, VCFLineFormatted.EndStr, VCFLineFormatted.Info["SVTYPE"],
					VCFLineFormatted.Info["SVLEN"], VCFLineFormatted.ID, samplePrint)
			}
		}
	}
}

func ReadVCFPOG(VCFLineRaw string, contigs *map[string]int, sampleNames *[]string) {
	var (
		dr  int
		dv  int
		vaf float64
		cdv int
		cid string
	)
	lineSplit := strings.Split(VCFLineRaw, "\t")
	VCFLineFormatted := new(vcf.VCF)
	VCFLineFormatted.Contig = lineSplit[0]
	if _, ok := (*contigs)[VCFLineFormatted.Contig]; ok {
		VCFPosInt, err := strconv.Atoi(lineSplit[1])
		if err != nil {
			fmt.Println("[FAILED] strconv.Atoi(lineSplit[1]")
			panic(err)
		}
		VCFLineFormatted.Pos = VCFPosInt
		VCFLineFormatted.Start = VCFPosInt
		VCFLineFormatted.ID = lineSplit[2]
		VCFLineFormatted.Ref = ""
		VCFLineFormatted.Alt = ""
		VCFLineFormatted.Quality = lineSplit[5]
		VCFLineFormatted.Filter = lineSplit[6]
		// split each key=value pair or flag
		info = make(map[string]string)
		for _, infoElem := range strings.Split(lineSplit[7], ";") {
			if strings.Contains(infoElem, "=") {
				infoKeyVal := strings.Split(infoElem, "=")
				info[infoKeyVal[0]] = infoKeyVal[1]
			} else {
				info[infoElem] = "flag"
			}
		}
		if info["SVTYPE"] == "BND" {
			VCFLineFormatted.End = VCFLineFormatted.Start + 1
			VCFLineFormatted.EndStr = lineSplit[4] // Alt
			info["SVLEN"] = "1"
		} else {
			end, err := strconv.Atoi(info["END"])
			if err != nil {
				fmt.Println("[FAILED] strconv.Atoi(info[\"END\"])")
				panic(err)
			}
			VCFLineFormatted.End = end
			VCFLineFormatted.EndStr = info["END"]
		}
		VCFLineFormatted.Info = info
		// we expect 2 samples: 1 cancer, 1 controls, in that order
		formatSplit := strings.Split(lineSplit[8], ":")
		sampleSV = make(map[string]map[string]string)
		for sidx, sample := range *sampleNames {
			sampleSVSplit := strings.Split(lineSplit[9+sidx], ":")
			sampleSV[sample] = make(map[string]string)
			for idx := range formatSplit {
				sampleSV[sample][formatSplit[idx]] = sampleSVSplit[idx]
			}
		}
		VCFLineFormatted.Samples = sampleSV
		/* Filters from the paper
		1. support vector (SUPP_VEC) 10 cancer + control
		2. AF >= 10% (per sample)
		3. DR+DV >= 10 (per sample)
		*/
		cancerSample := (*sampleNames)[0]
		controlSample := (*sampleNames)[1]
		if VCFLineFormatted.Info["SUPP_VEC"] == "10" {
			dr, err = strconv.Atoi(VCFLineFormatted.Samples[cancerSample]["DR"])
			dv, err = strconv.Atoi(VCFLineFormatted.Samples[cancerSample]["DV"])
			vaf = float64(dv) / float64(dr+dv)
			cdv, err = strconv.Atoi(VCFLineFormatted.Samples[controlSample]["DV"])
			cid = VCFLineFormatted.Samples[controlSample]["ID"]
			if dr+dv > minReadSomatic && vaf >= minVAFPOG {
				printOut := fmt.Sprintf("%s|%s|%0.3f:%d:%d", cancerSample,
					VCFLineFormatted.Samples[cancerSample]["GT"], vaf, dr, dv)
				if cdv == 0 && cid == "NULL" {
					fmt.Printf("%s\t%d\t%s\t%s\t%s\t%s\t%s\n", VCFLineFormatted.Contig,
						VCFLineFormatted.Start, VCFLineFormatted.EndStr, VCFLineFormatted.Info["SVTYPE"],
						VCFLineFormatted.Info["SVLEN"], VCFLineFormatted.ID, printOut)
				} else {
					fmt.Printf("*%s\t%d\t%s\t%s\t%s\t%s\t%s\n", VCFLineFormatted.Contig,
						VCFLineFormatted.Start, VCFLineFormatted.EndStr, VCFLineFormatted.Info["SVTYPE"],
						VCFLineFormatted.Info["SVLEN"], VCFLineFormatted.ID, printOut)
				}
			}
		}
	}
}
