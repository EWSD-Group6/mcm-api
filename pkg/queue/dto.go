package queue

import "mcm-api/pkg/common"

type ContributionCreatedPayload struct {
	ContributionId int64               `json:"contributionId"`
	UserId         int64               `json:"userId"`
	UserName       string              `json:"userName"`
	FacultyId      int64               `json:"facultyId"`
	User           common.LoggedInUser `json:"user"`
}

type ArticleUploadedMessage struct {
	Message
}

type ArticleUploadedPayload struct {
	ArticleId int                 `json:"articleId"`
	Link      string              `json:"link"`
	User      common.LoggedInUser `json:"user"`
}
