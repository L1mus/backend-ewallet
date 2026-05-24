package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/L1mus/backend-ewallet/internal/config"
	"github.com/L1mus/backend-ewallet/pkg"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// ── CREDENTIALS (hardcoded agar mudah testing) ───────────────────────────────
const (
	defaultPassword = "Password123" // semua user
	defaultPin      = "123456"      // semua user
)

// ── HELPERS ──────────────────────────────────────────────────────────────────
func strPtr(s string) *string { return &s }
func daysAgo(n int) time.Time { return time.Now().AddDate(0, 0, -n) }

// ── USER DATA ─────────────────────────────────────────────────────────────────
type seedUser struct {
	FullName   string
	Email      string
	Phone      *string
	PictureURL *string
	IsVerified bool
	Balance    float64
}

var users = []seedUser{
	// User 1 = main test user, data paling lengkap
	{FullName: "Budi Santoso", Email: "budi.santoso@email.com", Phone: strPtr("+6281111111111"), PictureURL: strPtr("https://api.dicebear.com/7.x/avataaars/svg?seed=Budi"), IsVerified: true, Balance: 5750000.00},
	{FullName: "Siti Aminah", Email: "siti.aminah@email.com", Phone: strPtr("+6281222222222"), PictureURL: strPtr("https://api.dicebear.com/7.x/avataaars/svg?seed=Siti"), IsVerified: true, Balance: 735600.65},
	{FullName: "Andi Wijaya", Email: "andi.wijaya@email.com", Phone: strPtr("+6281333333333"), PictureURL: strPtr("https://api.dicebear.com/7.x/avataaars/svg?seed=Andi"), IsVerified: true, Balance: 1191514.40},
	{FullName: "Dewi Lestari", Email: "dewi.lestari@email.com", Phone: strPtr("+6281444444444"), PictureURL: strPtr("https://api.dicebear.com/7.x/avataaars/svg?seed=Dewi"), IsVerified: false, Balance: 1551616.23},
	{FullName: "Eko Prasetyo", Email: "eko.prasetyo@email.com", Phone: strPtr("+6281555555555"), PictureURL: strPtr("https://api.dicebear.com/7.x/avataaars/svg?seed=Eko"), IsVerified: true, Balance: 1798470.77},
	{FullName: "Fajar Hidayat", Email: "fajar.hidayat@email.com", Phone: strPtr("+6281666666006"), PictureURL: strPtr("https://api.dicebear.com/7.x/avataaars/svg?seed=Fajar"), IsVerified: true, Balance: 805518.68},
	{FullName: "Gita Kusuma", Email: "gita.kusuma@email.com", Phone: strPtr("+6281666666007"), PictureURL: strPtr("https://api.dicebear.com/7.x/avataaars/svg?seed=Gita"), IsVerified: true, Balance: 409610.61},
	{FullName: "Hendra Putra", Email: "hendra.putra@email.com", Phone: strPtr("+6281666666008"), PictureURL: nil, IsVerified: false, Balance: 157975.64},
	{FullName: "Indah Pratama", Email: "indah.pratama@email.com", Phone: strPtr("+6281666666009"), PictureURL: strPtr("https://api.dicebear.com/7.x/avataaars/svg?seed=Indah"), IsVerified: true, Balance: 584545.16},
	{FullName: "Joko Santoso", Email: "joko.santoso@email.com", Phone: strPtr("+6281666666010"), PictureURL: strPtr("https://api.dicebear.com/7.x/avataaars/svg?seed=Joko"), IsVerified: true, Balance: 625355.95},
	{FullName: "Kartika Putri", Email: "kartika.putri@email.com", Phone: strPtr("+6281666666011"), PictureURL: strPtr("https://api.dicebear.com/7.x/avataaars/svg?seed=Kartika"), IsVerified: true, Balance: 1317475.66},
	{FullName: "Lukman Hakim", Email: "lukman.hakim@email.com", Phone: strPtr("+6281666666012"), PictureURL: nil, IsVerified: false, Balance: 202944.80},
	{FullName: "Mega Wati", Email: "mega.wati@email.com", Phone: strPtr("+6281666666013"), PictureURL: strPtr("https://api.dicebear.com/7.x/avataaars/svg?seed=Mega"), IsVerified: true, Balance: 404736.00},
	{FullName: "Nugroho Adi", Email: "nugroho.adi@email.com", Phone: strPtr("+6281666666014"), PictureURL: strPtr("https://api.dicebear.com/7.x/avataaars/svg?seed=Nugroho"), IsVerified: true, Balance: 1132083.68},
	{FullName: "Olivia Dewi", Email: "olivia.dewi@email.com", Phone: strPtr("+6281666666015"), PictureURL: strPtr("https://api.dicebear.com/7.x/avataaars/svg?seed=Olivia"), IsVerified: true, Balance: 331048.31},
	{FullName: "Prabowo Hadi", Email: "prabowo.hadi@email.com", Phone: strPtr("+6281666666016"), PictureURL: nil, IsVerified: false, Balance: 1791732.26},
	{FullName: "Rina Susanti", Email: "rina.susanti@email.com", Phone: strPtr("+6281666666017"), PictureURL: strPtr("https://api.dicebear.com/7.x/avataaars/svg?seed=Rina"), IsVerified: true, Balance: 1678342.61},
	{FullName: "Surya Dharma", Email: "surya.dharma@email.com", Phone: strPtr("+6281666666018"), PictureURL: strPtr("https://api.dicebear.com/7.x/avataaars/svg?seed=Surya"), IsVerified: true, Balance: 835569.17},
	{FullName: "Tika Rahayu", Email: "tika.rahayu@email.com", Phone: strPtr("+6281666666019"), PictureURL: strPtr("https://api.dicebear.com/7.x/avataaars/svg?seed=Tika"), IsVerified: true, Balance: 603233.85},
	{FullName: "Utomo Bagas", Email: "utomo.bagas@email.com", Phone: strPtr("+6281666666020"), PictureURL: nil, IsVerified: false, Balance: 1987809.32},
}

// ── PAYMENT CATEGORY ──────────────────────────────────────────────────────────
type seedPaymentCategory struct {
	Name string
}

var paymentCategories = []seedPaymentCategory{
	{Name: "Bank Transfer"},
	{Name: "E-Wallet"},
	{Name: "Convenience Store"},
	{Name: "Credit Card"},
	{Name: "Instant Payment"},
}

// ── PAYMENT METHOD ────────────────────────────────────────────────────────────
type seedPaymentMethod struct {
	CategoryID int // 1-based index into paymentCategories
	Name       string
	Code       string
	Fee        float64
	LogoURL    *string
}

