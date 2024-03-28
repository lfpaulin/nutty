package papers

import (
	"fmt"
	"nutty/config"
	"nutty/vcf"
	"strconv"
	"strings"
)

type SampleInterest struct {
	name    string
	dr      int
	dv      int
	vaf     float64
	gt      string
	statSV  string
	uniq    int
	uniqTxt string
}

const (
	minVAFGermline float64 = 0.25 // 25%
	minVAFMosaic   float64 = 0.05 // 5%
	maxVAFMosaic   float64 = 0.25 // < 25%
	minCoverage    int     = 10   // yeast genome is small we have high coverage
	fixHetVAFMin   float64 = 0.33 // for cases in which we have germline VAF but not GT (i.e../.)
	fixAltVAFMin   float64 = 0.66 // for cases in which we have germline VAF but not GT (i.e../.)
)

func YeastSpace(params *config.UserParam) {
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
			contigsVCF[contigName] = contigSize
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
			sampleNamesInfo := strings.Join(sampleNames, ", ")
			sampleNamesHeader := strings.Join(sampleNames, "\t")
			// Here goes the parser header
			fmt.Println("## Sample names: ", sampleNamesInfo)
			fmt.Printf("#CHROM\tSTART\tEND\tSVTYPE\tSVLEN\tID\t%s\n", sampleNamesHeader)
		case strings.Contains(line, "#"):
			//
		default:
			// each entry
			ReadVCFYeastSpace(line, &contigsVCF, &sampleNames, params.PaperAnalysis)
		}
	}
}

func ReadVCFYeastSpace(VCFLineRaw string, contigs *map[string]int, sampleNames *[]string, sampleInterest string) {
	lineSplit := strings.Split(VCFLineRaw, "\t")
	VCFLineFormatted := new(vcf.VCF)
	VCFLineFormatted.Contig = lineSplit[0]
	if _, ok := (*contigs)[VCFLineFormatted.Contig]; ok {
		VCFPosInt, err := strconv.Atoi(lineSplit[1])
		if err != nil {
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
		VCFLineFormatted.Info = info
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
		var suppVecUpdated = strings.Split(VCFLineFormatted.Info["SUPP_VEC"], "")
		var statusSV string
		var samplePrint string
		var printOut string
		var printPopulation []string
		sampleSV = make(map[string]map[string]string)
		formatSplit := strings.Split(lineSplit[8], ":")
		for sidx, sample := range *sampleNames {
			sampleSVSplit := strings.Split(lineSplit[9+sidx], ":")
			sampleSV[sample] = make(map[string]string)
			for idx := range formatSplit {
				sampleSV[sample][formatSplit[idx]] = sampleSVSplit[idx]
			}
		}
		VCFLineFormatted.Samples = sampleSV
		// sample of interest compare to
		dr, err = strconv.Atoi(VCFLineFormatted.Samples[sampleInterest]["DR"])
		dv, err = strconv.Atoi(VCFLineFormatted.Samples[sampleInterest]["DV"])
		gt = VCFLineFormatted.Samples[sampleInterest]["GT"]
		vaf = float64(dv) / float64(dr+dv)
		var sampleInterestValues = SampleInterest{
			name:    sampleInterest,
			dr:      dr,
			dv:      dv,
			vaf:     0.0,
			gt:      "",
			statSV:  "",
			uniq:    0,
			uniqTxt: "",
		}
		if minCoverage > dr+dv {
			sampleInterestValues.gt = "./."
			sampleInterestValues.vaf = 0.0
			sampleInterestValues.statSV = "undefined"
		} else {
			sampleInterestValues.vaf = vaf
			if vaf >= minVAFGermline {
				statusSV = "germline"
				if gt == "./." {
					if vaf > fixHetVAFMin && vaf <= fixAltVAFMin {
						gt = "0/1"
					} else if vaf > fixAltVAFMin {
						gt = "1/1"
					} else {
						//
					}
				}
			} else if vaf < maxVAFMosaic && vaf >= minVAFMosaic {
				statusSV = "mosaic"
				if gt == "./." {
					gt = "0/0"
				}
			} else if vaf < minVAFMosaic && vaf > 0.0 {
				statusSV = "lowVAF"
			} else if vaf == 0.0 {
				statusSV = "reference"
			} else {
				statusSV = "undefined"
			}
			sampleInterestValues.statSV = statusSV
			sampleInterestValues.gt = gt
		}
		// end sample of interest
		for sid, sampleName := range *sampleNames {
			dr, err = strconv.Atoi(VCFLineFormatted.Samples[sampleName]["DR"])
			dv, err = strconv.Atoi(VCFLineFormatted.Samples[sampleName]["DV"])
			gt = VCFLineFormatted.Samples[sampleName]["GT"]
			vaf = float64(dv) / float64(dr+dv)
			if minCoverage > dr+dv {
				statusSV = "undefined"
				suppVecUpdated[sid] = "0"
				gt = "./."
				vafString = "n/a"
			} else {
				if vaf >= minVAFGermline {
					statusSV = "germline"
					if gt == "./." {
						if vaf > fixHetVAFMin && vaf <= fixAltVAFMin {
							gt = "0/1"
						} else if vaf > fixAltVAFMin {
							gt = "1/1"
						} else {
							//
						}
					}
				} else if vaf < maxVAFMosaic && vaf >= minVAFMosaic {
					statusSV = "mosaic"
					if gt == "./." {
						gt = "0/0"
					}
				} else if vaf < minVAFMosaic && vaf > 0.0 {
					statusSV = "lowVAF"
				} else if vaf == 0.0 {
					statusSV = "reference"
				} else {
					statusSV = "undefined"
				}
				vafString = fmt.Sprintf("%0.3f", vaf)
			}
			var mainSampleList = strings.Split(sampleInterestValues.name, "_")
			var mainSampleName = strings.Join(mainSampleList[:2], "_")
			if !strings.Contains(sampleName, mainSampleName) {
				if sampleInterestValues.gt == gt && sampleInterestValues.statSV == statusSV {
					sampleInterestValues.uniq += 1
				}
			}
			printOut = fmt.Sprintf("%s|%s|%d|%d|%s", gt, vafString, dr, dv, statusSV)
			printPopulation = append(printPopulation, printOut)
		}
		samplePrint = strings.Join(printPopulation, "\t")
		if sampleInterestValues.uniq == 0 {
			sampleInterestValues.uniqTxt = "*candidate*"
		}
		fmt.Printf("%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\n", VCFLineFormatted.Contig, VCFLineFormatted.Start,
			VCFLineFormatted.EndStr, VCFLineFormatted.Info["SVTYPE"], VCFLineFormatted.Info["SVLEN"],
			VCFLineFormatted.ID, samplePrint, sampleInterestValues.uniqTxt)
	}
}
