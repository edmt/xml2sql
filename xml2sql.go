package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"archive/zip"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Static("/app", "./app")
	router.POST("/upload", func(c *gin.Context) {
		empresaId := c.Request.FormValue("empresaID")
		origin    := c.Request.FormValue("origin")

        file, fileHeader, _ := c.Request.FormFile("file")
        if strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".zip") {
        	content, _ := ioutil.ReadAll(file)
        	ioutil.WriteFile("./tmp/" + fileHeader.Filename, content, 0777)
        	r, _ := zip.OpenReader("./tmp/" + fileHeader.Filename)
        	defer r.Close()

			for _, f := range r.File {
				if strings.HasSuffix(strings.ToLower(f.Name), ".xml") {
				    rc, _ := f.Open()
				    xml, _ := ioutil.ReadAll(rc)
				    processXml(empresaId, origin, xml)
				}
			}

        } else {
        	xml, _ := ioutil.ReadAll(file)
        	processXml(empresaId, origin, xml)
        }

	})
	router.Run(":8080")
}

func processXml(empresaId, origin string, xml []byte) {
	// Save with new name
	cfdi := parseXml(xml)
	newName := buildFilename(cfdi)
	path := "output/" + buildDirectoryPath(cfdi)
	os.MkdirAll(path, 0777)
	ioutil.WriteFile(path+"/"+newName, xml, 0777)

	// Append to pool.csv file
	includeHeader := false
	if _, err := os.Stat("output/pool.csv"); os.IsNotExist(err) {
		includeHeader = true
	}
	poolf, _ := os.OpenFile("output/pool.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0660)
	if includeHeader {
		poolf.WriteString("RFCEmisor,RFCReceptor,UUID\n")
	}
	poolf.WriteString(cfdi.Emisor.RFC + "," + cfdi.Receptor.RFC + "," + cfdi.Complemento.TimbreFiscalDigital.UUID + "\n")
	defer poolf.Close()

	// Append to inserts.sql file
	insertsf, _ := os.OpenFile("output/inserts.sql", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0660)
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
