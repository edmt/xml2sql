package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Doc struct {
	XMLName         xml.Name        `xml:"Comprobante"`
	Tipo            string          `xml:"tipoDeComprobante,attr"`
	Version         string          `xml:"version,attr"`
	Serie           string          `xml:"serie,attr"`
	Folio           string          `xml:"folio,attr"`
	Fecha           string          `xml:"fecha,attr"`
	Moneda          string          `xml:"Moneda,attr"`
	TipoCambio      string          `xml:"TipoCambio,attr"`
	Total           string          `xml:"total,attr"`
	SubTotal        string          `xml:"subTotal,attr"`
	MetodoDePago    string          `xml:"metodoDePago,attr"`
	LugarExpedicion string          `xml:"LugarExpedicion,attr"`
	NoCertificado   string          `xml:"noCertificado,attr"`
	Emisor          CFDIEmisor      `xml:"Emisor"`
	Receptor        CFDIReceptor    `xml:"Receptor"`
	Conceptos       []CFDIConcepto  `xml:"Conceptos>Concepto"`
	Impuestos       CFDIImpuestos   `xml:"Impuestos"`
	Complemento     CFDIComplemento `xml:"Complemento"`
	Addenda         CFDIAddenda     `xml:"Addenda"`
}

type CFDIImpuestos struct {
	XMLName   xml.Name      `xml:"Impuestos"`
	Total     string        `xml:"totalImpuestosTrasladados,attr"`
	Traslados CFDITraslados `xml:"Traslados"`
}

type CFDITraslados struct {
	XMLName  xml.Name     `xml:"Traslados"`
	Traslado CFDITraslado `xml:"Traslado"`
}

type CFDITraslado struct {
	XMLName xml.Name `xml:"Traslado"`
	Importe string   `xml:"importe,attr"`
}

type CFDIAddenda struct {
	XMLName            xml.Name               `xml:"Addenda"`
	AddendaBuzonFiscal AddendaBuzonFiscalNode `xml:"AddendaBuzonFiscal"`
}

type AddendaBuzonFiscalNode struct {
	XMLName xml.Name `xml:"AddendaBuzonFiscal"`
	CFD     CFDNode  `xml:"CFD"`
}

type CFDNode struct {
	XMLName xml.Name `xml:"CFD"`
	RefID   string   `xml:"refID,attr"`
}

type CFDIEmisor struct {
	XMLName         xml.Name                  `xml:"Emisor"`
	RFC             string                    `xml:"rfc,attr"`
	Nombre          string                    `xml:"nombre,attr"`
	DomicilioFiscal EmisorDomicilioFiscalNode `xml:"DomicilioFiscal"`
}

type EmisorDomicilioFiscalNode struct {
	XMLName   xml.Name `xml:"DomicilioFiscal"`
	Municipio string   `xml:"municipio,attr"`
	Estado    string   `xml:"estado,attr"`
}

type CFDIReceptor struct {
	XMLName xml.Name `xml:"Receptor"`
	RFC     string   `xml:"rfc,attr"`
	Nombre  string   `xml:"nombre,attr"`
}

type CFDIConcepto struct {
	XMLName          xml.Name `xml:"Concepto"`
	Descripcion      string   `xml:"descripcion,attr"`
	NoIdentificacion string   `xml:"noIdentificacion,attr"`
	Cantidad         string   `xml:"cantidad,attr"`
	Unidad           string   `xml:"unidad,attr"`
	ValorUnitario    string   `xml:"valorUnitario,attr"`
	Importe          string   `xml:"importe,attr"`
}

type CFDIComplemento struct {
	XMLName             xml.Name               `xml:"Complemento"`
	TimbreFiscalDigital TFDTimbreFiscalDigital `xml:"TimbreFiscalDigital"`
	Nomina              NominaNomina           `xml:"Nomina"`
}

type TFDTimbreFiscalDigital struct {
	XMLName           xml.Name `xml:"TimbreFiscalDigital"`
	NumeroCertificado string   `xml:"noCertificadoSAT,attr"`
	FechaTimbrado     string   `xml:"FechaTimbrado,attr"`
	UUID              string   `xml:"UUID,attr"`
}

type NominaNomina struct {
	XMLName          xml.Name `xml:"Nomina"`
	FechaInicialPago string   `xml:"FechaInicialPago,attr"`
	FechaFinalPago   string   `xml:"FechaFinalPago,attr"`
}

func (d Doc) NumeroDeFactura() string {
	return fmt.Sprintf("%s-%s", d.Serie, d.Folio)
}

func (tfd TFDTimbreFiscalDigital) FechaTimbre() string {
	layout := "2006-01-02T15:04:05"
	t, err := time.Parse(layout, tfd.FechaTimbrado)

	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

func parseXml(doc []byte) Doc {
	var query Doc
	xml.Unmarshal(doc, &query)
	return query
}

func clean(value string) string {
	return strings.Replace(strings.Replace(value, "\t", "", -1), "\n", "", -1)
}

func EncodeAsRows(path string) []string {
	file, err := os.Open(path)

	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	rawContent, _ := ioutil.ReadAll(file)
	cfdi := parseXml(rawContent)

	var records []string
	var record = []string{
		clean(cfdi.Complemento.TimbreFiscalDigital.NumeroCertificado),
		clean(cfdi.Emisor.RFC),
		clean(cfdi.Emisor.Nombre),
		clean(cfdi.Emisor.DomicilioFiscal.Municipio),
		clean(cfdi.Emisor.DomicilioFiscal.Estado),
		clean(cfdi.Receptor.RFC),
		clean(cfdi.Receptor.Nombre),
		clean(cfdi.LugarExpedicion),
		clean(cfdi.Complemento.TimbreFiscalDigital.FechaTimbre()),
		clean(cfdi.Total),
		clean(cfdi.Moneda),
		clean(cfdi.Complemento.TimbreFiscalDigital.UUID),
	}
	records = append(records, strings.Join(record, "\t"))
	return records
}

func EncodeHeaders() string {
	var headerList = []string{
		"Certificado",
		"EmisorRFC",
		"EmisorRazonSocial",
		"EmisorMunicipio",
		"EmisorEstado",
		"ReceptorRFC",
		"ReceptorRazonSocial",
		"LugarDeExpedicion",
		"FechaTimbrado",
		"MontoTotal",
		"Moneda",
		"UUID",
	}
	return strings.Join(headerList, "\t")
}
