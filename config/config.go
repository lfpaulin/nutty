package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const version string = "0.2"

type UserParam struct {
	SubCMD         string
	VCF            string
	MinSupp        int
	MinSize        int
	MinContigLen   int
	FilerGT        string
	FilerBy        string
	FixGT          bool
	MinReadFixGT   int
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
	MinReadsMosaic int
	QualMinReads   int
	RmQC           bool
	InfoTag        string
	NoHeader       bool
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
		minReadFixGT   int
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
		minReadsMosaic int
		qualMinReads   int
		rmQC           bool
		infoTag        string
		noHeader       bool
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
		"    ghost\n" +
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
		cmdSVParse.IntVar(&minReadsMosaic, "min-mosaic-reads", 3, "Minimum number of reads to consider a SV mosaic, default = 3")
		cmdSVParse.StringVar(&filerGT, "filer-gt", "none", "Removed genotypes from output")
		cmdSVParse.StringVar(&filerBy, "filer-by", "none:none", "Filter by some parameter:value")
		cmdSVParse.BoolVar(&fixGT, "fix-gt", false, "Updated GT field based on VAF from SAMPLE column")
		cmdSVParse.IntVar(&minReadFixGT, "min-read-fix-gt", 5, "Min number of reads to update the GT")
		cmdSVParse.IntVar(&qualMinReads, "min-reads-qual", 5, "Min number of reads to be considered high quality SV")
		cmdSVParse.BoolVar(&rmQC, "rm-low-qc", false, "Remove entries with low QC")
		cmdSVParse.BoolVar(&saveRNames, "save-rnames", false, "Save RNAMES from INFO filed, if present")
		cmdSVParse.BoolVar(&noHeader, "no-header", false, "Do not print header of file")
		cmdSVParse.StringVar(&infoTag, "info-tag", "none", "Extracts tag from the INFO field")
		cmdSVParse.BoolVar(&asBED, "as-bed", false, "Output in BED format")
		cmdSVParse.BoolVar(&asDev, "as-dev", false, "Output extra logging")
		err := cmdSVParse.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
		if asDev {
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
			MinReadsMosaic: minReadsMosaic,
			FilerGT:        filerGT,
			FilerBy:        filerBy,
			FixGT:          fixGT,
			MinReadFixGT:   minReadFixGT,
			QualMinReads:   qualMinReads,
			RmQC:           rmQC,
			InfoTag:        infoTag,
			SaveRNames:     saveRNames,
			NoHeader:       noHeader,
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
		cmdPopParse.IntVar(&minReadsMosaic, "min-mosaic-reads", 3, "Minimum number of reads to consider a SV mosaic, default = 3")
		cmdPopParse.BoolVar(&uniq, "uniq", false, "Show only those that appear in a single individual (from SUPP_VEC)")
		cmdPopParse.BoolVar(&fixGT, "fix-gt", false, "Updated GT field based on VAF from SAMPLE column")
		cmdPopParse.IntVar(&minReadFixGT, "min-read-fix-gt", 5, "Min number of reads to update the GT")
		cmdPopParse.BoolVar(&fixSuppVec, "fix-suppvec", false, "Fix the SUPP_VEC based on GT/read counts")
		cmdPopParse.IntVar(&qualMinReads, "min-reads-qual", 5, "Min number of reads to be considered high quality SV")
		cmdPopParse.BoolVar(&rmQC, "rm-low-qc", false, "Remove entries with low QC")
		cmdPopParse.BoolVar(&outputVCF, "output-vcf", false, "Output is VCF")
		cmdPopParse.BoolVar(&asBED, "as-bed", false, "Output in BED format")
		cmdPopParse.BoolVar(&onlyGT, "only-gt", false, "Only prints the GT for each sample")
		cmdPopParse.BoolVar(&noHeader, "no-header", false, "Do not print header of file")
		cmdPopParse.BoolVar(&asDev, "as-dev", false, "Output extra logging")
		err := cmdPopParse.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
		if asDev {
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
			MinReadsMosaic: minReadsMosaic,
			Uniq:           uniq,
			AsBED:          asBED,
			FixSuppVec:     fixSuppVec,
			FixGT:          fixGT,
			MinReadFixGT:   minReadFixGT,
			QualMinReads:   qualMinReads,
			RmQC:           rmQC,
			OutputVCF:      outputVCF,
			OnlyGT:         onlyGT,
			NoHeader:       noHeader,
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
		if asDev {
			fmt.Printf("## CMD: nutty cancer %s \n", FlagsState(cmdCancerParse))
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
	case "ghost":
		cmdSpcPopParse := flag.NewFlagSet("pop", flag.ExitOnError)
		cmdSpcPopParse.StringVar(&vcf, "vcf", "stdin", "VCF file to read")
		cmdSpcPopParse.IntVar(&minSupp, "min-supp", 1, "Min. support for the SV calls (from SUPP_VEC), default = 1")
		cmdSpcPopParse.IntVar(&minSize, "min-size", 50, "Min. absolute size of the event (except for BDN), default = 1")
		cmdSpcPopParse.BoolVar(&uniq, "uniq", false, "Show only those that appear in a single individual (from SUPP_VEC)")
		cmdSpcPopParse.BoolVar(&asBED, "as-bed", false, "Output in BED format")
		cmdSpcPopParse.BoolVar(&onlyGT, "only-gt", false, "Only prints the GT for each sample")
		cmdSpcPopParse.BoolVar(&noHeader, "no-header", false, "Do not print header of file")
		cmdSpcPopParse.BoolVar(&asDev, "as-dev", false, "Output extra logging")
		err := cmdSpcPopParse.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
		if asDev {
			fmt.Printf("## CMD: nutty ghost %s \n", FlagsState(cmdSpcPopParse))
		}
		return UserParam{
			SubCMD:   subCMD,
			VCF:      vcf,
			MinSupp:  minSupp,
			MinSize:  minSize,
			Uniq:     uniq,
			AsBED:    asBED,
			OnlyGT:   onlyGT,
			NoHeader: noHeader,
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
		if asDev {
			fmt.Printf("## CMD: nutty paper %s \n", FlagsState(cmdPapers))
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
