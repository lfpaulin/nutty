package utils

import (
	"fmt"
	"nutty/config"
	"nutty/vcf"
	"strconv"
	"strings"
)

func ParsePopSpc(params *config.UserParam) {
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
		if strings.Contains(line, "#") {
			VCFHeader(&line, params)
		} else {
			// each entry
			ReadVCFPopEntrySpc(&line, &contigsVCF, &sampleNames, params)
		}
	}
}

func ReadVCFPopEntrySpc(VCFLineRaw *string, contigs *map[string]int, sampleNames *[]string, userParams *config.UserParam) {
	var lineSplit = strings.Split(*VCFLineRaw, "\t")
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
		VCFRecord.Ref = lineSplit[indexRef]
		VCFRecord.Alt = lineSplit[indexAlt]
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
		end, err := strconv.Atoi(info["END"])
		if err != nil {
			fmt.Println("[FAILED] strconv.Atoi(info[\"END\"])")
			panic(err)
		}
		VCFRecord.End = end
		VCFRecord.EndStr = info["END"]
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
		var samplePrint string
		var printPopulation []string
		var sampleNameUniq string
		sampleSV = make(map[string]map[string]string)
		var formatSplit = strings.Split(VCFRecord.Format, ":")
		for sampleIndex, sample := range *sampleNames {
			sampleSVSplit := strings.Split(lineSplit[indexSamples+sampleIndex], ":")
			sampleSV[sample] = make(map[string]string)
			for idx := range formatSplit {
				sampleSV[sample][formatSplit[idx]] = sampleSVSplit[idx]
			}
		}
		VCFRecord.Samples = sampleSV
		if userParams.Uniq && suppVecSum == 1 {
			sampleNameUniq = (*sampleNames)[suppVecUniq]
		}
		for _, sampleName := range *sampleNames {
			if (userParams.Uniq && suppVecSum == 1 && sampleName == sampleNameUniq) || !userParams.Uniq {
				svIDMergeCount = len(strings.Split(VCFRecord.Samples[sampleName]["ID"], ","))
				gt = VCFRecord.Samples[sampleName]["GT"]
				printPopulation = append(printPopulation, gt)
			}
		}
		// for printing
		if (userParams.Uniq && suppVecSum == 1) || !userParams.Uniq {
			samplePrint = strings.Join(printPopulation, "\t")
			if userParams.AsBED {
				VCFRecord.PrintBED()
			} else {
				VCFRecord.PrintParsed(&samplePrint)
			}
		}
	}
}
