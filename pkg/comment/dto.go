package comment

import (
	"mcm-api/pkg/common"
	"mcm-api/pkg/user"
	"time"
)

type IndexQuery struct {
	ContributionId int64 `query:"contributionId" validate:"required"`
	common.CursorQuery
}

type CommentRes struct {
	Id      string            `json:"id"`
	User    user.UserResponse `json:"user"`
	Content string            `json:"content"`
	common.TrackTime
}

type CommentCreateReq struct {
	ContributionId int64  `json:"contributionId"`
	Content        string `json:"content"`
}

type CommentUpdateReq struct {
	Content string `json:"content"`
}

type CursorPayload struct {
	Id        string
	CreatedAt time.Time
}
