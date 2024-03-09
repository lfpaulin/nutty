package config

import (
	"flag"
	"fmt"
	"os"
)

type UserParam struct {
	SubCMD   string
	VCF      string
	MinSupp  int
	MinSize  int
	FilerGT  string
	FilerBy  string
	AsBED    bool
	AsDev    bool
	Uniq     bool
	Cancer   string
	Mosaic   bool
	Germline bool
}

func GetParams() UserParam {
	// SV parsing for single sample VCF

	var (
		subCMD   string
		vcf      string
		minSupp  int
		minSize  int
		filerGT  string
		filerBy  string
		asBED    bool
		asDev    bool
		uniq     bool
		cancer   string
		mosaic   bool
		germline bool
	)

	if len(os.Args) < 2 {
		fmt.Println("Please specify a subcommand and parameters")
		os.Exit(1)
	}
	subCMD = os.Args[1]
	switch subCMD {
	case "sv":
		cmdSVParse := flag.NewFlagSet("sv", flag.ExitOnError)
		cmdSVParse.StringVar(&vcf, "vcf", "none", "VCF file to read")
		cmdSVParse.IntVar(&minSupp, "min-supp", 1, "Min. read support for the SV calls, default = 1")
		cmdSVParse.IntVar(&minSize, "min-size", 1, "Min. SV size, default = 1, in case of BND this is skipped")
		cmdSVParse.StringVar(&filerGT, "filer-gt", "none", "Removed genotypes from output")
		cmdSVParse.StringVar(&filerBy, "filer-by", "none", "Filter by some parameter:value")
		cmdSVParse.BoolVar(&asBED, "as-bed", false, "")
		cmdSVParse.BoolVar(&asDev, "as-dev", false, "")
		cmdSVParse.Parse(os.Args[2:])
		fmt.Printf("CMD: snf2_parser sv --vcf %s --min-sup %d --min-size %d --filer-gt %s --filer-by %s "+
			"--as-bed %t --as-dev %t \n", vcf, minSupp, minSize, filerGT, filerBy, asBED, asDev)
		return UserParam{
			SubCMD:  subCMD,
			VCF:     vcf,
			MinSupp: minSupp,
			MinSize: minSize,
			FilerGT: filerGT,
			FilerBy: filerBy,
			AsBED:   asBED,
			AsDev:   asDev,
		}
	case "pop":
		cmdPopParse := flag.NewFlagSet("pop", flag.ExitOnError)
		cmdPopParse.StringVar(&vcf, "vcf", "none", "VCF file to read")
		cmdPopParse.IntVar(&minSupp, "min-supp", 1, "Min. support for the SV calls (from SUPP_VEC), default = 1")
		cmdPopParse.IntVar(&minSize, "min-size", 1, "Min. absolute size of the event (except for BDN), default = 1")
		cmdPopParse.BoolVar(&uniq, "uniq", false, "Show only those that appear in a single individual (from SUPP_VEC)")
		cmdPopParse.BoolVar(&asDev, "as-dev", false, "")
		cmdPopParse.Parse(os.Args[2:])
		fmt.Printf("CMD: snf2_parser pop --vcf %s --min-supp %d --min-size %d --uniq %t --as-dev %t \n",
			vcf, minSupp, minSize, uniq, asDev)
		return UserParam{
			SubCMD:  subCMD,
			VCF:     vcf,
			MinSupp: minSupp,
			MinSize: minSize,
			AsBED:   asBED,
			AsDev:   asDev,
		}
	case "cancer":
		cmdCancerParse := flag.NewFlagSet("cancer", flag.ExitOnError)
		cmdCancerParse.StringVar(&vcf, "vcf", "none", "VCF file to read")
		cmdCancerParse.BoolVar(&uniq, "uniq", false, "Show only those that appear in a single individual (from SUPP_VEC)")
		cmdCancerParse.StringVar(&cancer, "cancer", "none", "SUPP_VEC if the cancer samples")
		cmdCancerParse.BoolVar(&mosaic, "mosaic", false, "Show mosaic calls (5% <= VAF <= 25%")
		cmdCancerParse.BoolVar(&germline, "germline", false, "Show germline calls (VAF>=25%")
		cmdCancerParse.BoolVar(&asDev, "as-dev", false, "")
		cmdCancerParse.Parse(os.Args[2:])
		fmt.Printf("CMD: snf2_parser cancer --vcf %s --uniq %t --cancer %s --mosaic %t --germline %t --as-dev %t\n",
			vcf, uniq, cancer, mosaic, germline, asDev)
		return UserParam{
			SubCMD:   subCMD,
			VCF:      vcf,
			Uniq:     uniq,
			Cancer:   cancer,
			Mosaic:   mosaic,
			Germline: germline,
			AsDev:    asDev,
		}
	default:
		fmt.Printf("Unknown subcommand: %s\n", subCMD)
		os.Exit(1)
	}
	return UserParam{
		SubCMD: "error",
	}
}
