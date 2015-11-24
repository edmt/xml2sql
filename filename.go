package main

import (
	"fmt"
	"strings"
)

func buildFilename(d Doc) string {
	return fmt.Sprintf("%s_%s%s_%s.xml",
		d.Emisor.RFC,
		d.Serie,
		d.Complemento.TimbreFiscalDigital.UUID,
		strings.Replace(substr(d.Fecha, 0, 10), "-", "", -1),
	)
}

func buildDirectoryPath(d Doc) string {
	return strings.Replace(substr(d.Fecha, 0, 10), "-", "/", -1)
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}
