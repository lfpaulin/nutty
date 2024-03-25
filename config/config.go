package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const version string = "0.1"

type UserParam struct {
	SubCMD         string
	VCF            string
	MinSupp        int
	MinSize        int
	MinContigLen   int
	FilerGT        string
	FilerBy        string
	FixGT          bool
	OnlyGT         bool
	SaveRNames     bool
	OutputVCF      bool
	AsBED          bool
	AsDev          bool
	Uniq           bool
	FixSuppVec     bool
	Cancer         string
	Germline       bool
	MinVAFGermline float64
	Mosaic         bool
	MinVAFMosaic   float64
	MaxVAFMosaic   float64
	InfoTag        string
	PaperID        string
	PaperAnalysis  string
	Help           string
}

func GetParams() UserParam {
	// SV parsing for single sample VCF

	var (
		subCMD         string
		vcf            string
		minSupp        int
		minSize        int
		minContigLen   int
		filerGT        string
		filerBy        string
		fixGT          bool
		onlyGT         bool
		outputVCF      bool
		asBED          bool
		saveRNames     bool
		asDev          bool
		uniq           bool
		fixSuppVec     bool
		cancer         string
		germline       bool
		minVAFGermline float64
		mosaic         bool
		minVAFMosaic   float64
		maxVAFMosaic   float64
		infotag        string
		help           string
		paperID        string
		paperAnalysis  string
	)

	help = "----------------------------------------\n" +
		"Nutty: a VCF parser for Sniffles2\n" +
		"  Usage: nutty <subcommand> <options>\n  Available subcommands:\n" +
		"    sv\n" +
		"    pop\n" +
		"    cancer\n" +
		"    paper\n" +
		"    ~~~~~~~~~~\n" +
		"    help\n" +
		"    version\n" +
		"----------------------------------------"

	if len(os.Args) < 2 {
		fmt.Println("Please specify a subcommand and parameters")
		fmt.Println(help)
		os.Exit(1)
	}
	// var cmdUsed []string
	subCMD = os.Args[1]
	switch subCMD {
	case "help":
		fmt.Println(help)
		return UserParam{
			SubCMD: subCMD,
		}
	case "version":
		fmt.Println("nutty version ", version)
		return UserParam{
			SubCMD: subCMD,
		}
	case "sv":
		cmdSVParse := flag.NewFlagSet("sv", flag.ExitOnError)
		cmdSVParse.StringVar(&vcf, "vcf", "stdin", "VCF file to read")
		cmdSVParse.IntVar(&minSupp, "min-supp", 10, "Min. read support for the SV calls, default = 10")
		cmdSVParse.IntVar(&minSize, "min-size", 50, "Min. SV size, default = 1, in case of BND this is skipped")
		cmdSVParse.IntVar(&minContigLen, "min-contig-len", 2000000, "Min. Contig/Chromosome size to be used, default = 2Mb")
		cmdSVParse.Float64Var(&minVAFGermline, "min-vaf-germline", 0.25, "Min. VAF considered for a germline SVs, default = 25%")
		cmdSVParse.Float64Var(&minVAFMosaic, "min-vaf-mosaic", 0.05, "Min. VAF considered for a mosaic SVs, default = 5%")
		cmdSVParse.Float64Var(&maxVAFMosaic, "max-vaf-mosaic", 0.25, "Max. VAF considered for a mosaic SVs, default < 25%")
		cmdSVParse.StringVar(&filerGT, "filer-gt", "none", "Removed genotypes from output")
		cmdSVParse.StringVar(&filerBy, "filer-by", "none:none", "Filter by some parameter:value")
		cmdSVParse.BoolVar(&fixGT, "fix-gt", false, "")
		cmdSVParse.BoolVar(&saveRNames, "save-rnames", false, "")
		cmdSVParse.StringVar(&infotag, "info-tag", "none", "Extracts tag from the INFO field")
		cmdSVParse.BoolVar(&asBED, "as-bed", false, "")
		cmdSVParse.BoolVar(&asDev, "as-dev", false, "")
		err := cmdSVParse.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
		if asDev{
			fmt.Printf("## CMD: nutty pop %s \n", FlagsState(cmdSVParse))
		}
		return UserParam{
			SubCMD:         subCMD,
			VCF:            vcf,
			MinSupp:        minSupp,
			MinSize:        minSize,
			MinContigLen:   minContigLen,
			MinVAFGermline: minVAFGermline,
			MinVAFMosaic:   minVAFMosaic,
			MaxVAFMosaic:   maxVAFMosaic,
			FilerGT:        filerGT,
			FilerBy:        filerBy,
			FixGT:          fixGT,
			InfoTag:        infotag,
			SaveRNames:     saveRNames,
			AsBED:          asBED,
			AsDev:          asDev,
		}
	case "pop":
		cmdPopParse := flag.NewFlagSet("pop", flag.ExitOnError)
		cmdPopParse.StringVar(&vcf, "vcf", "stdin", "VCF file to read")
		cmdPopParse.IntVar(&minSupp, "min-supp", 1, "Min. support for the SV calls (from SUPP_VEC), default = 1")
		cmdPopParse.IntVar(&minSize, "min-size", 50, "Min. absolute size of the event (except for BDN), default = 1")
		cmdPopParse.Float64Var(&minVAFGermline, "min-vaf-germline", 0.25, "Min. VAF considered for a germline SVs, default = 25%")
		cmdPopParse.Float64Var(&minVAFMosaic, "min-vaf-mosaic", 0.05, "Min. VAF considered for a mosaic SVs, default = 5%")
		cmdPopParse.Float64Var(&maxVAFMosaic, "max-vaf-mosaic", 0.25, "Max. VAF considered for a mosaic SVs, default < 25%")
		cmdPopParse.BoolVar(&uniq, "uniq", false, "Show only those that appear in a single individual (from SUPP_VEC)")
		cmdPopParse.BoolVar(&fixGT, "fix-gt", false, "")
		cmdPopParse.BoolVar(&fixSuppVec, "fix-suppvec", false, "")
		cmdPopParse.BoolVar(&outputVCF, "output-vcf", false, "")
		cmdPopParse.BoolVar(&asBED, "as-bed", false, "")
		cmdPopParse.BoolVar(&onlyGT, "only-gt", false, "")
		cmdPopParse.BoolVar(&asDev, "as-dev", false, "")
		err := cmdPopParse.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
		if asDev{
			fmt.Printf("## CMD: nutty pop %s \n", FlagsState(cmdPopParse))
		}
		return UserParam{
			SubCMD:         subCMD,
			VCF:            vcf,
			MinSupp:        minSupp,
			MinSize:        minSize,
			MinVAFGermline: minVAFGermline,
			MinVAFMosaic:   minVAFMosaic,
			MaxVAFMosaic:   maxVAFMosaic,
			Uniq:           uniq,
			AsBED:          asBED,
			FixSuppVec:     fixSuppVec,
			FixGT:          fixGT,
			OutputVCF:      outputVCF,
			OnlyGT:         onlyGT,
			AsDev:          asDev,
		}
	case "cancer":
		cmdCancerParse := flag.NewFlagSet("cancer", flag.ExitOnError)
		cmdCancerParse.StringVar(&vcf, "vcf", "stdin", "VCF file to read")
		cmdCancerParse.BoolVar(&uniq, "uniq", false, "Show only those that appear in a single individual (from SUPP_VEC)")
		cmdCancerParse.StringVar(&cancer, "cancer", "none", "SUPP_VEC if the cancer samples")
		cmdCancerParse.BoolVar(&mosaic, "mosaic", false, "Show mosaic calls (5% <= VAF <= 25%")
		cmdCancerParse.BoolVar(&germline, "germline", false, "Show germline calls (VAF>=25%")
		cmdCancerParse.BoolVar(&asDev, "as-dev", false, "")
		err := cmdCancerParse.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
		if asDev{
			fmt.Printf("## CMD: nutty pop %s \n", FlagsState(cmdCancerParse))
		}
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
		cmdPapers.StringVar(&vcf, "vcf", "stdin", "VCF file to read")
		cmdPapers.StringVar(&paperAnalysis, "paper-analysis", "none", "Analysis from the paper")
		err := cmdPapers.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
		if asDev{
			fmt.Printf("## CMD: nutty pop %s \n", FlagsState(cmdPapers))
		}
		return UserParam{
			SubCMD:        subCMD,
			VCF:           vcf,
			PaperID:       paperID,
			PaperAnalysis: paperAnalysis,
		}
	default:
		fmt.Printf("[CONF] Unknown subcommand: %s\n", subCMD)
		fmt.Println(help)
		os.Exit(1)
	}
	return UserParam{}
}

func FlagsState(fs *flag.FlagSet) string {
	var flaStatus []string
	fs.VisitAll(func(f *flag.Flag) {
		flaStatus = append(flaStatus, fmt.Sprintf("--%s %s", f.Name, f.Value))
	})
	return strings.Join(flaStatus, "  ")
}
