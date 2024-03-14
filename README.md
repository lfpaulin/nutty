# Nutty - a [Sniffles](https://github.com/fritzsedlazeck/Sniffles) companion app

## Required programs (go)

Setup a conda environment for Nutty
```bash
conda create --name nutty python=3.10
conda install anaconda::go
conda activate nutty
```
Alternatively, you can use download go from [https://go.dev/](https://go.dev/)

## How to run

Nutty can be used as without compilation using the following command:
```bash
go run main.go help
```
Or it can be compiled first
```bash
go run build
```
and then run

```bash
nutty help
```

---

## Examples
### Run Nutty simple sample parser


```bash
# not compressed VCF
nutty sv --vcf sv.vcf

# compressed VCF
nutty sv --vcf sv.vcf.gz

# to get help with the parameters
nutty sv --help

```


### Papers
In order to get reproducible results we included a `paper` subcomand that takes as input:
- Paper id
- VCF input
- Analysis flag
  and will output the output used for the analysis

```bash
nutty paper --cancer-t2t --vcf colo829.vcf.gz --colo829
```
This will output:
- Supplementary Table 3B if input is colo829_grch38.vcf.gz
- Supplementary Table 8B if input is colo829_t2t.vcf.gz
- Supplementary Table 11B if input is colo829_lifted.vcf.gz
