package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/L1mus/backend-ewallet/internal/appError"
	"github.com/L1mus/backend-ewallet/internal/dto"
	"github.com/L1mus/backend-ewallet/internal/repository"
	"github.com/L1mus/backend-ewallet/pkg"
	"github.com/jackc/pgx/v5"
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
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.GetUserProfileDTO{}, appError.UserNotFound
		}
		return dto.GetUserProfileDTO{}, err
	}
	return dto.GetUserProfileDTO{
		Id:                data.Id,
		FullName:          data.FullName,
		Email:             data.Email,
		Phone:             data.Phone,
		ProfilePictureURL: data.ProfilePictureURL,
	}, nil
}

func (s *UserService) GetUserDashboard(ctx context.Context, id int) (dto.GetUserDashboardDTO, error) {
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

func (s *UserService) FindReceiver(ctx context.Context, id int, req dto.ReceiverQuery) ([]dto.FindReceiverDTO, dto.PaginationMetaData, error) {
	data, err := s.userRepository.GetReceiver(ctx, id, req)
	if err != nil {
		return nil, dto.PaginationMetaData{}, err
	}

	if len(data) == 0 {
		return []dto.FindReceiverDTO{}, dto.PaginationMetaData{}, nil
	}

	totalData := data[0].TotalCount
	limit := 10
	totalPage := int(math.Ceil(float64(totalData) / float64(limit)))

	page, err := strconv.Atoi(req.Page)
	if err != nil {
		return nil, dto.PaginationMetaData{}, err
	}

	var users []dto.FindReceiverDTO
	prevLink := fmt.Sprintf("users/?search=%s&page=%s", req.Search, page-1)
	nextLink := fmt.Sprintf("users/?search=%s&page=%s", req.Search, page+1)
	for _, user := range data {
		users = append(users, dto.FindReceiverDTO{
			Id:                user.Id,
			FullName:          user.FullName,
			Phone:             user.Phone,
			ProfilePictureUrl: user.ProfilePictureUrl,
			IsVerified:        user.IsVerified,
		})
	}
	metaDataPAgination := dto.PaginationMetaData{
		TotalPages: totalPage,
		TotalData:  totalData,
		NextLink:   nextLink,
		PrevLink:   prevLink,
	}
	return users, metaDataPAgination, nil
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

func (s *UserService) GetTransactionHistory(ctx context.Context, id int, req dto.TransactionHistoryQuery) ([]dto.GetTransactionHistoryDTO, dto.PaginationMetaData, error) {
	data, err := s.userRepository.GetTransactionHistory(ctx, id, req)
	if err != nil {
		return nil, dto.PaginationMetaData{}, err
	}

	if len(data) == 0 {
		return []dto.GetTransactionHistoryDTO{}, dto.PaginationMetaData{}, nil
	}

	totalData := data[0].TotalCount
	limit := 10
	totalPage := int(math.Ceil(float64(totalData) / float64(limit)))

	page, err := strconv.Atoi(req.Page)
	if err != nil {
		return nil, dto.PaginationMetaData{}, err
	}

	var users []dto.GetTransactionHistoryDTO
	prevLink := fmt.Sprintf("transactions/?search=%s&page=%s", req.Search, page-1)
	nextLink := fmt.Sprintf("transactions/?search=%s&page=%s", req.Search, page+1)
	for _, user := range data {
		users = append(users, dto.GetTransactionHistoryDTO{
			TransactionID:     user.TransactionID,
			Amount:            user.Amount,
			Type:              user.Type,
			ActivityType:      user.ActivityType,
			Status:            user.Status,
			CreatedAt:         user.CreatedAt,
			Description:       user.Description,
			ReceiverName:      user.ReceiverName,
			PaymentMethodName: user.PaymentMethodName,
		})
	}
	metaDataPAgination := dto.PaginationMetaData{
		TotalPages: totalPage,
		TotalData:  totalData,
		NextLink:   nextLink,
		PrevLink:   prevLink,
	}
	return users, metaDataPAgination, nil
}

func (s *UserService) EditProfile(ctx context.Context, id int, req dto.EditProfileRequest) error {
	if req.Phone != nil && *req.Phone != "" {
		taken, err := s.userRepository.CheckPhoneTaken(ctx, id, *req.Phone)
		if err != nil {
			return err
		}
		if taken {
			return appError.PhoneAlreadyExists
		}
	}

	return s.userRepository.UpdateProfile(ctx, id, req)
}

func (s *UserService) EditPin(ctx context.Context, id int, req dto.EditPinRequest) error {
	data, err := s.userRepository.GetPin(ctx, id)
	if err != nil {
		return err
	}

	if data.HashPin == nil {
		return appError.EmptyPin
	}

	var hc pkg.HashConfig
	if err := hc.Compare(req.CurrentPin, *data.HashPin); err != nil {
		return appError.WrongPin
	}

	hc.UseRecommended()
	hashNewPin := hc.GenHash(req.NewPin)

	return s.userRepository.UpdatePin(ctx, id, hashNewPin)
}

func (s *UserService) EditPassword(ctx context.Context, id int, req dto.EditPasswordRequest) error {
	data, err := s.userRepository.GetHashPassword(ctx, id)
	if err != nil {
		return err
	}

	var hc pkg.HashConfig
	if err := hc.Compare(req.CurrentPassword, data.HashPassword); err != nil {
		return appError.WrongPassword
	}

	hc.UseRecommended()
	hashNewPassword := hc.GenHash(req.NewPassword)

	return s.userRepository.UpdatePassword(ctx, id, hashNewPassword)
}

func (s *UserService) UploadProfilePicture(ctx context.Context, userID int, fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader.Size > 2*1024*1024 {
		return "", appError.FileTooLarge
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := make([]byte, 512)
	if _, err := file.Read(buf); err != nil {
		return "", err
	}

	mimeType := http.DetectContentType(buf)

	allowedTypes := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
	}
	ext, allowed := allowedTypes[mimeType]
	if !allowed {
		return "", appError.FileTypeNotAllowed
	}

	oldProfile, err := s.userRepository.GetUserProfile(ctx, userID)
	if err == nil && oldProfile.ProfilePictureURL != nil && *oldProfile.ProfilePictureURL != "" {
		oldPath := filepath.Join("public", *oldProfile.ProfilePictureURL)
		err := os.Remove(oldPath)
		if err != nil {
			return "", err
		}
	}

	filename := fmt.Sprintf("user_%d_%d%s", userID, time.Now().UnixNano(), ext)
	savePath := filepath.Join("public", "img", "profiles", filename)

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	dst, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	publicURL := fmt.Sprintf("/img/profiles/%s", filename)
	if err := s.userRepository.UpdateProfilePictureURL(ctx, userID, publicURL); err != nil {
		os.Remove(savePath)
		return "", err
	}

	return publicURL, nil
}
