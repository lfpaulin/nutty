package utils

import (
	"sniffles2_helper_go/config"
	"sniffles2_helper_go/vcf"
)

func ParseSV(params *config.UserParam) {
	println("hey ", params.SubCMD)
	vcf := vcf.VCF{
		"contig": "1",
	}
}
