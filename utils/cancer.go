package utils

import (
	"bufio"
	"cmd/vendor/golang.org/x/sys/unix"
	"compress/gzip"
	"log"
	"os"
	"sniffles2_helper_go/config"
	"strings"
)

func ParseCancer(params *config.UserParam) {
	println("hey ", params.SubCMD)
	VCFHandler, err := os.Open(params.VCF)
	if err != nil {
		log.Fatal(err)
	}
	defer VCFHandler.Close()
	isGZ := strings.Contains(params.VCF, "gz")
	if isGZ {
		VCFRead, err := gzip.NewReader(VCFHandler)
		if err != nil {
			log.Fatal(err)
		}
		defer VCFRead.Close()
	} else {
		VCFRead := bufio.NewScanner(VCFHandler)
		if err != nil {
			log.Fatal(err)
		}
		defer VCFRead.Close()
	}
}

func ReaderGZ(*os.File fileHandler){

}
