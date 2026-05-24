package main

import (
	"context"
	"fmt"
	"log"

	"github.com/L1mus/backend-ewallet/internal/config"
	"github.com/L1mus/backend-ewallet/pkg"
	"github.com/joho/godotenv"
)

// helper agar bisa tulis phone sebagai pointer string
func strPtr(s string) *string { return &s }

type seedUser struct {
	FullName   string
	Email      string
	Phone      *string
	IsVerified bool
	Balance    float64
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env")
	}

	ctx := context.Background()

	db, err := config.ConnectPsql()
	if err != nil {
		log.Fatalf("DB connection error: %s", err.Error())
	}
	defer db.Close()

	// Semua user dev punya password yang sama: "Password123"
	var hc pkg.HashConfig
	hc.UseRecommended()
	hashPassword := hc.GenHash("Password123")

	log.Printf("Password hash generated (length: %d)", len(hashPassword))

	// DATA USER
	users := []seedUser{
		{FullName: "Budi Santoso", Email: "budi.santoso@email.com", Phone: strPtr("+6281111111111"), IsVerified: true, Balance: 5750000.00},
		{FullName: "Siti Aminah", Email: "siti.aminah@email.com", Phone: strPtr("+6281222222222"), IsVerified: true, Balance: 735600.65},
		{FullName: "Andi Wijaya", Email: "andi.wijaya@email.com", Phone: strPtr("+6281333333333"), IsVerified: true, Balance: 1191514.40},
		{FullName: "Dewi Lestari", Email: "dewi.lestari@email.com", Phone: strPtr("+6281444444444"), IsVerified: false, Balance: 1551616.23},
		{FullName: "Eko Prasetyo", Email: "eko.prasetyo@email.com", Phone: strPtr("+6281555555555"), IsVerified: true, Balance: 1798470.77},
		{FullName: "Fajar Hidayat", Email: "fajar.saputra@email.com", Phone: strPtr("+6281666666006"), IsVerified: true, Balance: 805518.68},
		{FullName: "Gita Kusuma", Email: "gita.pratama@email.com", Phone: strPtr("+6281666666007"), IsVerified: true, Balance: 409610.61},
		{FullName: "Hendra Putri", Email: "hendra.pratama@email.com", Phone: strPtr("+6281666666008"), IsVerified: false, Balance: 157975.64},
		{FullName: "Indah Pratama", Email: "indah.utami@email.com", Phone: strPtr("+6281666666009"), IsVerified: true, Balance: 584545.16},
		{FullName: "Joko Kusuma", Email: "joko.utami@email.com", Phone: strPtr("+6281666666010"), IsVerified: true, Balance: 625355.95},
		{FullName: "Kartika Kusuma", Email: "kartika.saputra@email.com", Phone: strPtr("+6281666666011"), IsVerified: true, Balance: 1317475.66},
		{FullName: "Luky Pratama", Email: "luky.wijaya@email.com", Phone: strPtr("+6281666666012"), IsVerified: false, Balance: 202944.80},
		{FullName: "Mega Pratama", Email: "mega.saputra@email.com", Phone: strPtr("+6281666666013"), IsVerified: true, Balance: 404736.00},
		{FullName: "Nugroho Putri", Email: "nugroho.utami@email.com", Phone: strPtr("+6281666666014"), IsVerified: true, Balance: 1132083.68},
		{FullName: "Olivia Kusuma", Email: "olivia.pratama@email.com", Phone: strPtr("+6281666666015"), IsVerified: true, Balance: 331048.31},
		{FullName: "Prabowo Kusuma", Email: "prabowo.wijaya@email.com", Phone: strPtr("+6281666666016"), IsVerified: false, Balance: 1791732.26},
		{FullName: "Rina Hidayat", Email: "rina.utami@email.com", Phone: strPtr("+6281666666017"), IsVerified: true, Balance: 1678342.61},
		{FullName: "Surya Hidayat", Email: "surya.saputra@email.com", Phone: strPtr("+6281666666018"), IsVerified: true, Balance: 835569.17},
		{FullName: "Tika Putri", Email: "tika.pratama@email.com", Phone: strPtr("+6281666666019"), IsVerified: true, Balance: 603233.85},
		{FullName: "Utomo Pratama", Email: "utomo.utami@email.com", Phone: strPtr("+6281666666020"), IsVerified: false, Balance: 1987809.32},
	}

	// TRUNCATE SEMUA TABLE
	log.Println("Truncating tables...")
	truncateOrder := []string{
		"topup_details",
		"transfer_details",
		"transactions",
		"favorite_contacts",
		"newsletter",
		"reviews",
		"forgot_password",
		"oauth_user",
		"wallet",
		"users",
	}
	for _, table := range truncateOrder {
		// RESTART IDENTITY → reset auto-increment sequence
		// CASCADE
		sql := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)
		if _, err := db.Exec(ctx, sql); err != nil {
			log.Fatalf("Error truncating %s: %s", table, err.Error())
		}
	}
	log.Println("All tables truncated")

	//INSERT USER + WALLET dalam satu DB transaction per user
	log.Println("Seeding users...")
	for i, u := range users {
		tx, err := db.Begin(ctx)
		if err != nil {
			log.Fatalf("Begin tx error: %s", err.Error())
		}

		// Insert user
		var userID int
		err = tx.QueryRow(ctx,
			`INSERT INTO users (full_name, email, hash_password, phone, is_verified)
			 VALUES ($1, $2, $3, $4, $5)
			 RETURNING id`,
			u.FullName, u.Email, hashPassword, u.Phone, u.IsVerified,
		).Scan(&userID)
		if err != nil {
			tx.Rollback(ctx)
			log.Fatalf("Error inserting user %s: %s", u.Email, err.Error())
		}

		// Insert wallet
		_, err = tx.Exec(ctx,
			`INSERT INTO wallet (user_id, balance, updated_at) VALUES ($1, $2, NOW())`,
			userID, u.Balance,
		)
		if err != nil {
			tx.Rollback(ctx)
			log.Fatalf("Error inserting wallet for user %s: %s", u.Email, err.Error())
		}

		if err := tx.Commit(ctx); err != nil {
			log.Fatalf("Commit error: %s", err.Error())
		}

		log.Printf("[%d/%d] ✓ %s (id: %d, balance: %.2f)",
			i+1, len(users), u.FullName, userID, u.Balance)
	}

	// SELESAI
	log.Println("")
	log.Println("========================================")
	log.Println("Seeding complete!")
	log.Println("Password semua user : Password123")
	log.Println("========================================")
	log.Printf("Total users seeded  : %d", len(users))
}
