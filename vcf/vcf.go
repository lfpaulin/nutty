package vcf

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

/*
  INFO example
	"SVTYPE": "DEL",
	"SVLEN": "-100",

  SAMPLES example
	"SAMPLE1": {
	}
*/
