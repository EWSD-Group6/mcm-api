package contribution

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"mcm-api/pkg/common"
	"mcm-api/pkg/enforcer"
)

type IndexQuery struct {
	common.PaginateQuery
	FacultyId             *int   `json:"facultyId"`
	StudentId             *int   `json:"studentId"`
	ContributionSessionId *int   `json:"contributionSessionId"`
	Status                Status `json:"status" enums:"accepted,rejected,reviewing"`
}

type ContributionRes struct {
	Id                  int     `json:"id"`
	User                UserRes `json:"user"`
	ContributeSessionId int     `json:"contributeSessionId"`
	ArticleId           *int    `json:"articleId"`
	Title               string  `json:"title"`
	Description         string  `json:"description"`
	Status              Status  `json:"status"`
	common.TrackTime
}

type UserRes struct {
	Id        int           `json:"id"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	FacultyId *int          `json:"facultyId"`
	Role      enforcer.Role `json:"role"`
}

type ImageRes struct {
	Key   string `json:"key"`
	Title string `json:"title"`
	Link  string `json:"link"`
}

type ArticleReq struct {
	Link string `json:"link"`
}

func (r *ArticleReq) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Link, validation.Required, validation.Length(10, 255)),
	)
}

type ContributionCreateReq struct {
	Article     *ArticleReq      `json:"article"`
	Images      []ImageCreateReq `json:"images"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
}

func (r *ContributionCreateReq) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Article, validation.Required.When(r.Images == nil)),
		validation.Field(&r.Images, validation.Required.When(r.Article == nil)),
		validation.Field(&r.Title, validation.Required, validation.Length(10, 255)),
		validation.Field(&r.Description, validation.Length(15, 512)),
	)
}

type ImageCreateReq struct {
	Key   string `json:"key"`
	Title string `json:"title"`
}

func (c ImageCreateReq) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Key, validation.Required),
		validation.Field(&c.Title, validation.Length(0, 200)),
	)
}

type ContributionUpdateReq struct {
	Article     *ArticleReq      `json:"article"`
	Images      []ImageCreateReq `json:"images"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
}

func (r *ContributionUpdateReq) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Article, validation.Required.When(r.Images == nil)),
		validation.Field(&r.Images, validation.Required.When(r.Article == nil)),
		validation.Field(&r.Title, validation.Required, validation.Length(10, 255)),
		validation.Field(&r.Description, validation.Length(15, 512)),
	)
}

type ContributionStatusReq struct {
	Status Status `json:"status"`
}

func (c ContributionStatusReq) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Status,
			validation.Required,
			validation.In(Accepted, Reviewing, Rejected),
		),
	)
}

type PaginateComposition struct {
	common.PaginateResponse
	Data []ContributionRes `json:"data"`
}
