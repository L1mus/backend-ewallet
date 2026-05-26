package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
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

func (s *UserService) FindReceiver(ctx context.Context, id int, req dto.PageQuery) ([]dto.FindReceiverDTO, dto.PaginationMetaData, error) {
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
	var page int

	if req.Page == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(req.Page)
		if err != nil {
			return nil, dto.PaginationMetaData{}, err
		}
	}

	var users []dto.FindReceiverDTO
	var prevLink string
	var nextLink string
	if page == 1 {
		prevLink = ""
	} else {
		prevLink = fmt.Sprintf("http://localhost:8080/users/transfer?search=%s&page=%d", req.Search, page-1)

	}
	if page == totalPage {
		nextLink = ""
	} else {
		nextLink = fmt.Sprintf("http://localhost:8080/users/transfer?search=%s&page=%d", req.Search, page+1)
	}
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

func (s *UserService) GetTransactionHistory(ctx context.Context, id int, req dto.PageQuery) ([]dto.GetTransactionHistoryDTO, dto.PaginationMetaData, error) {
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
	var page int

	if req.Page == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(req.Page)
		if err != nil {
			return nil, dto.PaginationMetaData{}, err
		}
	}

	var users []dto.GetTransactionHistoryDTO
	var prevLink string
	var nextLink string
	if page == 1 {
		prevLink = ""
	} else {
		prevLink = fmt.Sprintf("http://localhost:8080/users/transfer?search=%s&page=%d", req.Search, page-1)

	}
	if page == totalPage {
		nextLink = ""
	} else {
		nextLink = fmt.Sprintf("http://localhost:8080/users/transfer?search=%s&page=%d", req.Search, page+1)
	}
	for _, user := range data {
		users = append(users, dto.GetTransactionHistoryDTO{
			TransactionID:     user.TransactionID,
			Amount:            user.Amount,
			ProfilePictureUrl: user.ProfilePictureUrl,
			Phone:             user.Phone,
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
	/*
		validasi ukuran
		buat format file apa saja yang bisa di unggah
		cek format content
		hapus source lama
		buat format penamaan file yang di unggah dan
		save ke folder yang dibuat
		buat file path gambar
	*/

	if fileHeader.Size > 1*1024*1024 {
		return "", appError.FileTooLarge
	}
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Println("close error :", err)
		}
	}(file)

	buf := make([]byte, 512)
	if _, err := file.Read(buf); err != nil {
		return "", err
	}

	mimeType := http.DetectContentType(buf)

	allowedTypes := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/webp": ".webp",
	}
	ext, allowed := allowedTypes[mimeType]
	if !allowed {
		return "", appError.FileTypeNotAllowed
	}

	oldProfile, err := s.userRepository.GetUserProfile(ctx, userID)
	log.Println("old file", oldProfile.ProfilePictureURL, "errornya apa ya?", err)
	if err == nil && oldProfile.ProfilePictureURL != nil && *oldProfile.ProfilePictureURL != "" {
		oldPath := filepath.Join("public", *oldProfile.ProfilePictureURL)
		err := os.Remove(oldPath)
		if err != nil {
			return "", err
		}
	}

	filename := fmt.Sprintf("user_%d_%d%s", userID, time.Now().UnixNano(), ext)
	savePath := filepath.Join("public", "img", filename)

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	dst, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			log.Println("close error :", err)
		}
	}(dst)
	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	publicURL := fmt.Sprintf("/img/profiles/%s", filename)
	if err := s.userRepository.UpdateProfilePictureURL(ctx, userID, publicURL); err != nil {
		err := os.Remove(savePath)
		if err != nil {
			return "", err
		}
		return "", err
	}

	return publicURL, nil
}
