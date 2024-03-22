package utils

import (
	"fmt"
	"nutty/config"
	"nutty/vcf"
	"strconv"
	"strings"
)

func ParsePop(params *config.UserParam) {
	VCFReader := vcf.VCFReaderMaker(params.VCF)
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
				panic(err)
			}
			if contigSize > params.MinContigLen {
				contigsVCF[contigName] = contigSize
				if params.OutputVCF {
					fmt.Println(line)
				}
			}
		case strings.Contains(line, "##") && strings.Contains(line, "INFO"):
			infoMatch := vcf.HeaderRegex(line, "info")
			infoVCF = append(infoVCF, infoMatch[1])
			if params.OutputVCF {
				fmt.Println(line)
			}
		case strings.Contains(line, "##") && strings.Contains(line, "FORMAT"):
			formatMatch := vcf.HeaderRegex(line, "format")
			formatVCF = append(formatVCF, formatMatch[1])
			if params.OutputVCF {
				fmt.Println(line)
			}
		case strings.Contains(line, "#CHROM"):
			lineSplit := strings.Split(line, "\t")
			for _, sample := range lineSplit[9:] {
				sampleNames = append(sampleNames, sample)
			}
			sampleNamesInfo := strings.Join(sampleNames, ", ")
			sampleNamesHeader := strings.Join(sampleNames, "\t")
			if params.OutputVCF {
				fmt.Println(line)
			}
			// Here goes the parser header
			if !params.AsBED && !params.OutputVCF {
				fmt.Println("## Sample names: ", sampleNamesInfo)
				if !params.Uniq {
					fmt.Printf("#CHROM\tSTART\tEND\tSVTYPE\tSVLEN\tID\t%s\n", sampleNamesHeader)
				} else {
					fmt.Println("#CHROM\tSTART\tEND\tSVTYPE\tSVLEN\tID")
				}
			}
		case strings.Contains(line, "#"):
			if params.OutputVCF {
				fmt.Println(line)
			}
		default:
			// each entry
			ReadVCFPopEntry(line, &contigsVCF, &sampleNames, params)
		}
	}
}

func ReadVCFPopEntry(VCFLineRaw string, contigs *map[string]int, sampleNames *[]string, userParams *config.UserParam) {
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
		var suppVecSum = 0
		var suppVecUniq = 0
		var suppVecArray = strings.Split(VCFLineFormatted.Info["SUPP_VEC"], "")
		for suppVecIdx, suppVecElem := range suppVecArray {
			suppVecVal, err := strconv.Atoi(suppVecElem)
			if err != nil {
				panic(err)
			}
			suppVecSum += suppVecVal
			// Only works for unique SVs
			if suppVecVal == 1 {
				suppVecUniq = suppVecIdx
			}
		}
		var statusSV string
		var samplePrint string
		var printOut string
		var printPopulation []string
		var sampleNameUniq string
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
		if userParams.Uniq && suppVecSum == 1 {
			sampleNameUniq = (*sampleNames)[suppVecUniq]
		}
		for sid, sampleName := range *sampleNames {
			if (userParams.Uniq && suppVecSum == 1 && sampleName == sampleNameUniq) || !userParams.Uniq {
				dr, err = strconv.Atoi(VCFLineFormatted.Samples[sampleName]["DR"])
				dv, err = strconv.Atoi(VCFLineFormatted.Samples[sampleName]["DV"])
				gt = VCFLineFormatted.Samples[sampleName]["GT"]
				vaf = float64(dv) / float64(dr+dv)
				if minCoverage > dr+dv {
					statusSV = "undefined"
					gt = "./."
					vafString = "n/a"
					suppVecArray[sid] = "0"
				} else {
					if vaf >= userParams.MinVAFMosaic {
						statusSV = "germline"
						suppVecArray[sid] = "1"
						if gt == "./." && userParams.FixGT {
							if vaf > fixHetVAFMin && vaf <= fixAltVAFMin {
								gt = "0/1"
							} else if vaf > fixAltVAFMin {
								gt = "1/1"
							} else {
								//
							}
						}
					} else if vaf < userParams.MaxVAFMosaic && vaf >= userParams.MinVAFMosaic {
						statusSV = "mosaic"
						suppVecArray[sid] = "1"
						if gt == "./." && userParams.FixGT {
							gt = "0/0"
						}
					} else if vaf < userParams.MinVAFMosaic && vaf > 0.0 {
						statusSV = "lowVAF"
						suppVecArray[sid] = "0"
					} else if vaf == 0.0 {
						statusSV = "reference"
						suppVecArray[sid] = "0"
					} else {
						statusSV = "undefined"
						suppVecArray[sid] = "0"
					}
					vafString = fmt.Sprintf("%0.3f", vaf)
				}
				if userParams.OnlyGT {
					if statusSV == "undefined" {
						printOut = fmt.Sprintf("%s*", gt)
					} else {
						printOut = fmt.Sprintf("%s", gt)
					}
				} else {
					printOut = fmt.Sprintf("%s|%s|%d|%d|%s", gt, vafString, dr, dv, statusSV)
				}
				printPopulation = append(printPopulation, printOut)
			}
		}
		// for printing
		if (userParams.Uniq && suppVecSum == 1) || !userParams.Uniq {
			var suppVecOutput string
			if userParams.FixSuppVec {
				suppVecOutput = fmt.Sprintf("%s|%s", strings.Join(suppVecArray, ""), VCFLineFormatted.Info["SUPP_VEC"])
			} else {
				suppVecOutput = VCFLineFormatted.Info["SUPP_VEC"]
			}
			if userParams.OutputVCF {
				// make info, make samples, make vcf line
			} else {
				samplePrint = strings.Join(printPopulation, "\t")
				fmt.Printf("%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\n", VCFLineFormatted.Contig,
					VCFLineFormatted.Start, VCFLineFormatted.EndStr, VCFLineFormatted.Info["SVTYPE"],
					VCFLineFormatted.Info["SVLEN"], suppVecOutput, VCFLineFormatted.ID, samplePrint)
			}
		}
	}
}
