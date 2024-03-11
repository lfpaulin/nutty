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
    contigRegex, err := regexp.Compile("##contig=<ID=(.*),length=([0-9]+)>")
    if err != nil { panic(err) }
    infoRegex, err := regexp.Compile("##INFO=<ID=(.*),Number=.*,Type=.*,Description=\".*\">")
    if err != nil { panic(err) }
    formatRegex, err := regexp.Compile("<ID=(.*),Number=.*,Type=.*,Description=\".*\">")
    if err != nil { panic(err) }
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
            lineSplit := strings.Split(line, "\t")
            sampleName := lineSplit[9]
            // Here goes the parser header
            if !params.AsBED{
                fmt.Println("##Sample name: ", sampleName)
                fmt.Println(vcf.HeaderOut)
            }
        case strings.Contains(line, "#"):
            //
        default:
            // each entry
            ReadVCFEntry(line, params)
        }
    }
}


func ReadVCFEntry(VCFLineRaw string, userParams *config.UserParam) {
    var (
        dr int
        dv int
        vaf float64
    )
    lineSplit := strings.Split(VCFLineRaw, "\t")
    VCFLineFormated := new(vcf.VCF)
    VCFLineFormated.Contig = lineSplit[0]
    VCFPosInt, err := strconv.Atoi(lineSplit[1])
    if err != nil {
        panic(err)
    }
    VCFLineFormated.Pos = VCFPosInt 
    VCFLineFormated.ID = lineSplit[2]
    VCFLineFormated.Ref = lineSplit[3]
    VCFLineFormated.Alt = lineSplit[4]
    VCFLineFormated.Quality = lineSplit[5]
    VCFLineFormated.Filter = lineSplit[6]
    // split each key=value pair or flag
    info := make(map[string]string)
    for _, infoElem := range strings.Split(lineSplit[7], ";") {
        if strings.Contains(infoElem, "="){
            info_key_val := strings.Split(infoElem, "=")
            info[info_key_val[0]] = info_key_val[1]

        } else {
            info[infoElem] = "flag"
        }
    }
    VCFLineFormated.Info = info
    // we expect only one sample here, so we only use one
    sampleSV := make(map[string]string)
    formatSplit := strings.Split(lineSplit[8], ":")
    sampleSVSplit := strings.Split(lineSplit[9], ":")
    for idx := range formatSplit {
        sampleSV[formatSplit[idx]] = sampleSVSplit[idx]
    }
    VCFLineFormated.Samples = sampleSV
    // TODO: here go the output from the parser
    // #CONTTIG\tSTART\tEND\tSVTYPE\tSVLEN\tGT\tAF\tREFC\tALTC\tID
    dr, err = strconv.Atoi(VCFLineFormated.Samples["DR"])
    dv, err = strconv.Atoi(VCFLineFormated.Samples["DV"])
    vaf = float64(dv)/float64(dr+dv)*100
    fmt.Printf("%s\t%d\t%s\t%s\t%s\t%s\t%0.3f\t%d\t%d\t%s\n", VCFLineFormated.Contig, VCFLineFormated.Pos, 
        VCFLineFormated.Info["END"], VCFLineFormated.Info["SVTYPE"], VCFLineFormated.Info["SVLEN"],
        VCFLineFormated.Samples["GT"], vaf, dr, dv, VCFLineFormated.ID)
}