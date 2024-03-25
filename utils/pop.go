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
	var lineSplit = strings.Split(VCFLineRaw, "\t")
	var contig = lineSplit[indexChrom]
	if _, ok := (*contigs)[contig]; ok {
		var VCFRecord = new(vcf.VCF)
		VCFRecord.Contig = contig
		VCFPosInt, err := strconv.Atoi(lineSplit[indexPos])
		if err != nil {
			panic(err)
		}
		VCFRecord.Pos = VCFPosInt
		VCFRecord.Start = VCFPosInt
		VCFRecord.ID = lineSplit[indexID]
		VCFRecord.Ref = ""
		VCFRecord.Alt = ""
		VCFRecord.Quality = lineSplit[indexQual]
		VCFRecord.Filter = lineSplit[indexFilter]
		VCFRecord.Format = lineSplit[indexFormat]
		// split each key=value pair or flag
		info = make(map[string]string)
		for _, infoElem := range strings.Split(lineSplit[indexInfo], ";") {
			if strings.Contains(infoElem, "=") {
				infoKeyVal := strings.Split(infoElem, "=")
				info[infoKeyVal[0]] = infoKeyVal[1]
			} else {
				info[infoElem] = "flag"
			}
		}
		VCFRecord.Info = info
		if info["SVTYPE"] == "BND" {
			VCFRecord.End = VCFRecord.Start + 1
			VCFRecord.EndStr = lineSplit[indexAlt] // Alt
			info["SVLEN"] = "1"
		} else {
			end, err := strconv.Atoi(info["END"])
			if err != nil {
				fmt.Println("[FAILED] strconv.Atoi(info[\"END\"])")
				panic(err)
			}
			VCFRecord.End = end
			VCFRecord.EndStr = info["END"]
		}
		var suppVecSum = 0
		var suppVecUniq = 0
		var suppVecArrayUpdate = strings.Split(VCFRecord.Info["SUPP_VEC"], "")
		for suppVecIdx, suppVecElem := range suppVecArrayUpdate {
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
		var nPop int = len(*sampleNames)
		var printPopulation = make([]string, nPop, nPop)
		var sampleNameUniq string
		sampleSV = make(map[string]map[string]string)
		var formatSplit = strings.Split(VCFRecord.Format, ":")
		for sidx, sample := range *sampleNames {
			sampleSVSplit := strings.Split(lineSplit[indexSamples+sidx], ":")
			sampleSV[sample] = make(map[string]string)
			for idx := range formatSplit {
				sampleSV[sample][formatSplit[idx]] = sampleSVSplit[idx]
			}
		}
		VCFRecord.Samples = sampleSV
		if userParams.Uniq && suppVecSum == 1 {
			sampleNameUniq = (*sampleNames)[suppVecUniq]
		}
		for sid, sampleName := range *sampleNames {
			if (userParams.Uniq && suppVecSum == 1 && sampleName == sampleNameUniq) || !userParams.Uniq {
				dr, err = strconv.Atoi(VCFRecord.Samples[sampleName]["DR"])
				dv, err = strconv.Atoi(VCFRecord.Samples[sampleName]["DV"])
				gt = VCFRecord.Samples[sampleName]["GT"]
				vaf = float64(dv) / float64(dr+dv)
				if userParams.FixGT {
					if minCoverage > dr+dv {
						statusSV = "undefined"
						gt = "./."
						vafString = "n/a"
						suppVecArrayUpdate[sid] = "0"
					} else {
						if vaf >= userParams.MinVAFMosaic {
							statusSV = "germline"
							suppVecArrayUpdate[sid] = "1"
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
							suppVecArrayUpdate[sid] = "1"
							if gt == "./." && userParams.FixGT {
								gt = "0/0"
							}
						} else if vaf < userParams.MinVAFMosaic && vaf > 0.0 {
							statusSV = "lowVAF"
							suppVecArrayUpdate[sid] = "0"
						} else if vaf == 0.0 {
							statusSV = "reference"
							suppVecArrayUpdate[sid] = "0"
						} else {
							statusSV = "undefined"
							suppVecArrayUpdate[sid] = "0"
						}
						vafString = fmt.Sprintf("%0.3f", vaf)
					}
				}
				VCFRecord.Samples[sampleName]["GT"] = gt
				if _, ok := VCFRecord.Info["AF"]; ok {
					VCFRecord.Info["AF"] = vafString
				}
				if userParams.OnlyGT {
					if statusSV == "undefined" {
						printOut = fmt.Sprintf("%s*", gt)
					} else {
						printOut = fmt.Sprintf("%s", gt)
					}
				} else {
					if userParams.OutputVCF{
						// GT:GQ:DR:DV
						var nFormat = len(formatSplit)
						var sampleVCFCol = make([]string, nFormat, nFormat)
						for _, fromatKey := range formatSplit {
							sampleVCFCol = append(sampleVCFCol, VCFRecord.Samples[sampleName][fromatKey])
						}
						printOut = strings.Join(sampleVCFCol, ":")
					} else {
						printOut = fmt.Sprintf("%s|%s|%d|%d|%s", gt, vafString, dr, dv, statusSV)
					}
				}
				printPopulation = append(printPopulation, printOut)
			}
		}
		// for printing
		if (userParams.Uniq && suppVecSum == 1) || !userParams.Uniq {
			if userParams.FixSuppVec {
				VCFRecord.Info["SUPP_VEC"] = fmt.Sprintf("%s", strings.Join(suppVecArrayUpdate, ""))
			}
			samplePrint = strings.Join(printPopulation, "\t")
			if userParams.OutputVCF {
				// #CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO	FORMAT
				var outputVCFINFO string = lineSplit[indexInfo]
				// make info, make samples, make vcf line
				VCFRecord.PrintVCF(&lineSplit[indexRef], &lineSplit[indexAlt], &outputVCFINFO, &samplePrint)
			} else if userParams.AsBED {
				VCFRecord.PrintBED()
			} else {
				VCFRecord.PrintParsed(&samplePrint)
			}
		}
	}
}
