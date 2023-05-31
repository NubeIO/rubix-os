package interfaces

type GenerateViewTemplate struct {
	ViewUUID string `json:"view_uuid"`
	Name     string `json:"name"`
}

type AssignViewTemplate struct {
	ViewUUID         string `json:"view_uuid"`
	ViewTemplateUUID string `json:"view_template_uuid"`
	HostUUID         string `json:"host_uuid"`
}
