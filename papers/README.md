# Paper
Analysis to reproduce results from papers that use Sniffles2

## Cancer/T2T

We hypothesized that use of a completed human reference genome (CHM13-T2T) would improve somatic SV calling. 
We assessed the current SV benchmark set for COLO829/BL across four replicates sequenced at different centers with different long-read technologies. 
Our work demonstrates new approaches to optimize somatic SV prioritization in cancer with potential improvements in other genetic diseases.

### COLO829 somatic SV filtering

```bash
go run main.go paper --vcf colo829_<version>.vcf.gz --paper-id cancer_t2t --paper-analysis colo829
```
This will output the first draft to create the following tables:
- Supplementary Table 3B if input is colo829_grch38.vcf.gz
- Supplementary Table 8B if input is colo829_t2t.vcf.gz
- Supplementary Table 11B if input is colo829_lifted.vcf.gz

### POG sampls somatic SV filtering

```bash
go run main.go paper --vcf pog.vcf.gz --paper-id cancer_t2t --paper-analysis pog
```
This will output the first draft to create the following tables:
- Supplementary Table 18 if input is pog044.vcf.gz
- Supplementary Table 19 if input is pog1022.vcf.gz
