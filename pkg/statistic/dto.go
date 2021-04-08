package statistic

import (
	"mcm-api/pkg/contribution"
	"time"
)

type AdminDashboard struct {
	ActiveUserCount           int64 `json:"activeUserCount"`
	DisableUserCount          int64 `json:"disableUserCount"`
	TotalContribution         int64 `json:"totalContribution"`
	TotalContributeSession    int64 `json:"totalContributeSession"`
	StudentCount              int64 `json:"studentCount"`
	MarketingManagerCount     int64 `json:"marketingManagerCount"`
	MarketingCoordinatorCount int64 `json:"marketingCoordinatorCount"`
	GuestCount                int64 `json:"guestCount"`
}

type FacultyContributionData struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type Session struct {
	Id               int       `json:"id"`
	OpenTime         time.Time `json:"openTime"`
	ClosureTime      time.Time `json:"closureTime"`
	FinalClosureTime time.Time `json:"finalClosureTime"`
}

type ContributionFacultyChart struct {
	Session *Session                   `json:"session"`
	Data    []*FacultyContributionData `json:"data"`
}

type ContributionFacultyChartQuery struct {
	Status *contribution.Status `query:"status" enums:"accepted,reviewing,rejected"`
}

type ContributionStudentChartQuery struct {
	Status *contribution.Status `query:"status" enums:"accepted,reviewing,rejected"`
}

type ContributionStudentData struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Count int    `json:"count"`
}

type ContributionStudentChart struct {
	Session *Session                   `json:"session"`
	Data    []*ContributionStudentData `json:"data"`
}

type ContributionSessionChart struct {
	Data []struct {
		Id                int       `json:"id"`
		OpenTime          time.Time `json:"openTime"`
		ClosureTime       time.Time `json:"closureTime"`
		FinalClosureTime  time.Time `json:"finalClosureTime"`
		ContributionCount int       `json:"contributionCount"`
	} `json:"data"`
}

type StudentSessionChart struct {
	Data []struct {
		Id               int       `json:"id"`
		OpenTime         time.Time `json:"openTime"`
		ClosureTime      time.Time `json:"closureTime"`
		FinalClosureTime time.Time `json:"finalClosureTime"`
		StudentCount     int       `json:"contributionCount"`
	} `json:"data"`
}

type StudentFacultyChart struct {
	Data []struct {
		FacultyName  string `json:"facultyName"`
		StudentCount int    `json:"studentCount"`
	} `json:"data"`
}
