package request


type CreateDocumentRequest struct {
	VisaApplicationID string `json:"visa_application_id" validate:"required,len=26"`
	FileType          string `json:"file_type" validate:"required,oneof=passport_photo international_passport bank_statement signed_form transcript cv sop"`
	FilePath          string `json:"file_path" validate:"required,url"`
}