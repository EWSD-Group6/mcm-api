package article

import (
	"mcm-api/pkg/common"
)

type ArticleRes struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Versions    []Version `json:"versions"`
	common.TrackTime
}

type ArticleReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
}
