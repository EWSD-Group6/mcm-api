package notification

import (
	"mcm-api/config"
	"testing"
)

func TestService_SendNewContributionEmail(t *testing.T) {
	service := InitializeService(&config.Config{SesSenderEmail: "noreply@devstack.cloud"})
	err := service.SendNewContributionEmail(&Destination{ToAddresses: []string{
		"superquanganh@gmail.com",
	}}, &TemplateNewContributionPayLoad{
		Name:        "Teacher 1",
		StudentName: "User 1",
		Link:        "https://google.com",
	})
	if err != nil {
		t.Error(err)
	}
}
