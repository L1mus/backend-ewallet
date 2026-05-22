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
		return nil, appError.InvalidPageNumber
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
	return users, totalData, totalPage, nil
}

func (s *UserService) GetTransactionReport(ctx context.Context, id int, period string) ([]dto.GetTransactionReportDTO, error) {
	data, err := s.userRepository.GetTransactionReport(ctx, id, period)
	if err != nil {
		return nil, err
	}
	var transactions []dto.GetTransactionReportDTO
	for _, transaction := range data {
		transactions = append(transactions, dto.GetTransactionReportDTO{
			Period:       transaction.Period,
			TotalIncome:  transaction.TotalIncome,
			TotalExpense: transaction.TotalExpense,
		})
	}
	return transactions, nil
}

func (s *UserService) CheckPin(ctx context.Context, id int) error {
	data, err := s.userRepository.GetPin(ctx, id)
	if err != nil {
		return err
	}
	if data.HashPin == nil {
		return appError.EmptyPin
	}
	return nil
}

func (s *UserService) GetTransactionHistory(ctx context.Context, id int, search string, limit int8, page int8) {
	//offset := (page - 1) * limit
	//data, err := s.userRepository.GetTransactionHistory(ctx, id, search, limit, offset)
	//if err != nil {
	//	return
	//}
}
