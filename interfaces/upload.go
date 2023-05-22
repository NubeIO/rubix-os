package interfaces

import "mime/multipart"

type EdgeUploadResponse struct {
	Destination string `json:"destination"`
	File        string `json:"file"`
	Size        string `json:"size"`
	UploadTime  string `json:"upload_time"`
}

type FileUpload struct {
	Arch    string `json:"arch" binding:"required"`
	Version string `json:"version" binding:"required"`
	File    string `json:"file"`
}

type AppUpload struct {
	Name                            string `json:"name" binding:"required"`
	Arch                            string `json:"arch" binding:"required"`
	Version                         string `json:"version" binding:"required"`
	MoveExtractedFileToNameApp      bool   `json:"move_extracted_file_to_name_app" binding:"required"`
	MoveOneLevelInsideFileToOutside bool   `json:"move_one_level_inside_file_to_outside" binding:"required"`
}

type AppUninstall struct {
	Name    string `json:"name" binding:"required"`
	Version string `json:"version" binding:"required"`
}

type Upload struct {
	Name                            string                `json:"name"`
	Version                         string                `json:"version"`
	Product                         string                `json:"product"`
	Arch                            string                `json:"arch"`
	DoNotValidateArch               bool                  `json:"do_not_validate_arch"`
	MoveExtractedFileToNameApp      bool                  `json:"move_extracted_file_to_name_app"`
	MoveOneLevelInsideFileToOutside bool                  `json:"move_one_level_inside_file_to_outside"`
	File                            *multipart.FileHeader `json:"file"`
}

type UploadResponse struct {
	FileName     string `json:"file_name,omitempty"`
	TmpFile      string `json:"tmp_file,omitempty"`
	UploadedFile string `json:"uploaded_file,omitempty"`
}
