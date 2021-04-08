package statistic

import (
	"context"
	"gorm.io/gorm"
	"mcm-api/pkg/contributesession"
	"mcm-api/pkg/contribution"
	"mcm-api/pkg/enforcer"
	"mcm-api/pkg/user"
)

type repository struct {
	db *gorm.DB
}

func InitializeRepository(db *gorm.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r repository) countActiveUser(ctx context.Context) (int64, error) {
	var result int64
	db := r.db.WithContext(ctx).Table("users").Where("status = ?", user.UserActive).Count(&result)
	if db.Error != nil {
		return 0, nil
	}
	return result, nil
}

func (r repository) countDisableUser(ctx context.Context) (int64, error) {
	var result int64
	db := r.db.WithContext(ctx).Table("users").Where("status = ?", user.UserDisable).Count(&result)
	if db.Error != nil {
		return 0, nil
	}
	return result, nil
}

func (r repository) totalContributions(ctx context.Context) (int64, error) {
	var result int64
	db := r.db.WithContext(ctx).Model(&contribution.Entity{}).Count(&result)
	if db.Error != nil {
		return 0, nil
	}
	return result, nil
}

func (r repository) totalContributeSessions(ctx context.Context) (int64, error) {
	var result int64
	db := r.db.WithContext(ctx).Model(&contributesession.Entity{}).Count(&result)
	if db.Error != nil {
		return 0, nil
	}
	return result, nil
}

type countByRole struct {
	Role  enforcer.Role
	Total int64
}

func (r repository) countUserByRole(ctx context.Context) (map[enforcer.Role]int64, error) {
	var result []*countByRole
	db := r.db.WithContext(ctx).Model(&user.Entity{}).Select("role, count(id) as total").
		Group("role").Find(&result)
	if db.Error != nil {
		return nil, db.Error
	}
	var resultMap = map[enforcer.Role]int64{
		enforcer.Guest:                0,
		enforcer.Student:              0,
		enforcer.MarketingCoordinator: 0,
		enforcer.MarketingManager:     0,
		enforcer.Administrator:        0,
	}
	for _, v := range result {
		resultMap[v.Role] = v.Total
	}
	return resultMap, nil
}

func (r repository) countContributionGroupByFaculty(ctx context.Context, sessionId *int, status *contribution.Status) ([]*FacultyContributionData, error) {
	var result []*FacultyContributionData
	db := r.db.WithContext(ctx).Model(&contribution.Entity{}).
		Select("faculties.id as id, faculties.name as name, count(contributions.id) as count").
		Joins("left join users on contributions.user_id = users.id").
		Joins("left join faculties on users.faculty_id = faculties.id").
		Group("faculties.id")
	if status != nil {
		db.Where("contributions.status = ?", *status)
	}
	if sessionId != nil {
		db.Where("contributions.contribute_session_id = ?", *sessionId)
	}
	db = db.Find(&result)
	if db.Error != nil {
		return nil, db.Error
	}
	return result, nil
}

func (r repository) countContributionGroupByStudent(ctx context.Context, sessionId *int, status *contribution.Status) ([]*ContributionStudentData, error) {
	var result []*ContributionStudentData
	db := r.db.WithContext(ctx).Model(&contribution.Entity{}).
		Select("users.id, users.name, users.email, count(*) as count").
		Joins("left join users on contributions.user_id = users.id").
		Group("users.id").
		Order("count desc").
		Limit(100)
	if status != nil {
		db.Where("contributions.status = ?", *status)
	}
	if sessionId != nil {
		db.Where("contributions.contribute_session_id = ?", *sessionId)
	}
	db = db.Find(&result)
	if db.Error != nil {
		return nil, db.Error
	}
	return result, nil
}
