package article

import (
	"mcm-api/pkg/common"
	"time"
)

type ArticleRes struct {
	Id       int           `json:"id"`
	Versions []*VersionRes `json:"versions"`
	common.TrackTime
}

type VersionRes struct {
	Id              int       `json:"id"`
	Hash            string    `json:"hash"`
	ArticleId       int       `json:"articleId"`
	LinkOriginal    string    `json:"linkOriginal,omitempty"`
	LinkOriginalCdn string    `json:"linkOriginalCdn,omitempty"`
	LinkPdf         string    `json:"linkPdf,omitempty"`
	LinkPdfCdn      string    `json:"linkPdfCdn,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
}

type ArticleReq struct {
	Link string `json:"link"`
}
