package vcf

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"
	"strings"
	// "fmt"
)

/*
  INFO example
    "SVTYPE": "DEL",
    "SVLEN": "-100",

  SAMPLES example
    "SAMPLE1": {
    }
*/

type VCF struct {
    Contig  string
    Pos     int
    ID      string
    Ref     string
    Alt     string
    Quality string
    Filter  string
    Info    map[string]string
    Samples map[string]string
}

// For file reading
type FileScanner struct {
    io.Closer
    *bufio.Scanner
}


const HeaderOut string = "#CONTTIG\tSTART\tEND\tSVTYPE\tSVLEN\tGT\tVAF\tREFC\tALTC\tID"
const maxCapacity = 512*1024
var buf = make([]byte, maxCapacity)

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
        // scanner.Buffer(buf, maxCapacity)
        return &FileScanner{VCFHandler, scanner}
    } else {
        scanner := bufio.NewScanner(VCFHandler)
        // scanner.Buffer(buf, maxCapacity)
        return &FileScanner{VCFHandler, scanner}
    }
}