var paymentMethods = []seedPaymentMethod{
	{CategoryID: 1, Name: "BCA Virtual Account", Code: "BCAVA", Fee: 1000, LogoURL: strPtr("https://logo.url/bcava.png")},
	{CategoryID: 1, Name: "Mandiri Virtual Account", Code: "MANDIRIVA", Fee: 1000, LogoURL: strPtr("https://logo.url/mandiriva.png")},
	{CategoryID: 1, Name: "BNI Virtual Account", Code: "BNIVA", Fee: 1000, LogoURL: strPtr("https://logo.url/bniva.png")},
	{CategoryID: 1, Name: "BRI Virtual Account", Code: "BRIVA", Fee: 1000, LogoURL: strPtr("https://logo.url/briva.png")},
	{CategoryID: 1, Name: "Permata Virtual Account", Code: "PERMATAVA", Fee: 1000, LogoURL: strPtr("https://logo.url/permatava.png")},
	{CategoryID: 2, Name: "GoPay", Code: "GOPAY", Fee: 1500, LogoURL: strPtr("https://logo.url/gopay.png")},
	{CategoryID: 2, Name: "OVO", Code: "OVO", Fee: 1500, LogoURL: strPtr("https://logo.url/ovo.png")},
	{CategoryID: 2, Name: "Dana", Code: "DANA", Fee: 1000, LogoURL: strPtr("https://logo.url/dana.png")},
	{CategoryID: 2, Name: "LinkAja", Code: "LINKAJA", Fee: 1000, LogoURL: strPtr("https://logo.url/linkaja.png")},
	{CategoryID: 2, Name: "ShopeePay", Code: "SPAY", Fee: 2000, LogoURL: strPtr("https://logo.url/spay.png")},
	{CategoryID: 3, Name: "Indomaret", Code: "IDM", Fee: 2500, LogoURL: strPtr("https://logo.url/idm.png")},
	{CategoryID: 3, Name: "Alfamart", Code: "ALFA", Fee: 2500, LogoURL: strPtr("https://logo.url/alfa.png")},
	{CategoryID: 4, Name: "Visa Card", Code: "VISA", Fee: 5000, LogoURL: strPtr("https://logo.url/visa.png")},
	{CategoryID: 4, Name: "MasterCard", Code: "MC", Fee: 5000, LogoURL: strPtr("https://logo.url/mc.png")},
	{CategoryID: 4, Name: "JCB Card", Code: "JCB", Fee: 6000, LogoURL: strPtr("https://logo.url/jcb.png")},
	{CategoryID: 5, Name: "QRIS", Code: "QRIS", Fee: 0, LogoURL: strPtr("https://logo.url/qris.png")},
	{CategoryID: 5, Name: "BI-Fast", Code: "BIFAST", Fee: 2500, LogoURL: strPtr("https://logo.url/bifast.png")},
	{CategoryID: 1, Name: "CIMB Niaga VA", Code: "CIMBVA", Fee: 1000, LogoURL: strPtr("https://logo.url/cimbva.png")},
	{CategoryID: 2, Name: "Sakuku", Code: "SAKUKU", Fee: 1000, LogoURL: strPtr("https://logo.url/sakuku.png")},
	{CategoryID: 3, Name: "Lawson", Code: "LAWSON", Fee: 2500, LogoURL: strPtr("https://logo.url/lawson.png")},
}

// ── TRANSACTION DATA ──────────────────────────────────────────────────────────
type txSeed struct {
	UserID       int
	Amount       float64
	TxType       string // "income" | "expense"
	ActivityType string // "topup" | "transfer"
	DaysAgo      int
	// untuk transfer
	ReceiverID  int
	Description string
	// untuk topup
	PaymentMethodID int
	DeliveryFee     float64
}

