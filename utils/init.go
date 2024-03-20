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

var (
	dr        int
	dv        int
	vaf       float64
	gt        string
	vafString string
)

const (
	fixHetVAFMin float64 = 0.33
	fixAltVAFMin float64 = 0.66
	minCoverage  int     = 5 // do not trust anything below this number of reads
)
