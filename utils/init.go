package utils

var (
	formatVCF   []string
	sampleName  string
	sampleNames []string
	info        map[string]string
	sampleSV    map[string]map[string]string
)
var contigsVCF = make(map[string]int)
var infoVCF = make(map[string]string)

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
	// #CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO	FORMAT
	indexChrom   int = 0
	indexPos     int = 1
	indexID      int = 2
	indexRef     int = 3
	indexAlt     int = 4
	indexQual    int = 5
	indexFilter  int = 6
	indexInfo    int = 7
	indexFormat  int = 8
	indexSamples int = 9
)
