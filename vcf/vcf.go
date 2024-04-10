package vcf

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	// "fmt"
)

type VCF struct {
	Contig  string
	Pos     int
	ID      string
	Ref     string
	Alt     string
	Quality string
	Filter  string
	Format  string
	Info    map[string]string
	Samples map[string]map[string]string
	Start   int
	End     int
	EndStr  string
}

func (vcf *VCF) PrintVCF(ref *string, alt *string, info *string, sample *string) {
	fmt.Printf("%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", vcf.Contig, vcf.Pos, vcf.ID, *ref, *alt,
		vcf.Quality, vcf.Filter, *info, vcf.Format, *sample)
}

func (vcf *VCF) PrintBED() {
	fmt.Printf("%s\t%d\t%s\t%s\t%s:%s\n", vcf.Contig, vcf.Start, vcf.EndStr, vcf.ID,
		vcf.Info["SVTYPE"], vcf.Info["SVLEN"])
}

func (vcf *VCF) PrintParsedPop(sample *string) {
	fmt.Printf("%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\n", vcf.Contig, vcf.Start, vcf.EndStr,
		vcf.Info["SVTYPE"], vcf.Info["SVLEN"], vcf.Info["SUPP_VEC"], vcf.ID, *sample)
}

func (vcf *VCF) PrintSpectrePop(sample *string) {
	fmt.Printf("%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", vcf.Contig, vcf.Start, vcf.EndStr, vcf.Info["SVTYPE"],
		vcf.Info["SVLEN"], vcf.Info["SUPP_VEC"], vcf.Info["CN"], vcf.ID, *sample)
}

// FileScanner used for file reading
type FileScanner struct {
	io.Closer
	*bufio.Scanner
}

// VAFHomRef and VAFHomAlt values are currently fixed and should be changed in the future
const VAFHomRef float64 = 0.25
const VAFHomAlt float64 = 0.75

const maxCapacity = 512 * 1024

var buf = make([]byte, maxCapacity)

func HeaderRegex(VCFLine string, headerTag string) []string {
	contigRegex, err := regexp.Compile("##contig=<ID=(.*),length=([0-9]+)>")
	if err != nil {
		panic(err)
	}
	infoRegex, err := regexp.Compile("##INFO=<ID=(.*),Number=.*,Type=(.*),Description=\".*\">")
	if err != nil {
		panic(err)
	}
	formatRegex, err := regexp.Compile("##FORMAT=<ID=(.*),Number=.*,Type=.*,Description=\".*\">")
	if err != nil {
		panic(err)
	}
	switch headerTag {
	case "contig":
		return contigRegex.FindStringSubmatch(VCFLine)
	case "info":
		return infoRegex.FindStringSubmatch(VCFLine)
	case "format":
		return formatRegex.FindStringSubmatch(VCFLine)
	}
	return nil
}

func ReadVCF(VCFFile string) *FileScanner {
	// File handler
	if _, err := os.Stat(VCFFile); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("[ERROR]: %v. File %s does not exists\n", err, VCFFile)
		os.Exit(1)
	}
	VCFHandler, err := os.Open(VCFFile)
	if err != nil {
		panic(err)
	}
	// check if compressed
	isGZ := strings.Contains(VCFFile, "gz")
	if isGZ {
		VCFReadGZ, err := gzip.NewReader(VCFHandler)
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(VCFReadGZ)
		scanner.Buffer(buf, maxCapacity)
		return &FileScanner{VCFHandler, scanner}
	} else {
		scanner := bufio.NewScanner(VCFHandler)
		scanner.Buffer(buf, maxCapacity)
		return &FileScanner{VCFHandler, scanner}
	}
}

func ReadVCFStdin() *FileScanner {
	// File from stdin
	var fileIn io.Closer
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(buf, maxCapacity)
	return &FileScanner{fileIn, scanner}
}

func ReaderMaker(VCFFile string) *FileScanner {
	if VCFFile == "-" || VCFFile == "stdin" {
		return ReadVCFStdin()
	} else {
		return ReadVCF(VCFFile)
	}
}
