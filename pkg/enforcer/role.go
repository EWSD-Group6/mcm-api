package enforcer

type Role string

const (
	Administrator        Role = "admin"
	MarketingManager     Role = "marketing_manager"
	MarketingCoordinator Role = "marketing_coordinator"
	Student              Role = "student"
	Guest                Role = "guest"
)

var permissionsOfRole map[Role][]Permission

func init() {
	permissionsOfRole = make(map[Role][]Permission)

	addPermissions(Administrator,
		ReadUser,
		UpdateUser,
		CreateUser,
		DeleteUser,

		ReadFaculty,
		UpdateFaculty,
		CreateFaculty,
		DeleteFaculty,

		ReadContributeSession,
		UpdateContributeSession,
		CreateContributeSession,
		DeleteContributeSession,

		ReadSystemData,
		UpdateSystemData,
	)

	addPermissions(MarketingManager,
		ReadContribution,
		ReadContributeSession,
		ExportContributeSession,
		ReadStatistic,

		ReadFaculty,

		ReadContributeSession,
	)

	addPermissions(MarketingCoordinator,
		ReadContribution,
		UpdateContributionStatus,

		ReadComment,
		UpdateComment,
		CreateComment,
		DeleteComment,

		ReadFaculty,

		ReadContributeSession,
	)

	addPermissions(Student,
		ReadContribution,
		UpdateContribution,
		CreateContribution,
		DeleteContribution,

		CreateMedia,

		ReadComment,
		UpdateComment,
		CreateComment,
		DeleteComment,

		ReadSystemData,

		ReadFaculty,

		ReadContributeSession,
	)

	addPermissions(Guest, ReadContribution)
}

func addPermissions(role Role, permission ...Permission) {
	for _, p := range permission {
		permissionsOfRole[role] = append(permissionsOfRole[role], p)
	}
}

func Can(role Role, permission Permission) bool {
	for _, v := range permissionsOfRole[role] {
		if v == permission {
			return true
		}
	}
	return false
}
