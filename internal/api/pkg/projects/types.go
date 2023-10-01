package projects

type Project struct {
	Name   string `json:"name,omitempty"`
	WMSurl string `json:"wms,omitempty"`
}
