package service

import (
	"github.com/golang/mock/gomock"
	_repo "github.com/lantonster/askme/internal/repo"
	_service "github.com/lantonster/askme/internal/service"
	mockrepo "github.com/lantonster/askme/mock/repo"
	mockservice "github.com/lantonster/askme/mock/service"
)

var (
	mockActivityService *mockservice.MockActivityService
	mockAuthService     *mockservice.MockAuthService
	mockConfigService   *mockservice.MockConfigService
	mockEmailService    *mockservice.MockEmailService
	mockSiteInfoService *mockservice.MockSiteInfoService
	mockUploadsService  *mockservice.MockUploadsService
	mockUserService     *mockservice.MockUserService

	mockActivityRepo *mockrepo.MockActivityRepo
	mockAuthRepo     *mockrepo.MockAuthRepo
	mockConfigRepo   *mockrepo.MockConfigRepo
	mockEmailRepo    *mockrepo.MockEmailRepo
	mockSiteInfoRepo *mockrepo.MockSiteInfoRepo
	mockUserRepo     *mockrepo.MockUserRepo
)

func setupMock(ctrl *gomock.Controller) *_repo.Repo {
	mockActivityService = mockservice.NewMockActivityService(ctrl)
	mockAuthService = mockservice.NewMockAuthService(ctrl)
	mockConfigService = mockservice.NewMockConfigService(ctrl)
	mockEmailService = mockservice.NewMockEmailService(ctrl)
	mockSiteInfoService = mockservice.NewMockSiteInfoService(ctrl)
	mockUploadsService = mockservice.NewMockUploadsService(ctrl)
	mockUserService = mockservice.NewMockUserService(ctrl)
	_service.NewService(
		mockActivityService,
		mockAuthService,
		mockConfigService,
		mockEmailService,
		mockSiteInfoService,
		mockUploadsService,
		mockUserService,
	)

	mockActivityRepo = mockrepo.NewMockActivityRepo(ctrl)
	mockAuthRepo = mockrepo.NewMockAuthRepo(ctrl)
	mockConfigRepo = mockrepo.NewMockConfigRepo(ctrl)
	mockEmailRepo = mockrepo.NewMockEmailRepo(ctrl)
	mockSiteInfoRepo = mockrepo.NewMockSiteInfoRepo(ctrl)
	mockUserRepo = mockrepo.NewMockUserRepo(ctrl)
	return _repo.NewRepo(
		mockActivityRepo,
		mockAuthRepo,
		mockConfigRepo,
		mockEmailRepo,
		mockSiteInfoRepo,
		mockUserRepo,
	)
}
