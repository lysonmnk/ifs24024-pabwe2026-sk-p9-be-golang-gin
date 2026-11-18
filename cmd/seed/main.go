package main

import (
	"log"
	"time"

	"app/internal/config"
	"app/internal/database"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi: %v", err)
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Gagal inisialisasi basis data: %v", err)
	}

	// 1. Pastikan Author User ada
	var author database.User
	err = db.Where("username = ?", "penulis1").First(&author).Error
	if err != nil {
		hashedPass, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		author = database.User{
			Name:     "Lyson Manurung (Penulis)",
			Username: "penulis1",
			Email:    "penulis1@delcom.org",
			Password: string(hashedPass),
			Role:     "author",
		}
		if err := db.Create(&author).Error; err != nil {
			log.Fatalf("Gagal membuat author seeder: %v", err)
		}
		log.Printf("Seeded Author: username=%s, password=password123", author.Username)
	}

	// 2. Sample Posts
	now := time.Now()
	samplePosts := []database.Post{
		{
			UserID:      author.ID,
			Title:       "Membangun REST API Performa Tinggi dengan Golang Gin",
			Slug:        "membangun-rest-api-performa-tinggi-dengan-golang-gin",
			Category:    "Teknologi",
			Excerpt:     "Panduan komprehensif merancang REST API yang cepat, andal, dan modular menggunakan framework Gin di bahasa pemrograman Go.",
			ReadingTime: 4,
			IsPublished: true,
			PublishedAt: &now,
			Content: `# Membangun REST API Performa Tinggi dengan Golang Gin

Golang dan Gin merupakan kombinasi yang sangat populer dalam membangun sistem backend modern. Dengan fitur goroutine yang ringan dan routing berkecepatan tinggi, Gin mampu menangani puluhan ribu permintaan per detik dengan latensi minimal.

## Keunggulan Utama Gin
1. **Performa Ekstrem**: Routing berbasis *radix tree* yang sangat efisien dalam konsumsi memori.
2. **Middleware Terstruktur**: Memudahkan penerapan proteksi otentikasi JWT, CORS, logging, dan penanganan error.
3. **Validasi Request Otomatis**: Mendukung binding JSON dan validasi struct tag bawaan yang kuat.

` + "```go" + `
package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func main() {
    r := gin.Default()
    r.GET("/api/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "pong"})
    })
    r.Run(":8081")
}
` + "```" + `

## Arsitektur Modular
Dengan memisahkan lapisan **Handler**, **Service**, dan **Repository**, kode menjadi mudah dipelihara dan dapat diuji secara terisolasi menggunakan unit test.
`,
		},
		{
			UserID:      author.ID,
			Title:       "Desain Arsitektur Multi-Tier pada Web Modern",
			Slug:        "desain-arsitektur-multi-tier-pada-web-modern",
			Category:    "Pendidikan",
			Excerpt:     "Memahami pemisahan tanggung jawab antara tier Publik, Penulis (Author), dan Administrator dalam aplikasi web skala enterprise.",
			ReadingTime: 5,
			IsPublished: true,
			PublishedAt: &now,
			Content: `# Desain Arsitektur Multi-Tier pada Web Modern

Arsitektur aplikasi web modern memerlukan batasan akses yang tegas antar peran pengguna guna menjaga integritas dan keamanan data.

## Tiga Tier Utama
- **Public Tier**: Dapat diakses tanpa autentikasi. Menampilkan artikel kurasi, fungsi pencarian, dan penjelajahan kategori.
- **Author Tier**: Memerlukan akun penulis. Menyediakan dashboard statistik performa, pengelolaan draf, serta integrasi pembuat draf berbantuan kecerdasan buatan.
- **Admin Tier**: Khusus administrator untuk memantau direktori seluruh pengguna dan mengawasi arsip artikel yang ada dalam sistem.

Setiap tier dilindungi oleh JWT middleware yang memvalidasi hak akses secara stateless pada setiap endpoint API.
`,
		},
		{
			UserID:      author.ID,
			Title:       "Penerapan AI Generator dalam Alur Publikasi Konten",
			Slug:        "penerapan-ai-generator-dalam-alur-publikasi-konten",
			Category:    "Bisnis",
			Excerpt:     "Bagaimana asisten kecerdasan buatan dengan Google Search Grounding mempercepat riset dan penulisan artikel berkualitas.",
			ReadingTime: 3,
			IsPublished: true,
			PublishedAt: &now,
			Content: `# Penerapan AI Generator dalam Alur Publikasi Konten

Kehadiran asisten AI bukan untuk menggantikan penulis, melainkan berfungsi sebagai *co-pilot* riset yang mempercepat pengumpulan referensi dan pembuatan kerangka tulisan.

## Mengapa Google Search Grounding Penting?
Model AI tradisional sering mengalami halusinasi data karena batasan *knowledge cutoff*. Dengan mengaktifkan pencarian web langsung:
- Informasi yang disajikan selalu aktual dan faktual.
- Sitasi sumber dapat ditelusuri langsung oleh pembaca.
- Penulis dapat memverifikasi klaim dan statistik secara instan.
`,
		},
		{
			UserID:      author.ID,
			Title:       "Panduan Praktis Clean Code dan Otomasi Testing",
			Slug:        "panduan-praktis-clean-code-dan-otomasi-testing",
			Category:    "Teknologi",
			Excerpt:     "Menerapkan prinsip pengujian otomatis hingga mencapai 100% code coverage menggunakan Vitest dan React Testing Library.",
			ReadingTime: 4,
			IsPublished: true,
			PublishedAt: &now,
			Content: `# Panduan Praktis Clean Code dan Otomasi Testing

Kualitas perangkat lunak tidak hanya diukur dari apakah fiturnya berjalan, tetapi juga dari seberapa andal kode tersebut saat dilakukan perubahan di masa depan.

## Metrik Pengujian Kritis
- **Statement Coverage**: Memastikan seluruh baris instruksi dieksekusi selama pengujian.
- **Branch Coverage**: Menguji setiap jalur percabangan logika if-else dan fallback.
- **Function Coverage**: Memastikan seluruh fungsi dan callback telah diverifikasi.
- **Line Coverage**: Menjamin tidak ada baris logika yang terlewati.

Mencapai 100% test coverage memberikan keyakinan tinggi bagi tim sebelum merilis aplikasi ke tahap produksi.
`,
		},
		{
			UserID:      author.ID,
			Title:       "Eksplorasi Komponen UI Berbasis Desain Modular",
			Slug:        "eksplorasi-komponen-ui-berbasis-desain-modular",
			Category:    "Desain",
			Excerpt:     "Membangun modal fullscreen, toast notifications, dan tema gelap/terang yang ramah pengguna.",
			ReadingTime: 2,
			IsPublished: false,
			Content: `# Eksplorasi Komponen UI Berbasis Desain Modular

Draf artikel mengenai prinsip desain antarmuka responsif dan aksesibel pada aplikasi SPA. Saat ini artikel masih berstatus draf di workspace penulis.
`,
		},
	}

	for _, p := range samplePosts {
		var existing database.Post
		if err := db.Where("slug = ?", p.Slug).First(&existing).Error; err != nil {
			if err := db.Create(&p).Error; err != nil {
				log.Printf("Gagal seed post '%s': %v", p.Title, err)
			} else {
				log.Printf("Seeded Post: [%s] %s (Published: %v)", p.Category, p.Title, p.IsPublished)
			}
		}
	}

	log.Println("✅ Database seeding completed successfully!")
}