// Transactions untuk USER 1 (Budi Santoso) — komprehensif untuk semua test
// Rentang 12 bulan untuk test report endpoint
var user1Transactions = []txSeed{
	// ── BULAN INI (days 1-30) ──────────────────────────────────────────────
	{Amount: 1000000, TxType: "income", ActivityType: "topup", DaysAgo: 1, PaymentMethodID: 1, DeliveryFee: 1000},
	{Amount: 500000, TxType: "income", ActivityType: "topup", DaysAgo: 2, PaymentMethodID: 6, DeliveryFee: 1500},
	{Amount: 400000, TxType: "expense", ActivityType: "transfer", DaysAgo: 3, ReceiverID: 2, Description: "Bayar makan siang bareng"},
	{Amount: 700000, TxType: "income", ActivityType: "topup", DaysAgo: 4, PaymentMethodID: 8, DeliveryFee: 1000},
	{Amount: 800000, TxType: "expense", ActivityType: "transfer", DaysAgo: 5, ReceiverID: 3, Description: "Patungan beli kado ulang tahun"},
	{Amount: 300000, TxType: "income", ActivityType: "topup", DaysAgo: 6, PaymentMethodID: 16, DeliveryFee: 0},
	{Amount: 150000, TxType: "expense", ActivityType: "transfer", DaysAgo: 7, ReceiverID: 4, Description: "Bayar bensin"},
	{Amount: 500000, TxType: "income", ActivityType: "topup", DaysAgo: 8, PaymentMethodID: 1, DeliveryFee: 1000},
	{Amount: 300000, TxType: "expense", ActivityType: "transfer", DaysAgo: 9, ReceiverID: 5, Description: "Ongkir barang titipan"},
	{Amount: 500000, TxType: "income", ActivityType: "topup", DaysAgo: 10, PaymentMethodID: 6, DeliveryFee: 1500},
	{Amount: 200000, TxType: "expense", ActivityType: "transfer", DaysAgo: 11, ReceiverID: 2, Description: "Patungan arisan"},
	{Amount: 500000, TxType: "income", ActivityType: "topup", DaysAgo: 12, PaymentMethodID: 7, DeliveryFee: 1500},
	{Amount: 75000, TxType: "expense", ActivityType: "transfer", DaysAgo: 13, ReceiverID: 6, Description: "Pulsa darurat"},
	{Amount: 250000, TxType: "income", ActivityType: "topup", DaysAgo: 14, PaymentMethodID: 16, DeliveryFee: 0},
	{Amount: 100000, TxType: "expense", ActivityType: "transfer", DaysAgo: 15, ReceiverID: 7, Description: "Ganti uang parkir"},
	{Amount: 1200000, TxType: "income", ActivityType: "topup", DaysAgo: 16, PaymentMethodID: 2, DeliveryFee: 1000},
	{Amount: 600000, TxType: "expense", ActivityType: "transfer", DaysAgo: 17, ReceiverID: 3, Description: "Bayar kontrakan bulan ini"},
	{Amount: 450000, TxType: "income", ActivityType: "topup", DaysAgo: 18, PaymentMethodID: 8, DeliveryFee: 1000},
	{Amount: 800000, TxType: "income", ActivityType: "topup", DaysAgo: 19, PaymentMethodID: 3, DeliveryFee: 1000},
	{Amount: 250000, TxType: "expense", ActivityType: "transfer", DaysAgo: 20, ReceiverID: 9, Description: "Belanja keperluan dapur"},
	// ── 2 BULAN LALU (days 31-60) ─────────────────────────────────────────
	{Amount: 3000000, TxType: "income", ActivityType: "topup", DaysAgo: 25, PaymentMethodID: 1, DeliveryFee: 1000},
	{Amount: 1500000, TxType: "expense", ActivityType: "transfer", DaysAgo: 26, ReceiverID: 2, Description: "Transfer bulanan keluarga"},
	{Amount: 500000, TxType: "expense", ActivityType: "transfer", DaysAgo: 28, ReceiverID: 10, Description: "Bayar langganan gym"},
	{Amount: 4000000, TxType: "income", ActivityType: "topup", DaysAgo: 35, PaymentMethodID: 1, DeliveryFee: 1000},
	{Amount: 1200000, TxType: "expense", ActivityType: "transfer", DaysAgo: 40, ReceiverID: 2, Description: "Patungan liburan"},
	{Amount: 800000, TxType: "expense", ActivityType: "transfer", DaysAgo: 48, ReceiverID: 3, Description: "Biaya servis motor"},
	{Amount: 500000, TxType: "income", ActivityType: "topup", DaysAgo: 52, PaymentMethodID: 16, DeliveryFee: 0},
	{Amount: 350000, TxType: "expense", ActivityType: "transfer", DaysAgo: 55, ReceiverID: 4, Description: "Bayar makan keluarga"},
	{Amount: 600000, TxType: "income", ActivityType: "topup", DaysAgo: 58, PaymentMethodID: 6, DeliveryFee: 1500},
	// ── 3 BULAN LALU (days 61-90) ─────────────────────────────────────────
	{Amount: 4000000, TxType: "income", ActivityType: "topup", DaysAgo: 65, PaymentMethodID: 1, DeliveryFee: 1000},
	{Amount: 1200000, TxType: "expense", ActivityType: "transfer", DaysAgo: 70, ReceiverID: 2, Description: "Transfer rutin bulanan"},
	{Amount: 800000, TxType: "expense", ActivityType: "transfer", DaysAgo: 78, ReceiverID: 3, Description: "Bayar utang teman"},
	{Amount: 500000, TxType: "income", ActivityType: "topup", DaysAgo: 82, PaymentMethodID: 16, DeliveryFee: 0},
	{Amount: 200000, TxType: "expense", ActivityType: "transfer", DaysAgo: 85, ReceiverID: 5, Description: "Beli token listrik"},
	{Amount: 750000, TxType: "income", ActivityType: "topup", DaysAgo: 88, PaymentMethodID: 7, DeliveryFee: 1500},
	// ── 4 BULAN LALU (days 91-120) ────────────────────────────────────────
	{Amount: 4000000, TxType: "income", ActivityType: "topup", DaysAgo: 95, PaymentMethodID: 1, DeliveryFee: 1000},
	{Amount: 1200000, TxType: "expense", ActivityType: "transfer", DaysAgo: 100, ReceiverID: 2, Description: "Transfer bulanan"},
	{Amount: 800000, TxType: "expense", ActivityType: "transfer", DaysAgo: 108, ReceiverID: 3, Description: "Bayar cicilan"},
	{Amount: 500000, TxType: "income", ActivityType: "topup", DaysAgo: 112, PaymentMethodID: 16, DeliveryFee: 0},
	{Amount: 400000, TxType: "expense", ActivityType: "transfer", DaysAgo: 115, ReceiverID: 6, Description: "Patungan kado"},
	{Amount: 900000, TxType: "income", ActivityType: "topup", DaysAgo: 118, PaymentMethodID: 4, DeliveryFee: 1000},
	// ── 5 BULAN LALU (days 121-150) ───────────────────────────────────────
	{Amount: 4000000, TxType: "income", ActivityType: "topup", DaysAgo: 125, PaymentMethodID: 1, DeliveryFee: 1000},
	{Amount: 1200000, TxType: "expense", ActivityType: "transfer", DaysAgo: 130, ReceiverID: 2, Description: "Transfer rutin"},
	{Amount: 800000, TxType: "expense", ActivityType: "transfer", DaysAgo: 138, ReceiverID: 3, Description: "Bayar tagihan"},
	{Amount: 500000, TxType: "income", ActivityType: "topup", DaysAgo: 142, PaymentMethodID: 16, DeliveryFee: 0},
	{Amount: 300000, TxType: "expense", ActivityType: "transfer", DaysAgo: 145, ReceiverID: 7, Description: "Ongkos titip"},
	// ── 6 BULAN LALU (days 151-180) ───────────────────────────────────────
	{Amount: 4000000, TxType: "income", ActivityType: "topup", DaysAgo: 155, PaymentMethodID: 1, DeliveryFee: 1000},
	{Amount: 1200000, TxType: "expense", ActivityType: "transfer", DaysAgo: 160, ReceiverID: 2, Description: "Transfer bulanan"},
	{Amount: 800000, TxType: "expense", ActivityType: "transfer", DaysAgo: 168, ReceiverID: 3, Description: "Cicilan motor"},
	{Amount: 500000, TxType: "income", ActivityType: "topup", DaysAgo: 172, PaymentMethodID: 16, DeliveryFee: 0},
	{Amount: 1500000, TxType: "income", ActivityType: "topup", DaysAgo: 175, PaymentMethodID: 2, DeliveryFee: 1000},
	// ── 7 BULAN LALU (days 181-210) ───────────────────────────────────────
	{Amount: 4000000, TxType: "income", ActivityType: "topup", DaysAgo: 185, PaymentMethodID: 1, DeliveryFee: 1000},
	{Amount: 1200000, TxType: "expense", ActivityType: "transfer", DaysAgo: 190, ReceiverID: 2, Description: "Transfer rutin"},
	{Amount: 800000, TxType: "expense", ActivityType: "transfer", DaysAgo: 198, ReceiverID: 3, Description: "Bayar kontrakan"},
	{Amount: 500000, TxType: "income", ActivityType: "topup", DaysAgo: 202, PaymentMethodID: 16, DeliveryFee: 0},
	{Amount: 600000, TxType: "expense", ActivityType: "transfer", DaysAgo: 207, ReceiverID: 5, Description: "Patungan kawinan"},
	// ── 8 BULAN LALU (days 211-240) ───────────────────────────────────────
	{Amount: 4000000, TxType: "income", ActivityType: "topup", DaysAgo: 215, PaymentMethodID: 1, DeliveryFee: 1000},
	{Amount: 1200000, TxType: "expense", ActivityType: "transfer", DaysAgo: 220, ReceiverID: 2, Description: "Transfer bulanan"},
	{Amount: 800000, TxType: "expense", ActivityType: "transfer", DaysAgo: 228, ReceiverID: 3, Description: "Biaya pendidikan"},
	{Amount: 500000, TxType: "income", ActivityType: "topup", DaysAgo: 232, PaymentMethodID: 16, DeliveryFee: 0},
	// ── 9 BULAN LALU (days 241-270) ───────────────────────────────────────
	{Amount: 4000000, TxType: "income", ActivityType: "topup", DaysAgo: 245, PaymentMethodID: 1, DeliveryFee: 1000},
	{Amount: 1200000, TxType: "expense", ActivityType: "transfer", DaysAgo: 250, ReceiverID: 2, Description: "Transfer rutin"},
	{Amount: 800000, TxType: "expense", ActivityType: "transfer", DaysAgo: 258, ReceiverID: 3, Description: "Cicilan kendaraan"},
	{Amount: 500000, TxType: "income", ActivityType: "topup", DaysAgo: 262, PaymentMethodID: 16, DeliveryFee: 0},
	{Amount: 450000, TxType: "expense", ActivityType: "transfer", DaysAgo: 268, ReceiverID: 9, Description: "Bayar iuran"},
	// ── 10 BULAN LALU (days 271-300) ──────────────────────────────────────
	{Amount: 4000000, TxType: "income", ActivityType: "topup", DaysAgo: 275, PaymentMethodID: 1, DeliveryFee: 1000},
	{Amount: 1200000, TxType: "expense", ActivityType: "transfer", DaysAgo: 280, ReceiverID: 2, Description: "Transfer bulanan"},
	{Amount: 800000, TxType: "expense", ActivityType: "transfer", DaysAgo: 288, ReceiverID: 3, Description: "Bayar utang"},
	{Amount: 500000, TxType: "income", ActivityType: "topup", DaysAgo: 292, PaymentMethodID: 16, DeliveryFee: 0},
	// ── 11 BULAN LALU (days 301-330) ──────────────────────────────────────
	{Amount: 4000000, TxType: "income", ActivityType: "topup", DaysAgo: 305, PaymentMethodID: 1, DeliveryFee: 1000},
	{Amount: 1200000, TxType: "expense", ActivityType: "transfer", DaysAgo: 310, ReceiverID: 2, Description: "Transfer rutin"},
	{Amount: 800000, TxType: "expense", ActivityType: "transfer", DaysAgo: 318, ReceiverID: 3, Description: "Bayar asuransi"},
	{Amount: 500000, TxType: "income", ActivityType: "topup", DaysAgo: 322, PaymentMethodID: 16, DeliveryFee: 0},
	{Amount: 700000, TxType: "income", ActivityType: "topup", DaysAgo: 328, PaymentMethodID: 8, DeliveryFee: 1000},
	// ── 12 BULAN LALU (days 331-365) ──────────────────────────────────────
	{Amount: 4000000, TxType: "income", ActivityType: "topup", DaysAgo: 335, PaymentMethodID: 1, DeliveryFee: 1000},
	{Amount: 1200000, TxType: "expense", ActivityType: "transfer", DaysAgo: 340, ReceiverID: 2, Description: "Transfer akhir tahun"},
	{Amount: 800000, TxType: "expense", ActivityType: "transfer", DaysAgo: 348, ReceiverID: 3, Description: "Bayar tagihan tahunan"},
	{Amount: 500000, TxType: "income", ActivityType: "topup", DaysAgo: 352, PaymentMethodID: 16, DeliveryFee: 0},
	{Amount: 2000000, TxType: "income", ActivityType: "topup", DaysAgo: 360, PaymentMethodID: 1, DeliveryFee: 1000},
}

