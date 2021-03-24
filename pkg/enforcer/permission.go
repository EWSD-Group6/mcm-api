package enforcer

type Permission int

const (
	ReadUser Permission = iota
	CreateUser
	UpdateUser
	DeleteUser

	ReadFaculty
	CreateFaculty
	UpdateFaculty
	DeleteFaculty

	ReadContributeSession
	CreateContributeSession
	UpdateContributeSession
	DeleteContributeSession
	ExportContributeSession

	CreateMedia

	ReadContribution
	CreateContribution
	UpdateContribution
	DeleteContribution
	UpdateContributionStatus

	ReadComment
	CreateComment
	UpdateComment
	DeleteComment

	ReadSystemData
	UpdateSystemData

	ReadStatistic
)
