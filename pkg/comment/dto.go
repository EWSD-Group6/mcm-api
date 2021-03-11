package comment

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"mcm-api/pkg/common"
	"mcm-api/pkg/user"
	"time"
)

type IndexQuery struct {
	ContributionId int `query:"contributionId" validate:"required"`
	common.CursorQuery
}

type CommentRes struct {
	Id      string            `json:"id"`
	User    user.UserResponse `json:"user"`
	Content string            `json:"content"`
	Edited  bool              `json:"edited"`
	common.TrackTime
}

type CommentCreateReq struct {
	ContributionId int    `json:"contributionId"`
	Content        string `json:"content"`
}

func (c CommentCreateReq) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ContributionId, validation.Required),
		validation.Field(&c.Content, validation.Required, validation.Length(0, 500)),
	)
}

type CommentUpdateReq struct {
	Content string `json:"content"`
}

func (c CommentUpdateReq) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Content, validation.Required, validation.Length(0, 500)),
	)
}

type CursorPayload struct {
	Id        string
	CreatedAt time.Time
}
