package papers

import (
	"fmt"
	"nutty/config"
	"nutty/vcf"
	"strconv"
	"strings"
)

func GIAB(params *config.UserParam) {
	VCFReader := vcf.ReaderMaker(params.VCF)
	if params.VCF != "-" && params.VCF != "stdin" {
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
			if contigSize > minContigLen {
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
			fmt.Println("#CHROM\tSTART\tSVTYPE\tSVLEN\tID\tTUMOR\tCONTROL")
		case strings.Contains(line, "#"):
			//
		default:
			// each entry
			switch params.PaperAnalysis {
			case "hg008":
				ReadVCFHG008(line, &contigsVCF, &sampleNames, params)
			default:
				fmt.Printf("Paper Analysis %s not known.. exiting", params.PaperAnalysis)
			}
		}
	}
}

func SuppIsTumor(SuppVec string) bool {
	/*
		Tumor/Control sample that has six samples three control (111000) and three tumor (000111):
		We ought to get tumor-only/somatic SVs
	*/
	// indices
	tumorIndex := 3
	isTumorOnly := false
	// check tumor first
	for _, supp := range strings.Split(SuppVec, "")[tumorIndex:] {
		if supp == "1" {
			isTumorOnly = true
		}
	}
	// look in controls next
	for _, supp := range strings.Split(SuppVec, "")[:tumorIndex] {
		if supp == "1" {
			isTumorOnly = false
			break
		}
	}
	return isTumorOnly
}

func ReadVCFHG008(VCFLineRaw string, contigs *map[string]int, sampleNames *[]string, params *config.UserParam) {
	minReadsAlt := 3
	minReadsTotal := 10
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
		1. support vector (SUPP_VEC) 000111 (and any combination within the tumors)
		2. no reads in the control
		*/
		tumorSamples := 3
		printTumor := make([]string, 0)
		printOut := ""
		printPass := 0
		if SuppIsTumor(VCFLineFormatted.Info["SUPP_VEC"]) {
			for _, sampleName := range (*sampleNames)[tumorSamples:] {
				gt = VCFLineFormatted.Samples[sampleName]["GT"]
				dr, err = strconv.Atoi(VCFLineFormatted.Samples[sampleName]["DR"])
				dv, err = strconv.Atoi(VCFLineFormatted.Samples[sampleName]["DV"])
				vaf = float64(dv) / float64(dr+dv)
				if dv >= minReadsAlt && gt == "1/1" && dr+dv > minReadsTotal {
					printOut = fmt.Sprintf("%s|%s|%0.3f:%d:%d", sampleName, gt, vaf, dr, dv)
					printTumor = append(printTumor, printOut)
					printPass += 1
				} else {
					printTumor = append(printTumor, sampleName)
				}
			}
			if printPass == 3 {
				samplePrint := strings.Join(printTumor, ", ")
				fmt.Printf("%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\n", VCFLineFormatted.Contig, VCFLineFormatted.Start,
					VCFLineFormatted.EndStr, VCFLineFormatted.Info["SVTYPE"], VCFLineFormatted.Info["SVLEN"],
					VCFLineFormatted.Info["SUPP_VEC"], VCFLineFormatted.ID, samplePrint)
			}
		}
	}
}
