package utils

import (
	"fmt"
	"nutty/config"
	"nutty/vcf"
	"strconv"
	"strings"
)

var sampleNames []string

func ParsePop(params *config.UserParam) {
	VCFReader := vcf.ReadVCF(params.VCF)
	defer func(VCFReader *vcf.FileScanner) {
		err := VCFReader.Close()
		if err != nil {
			panic(err)
		}
	}(VCFReader)
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
			if !params.AsBED {
				fmt.Println("##Sample names: ", sampleNamesHeader)
				fmt.Println(vcf.HeaderOut)
			}
		case strings.Contains(line, "#"):
			//
		default:
			// each entry
			ReadVCFPopEntry(line, &contigsVCF, &sampleNames, params)
		}
	}
}

func ReadVCFPopEntry(VCFLineRaw string, contigs *map[string]int, sampleNames *[]string, userParams *config.UserParam) {
	/*var (
		dr       int
		dv       int
		vaf      float64
		vafPrint float64
	)*/
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
				info[infoKeyVal[0]] = infoKeyVal[1]
			} else {
				info[infoElem] = "flag"
			}
		}
		VCFLineFormatted.Info = info
		/*
			// Aqui
			if dr+dv >= userParams.MinSupp {
				fmt.Printf("%s\t%d\t%s\t%s\t%s\t%s\t%0.3f\t%d\t%d\t%s\n", VCFLineFormatted.Contig,
					VCFLineFormatted.Pos, VCFLineFormatted.Info["END"], VCFLineFormatted.Info["SVTYPE"],
					VCFLineFormatted.Info["SVLEN"], gt, vafPrint, dr, dv, VCFLineFormatted.ID)
			}

			// we expect only one sample here, so we only use one
			sampleSV := make(map[string]map[string]string)
			formatSplit := strings.Split(lineSplit[8], ":")
			sampleSVSplit := strings.Split(lineSplit[9], ":")
			for sidx := range *sampleNames {
				for idx := range formatSplit {
					sampleSV[(*sampleNames)[sidx]][formatSplit[idx]] = sampleSVSplit[idx]
				}
			}
			VCFLineFormatted.Samples = sampleSV
			// #CONTTIG\tSTART\tEND\tSVTYPE\tSVLEN\tGT\tAF\tREF\tALT\tID
			dr, err = strconv.Atoi(VCFLineFormatted.Samples["DR"])
			dv, err = strconv.Atoi(VCFLineFormatted.Samples["DV"])
			vaf = float64(dv) / float64(dr+dv)
			vafPrint = vaf * 100
			// Fix GT
			gt := VCFLineFormatted.Samples["GT"]
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
		*/
	}
}
