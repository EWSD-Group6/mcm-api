package queue

import (
	"mcm-api/pkg/enforcer"
)

type ContributionCreatedPayload struct {
	ContributionId int                   `json:"contributionId"`
	UserId         int                   `json:"userId"`
	UserName       string                `json:"userName"`
	FacultyId      int                   `json:"facultyId"`
	User           enforcer.LoggedInUser `json:"user"`
}

type ArticleUploadedPayload struct {
	ArticleId int                   `json:"articleId"`
	Link      string                `json:"link"`
	User      enforcer.LoggedInUser `json:"user"`
}

type ExportContributeSessionPayload struct {
	ContributeSessionId int `json:"contributeSessionId"`
}
