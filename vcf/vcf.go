package vcf

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"
	"regexp"
	"strings"
	// "fmt"
)

// VAFHomRef and VAFHomAlt values are currently fixed and should be changed in the future
const VAFHomRef float64 = 0.25
const VAFHomAlt float64 = 0.75

type VCF struct {
	Contig  string
	Pos     int
	ID      string
	Ref     string
	Alt     string
	Quality string
	Filter  string
	Info    map[string]string
	Samples map[string]map[string]string
}

// FileScanner used for file reading
type FileScanner struct {
	io.Closer
	*bufio.Scanner
}

const HeaderOut string = "#CONTTIG\tSTART\tEND\tSVTYPE\tSVLEN\tGT\tVAF\tREFC\tALTC\tID"
const maxCapacity = 512 * 1024

var buf = make([]byte, maxCapacity)

func HeaderRegex(VCFLine string, headerTag string) []string {
	contigRegex, err := regexp.Compile("##contig=<ID=(.*),length=([0-9]+)>")
	if err != nil {
		panic(err)
	}
	infoRegex, err := regexp.Compile("##INFO=<ID=(.*),Number=.*,Type=.*,Description=\".*\">")
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