// Transactions untuk user lain — data baseline
type otherUserTx struct {
	UserID          int
	Amount          float64
	TxType          string
	ActivityType    string
	DaysAgo         int
	ReceiverID      int
	Description     string
	PaymentMethodID int
	DeliveryFee     float64
}

var otherTransactions = []otherUserTx{
	// User 2 - Siti: topup dan transfer ke user 1
	{UserID: 2, Amount: 500000, TxType: "income", ActivityType: "topup", DaysAgo: 33, PaymentMethodID: 3, DeliveryFee: 1000},
	{UserID: 2, Amount: 200000, TxType: "expense", ActivityType: "transfer", DaysAgo: 15, ReceiverID: 1, Description: "Bayar titipan"},
	// User 3 - Andi
	{UserID: 3, Amount: 1000000, TxType: "income", ActivityType: "topup", DaysAgo: 11, PaymentMethodID: 1, DeliveryFee: 1000},
	{UserID: 3, Amount: 300000, TxType: "expense", ActivityType: "transfer", DaysAgo: 8, ReceiverID: 1, Description: "Bayar hutang"},
	// User 4 - Dewi
	{UserID: 4, Amount: 2000000, TxType: "income", ActivityType: "topup", DaysAgo: 20, PaymentMethodID: 2, DeliveryFee: 1000},
	{UserID: 4, Amount: 500000, TxType: "expense", ActivityType: "transfer", DaysAgo: 18, ReceiverID: 1, Description: "Transfer ke Budi"},
	// User 5 - Eko
	{UserID: 5, Amount: 750000, TxType: "income", ActivityType: "topup", DaysAgo: 48, PaymentMethodID: 4, DeliveryFee: 1000},
	{UserID: 5, Amount: 250000, TxType: "expense", ActivityType: "transfer", DaysAgo: 30, ReceiverID: 1, Description: "Bayar iuran"},
	// User 6 - Fajar
	{UserID: 6, Amount: 500000, TxType: "income", ActivityType: "topup", DaysAgo: 40, PaymentMethodID: 6, DeliveryFee: 1500},
	{UserID: 6, Amount: 100000, TxType: "expense", ActivityType: "transfer", DaysAgo: 35, ReceiverID: 2, Description: "Ongkir titipan"},
	// User 7 - Gita
	{UserID: 7, Amount: 300000, TxType: "income", ActivityType: "topup", DaysAgo: 20, PaymentMethodID: 16, DeliveryFee: 0},
	{UserID: 7, Amount: 150000, TxType: "expense", ActivityType: "transfer", DaysAgo: 15, ReceiverID: 1, Description: "Bayar nasi padang"},
	// User 8 - Hendra
	{UserID: 8, Amount: 200000, TxType: "income", ActivityType: "topup", DaysAgo: 11, PaymentMethodID: 5, DeliveryFee: 1000},
	// User 9 - Indah
	{UserID: 9, Amount: 600000, TxType: "income", ActivityType: "topup", DaysAgo: 29, PaymentMethodID: 7, DeliveryFee: 1500},
	{UserID: 9, Amount: 250000, TxType: "expense", ActivityType: "transfer", DaysAgo: 18, ReceiverID: 1, Description: "Patungan bayar makan"},
	// User 10 - Joko
	{UserID: 10, Amount: 800000, TxType: "income", ActivityType: "topup", DaysAgo: 29, PaymentMethodID: 4, DeliveryFee: 1000},
	// User 11 - Kartika
	{UserID: 11, Amount: 1500000, TxType: "income", ActivityType: "topup", DaysAgo: 43, PaymentMethodID: 1, DeliveryFee: 1000},
	{UserID: 11, Amount: 500000, TxType: "expense", ActivityType: "transfer", DaysAgo: 22, ReceiverID: 1, Description: "Bayar sewa kos"},
	// User 12 - Lukman
	{UserID: 12, Amount: 300000, TxType: "income", ActivityType: "topup", DaysAgo: 34, PaymentMethodID: 11, DeliveryFee: 2500},
	// User 13 - Mega
	{UserID: 13, Amount: 500000, TxType: "income", ActivityType: "topup", DaysAgo: 35, PaymentMethodID: 8, DeliveryFee: 1000},
	{UserID: 13, Amount: 200000, TxType: "expense", ActivityType: "transfer", DaysAgo: 12, ReceiverID: 1, Description: "Bayar parkiran"},
	// User 14 - Nugroho
	{UserID: 14, Amount: 1200000, TxType: "income", ActivityType: "topup", DaysAgo: 35, PaymentMethodID: 1, DeliveryFee: 1000},
	{UserID: 14, Amount: 400000, TxType: "expense", ActivityType: "transfer", DaysAgo: 7, ReceiverID: 2, Description: "Transfer ke Siti"},
	// User 15 - Olivia
	{UserID: 15, Amount: 400000, TxType: "income", ActivityType: "topup", DaysAgo: 34, PaymentMethodID: 16, DeliveryFee: 0},
	{UserID: 15, Amount: 100000, TxType: "expense", ActivityType: "transfer", DaysAgo: 7, ReceiverID: 1, Description: "Patungan kopi"},
	// User 16 - Prabowo
	{UserID: 16, Amount: 2000000, TxType: "income", ActivityType: "topup", DaysAgo: 34, PaymentMethodID: 2, DeliveryFee: 1000},
	// User 17 - Rina
	{UserID: 17, Amount: 200000, TxType: "income", ActivityType: "topup", DaysAgo: 34, PaymentMethodID: 5, DeliveryFee: 1000},
	{UserID: 17, Amount: 50000, TxType: "expense", ActivityType: "transfer", DaysAgo: 23, ReceiverID: 1, Description: "Bayar ojek"},
	// User 18 - Surya
	{UserID: 18, Amount: 1000000, TxType: "income", ActivityType: "topup", DaysAgo: 34, PaymentMethodID: 5, DeliveryFee: 1000},
	// User 19 - Tika
	{UserID: 19, Amount: 800000, TxType: "income", ActivityType: "topup", DaysAgo: 26, PaymentMethodID: 1, DeliveryFee: 1000},
	{UserID: 19, Amount: 300000, TxType: "expense", ActivityType: "transfer", DaysAgo: 12, ReceiverID: 2, Description: "Transfer ke Siti"},
	// User 20 - Utomo
	{UserID: 20, Amount: 2500000, TxType: "income", ActivityType: "topup", DaysAgo: 26, PaymentMethodID: 1, DeliveryFee: 1000},
}

