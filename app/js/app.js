var configureDropzone = function() {
  Dropzone.options.myAwesomeDropzone = {
    paramName: "file",
    maxFilesize: 10, // MB
    uploadMultiple: false,
    parallelUploads: 10,
    acceptedFiles: ".xml,.XML,.zip,.ZIP",
    url: "/upload",
    dictDefaultMessage: 'Arrastra aquí o selecciona los archivos',
    dictFallbackMessage: 'Tu navegador no soporta drag & drop',
    dictInvalidFileType: 'Tipo de archivo no soportado',
    dictFileTooBig: 'Archivo demasiado grande',
    dictResponseError: 'Ha ocurrido un error',
    dictMaxFilesExceeded: 'Demasiados archivos simultáneos',
    init: function() {
      this.on('drop', function(e) {
        this.removeAllFiles();
      });
      
      this.on('sending', function(file, xhr, formData) {
        formData.append("empresaID", $('#empresaId').val());
        formData.append("origin", $('#origin').val());
      })
    }
  }
};

$(document).ready(function() {
  configureDropzone();
});
