package main

import (
	"fmt"
	"strings"
)

func formatAsInsert(d Doc, name, refId, empresaId, origen, hash string) string {
	const template = `
insert into CFD (
    idInternal,        emisor,
    expedido,          fechaCancelacion,
    fechaGeneracion,   fechaValidacion,
    folio,             idRemision,
    idSucursal,        montoTotal,
    noCertificado,     nombreArchivo,
    razonNoValido,     razonSocial,
    rfc,               serie,
    tipoDeComprobante, tipoDocumento,
    valido,            version,
    vigente,           leido,
    Empresa_Id,        descuento,
    subTotal,          totalImpuestosRetenidos,
    totalIVATrasladado,refID,
    totalImpuestosTrasladados, InfoAduanera_id,
    descargas,         numTimbre,
    fechaTimbrado,     origen,
    Moneda_id,         hash,
    AcuseCancela_id,   EstatusCfd_id,
    Validacion_id,     fecha_entregaBF,
    Unidad_Negocio_Id, descripcionComercial
)
values (
     %s,   '%s',
     %s,    %s,
    '%s',   %s,
    '%s',   %s,
     %s,    %s,
    '%s',  '%s',
     %s,   '%s',
    '%s',   %s,
    '%s',  '%s',
     %s,   '%s',
     %s,   %s,
    '%s',  %s,
     %s,   %s,
     %s,   %s,
     %s,   %s,
     %s,  '%s',
    '%s',  %s,
     %s,  '%s',
     %s,   %s,
     %s,   %s,
     %s,   %s,
);

`
	return fmt.Sprintf(template,
		"REPLACE(newid(),'-','')", d.Emisor.Nombre,
		"1", "null",
		strings.Replace(d.Fecha, "T", " ", -1), "null",
		d.Folio, "null",
		"null", d.Total,
		d.NoCertificado, "/"+buildDirectoryPath(d)+"/"+name,
		"null", d.Receptor.Nombre,
		d.Receptor.RFC, "null",
		d.Tipo, d.Tipo,
		"0", d.Version,
		"1", "0",
		empresaId, "null",
		d.SubTotal, "null",
		"null", refId,
		d.Impuestos.Total, "null",
		"0", d.Complemento.TimbreFiscalDigital.UUID,
		strings.Replace(d.Complemento.TimbreFiscalDigital.FechaTimbrado, "T", " ", -1), origen,
		"null", hash,
		"null", "null",
		"null", "null",
		"null", "null",
	)
}
