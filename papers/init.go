package papers

var (
	dr          int
	dv          int
	vaf         float64
	vafString   string
	gt          string
	infoVCF     []string
	formatVCF   []string
	sampleNames []string
	info        map[string]string
	sampleSV    map[string]map[string]string
)

var contigsVCF = make(map[string]int)
