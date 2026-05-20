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

func (s *UserService) FindReceiver(ctx context.Context, id int, search string, limit int, page int) ([]dto.FindReceiverDTO, error) {
	if page <= 0 || limit <= 0 {
		return nil, appError.InvalidEmailFormat
	}

	offset := (page - 1) * limit
	data, err := s.userRepository.FindReceiver(ctx, id, search, limit, offset)
	if err != nil {
		return nil, err
	}
	var users []dto.FindReceiverDTO
	for _, user := range data {
		users = append(users, dto.FindReceiverDTO{
			Id:                user.Id,
			FullName:          user.FullName,
			Phone:             user.Phone,
			ProfilePictureUrl: user.ProfilePictureUrl,
			IsVerified:        user.IsVerified,
		})
	}
	return users, nil
}