// ── REVIEW DATA ───────────────────────────────────────────────────────────────
type seedReview struct {
	UserID      int
	Rating      int
	Description string
	DaysAgo     int
}

var reviews = []seedReview{
	{UserID: 7, Rating: 4, Description: "Suka dengan tampilan aplikasinya, mudah digunakan", DaysAgo: 1},
	{UserID: 3, Rating: 5, Description: "Aplikasi transfer paling mudah yang pernah saya pakai!", DaysAgo: 2},
	{UserID: 17, Rating: 4, Description: "Sangat membantu mengatur keuangan bulanan", DaysAgo: 3},
	{UserID: 1, Rating: 5, Description: "Transfernya cepat dan aman, sangat recommended", DaysAgo: 4},
	{UserID: 6, Rating: 3, Description: "Bagus tapi loading agak lambat kadang", DaysAgo: 5},
	{UserID: 5, Rating: 4, Description: "Fitur laporan keuangan sangat berguna", DaysAgo: 6},
	{UserID: 8, Rating: 4, Description: "Mudah digunakan untuk pembayaran sehari-hari", DaysAgo: 7},
	{UserID: 2, Rating: 5, Description: "Proses topup sangat cepat, tidak sampai 1 menit", DaysAgo: 8},
	{UserID: 9, Rating: 4, Description: "Tampilan bersih dan tidak membingungkan", DaysAgo: 9},
	{UserID: 16, Rating: 3, Description: "Sudah cukup bagus, harapkan fitur split bill", DaysAgo: 10},
	{UserID: 11, Rating: 4, Description: "Aplikasi stabil, jarang error", DaysAgo: 11},
	{UserID: 18, Rating: 5, Description: "Terbaik dibanding e-wallet lain yang pernah saya coba", DaysAgo: 12},
	{UserID: 13, Rating: 4, Description: "Support-nya responsif dan membantu", DaysAgo: 13},
	{UserID: 14, Rating: 5, Description: "Fitur riwayat transaksi lengkap banget", DaysAgo: 14},
	{UserID: 4, Rating: 3, Description: "Perlu peningkatan di notifikasi transaksi", DaysAgo: 15},
	{UserID: 19, Rating: 4, Description: "Pengiriman uang ke luar negeri kapan ada?", DaysAgo: 16},
	{UserID: 10, Rating: 5, Description: "Saldo selalu aman dan akurat", DaysAgo: 17},
	{UserID: 12, Rating: 4, Description: "Fitur pencarian kontak sangat membantu", DaysAgo: 18},
	{UserID: 15, Rating: 2, Description: "Pernah gagal transfer 2x, tapi uang kembali", DaysAgo: 19},
	{UserID: 20, Rating: 3, Description: "Cukup oke, tapi masih ada bug kecil", DaysAgo: 20},
}

// ── NEWSLETTER DATA ───────────────────────────────────────────────────────────
type seedNewsletter struct {
	UserID  int
	Email   string
	Status  string
	DaysAgo int
}

