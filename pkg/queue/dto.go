package queue

import "mcm-api/pkg/common"

type ContributionCreatedPayload struct {
	ContributionId int                 `json:"contributionId"`
	UserId         int                 `json:"userId"`
	UserName       string              `json:"userName"`
	FacultyId      int                 `json:"facultyId"`
	User           common.LoggedInUser `json:"user"`
}

type ArticleUploadedMessage struct {
	Message
}

type ArticleUploadedPayload struct {
	ArticleId int                 `json:"articleId" mapstructure:"articleId"`
	Link      string              `json:"link"`
	User      common.LoggedInUser `json:"user"`
}
