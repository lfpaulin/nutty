package utils

import (
	"fmt"
	"nutty/config"
	"nutty/vcf"
	"strings"
)

func ParseCancer(params *config.UserParam) {
	VCFReader := vcf.ReadVCF(params.VCF)
	defer func(VCFReader *vcf.FileScanner) {
		err := VCFReader.Close()
		if err != nil {
			panic(err)
		}
	}(VCFReader)
	for VCFReader.Scan() {
		line := strings.TrimSpace(VCFReader.Text())
		switch {
		case strings.Contains(line, "##") && strings.Contains(line, "contig"):
			fmt.Println("contig => ", line)
		case strings.Contains(line, "##") && strings.Contains(line, "INFO"):
			fmt.Println("info => ", line)
		case strings.Contains(line, "##") && strings.Contains(line, "FORMAT"):
			fmt.Println("format => ", line)
		default:
			//
		}
	}
}