var newsletters = []seedNewsletter{
	{UserID: 1, Email: "newsletter.budi@example.com", Status: "active", DaysAgo: 5},
	{UserID: 2, Email: "newsletter.siti@example.com", Status: "active", DaysAgo: 10},
	{UserID: 3, Email: "newsletter.andi@example.com", Status: "active", DaysAgo: 15},
	{UserID: 4, Email: "newsletter.dewi@example.com", Status: "unsubscribe", DaysAgo: 20},
	{UserID: 5, Email: "newsletter.eko@example.com", Status: "active", DaysAgo: 25},
	{UserID: 6, Email: "newsletter.fajar@example.com", Status: "active", DaysAgo: 30},
	{UserID: 7, Email: "newsletter.gita@example.com", Status: "active", DaysAgo: 35},
	{UserID: 8, Email: "newsletter.hendra@example.com", Status: "unsubscribe", DaysAgo: 40},
	{UserID: 9, Email: "newsletter.indah@example.com", Status: "active", DaysAgo: 45},
	{UserID: 10, Email: "newsletter.joko@example.com", Status: "active", DaysAgo: 50},
	{UserID: 11, Email: "newsletter.kartika@example.com", Status: "active", DaysAgo: 55},
	{UserID: 12, Email: "newsletter.lukman@example.com", Status: "unsubscribe", DaysAgo: 60},
	{UserID: 13, Email: "newsletter.mega@example.com", Status: "active", DaysAgo: 65},
	{UserID: 14, Email: "newsletter.nugroho@example.com", Status: "active", DaysAgo: 70},
	{UserID: 15, Email: "newsletter.olivia@example.com", Status: "active", DaysAgo: 75},
	{UserID: 16, Email: "newsletter.prabowo@example.com", Status: "unsubscribe", DaysAgo: 80},
	{UserID: 17, Email: "newsletter.rina@example.com", Status: "active", DaysAgo: 85},
	{UserID: 18, Email: "newsletter.surya@example.com", Status: "active", DaysAgo: 90},
	{UserID: 19, Email: "newsletter.tika@example.com", Status: "active", DaysAgo: 95},
	{UserID: 20, Email: "newsletter.utomo@example.com", Status: "unsubscribe", DaysAgo: 100},
}

// ── MAIN ──────────────────────────────────────────────────────────────────────
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

	// 1. Generate hashes
	var hc pkg.HashConfig
	hc.UseRecommended()
	hashPassword := hc.GenHash(defaultPassword)
	hashPin := hc.GenHash(defaultPin)
	log.Println("✓ Hash password & PIN generated")

	// 2. Truncate semua table (child → parent)
	truncateAll(ctx, db)

	// 3. Seed dalam urutan yang benar
	userIDs := seedUsers(ctx, db, hashPassword, hashPin)
	categoryIDs := seedPaymentCategories(ctx, db)
	methodIDs := seedPaymentMethods(ctx, db, categoryIDs)
	txIDs := seedUser1Transactions(ctx, db, userIDs, methodIDs)
	otherTxIDs := seedOtherTransactions(ctx, db, userIDs, methodIDs)
	seedTransferDetails(ctx, db, txIDs, otherTxIDs, userIDs)
	seedTopupDetails(ctx, db, txIDs, otherTxIDs, methodIDs)
	seedFavoriteContacts(ctx, db, userIDs)
	seedReviews(ctx, db, userIDs)
	seedNewsletters(ctx, db, userIDs)
	seedForgotPasswords(ctx, db, userIDs)
	seedOAuthUsers(ctx, db, userIDs)

	// 4. Ringkasan
	printSummary(userIDs)
}

// ── TRUNCATE ──────────────────────────────────────────────────────────────────
func truncateAll(ctx context.Context, db *pgxpool.Pool) {
	log.Println("Truncating all tables...")
	tables := []string{
		"topup_details", "transfer_details", "transactions",
		"favorite_contacts", "newsletter", "reviews",
		"forgot_password", "oauth_user", "wallet", "users",
		"payment_method", "category_payment_method",
	}
	for _, t := range tables {
		if _, err := db.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", t)); err != nil {
			log.Fatalf("Error truncating %s: %s", t, err.Error())
		}
	}
	log.Println("✓ All tables truncated")
}

