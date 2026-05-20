package service

import (
	"context"

	"github.com/L1mus/backend-ewallet/internal/appError"
	"github.com/L1mus/backend-ewallet/internal/dto"
	"github.com/L1mus/backend-ewallet/internal/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetUserProfile(ctx context.Context, id int) (dto.GetUserProfileDTO, error) {
	data, err := s.userRepository.GetUserProfile(ctx, id)
	if err != nil {
		return dto.GetUserProfileDTO{}, appError.UserNotFound
	}
	return dto.GetUserProfileDTO{
		Id:                data.Id,
		FullName:          data.FullName,
		Email:             data.Email,
		Phone:             data.Phone,
		ProfilePictureURL: data.ProfilePictureURL,
	}, nil
}

func (s *UserService) GetUserDashboad(ctx context.Context, id int) (dto.GetUserDashboardDTO, error) {
	data, err := s.userRepository.GetUserDashboard(ctx, id)
	if err != nil {
		return dto.GetUserDashboardDTO{}, err
	}
	return dto.GetUserDashboardDTO{
		Balance:       data.Balance,
		TotalIncome:   data.TotalIncome,
		TotalExpenses: data.TotalExpenses,
	}, nil
}
