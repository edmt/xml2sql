package main

import (
	"archive/zip"
	"crypto/sha1"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	log.Print(os.Getwd())
	router := gin.Default()
	router.Static("/app", "./app")
	router.POST("/upload", func(c *gin.Context) {
		empresaId := c.Request.FormValue("empresaID")
		origin := c.Request.FormValue("origin")

		file, fileHeader, err := c.Request.FormFile("file")
		logger(err)
		if strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".zip") {
			content, err := ioutil.ReadAll(file)
			logger(err)
			ioutil.WriteFile("./tmp/"+fileHeader.Filename, content, 0777)
			r, err := zip.OpenReader("./tmp/" + fileHeader.Filename)
			logger(err)
			defer r.Close()

			for _, f := range r.File {
				if strings.HasSuffix(strings.ToLower(f.Name), ".xml") {
					rc, err := f.Open()
					logger(err)
					xml, err := ioutil.ReadAll(rc)
					logger(err)
					processXml(empresaId, origin, xml)
				}
			}

		} else {
			xml, err := ioutil.ReadAll(file)
			logger(err)
			processXml(empresaId, origin, xml)
		}

	})
	router.Run(":8080")
}

func logger(err error) {
	if err != nil {
		log.Print(err)
	}
}

func processXml(empresaId, origin string, xml []byte) {
	cfdi := parseXml(xml)
	newName := buildFilename(cfdi)
	path := "./output/" + buildDirectoryPath(cfdi)
	log.Printf("Creating directory: %s\n", path)
	err := os.MkdirAll(path, 0777)
	logger(err)
	err = ioutil.WriteFile(path + "/" + newName, xml, 0777)
	logger(err)

	// Append to pool.csv file
	includeHeader := false
	if _, err := os.Stat("./output/pool.csv"); os.IsNotExist(err) {
		includeHeader = true
	}
	poolf, err := os.OpenFile("./output/pool.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0660)
	logger(err)
	if includeHeader {
		poolf.WriteString("RFCEmisor,RFCReceptor,UUID\n")
	}
	poolf.WriteString(cfdi.Emisor.RFC + "," + cfdi.Receptor.RFC + "," + cfdi.Complemento.TimbreFiscalDigital.UUID + "\n")
	defer poolf.Close()

	// Append to inserts.sql file
	insertsf, err := os.OpenFile("./output/inserts.sql", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0660)
	logger(err)
	insertsf.WriteString(formatAsInsert(cfdi, empresaId, origin, getHash(path+"/"+newName)))
	defer insertsf.Close()
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
