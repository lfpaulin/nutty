package vcf

import (
    "compress/gzip"
    "io"
    "bufio"
    "log"
    "os"
    "strings"
)



type VCFInfo struct {
    Info map[string]string
}

type VCFSamples struct {
    VCFSample map[string]map[string]string
}

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
    Info    VCFInfo
    Samples VCFSamples
}

// For file reading
type FileScanner struct {
    io.Closer
    *bufio.Scanner
}

func ReadVCF(VCFFile string) *FileScanner {
    // File handler
    VCFHandler, err := os.Open(VCFFile)
    if err != nil {
        log.Fatal(err)
    }
    defer VCFHandler.Close()
    // check if compressed
    isGZ := strings.Contains(VCFFile, "gz")
    if isGZ {
        if err != nil {
            log.Fatal(nil)
        }
        VCFReadGZ, err := gzip.NewReader(VCFHandler)
        if err != nil {
            log.Fatal(err)
        }
        // line by line
        scanner := bufio.NewScanner(VCFReadGZ)
        return &FileScanner{VCFHandler, scanner}
    } else {
        scanner := bufio.NewScanner(VCFHandler)
        return &FileScanner{VCFHandler, scanner}
    }
}
