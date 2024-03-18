package utils

var (
    infoVCF     []string
    formatVCF   []string
    sampleName  string
    sampleNames []string
    info        map[string]string
    sampleSV    map[string]map[string]string
)
var contigsVCF = make(map[string]int)
