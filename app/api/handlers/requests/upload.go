package requests

type UploadRequest struct {
	Name string `json:"name" form:"name"`
}
