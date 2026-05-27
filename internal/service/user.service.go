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
	"github.com/L1mus/backend-ewallet/internal/cache"
	"github.com/L1mus/backend-ewallet/internal/dto"
	"github.com/L1mus/backend-ewallet/internal/repository"
	"github.com/L1mus/backend-ewallet/pkg"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	userRepository *repository.UserRepository
	rdb            *redis.Client
}

func NewUserService(userRepository *repository.UserRepository, rdb *redis.Client) *UserService {
	return &UserService{
		userRepository: userRepository,
		rdb:            rdb,
	}
}

func (s *UserService) GetUserProfile(ctx context.Context, id int) (dto.GetUserProfileDTO, error) {
	rkey := fmt.Sprintf("user:profile:%d", id)

	var cachedProfile dto.GetUserProfileDTO

	found, err := cache.GetFromCache(ctx, s.rdb, rkey, &cachedProfile)
	if err == nil && found {
		log.Printf("Cache HIT for key: %s Retrieving data from Redis...", rkey)
		return cachedProfile, nil
	}
	if err != nil {
		log.Printf("Redis error while retrieving cache: %v", err)
	}

	log.Printf("Cache MISS for key: %s Fetching data from DB...", rkey)

	data, err := s.userRepository.GetUserProfile(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.GetUserProfileDTO{}, appError.UserNotFound
		}
		return dto.GetUserProfileDTO{}, err
	}

	profileDTO := dto.GetUserProfileDTO{
		Id:                data.Id,
		FullName:          data.FullName,
		Email:             data.Email,
		Phone:             data.Phone,
		ProfilePictureURL: data.ProfilePictureURL,
	}

	cacheTTL := 1 * time.Hour
	err = cache.SaveToCache(ctx, s.rdb, rkey, profileDTO, cacheTTL)
	if err != nil {
		log.Printf("Failed to save data: %v", err)
	}

	return profileDTO, nil
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
	page := 1

	if req.Page != "" {
		if p, err := strconv.Atoi(req.Page); err == nil && p > 0 {
			page = p
		}
	}

	prevLink := ""
	nextLink := ""
	if page > 1 {
		prevLink = fmt.Sprintf("http://localhost:8080/users/transfer?search=%s&page=%d", req.Search, page-1)
	}
	if page < totalPage {
		nextLink = fmt.Sprintf("http://localhost:8080/users/transfer?search=%s&page=%d", req.Search, page+1)
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
	metaDataPAgination := dto.PaginationMetaData{
		TotalPages: totalPage,
		TotalData:  totalData,
		NextLink:   nextLink,
		PrevLink:   prevLink,
	}
	return users, metaDataPAgination, nil
}

