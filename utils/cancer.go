package utils

import (
    "fmt"
    "strings"
    "sniffles2_helper_go/config"
    "sniffles2_helper_go/vcf"
)


func ParseCancer(params *config.UserParam) {
    VCFReader := vcf.ReadVCF(params.VCF)
    defer VCFReader.Close()
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
