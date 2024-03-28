package utils

import (
	"fmt"
	"nutty/config"
	"nutty/vcf"
	"strconv"
	"strings"
)

func ParseSV(params *config.UserParam) {
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
			parseBy := "none"
			if params.InfoTag != "none" {
				parseBy = "info_tag"
			}
			switch parseBy {
			case "info_tag":
				ReadVCFInfo(line, &contigsVCF, params.InfoTag)
			default:
				ReadVCFEntry(line, &contigsVCF, sampleName, params)
			}
		}
	}
}

func ReadVCFEntry(VCFLineRaw string, contigs *map[string]int, sampleName string, userParams *config.UserParam) {
	var (
		dr       int
		dv       int
		vaf      float64
		vafPrint float64
	)
	lineSplit := strings.Split(VCFLineRaw, "\t")
	VCFLineFormatted := new(vcf.VCF)
	VCFLineFormatted.Contig = lineSplit[0]
	if _, ok := (*contigs)[VCFLineFormatted.Contig]; ok {
		VCFPosInt, err := strconv.Atoi(lineSplit[1])
		if err != nil {
			panic(err)
		}
		VCFLineFormatted.Pos = VCFPosInt
		VCFLineFormatted.ID = lineSplit[2]
		VCFLineFormatted.Ref = ""
		VCFLineFormatted.Alt = ""
		VCFLineFormatted.Quality = lineSplit[5]
		VCFLineFormatted.Filter = lineSplit[6]
		// split each key=value pair or flag
		info := make(map[string]string)
		for _, infoElem := range strings.Split(lineSplit[7], ";") {
			if strings.Contains(infoElem, "=") {
				infoKeyVal := strings.Split(infoElem, "=")
				if infoKeyVal[0] == "RNAMES" && userParams.SaveRNames {
					info[infoKeyVal[0]] = infoKeyVal[1]
				} else if infoKeyVal[0] == "RNAMES" && !userParams.SaveRNames {
					//
				} else {
					info[infoKeyVal[0]] = infoKeyVal[1]
				}

			} else {
				info[infoElem] = "flag"
			}
		}
		VCFLineFormatted.Info = info
		// we expect only one sample here, so we only use one
		sampleSV = make(map[string]map[string]string)
		sampleSV[sampleName] = make(map[string]string)
		formatSplit := strings.Split(lineSplit[8], ":")
		sampleSVSplit := strings.Split(lineSplit[9], ":")
		for idx := range formatSplit {
			sampleSV[sampleName][formatSplit[idx]] = sampleSVSplit[idx]
		}
		VCFLineFormatted.Samples = sampleSV
		// #CONTTIG\tSTART\tEND\tSVTYPE\tSVLEN\tGT\tAF\tREF\tALT\tID
		dr, err = strconv.Atoi(VCFLineFormatted.Samples[sampleName]["DR"])
		dv, err = strconv.Atoi(VCFLineFormatted.Samples[sampleName]["DV"])
		vaf = float64(dv) / float64(dr+dv)
		vafPrint = vaf * 100
		// Fix GT
		gt := VCFLineFormatted.Samples[sampleName]["GT"]
		if userParams.FixGT && gt == "./." && dr+dv >= userParams.MinSupp {
			if vaf <= vcf.VAFHomRef {
				gt = "0/0"
			} else if vaf >= vcf.VAFHomAlt {
				gt = "1/1"
			} else if vaf > vcf.VAFHomRef && vaf < vcf.VAFHomAlt {
				gt = "0/1"
			} else {
				//
			}
		}
		if dr+dv >= userParams.MinSupp {
			fmt.Printf("%s\t%d\t%s\t%s\t%s\t%s\t%0.3f\t%d\t%d\t%s\n", VCFLineFormatted.Contig,
				VCFLineFormatted.Pos, VCFLineFormatted.Info["END"], VCFLineFormatted.Info["SVTYPE"],
				VCFLineFormatted.Info["SVLEN"], gt, vafPrint, dr, dv, VCFLineFormatted.ID)
		}
	}
}

func ReadVCFInfo(VCFLineRaw string, contigs *map[string]int, userInfoTag string) {
	lineSplit := strings.Split(VCFLineRaw, "\t")
	VCFLineFormatted := new(vcf.VCF)
	VCFLineFormatted.Contig = lineSplit[0]
	if _, ok := (*contigs)[VCFLineFormatted.Contig]; ok {
		VCFPosInt, err := strconv.Atoi(lineSplit[1])
		if err != nil {
			panic(err)
		}
		VCFLineFormatted.Pos = VCFPosInt
		VCFLineFormatted.ID = lineSplit[2]
		VCFLineFormatted.Ref = ""
		VCFLineFormatted.Alt = ""
		VCFLineFormatted.Quality = ""
		VCFLineFormatted.Filter = lineSplit[6]
		// split each key=value pair or flag
		info := make(map[string]string)
		for _, infoElem := range strings.Split(lineSplit[7], ";") {
			if strings.Contains(infoElem, "=") {
				infoKeyVal := strings.Split(infoElem, "=")
				info[infoKeyVal[0]] = infoKeyVal[1]
			} else {
				info[infoElem] = "flag"
			}
		}
		VCFLineFormatted.Info = info
		if _, ok := VCFLineFormatted.Info[userInfoTag]; ok {
			fmt.Printf("%s:%d\t%s\t%s\n", VCFLineFormatted.Contig, VCFLineFormatted.Pos,
				VCFLineFormatted.ID, VCFLineFormatted.Info[userInfoTag])
		}
	}
}
