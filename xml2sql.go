package main

import (
	"archive/zip"
	"crypto/sha1"
	"fmt"
	// "github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func processFromFileSystem() {
	db := readCSV("input/empresas.csv")
	log.Print(db)
	matches, _ := filepath.Glob("input/*.zip")
	for _, zip := range matches {
		processZip(zip, db)
	}
}

func processZip(path string, db map[string]CSVRecord) {
	r, err := zip.OpenReader(path)
	logger(err)
	defer r.Close()

	for _, f := range r.File {
		if strings.HasSuffix(strings.ToLower(f.Name), ".xml") {
			rc, err := f.Open()
			logger(err)
			xml, err := ioutil.ReadAll(rc)
			logger(err)
			processXml(xml, f.Name, db)
		}
	}
}

func main() {
	processFromFileSystem()
	// gin.SetMode(gin.ReleaseMode)
	// log.Print(os.Getwd())
	// router := gin.Default()
	// router.Static("/app", "./app")
	// router.POST("/upload", func(c *gin.Context) {
	// 	empresaId := c.Request.FormValue("empresaID")
	// 	origin := c.Request.FormValue("origin")

	// 	file, fileHeader, err := c.Request.FormFile("file")
	// 	logger(err)
	// 	if strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".zip") {
	// 		content, err := ioutil.ReadAll(file)
	// 		logger(err)
	// 		ioutil.WriteFile("./tmp/"+fileHeader.Filename, content, 0777)
	// 		r, err := zip.OpenReader("./tmp/" + fileHeader.Filename)
	// 		logger(err)
	// 		defer r.Close()

	// 		for _, f := range r.File {
	// 			if strings.HasSuffix(strings.ToLower(f.Name), ".xml") {
	// 				rc, err := f.Open()
	// 				logger(err)
	// 				xml, err := ioutil.ReadAll(rc)
	// 				logger(err)
	// 				processXml(empresaId, origin, xml)
	// 			}
	// 		}

	// 	} else {
	// 		xml, err := ioutil.ReadAll(file)
	// 		logger(err)
	// 		processXml(empresaId, origin, xml)
	// 	}

	// })
	// router.Run(":8080")
}

func logger(err error) {
	if err != nil {
		log.Print(err)
	}
}

func processXml(xml []byte, name string, db map[string]CSVRecord) {
	cfdi := parseXml(xml)

	empresaId := db[cfdi.Emisor.RFC].Id
	ambient := db[cfdi.Emisor.RFC].Ambient
	origin := getOrigin(ambient, cfdi)
	refId := getRefId(cfdi)

	newName := name
	path := "./tmp/" + buildDirectoryPath(cfdi)
	log.Printf("Creating directory: %s\n", path)
	err := os.MkdirAll(path, 0777)
	logger(err)
	err = ioutil.WriteFile(path+"/"+newName, xml, 0777)
	logger(err)

	// Append to pool.csv file
	includeHeader := false
	if _, err := os.Stat("./output/pool-" + ambient + ".csv"); os.IsNotExist(err) {
		includeHeader = true
	}
	poolf, err := os.OpenFile("./output/pool-" + ambient + ".csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0660)
	logger(err)
	if includeHeader {
		poolf.WriteString("RFCEmisor,RFCReceptor,UUID\n")
	}
	poolf.WriteString(cfdi.Emisor.RFC + "," + cfdi.Receptor.RFC + "," + cfdi.Complemento.TimbreFiscalDigital.UUID + "\n")
	defer poolf.Close()

	// Append to inserts.sql file
	insertsf, err := os.OpenFile("./output/inserts-" + ambient + ".sql", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0660)
	logger(err)
	insertsf.WriteString(formatAsInsert(cfdi, name, refId, empresaId, origin, getHash(path+"/"+newName)))
	defer insertsf.Close()
}

func getOrigin(ambient string, cfdi Doc) string {
	if ambient == "BF" {
		return `'Portal Web'`
	}
	if ambient == "OX" {
		return `'WebServiceCFDi'`
	}
	if ambient == "NF" {
		return `'WebServiceCFDi'`
	}
	if ambient == "TF" {
		if strings.ToUpper(cfdi.Complemento.Nomina.XMLName.Local) == "NOMINA" {
			return `'timbreNomina'`
		} else {
			return `'Timbrado'`
		}
	}
	if ambient == "CF" {
		if strings.ToUpper(cfdi.Complemento.Nomina.XMLName.Local) == "NOMINA" {
			return `'emisionNomina'`
		} else {
			return `'WebServiceCFDi'`
		}
	}
	return "NULL"
}

func getRefId(cfdi Doc) string {
	if strings.ToUpper(cfdi.Addenda.AddendaBuzonFiscal.CFD.XMLName.Local) == "CFD" && cfdi.Addenda.AddendaBuzonFiscal.CFD.RefID != "" {
		log.Print("Con refId")
		return "'" + cfdi.Addenda.AddendaBuzonFiscal.CFD.RefID + "'"
	} else {
		return "REPLACE(newid(),'-','')"
	}
}

func getHash(path string) string {
	app := "xsltproc"

	arg0 := "./app/xslt/cfd/_32/cadenaoriginal_3_2.xslt"

	out, err := exec.Command(app, arg0, path).Output()
	if err != nil {
		log.Fatal(err)
	}

	h := sha1.New()
	io.WriteString(h, string(out))
	return strings.Replace(fmt.Sprintf("% x", h.Sum(nil)), " ", "", -1)
}
