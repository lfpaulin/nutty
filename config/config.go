package config

import (
	"flag"
	"fmt"
	"os"
)

type UserParam struct {
	SubCMD        string
	VCF           string
	MinSupp       int
	MinSize       int
	MinContigLen  int
	FilerGT       string
	FilerBy       string
	FixGT         bool
	SaveRNames    bool
	AsBED         bool
	AsDev         bool
	Uniq          bool
	FixSuppVec    bool
	Cancer        string
	Mosaic        bool
	Germline      bool
	PaperID       string
	PaperAnalysis string
	Help          string
}

func GetParams() UserParam {
	// SV parsing for single sample VCF

	var (
		subCMD        string
		vcf           string
		minSupp       int
		minSize       int
		minContigLen  int
		filerGT       string
		filerBy       string
		fixGT         bool
		asBED         bool
		saveRNames    bool
		asDev         bool
		uniq          bool
		fixSuppVec    bool
		cancer        string
		mosaic        bool
		germline      bool
		help          string
		paperID       string
		paperAnalysis string
	)

	help = "----------------------------------------\n" +
		"Sniffles 2 helper\n" +
		"  Usage: snf2_parse <subcommand> <options>\n  Available subcommands:\n" +
		"    sv\n" +
		"    pop\n" +
		"    cancer\n" +
		"    paper\n" +
		"----------------------------------------"

	if len(os.Args) < 2 {
		fmt.Println("Please specify a subcommand and parameters")
		fmt.Println(help)
		os.Exit(1)
	}
	subCMD = os.Args[1]
	switch subCMD {
	case "help":
		fmt.Println(help)
	case "sv":
		cmdSVParse := flag.NewFlagSet("sv", flag.ExitOnError)
		cmdSVParse.StringVar(&vcf, "vcf", "none", "VCF file to read")
		cmdSVParse.IntVar(&minSupp, "min-supp", 10, "Min. read support for the SV calls, default = 10")
		cmdSVParse.IntVar(&minSize, "min-size", 50, "Min. SV size, default = 1, in case of BND this is skipped")
		cmdSVParse.IntVar(&minContigLen, "min-contig-len", 2000000, "Min. Contig/Chromosome size to be used, default = 2Mb")
		cmdSVParse.StringVar(&filerGT, "filer-gt", "none", "Removed genotypes from output")
		cmdSVParse.StringVar(&filerBy, "filer-by", "none:none", "Filter by some parameter:value")
		cmdSVParse.BoolVar(&fixGT, "fix-gt", false, "")
		cmdSVParse.BoolVar(&saveRNames, "save-rnames", false, "")
		cmdSVParse.BoolVar(&asBED, "as-bed", false, "")
		cmdSVParse.BoolVar(&asDev, "as-dev", false, "")
		err := cmdSVParse.Parse(os.Args[2:])
		if err != nil {
			panic(err)
			return UserParam{}
		}
		fmt.Printf("#CMD: snf2_parser sv --vcf %s --min-sup %d --min-size %d --filer-gt %s --filer-by %s "+
			"--fix-gt %t --as-bed %t --as-dev %t \n", vcf, minSupp, minSize, filerGT, filerBy, fixGT, asBED, asDev)
		return UserParam{
			SubCMD:       subCMD,
			VCF:          vcf,
			MinSupp:      minSupp,
			MinSize:      minSize,
			MinContigLen: minContigLen,
			FilerGT:      filerGT,
			FilerBy:      filerBy,
			FixGT:        fixGT,
			SaveRNames:   saveRNames,
			AsBED:        asBED,
			AsDev:        asDev,
		}
	case "pop":
		cmdPopParse := flag.NewFlagSet("pop", flag.ExitOnError)
		cmdPopParse.StringVar(&vcf, "vcf", "none", "VCF file to read")
		cmdPopParse.IntVar(&minSupp, "min-supp", 1, "Min. support for the SV calls (from SUPP_VEC), default = 1")
		cmdPopParse.IntVar(&minSize, "min-size", 50, "Min. absolute size of the event (except for BDN), default = 1")
		cmdPopParse.BoolVar(&uniq, "uniq", false, "Show only those that appear in a single individual (from SUPP_VEC)")
		cmdPopParse.BoolVar(&fixSuppVec, "fix-suppvec", false, "")
		cmdPopParse.BoolVar(&asDev, "as-dev", false, "")
		err := cmdPopParse.Parse(os.Args[2:])
		if err != nil {
			panic(err)
			return UserParam{}
		}
		fmt.Printf("#CMD: snf2_parser pop --vcf %s --min-supp %d --min-size %d --uniq %t --fix-suppvec %t --as-dev %t\n",
			vcf, minSupp, minSize, uniq, fixSuppVec, asDev)
		return UserParam{
			SubCMD:     subCMD,
			VCF:        vcf,
			MinSupp:    minSupp,
			MinSize:    minSize,
			Uniq:       uniq,
			FixSuppVec: fixSuppVec,
			AsDev:      asDev,
		}
	case "cancer":
		cmdCancerParse := flag.NewFlagSet("cancer", flag.ExitOnError)
		cmdCancerParse.StringVar(&vcf, "vcf", "none", "VCF file to read")
		cmdCancerParse.BoolVar(&uniq, "uniq", false, "Show only those that appear in a single individual (from SUPP_VEC)")
		cmdCancerParse.StringVar(&cancer, "cancer", "none", "SUPP_VEC if the cancer samples")
		cmdCancerParse.BoolVar(&mosaic, "mosaic", false, "Show mosaic calls (5% <= VAF <= 25%")
		cmdCancerParse.BoolVar(&germline, "germline", false, "Show germline calls (VAF>=25%")
		cmdCancerParse.BoolVar(&asDev, "as-dev", false, "")
		err := cmdCancerParse.Parse(os.Args[2:])
		if err != nil {
			panic(err)
			return UserParam{}
		}
		fmt.Printf("#CMD: snf2_parser cancer --vcf %s --uniq %t --cancer %s --mosaic %t --germline %t --as-dev %t\n",
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
	case "paper":
		cmdPapers := flag.NewFlagSet("paper", flag.ExitOnError)
		cmdPapers.StringVar(&paperID, "paper-id", "0", "Paper ID, see README")
		cmdPapers.StringVar(&vcf, "vcf", "none", "VCF file to read")
		cmdPapers.StringVar(&paperAnalysis, "paper-analysis", "none", "Analysis from the paper")
		err := cmdPapers.Parse(os.Args[2:])
		if err != nil {
			panic(err)
			return UserParam{}
		}
		fmt.Printf("#CMD: snf2_parser paper --paper-id %s --vcf %s --paper-analysis %s\n",
			paperID, vcf, paperAnalysis)
		return UserParam{
			SubCMD:        subCMD,
			VCF:           vcf,
			PaperID:       paperID,
			PaperAnalysis: paperAnalysis,
		}
	default:
		fmt.Printf("Unknown subcommand: %s\n", subCMD)
		fmt.Println(help)
		os.Exit(1)
	}
	return UserParam{
		SubCMD: "error",
	}
}