// ── SEED USERS ────────────────────────────────────────────────────────────────
func seedUsers(ctx context.Context, db *pgxpool.Pool, hashPassword, hashPin string) []int {
	log.Println("Seeding users + wallets...")
	ids := make([]int, 0, len(users))

	for i, u := range users {
		tx, _ := db.Begin(ctx)

		var id int
		err := tx.QueryRow(ctx,
			`INSERT INTO users (full_name, email, hash_password, hash_pin, phone, profile_picture_url, is_verified)
			 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			u.FullName, u.Email, hashPassword, hashPin, u.Phone, u.PictureURL, u.IsVerified,
		).Scan(&id)
		if err != nil {
			tx.Rollback(ctx)
			log.Fatalf("Error inserting user %s: %s", u.Email, err.Error())
		}

		if _, err := tx.Exec(ctx,
			`INSERT INTO wallet (user_id, balance, updated_at) VALUES ($1, $2, NOW())`,
			id, u.Balance,
		); err != nil {
			tx.Rollback(ctx)
			log.Fatalf("Error inserting wallet for %s: %s", u.Email, err.Error())
		}

		if err := tx.Commit(ctx); err != nil {
			log.Fatalf("Commit error: %s", err.Error())
		}

		ids = append(ids, id)
		log.Printf("  [%02d/20] ✓ %s (id:%d, balance:%.0f)", i+1, u.FullName, id, u.Balance)
	}
	log.Println("✓ Users & wallets seeded")
	return ids
}

// ── SEED PAYMENT CATEGORIES ───────────────────────────────────────────────────
func seedPaymentCategories(ctx context.Context, db *pgxpool.Pool) []int {
	log.Println("Seeding payment categories...")
	ids := make([]int, 0, len(paymentCategories))
	for _, c := range paymentCategories {
		var id int
		if err := db.QueryRow(ctx,
			`INSERT INTO category_payment_method (category_name) VALUES ($1) RETURNING id`,
			c.Name,
		).Scan(&id); err != nil {
			log.Fatalf("Error inserting category %s: %s", c.Name, err.Error())
		}
		ids = append(ids, id)
	}
	log.Printf("✓ %d payment categories seeded", len(ids))
	return ids
}

// ── SEED PAYMENT METHODS ──────────────────────────────────────────────────────
func seedPaymentMethods(ctx context.Context, db *pgxpool.Pool, categoryIDs []int) []int {
	log.Println("Seeding payment methods...")
	ids := make([]int, 0, len(paymentMethods))
	for _, m := range paymentMethods {
		catID := categoryIDs[m.CategoryID-1]
		var id int
		if err := db.QueryRow(ctx,
			`INSERT INTO payment_method (payment_category_id, name, code, fee, logo_url, is_active)
			 VALUES ($1, $2, $3, $4, $5, TRUE) RETURNING id`,
			catID, m.Name, m.Code, m.Fee, m.LogoURL,
		).Scan(&id); err != nil {
			log.Fatalf("Error inserting payment method %s: %s", m.Name, err.Error())
		}
		ids = append(ids, id)
	}
	log.Printf("✓ %d payment methods seeded", len(ids))
	return ids
}

// ── SEED USER 1 TRANSACTIONS ──────────────────────────────────────────────────
// Returns slice of (transactionID, txSeed index) untuk dipakai seed detail
type txRecord struct {
	ID   int
	Seed txSeed
}

func seedUser1Transactions(ctx context.Context, db *pgxpool.Pool, userIDs, methodIDs []int) []txRecord {
	log.Println("Seeding user 1 transactions (comprehensive)...")
	records := make([]txRecord, 0, len(user1Transactions))

	user1ID := userIDs[0] // Budi Santoso

	for _, s := range user1Transactions {
		seed := s
		seed.UserID = user1ID
		if seed.ReceiverID > 0 {
			seed.ReceiverID = userIDs[seed.ReceiverID-1]
		}

		var txID int
		if err := db.QueryRow(ctx,
			`INSERT INTO transactions (user_id, amount, type, activity_type, status, created_at)
			 VALUES ($1, $2, $3, $4, 'success', $5) RETURNING id`,
			seed.UserID, seed.Amount, seed.TxType, seed.ActivityType, daysAgo(seed.DaysAgo),
		).Scan(&txID); err != nil {
			log.Fatalf("Error inserting tx: %s", err.Error())
		}
		records = append(records, txRecord{ID: txID, Seed: seed})
	}
	log.Printf("✓ %d transactions seeded for user 1", len(records))
	return records
}

// ── SEED OTHER USER TRANSACTIONS ──────────────────────────────────────────────
type otherTxRecord struct {
	ID   int
	Seed otherUserTx
}

func seedOtherTransactions(ctx context.Context, db *pgxpool.Pool, userIDs, methodIDs []int) []otherTxRecord {
	log.Println("Seeding transactions for other users...")
	records := make([]otherTxRecord, 0, len(otherTransactions))

	for _, s := range otherTransactions {
		seed := s
		seed.UserID = userIDs[seed.UserID-1]
		if seed.ReceiverID > 0 {
			seed.ReceiverID = userIDs[seed.ReceiverID-1]
		}

		var txID int
		if err := db.QueryRow(ctx,
			`INSERT INTO transactions (user_id, amount, type, activity_type, status, created_at)
			 VALUES ($1, $2, $3, $4, 'success', $5) RETURNING id`,
			seed.UserID, seed.Amount, seed.TxType, seed.ActivityType, daysAgo(seed.DaysAgo),
		).Scan(&txID); err != nil {
			log.Fatalf("Error inserting other tx: %s", err.Error())
		}
		records = append(records, otherTxRecord{ID: txID, Seed: seed})
	}
	log.Printf("✓ %d transactions seeded for other users", len(records))
	return records
}

// ── SEED TRANSFER & TOPUP DETAILS ─────────────────────────────────────────────
func seedTransferDetails(ctx context.Context, db *pgxpool.Pool, user1Txs []txRecord, otherTxs []otherTxRecord, userIDs []int) {
	log.Println("Seeding transfer details...")
	count := 0

	for _, r := range user1Txs {
		if r.Seed.ActivityType != "transfer" {
			continue
		}
		if _, err := db.Exec(ctx,
			`INSERT INTO transfer_details (transaction_id, receiver_id, description) VALUES ($1, $2, $3)`,
			r.ID, r.Seed.ReceiverID, r.Seed.Description,
		); err != nil {
			log.Fatalf("Error inserting transfer_detail: %s", err.Error())
		}
		count++
	}

	for _, r := range otherTxs {
		if r.Seed.ActivityType != "transfer" {
			continue
		}
		if _, err := db.Exec(ctx,
			`INSERT INTO transfer_details (transaction_id, receiver_id, description) VALUES ($1, $2, $3)`,
			r.ID, r.Seed.ReceiverID, r.Seed.Description,
		); err != nil {
			log.Fatalf("Error inserting transfer_detail (other): %s", err.Error())
		}
		count++
	}
	log.Printf("✓ %d transfer details seeded", count)
}

func seedTopupDetails(ctx context.Context, db *pgxpool.Pool, user1Txs []txRecord, otherTxs []otherTxRecord, methodIDs []int) {
	log.Println("Seeding topup details...")
	count := 0

	for _, r := range user1Txs {
		if r.Seed.ActivityType != "topup" {
			continue
		}
		pmID := methodIDs[r.Seed.PaymentMethodID-1]
		total := r.Seed.Amount + r.Seed.DeliveryFee
		if _, err := db.Exec(ctx,
			`INSERT INTO topup_details (transaction_id, payment_method_id, order_amount, delivery_fee, tax_amount, total_amount)
			 VALUES ($1, $2, $3, $4, 0.00, $5)`,
			r.ID, pmID, r.Seed.Amount, r.Seed.DeliveryFee, total,
		); err != nil {
			log.Fatalf("Error inserting topup_detail: %s", err.Error())
		}
		count++
	}

	for _, r := range otherTxs {
		if r.Seed.ActivityType != "topup" {
			continue
		}
		pmID := methodIDs[r.Seed.PaymentMethodID-1]
		total := r.Seed.Amount + r.Seed.DeliveryFee
		if _, err := db.Exec(ctx,
			`INSERT INTO topup_details (transaction_id, payment_method_id, order_amount, delivery_fee, tax_amount, total_amount)
			 VALUES ($1, $2, $3, $4, 0.00, $5)`,
			r.ID, pmID, r.Seed.Amount, r.Seed.DeliveryFee, total,
		); err != nil {
			log.Fatalf("Error inserting topup_detail (other): %s", err.Error())
		}
		count++
	}
	log.Printf("✓ %d topup details seeded", count)
}

// ── SEED FAVORITE CONTACTS ────────────────────────────────────────────────────
func seedFavoriteContacts(ctx context.Context, db *pgxpool.Pool, userIDs []int) {
	log.Println("Seeding favorite contacts...")
	// User 1 simpan user 2,3,4,5 sebagai favorit
	// User 2 simpan user 1,3,4,5
	// dst
	pairs := [][2]int{
		{1, 2}, {1, 3}, {1, 4}, {1, 5},
		{2, 1}, {2, 3}, {2, 4}, {2, 5},
		{3, 1}, {3, 2}, {3, 4}, {3, 5},
		{4, 1}, {4, 2}, {4, 3}, {4, 5},
		{5, 1}, {5, 2}, {5, 3}, {5, 4},
	}
	for _, p := range pairs {
		uid := userIDs[p[0]-1]
		fid := userIDs[p[1]-1]
		if _, err := db.Exec(ctx,
			`INSERT INTO favorite_contacts (user_id, favorite_user_id) VALUES ($1, $2)`,
			uid, fid,
		); err != nil {
			log.Fatalf("Error inserting favorite_contact: %s", err.Error())
		}
	}
	log.Printf("✓ %d favorite contacts seeded", len(pairs))
}

// ── SEED REVIEWS ──────────────────────────────────────────────────────────────
func seedReviews(ctx context.Context, db *pgxpool.Pool, userIDs []int) {
	log.Println("Seeding reviews...")
	for _, r := range reviews {
		uid := userIDs[r.UserID-1]
		if _, err := db.Exec(ctx,
			`INSERT INTO reviews (user_id, rating, description, created_at) VALUES ($1, $2, $3, $4)`,
			uid, r.Rating, r.Description, daysAgo(r.DaysAgo),
		); err != nil {
			log.Fatalf("Error inserting review: %s", err.Error())
		}
	}
	log.Printf("✓ %d reviews seeded", len(reviews))
}

// ── SEED NEWSLETTERS ──────────────────────────────────────────────────────────
func seedNewsletters(ctx context.Context, db *pgxpool.Pool, userIDs []int) {
	log.Println("Seeding newsletters...")
	for _, n := range newsletters {
		uid := userIDs[n.UserID-1]
		if _, err := db.Exec(ctx,
			`INSERT INTO newsletter (user_id, email, status, created_at) VALUES ($1, $2, $3, $4)`,
			uid, n.Email, n.Status, daysAgo(n.DaysAgo),
		); err != nil {
			log.Fatalf("Error inserting newsletter: %s", err.Error())
		}
	}
	log.Printf("✓ %d newsletter entries seeded", len(newsletters))
}

// ── SEED FORGOT PASSWORDS ─────────────────────────────────────────────────────
func seedForgotPasswords(ctx context.Context, db *pgxpool.Pool, userIDs []int) {
	log.Println("Seeding forgot password tokens...")
	entries := []struct {
		UserIdx int
		Token   string
		IsUsed  bool
		DaysAgo int
	}{
		{7, "token_expired_001", true, 10},
		{11, "token_expired_002", true, 8},
		{3, "token_active_003", false, 1}, // masih aktif (1 jam ke depan)
		{15, "token_expired_004", true, 5},
		{2, "token_active_005", false, 0}, // baru saja dibuat
	}
	for _, e := range entries {
		uid := userIDs[e.UserIdx-1]
		createdAt := daysAgo(e.DaysAgo)
		expiredAt := createdAt.Add(1 * time.Hour)
		if _, err := db.Exec(ctx,
			`INSERT INTO forgot_password (user_id, token, is_used, created_at, expired_at) VALUES ($1,$2,$3,$4,$5)`,
			uid, e.Token, e.IsUsed, createdAt, expiredAt,
		); err != nil {
			log.Fatalf("Error inserting forgot_password: %s", err.Error())
		}
	}
	log.Printf("✓ %d forgot password tokens seeded", len(entries))
}

// ── SEED OAUTH USERS ──────────────────────────────────────────────────────────
func seedOAuthUsers(ctx context.Context, db *pgxpool.Pool, userIDs []int) {
	log.Println("Seeding OAuth users...")
	entries := []struct {
		UserIdx      int
		Provider     string
		ProviderUID  string
		AccessToken  string
		RefreshToken string
	}{
		{1, "google", "google_uid_001", "access_budi_001", "refresh_budi_001"},
		{2, "facebook", "fb_uid_002", "access_siti_002", "refresh_siti_002"},
		{3, "google", "google_uid_003", "access_andi_003", "refresh_andi_003"},
		{5, "google", "google_uid_005", "access_eko_005", "refresh_eko_005"},
		{6, "facebook", "fb_uid_006", "access_fajar_006", "refresh_fajar_006"},
		{9, "google", "google_uid_009", "access_indah_009", "refresh_indah_009"},
		{10, "facebook", "fb_uid_010", "access_joko_010", "refresh_joko_010"},
	}
	for _, e := range entries {
		uid := userIDs[e.UserIdx-1]
		if _, err := db.Exec(ctx,
			`INSERT INTO oauth_user (user_id, provider_name, provider_user_id, access_token, refresh_token)
			 VALUES ($1, $2, $3, $4, $5)`,
			uid, e.Provider, e.ProviderUID, e.AccessToken, e.RefreshToken,
		); err != nil {
			log.Fatalf("Error inserting oauth_user: %s", err.Error())
		}
	}
	log.Printf("✓ %d OAuth users seeded", len(entries))
}

// ── SUMMARY ───────────────────────────────────────────────────────────────────
func printSummary(userIDs []int) {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                   SEEDING COMPLETE!                          ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	fmt.Println("║  CREDENTIALS (semua user)                                    ║")
	fmt.Printf("║  Password : %-49s║\n", defaultPassword)
	fmt.Printf("║  PIN      : %-49s║\n", defaultPin)
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	fmt.Println("║  TEST ACCOUNTS                                               ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	for i, u := range users {
		line := fmt.Sprintf("  [%02d] %-20s │ %s", i+1, u.FullName, u.Email)
		fmt.Printf("║%-62s║\n", line)
	}
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	fmt.Println("║  DATA SUMMARY                                                ║")
	fmt.Println("║  • User 1 (Budi): 80+ transaksi, rentang 12 bulan           ║")
	fmt.Println("║  • Semua user: balance & PIN tersedia                        ║")
	fmt.Println("║  • 20 payment methods dari 5 kategori                        ║")
	fmt.Println("║  • Favorite contacts: user 1-5 saling terhubung             ║")
	fmt.Println("║  • 20 reviews & 20 newsletter entries                        ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	fmt.Println("║  ENDPOINT TEST GUIDE                                         ║")
	fmt.Println("║  POST /auth/register     → email baru                        ║")
	fmt.Println("║  POST /auth              → email user + Password123          ║")
	fmt.Println("║  GET  /users/profile     → butuh Bearer token                ║")
	fmt.Println("║  GET  /users/dashboard   → coba user 1 (data terlengkap)    ║")
	fmt.Println("║  GET  /users/report?period=month                             ║")
	fmt.Println("║  GET  /users/report?period=year                              ║")
	fmt.Println("║  GET  /users/transfer?search=budi&page=1                    ║")
	fmt.Println("║  GET  /users/transactions?page=1 (user 1: 80+ items)        ║")
	fmt.Println("║  PATCH /users/pin        → current_pin: 123456               ║")
	fmt.Println("║  PATCH /users/password   → current_password: Password123     ║")
	fmt.Println("║  POST /transactions/transfer → pin: 123456                   ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}
