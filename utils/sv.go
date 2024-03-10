package utils

import (
	"fmt"
	"regexp"
	"sniffles2_helper_go/config"
	"sniffles2_helper_go/vcf"
	"strconv"
	"strings"
)


var (
    infoVCF []string
    formatVCF []string
)
var contigsVCF = make(map[string]int)

func ParseSV(params *config.UserParam) {
    VCFReader := vcf.ReadVCF(params.VCF)
    defer VCFReader.Close()
    // header metadata needed
    contigRegex, _ := regexp.Compile("##contig=<ID=(.*),length=([0-9]+)>")
    infoRegex, _ := regexp.Compile("##INFO=<ID=(.*),Number=.*,Type=.*,Description=\".*\">")
    formatRegex, _ := regexp.Compile("<ID=(.*),Number=.*,Type=.*,Description=\".*\">")
    // contig size limit
    minContigSize := 1000000
    for VCFReader.Scan() {
        line := strings.TrimSpace(VCFReader.Text())
        switch {
        case strings.Contains(line, "##") && strings.Contains(line, "contig"):
            contigMatch := contigRegex.FindStringSubmatch(line)
            contigName := contigMatch[1]
            contigSize, err := strconv.Atoi(contigMatch[2])
            if err != nil {
                panic(err)
            }
            if int(contigSize) > minContigSize {
                contigsVCF[contigName] = contigSize
            }
        case strings.Contains(line, "##") && strings.Contains(line, "INFO"):
            infoMatch := infoRegex.FindStringSubmatch(line)
            infoVCF = append(infoVCF, infoMatch[1])
        case strings.Contains(line, "##") && strings.Contains(line, "FORMAT"):
            formatMatch := formatRegex.FindStringSubmatch(line)
            formatVCF = append(formatVCF, formatMatch[1])
        case strings.Contains(line, "#CHROM"):
            // ignore
        default:
            // each entry
        }
    }
}
