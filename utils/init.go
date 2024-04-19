package utils

import (
	"fmt"
	"nutty/config"
	"nutty/vcf"
	"strconv"
	"strings"
)

var (
	formatVCF   []string
	sampleName  string
	sampleNames []string
	info        map[string]string
	sampleSV    map[string]map[string]string
)
var contigsVCF = make(map[string]int)
var infoVCF = make(map[string]string)

var (
	dr             int
	dv             int
	vaf            float64
	gt             string
	vafString      string
	svIDMergeCount int
	// svIDMerge []string
)

const (
	fixHetVAFMin float64 = 0.33
	fixAltVAFMin float64 = 0.66
	// #CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO	FORMAT
	indexChrom   int = 0
	indexPos     int = 1
	indexID      int = 2
	indexRef     int = 3
	indexAlt     int = 4
	indexQual    int = 5
	indexFilter  int = 6
	indexInfo    int = 7
	indexFormat  int = 8
	indexSamples int = 9
)

func VCFHeader(lineHeaderVCF *string, userParams *config.UserParam) {
	switch {
	case strings.Contains(*lineHeaderVCF, "##") && strings.Contains(*lineHeaderVCF, "contig"):
		contigMatch := vcf.HeaderRegex(*lineHeaderVCF, "contig")
		contigName := contigMatch[1]
		contigSize, err := strconv.Atoi(contigMatch[2])
		if err != nil {
			panic(err)
		}
		if contigSize > userParams.MinContigLen {
			contigsVCF[contigName] = contigSize
			if userParams.OutputVCF {
				fmt.Println(*lineHeaderVCF)
			}
		}
	case strings.Contains(*lineHeaderVCF, "##") && strings.Contains(*lineHeaderVCF, "INFO="):
		infoMatch := vcf.HeaderRegex(*lineHeaderVCF, "info")
		infoVCF[infoMatch[1]] = infoMatch[2]
		if userParams.OutputVCF {
			fmt.Println(*lineHeaderVCF)
		}
	case strings.Contains(*lineHeaderVCF, "##") && strings.Contains(*lineHeaderVCF, "FORMAT="):
		formatMatch := vcf.HeaderRegex(*lineHeaderVCF, "format")
		formatVCF = append(formatVCF, formatMatch[1])
		if userParams.OutputVCF {
			fmt.Println(*lineHeaderVCF)
		}
	case strings.Contains(*lineHeaderVCF, "#CHROM"):
		lineHeaderVCFSplit := strings.Split(*lineHeaderVCF, "\t")
		for _, sample := range lineHeaderVCFSplit[9:] {
			sampleNames = append(sampleNames, sample)
		}
		sampleNamesInfo := strings.Join(sampleNames, ", ")
		sampleNamesHeader := strings.Join(sampleNames, "\t")
		if userParams.OutputVCF {
			fmt.Println(*lineHeaderVCF)
		}
		// Here we print the header of the parsed file
		if !userParams.AsBED && !userParams.OutputVCF && userParams.InfoTag == "" {
			switch userParams.SubCMD {
			case "sv":
				if userParams.NoHeader {
					//
				} else {
					fmt.Println("##Sample name: ", sampleNamesInfo)
					fmt.Println("#CONTIG\tSTART\tEND\tSVTYPE\tSVLEN\tGT\tVAF\tREFC\tALTC\tID")
				}
			case "pop":
				if userParams.NoHeader {
					//
				} else {
					fmt.Println("##Sample names: ", sampleNamesInfo)
					if !userParams.Uniq {
						fmt.Printf("#CONTIG\tSTART\tEND\tSVTYPE\tSVLEN\tSUPPVEC\tID\t%s\n", sampleNamesHeader)
					} else {
						fmt.Println("#CONTIG\tSTART\tEND\tSVTYPE\tSVLEN\tID")
					}
				}
			case "ghost":
				fmt.Printf("#CONTIG\tSTART\tEND\tSVTYPE\tSVLEN\tSUPPVEC\tCOPY_NUMBER\tID\t%s\n", sampleNamesHeader)
			default:
				//
			}
		}
	case strings.Contains(*lineHeaderVCF, "#"):
		if userParams.OutputVCF && !userParams.NoHeader {
			fmt.Println(*lineHeaderVCF)
		}
	default:
		//
	}
}
