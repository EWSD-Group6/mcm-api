package statistic

import (
	"context"
	"mcm-api/pkg/contributesession"
	"mcm-api/pkg/enforcer"
)

type Service struct {
	repository               *repository
	contributeSessionService *contributesession.Service
}

func InitializeService(
	repository *repository,
	contributeSessionService *contributesession.Service,
) *Service {
	return &Service{
		repository:               repository,
		contributeSessionService: contributeSessionService,
	}
}

func (s Service) adminDashboard(ctx context.Context) (*AdminDashboard, error) {
	byRole, err := s.repository.countUserByRole(ctx)
	if err != nil {
		return nil, err
	}
	countActive, err := s.repository.countActiveUser(ctx)
	if err != nil {
		return nil, err
	}
	countDisableUser, err := s.repository.countDisableUser(ctx)
	if err != nil {
		return nil, err
	}
	totalSessions, err := s.repository.totalContributeSessions(ctx)
	if err != nil {
		return nil, err
	}
	totalContributions, err := s.repository.totalContributions(ctx)
	if err != nil {
		return nil, err
	}
	return &AdminDashboard{
		ActiveUserCount:           countActive,
		DisableUserCount:          countDisableUser,
		TotalContribution:         totalContributions,
		TotalContributeSession:    totalSessions,
		StudentCount:              byRole[enforcer.Student],
		MarketingManagerCount:     byRole[enforcer.MarketingManager],
		MarketingCoordinatorCount: byRole[enforcer.MarketingCoordinator],
		GuestCount:                byRole[enforcer.Guest],
	}, nil
}

func (s Service) contributionFacultyChart(ctx context.Context, query *ContributionFacultyChartQuery) (*ContributionFacultyChart, error) {
	session, err := s.contributeSessionService.GetCurrentSession(ctx)
	if err != nil {
		return &ContributionFacultyChart{
			Session: nil,
			Data:    []*FacultyContributionData{},
		}, nil
	}
	groupByFaculty, err := s.repository.countContributionGroupByFaculty(ctx, &session.Id, query.Status)
	if err != nil {
		return nil, err
	}
	return &ContributionFacultyChart{
		Session: &Session{
			Id:               session.Id,
			OpenTime:         session.OpenTime,
			ClosureTime:      session.ClosureTime,
			FinalClosureTime: session.FinalClosureTime,
		},
		Data: groupByFaculty,
	}, nil
}

func (s Service) contributionStudentChart(ctx context.Context, query *ContributionStudentChartQuery) (*ContributionStudentChart, error) {
	session, err := s.contributeSessionService.GetCurrentSession(ctx)
	if err != nil {
		return &ContributionStudentChart{
			Session: nil,
			Data:    []*ContributionStudentData{},
		}, nil
	}
	groupBy, err := s.repository.countContributionGroupByStudent(ctx, &session.Id, query.Status)
	if err != nil {
		return nil, err
	}
	return &ContributionStudentChart{
		Session: &Session{
			Id:               session.Id,
			OpenTime:         session.OpenTime,
			ClosureTime:      session.ClosureTime,
			FinalClosureTime: session.FinalClosureTime,
		},
		Data: groupBy,
	}, nil
}
