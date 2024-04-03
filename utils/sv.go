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
			} else if params.AsBED {
				parseBy = "bed"
			} else {
				//
			}
			switch parseBy {
			case "info_tag":
				ReadVCFInfo(&line, &contigsVCF, params.InfoTag)
			case "bed":
				ReadVCF2BED(&line, &contigsVCF)
			default:
				ReadVCFEntry(&line, &contigsVCF, sampleName, params)
			}
		}
	}
}

func ReadVCFEntry(VCFLineRaw *string, contigs *map[string]int, sampleName string, userParams *config.UserParam) {
	var (
		dr       int
		dv       int
		vaf      float64
		vafPrint float64
	)
	lineSplit := strings.Split(*VCFLineRaw, "\t")
	VCFRecord := new(vcf.VCF)
	VCFRecord.Contig = lineSplit[0]
	if _, ok := (*contigs)[VCFRecord.Contig]; ok {
		VCFPosInt, err := strconv.Atoi(lineSplit[1])
		if err != nil {
			panic(err)
		}
		VCFRecord.Pos = VCFPosInt
		VCFRecord.ID = lineSplit[2]
		VCFRecord.Ref = ""
		VCFRecord.Alt = ""
		VCFRecord.Quality = lineSplit[5]
		VCFRecord.Filter = lineSplit[6]
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
		// we expect only one sample here, so we only use one
		sampleSV = make(map[string]map[string]string)
		sampleSV[sampleName] = make(map[string]string)
		formatSplit := strings.Split(lineSplit[8], ":")
		sampleSVSplit := strings.Split(lineSplit[9], ":")
		for idx := range formatSplit {
			sampleSV[sampleName][formatSplit[idx]] = sampleSVSplit[idx]
		}
		VCFRecord.Samples = sampleSV
		// #CONTTIG\tSTART\tEND\tSVTYPE\tSVLEN\tGT\tAF\tREF\tALT\tID
		dr, err = strconv.Atoi(VCFRecord.Samples[sampleName]["DR"])
		dv, err = strconv.Atoi(VCFRecord.Samples[sampleName]["DV"])
		vaf = float64(dv) / float64(dr+dv)
		vafPrint = vaf * 100
		// Fix GT
		gt := VCFRecord.Samples[sampleName]["GT"]
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
			fmt.Printf("%s\t%d\t%s\t%s\t%s\t%s\t%0.3f\t%d\t%d\t%s\n", VCFRecord.Contig,
				VCFRecord.Pos, VCFRecord.Info["END"], VCFRecord.Info["SVTYPE"],
				VCFRecord.Info["SVLEN"], gt, vafPrint, dr, dv, VCFRecord.ID)
		}
	}
}

func ReadVCFInfo(VCFLineRaw *string, contigs *map[string]int, userInfoTag string) {
	lineSplit := strings.Split(*VCFLineRaw, "\t")
	VCFRecord := new(vcf.VCF)
	VCFRecord.Contig = lineSplit[0]
	if _, ok := (*contigs)[VCFRecord.Contig]; ok {
		VCFPosInt, err := strconv.Atoi(lineSplit[1])
		if err != nil {
			panic(err)
		}
		VCFRecord.Pos = VCFPosInt
		VCFRecord.ID = lineSplit[2]
		VCFRecord.Ref = ""
		VCFRecord.Alt = ""
		VCFRecord.Quality = ""
		VCFRecord.Filter = lineSplit[6]
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
		VCFRecord.Info = info
		if _, ok := VCFRecord.Info[userInfoTag]; ok {
			fmt.Printf("%s:%d\t%s\t%s\n", VCFRecord.Contig, VCFRecord.Pos,
				VCFRecord.ID, VCFRecord.Info[userInfoTag])
		}
	}
}

func ReadVCF2BED(VCFLineRaw *string, contigs *map[string]int) {
	lineSplit := strings.Split(*VCFLineRaw, "\t")
	VCFRecord := new(vcf.VCF)
	VCFRecord.Contig = lineSplit[0]
	if _, ok := (*contigs)[VCFRecord.Contig]; ok {
		VCFPosInt, err := strconv.Atoi(lineSplit[1])
		if err != nil {
			panic(err)
		}
		VCFRecord.Pos = VCFPosInt
		VCFRecord.ID = lineSplit[2]
		// split each key=value pair or flag
		var end int = 0
		var endStr string
		var svtype string = ""
		for _, infoElem := range strings.Split(lineSplit[7], ";") {
			if strings.Contains(infoElem, "=") {
				infoKeyVal := strings.Split(infoElem, "=")
				switch infoKeyVal[0]{
				case "END":
					endStr = infoKeyVal[1]
				case "SVTYPE":
					svtype = infoKeyVal[1]
				default:
					//
				}
			}
		}
		if svtype == "BND" {
			end = VCFRecord.Pos + 1
		} else {
			end, err = strconv.Atoi(endStr)
			if err != nil {
				fmt.Println("[FAILED] strconv.Atoi(info -> END)")
				panic(err)
			}
		}
		fmt.Printf("%s\t%d\t%d\t%s|%s\n", VCFRecord.Contig, VCFRecord.Pos, end,
			VCFRecord.ID, svtype)
	}
}
