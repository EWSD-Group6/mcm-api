package media

import (
	"io"
	"mcm-api/pkg/enforcer"
)

type UploadType string

const (
	Document UploadType = "document"
	Image    UploadType = "image"
)

type UploadResult struct {
	Key string `json:"key"`
}

type FileUploadOriginalReq struct {
	File io.ReadSeeker
	Size int64
	Name string
	User *enforcer.LoggedInUser
}

type FileUploadPreviewReq struct {
	File io.Reader
	Name string
	User enforcer.LoggedInUser
}

type UploadQuery struct {
	Type UploadType `query:"type" enums:"document,image"`
}

type ContributionUploadReq struct {
	File                io.Reader
	ContributeSessionId int
}