func (s *UserService) GetTransactionReport(ctx context.Context, id int, period string) (dto.GetTransactionReportResponse, error) {
	type periodConfig struct {
		timePeriode string
		startDate   time.Time
		labelFormat string
	}

	now := time.Now()
	configs := map[string]periodConfig{
		"week": {
			timePeriode: "day",
			startDate:   now.AddDate(0, 0, -6),
			labelFormat: "Mon",
		},
		"month": {
			timePeriode: "day",
			startDate:   time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()),
			labelFormat: "2",
		},
		"year": {
			timePeriode: "month",
			startDate:   time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()),
			labelFormat: "Jan",
		},
	}

	cfg, ok := configs[period]
	if !ok {
		return dto.GetTransactionReportResponse{}, appError.InvalidPeriod
	}

	data, err := s.userRepository.GetTransactionReport(ctx, id, cfg.timePeriode, cfg.startDate)
	if err != nil {
		return dto.GetTransactionReportResponse{}, err
	}

	var totalIncome float32
	var totalExpense float32
	var chartData []dto.GetTransactionReportDTO

	for _, t := range data {
		label := ""
		if t.Period != nil {
			label = t.Period.Format(cfg.labelFormat)
		}
		totalIncome += t.TotalIncome
		totalExpense += t.TotalExpense

		chartData = append(chartData, dto.GetTransactionReportDTO{
			Period:       label,
			TotalIncome:  t.TotalIncome,
			TotalExpense: t.TotalExpense,
		})
	}

	if chartData == nil {
		chartData = []dto.GetTransactionReportDTO{}
	}

	return dto.GetTransactionReportResponse{
		Period: period,
		Summary: dto.ReportSummaryDTO{
			TotalIncome:  totalIncome,
			TotalExpense: totalExpense,
			NetAmount:    totalIncome - totalExpense,
		},
		Data: chartData,
	}, nil
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
	page := 1

	if req.Page != "" {
		if p, err := strconv.Atoi(req.Page); err == nil && p > 0 {
			page = p
		}
	}

	prevLink := ""
	nextLink := ""
	if page > 1 {
		prevLink = fmt.Sprintf("http://localhost:8080/users/transactions?search=%s&page=%d", req.Search, page-1)
	}
	if page < totalPage {
		nextLink = fmt.Sprintf("http://localhost:8080/users/transactions?search=%s&page=%d", req.Search, page+1)
	}

	var transactions []dto.GetTransactionHistoryDTO
	for _, transaction := range data {
		item := dto.GetTransactionHistoryDTO{
			TransactionID: transaction.TransactionID,
			Amount:        transaction.Amount,
			Type:          transaction.Type,
			ActivityType:  transaction.ActivityType,
			Status:        transaction.Status,
		}
		if transaction.ReceiverName != nil {
			item.ReceiverName = *transaction.ReceiverName
		}
		if transaction.Phone != nil {
			item.Phone = *transaction.Phone
		}
		if transaction.ProfilePictureUrl != nil {
			item.ProfilePictureUrl = *transaction.ProfilePictureUrl
		}
		transactions = append(transactions, item)
	}
	metaDataPAgination := dto.PaginationMetaData{
		TotalPages: totalPage,
		TotalData:  totalData,
		NextLink:   nextLink,
		PrevLink:   prevLink,
	}
	return transactions, metaDataPAgination, nil
}

func (s *UserService) EditProfile(ctx context.Context, id int, req dto.EditProfileRequest, fileHeader *multipart.FileHeader) error {
	newProfileURL := req.ProfilePictureURL
	if fileHeader != nil {
		/*
			validasi ukuran
			buat format file apa saja yang bisa di unggah
			cek format content
			hapus source lama
			buat format penamaan file yang di unggah dan
			save ke folder yang dibuat
			buat file path gambar
		*/
		if fileHeader.Size > 2*1024*1024 {
			return appError.FileTooLarge
		}

		file, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer func(file multipart.File) {
			_ = file.Close()
		}(file)

		buf := make([]byte, 512)
		if _, err := file.Read(buf); err != nil {
			return err
		}
		mimeType := http.DetectContentType(buf)

		allowedTypes := map[string]string{
			"image/jpeg": ".jpg",
			"image/png":  ".png",
			"image/webp": ".webp",
		}
		ext, allowed := allowedTypes[mimeType]
		if !allowed {
			return appError.FileTypeNotAllowed
		}

		oldProfile, err := s.userRepository.GetUserProfile(ctx, id)
		if err == nil && oldProfile.ProfilePictureURL != nil && *oldProfile.ProfilePictureURL != "" {
			oldPath := filepath.Join("public", *oldProfile.ProfilePictureURL)
			_ = os.Remove(oldPath)
		}

		filename := fmt.Sprintf("user_%d_%d%s", id, time.Now().UnixNano(), ext)
		savePath := filepath.Join("public", "img", filename)

		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return err
		}

		dst, err := os.Create(savePath)
		if err != nil {
			return err
		}
		defer func(dst *os.File) {
			_ = dst.Close()
		}(dst)

		if _, err := io.Copy(dst, file); err != nil {
			return err
		}

		newProfileURL = new(fmt.Sprintf("/img/%s", filename))
	}

	req.ProfilePictureURL = newProfileURL

	err := s.userRepository.UpdateProfile(ctx, id, req)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("user:profile:%d", id)
	_ = cache.DelFromCache(ctx, s.rdb, cacheKey)

	return nil
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
