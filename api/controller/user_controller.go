package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golang-api/api/responses"
	"golang-api/config"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// LoginUser adalah handler untuk proses login pengguna.
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user map[string]interface{}

	// Membaca data JSON dari body permintaan
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		responses.ErrorResponse(w, "Gagal membaca data pengguna dari permintaan", http.StatusBadRequest)
		return
	}

	// Mendapatkan username dan password dari data pengguna
	username, ok := user["username"].(string)
	if !ok {
		responses.ErrorResponse(w, "Username harus diisi", http.StatusBadRequest)
		return
	}

	password, ok := user["password"].(string)
	if !ok {
		responses.ErrorResponse(w, "Password harus diisi", http.StatusBadRequest)
		return
	}

	// Mengecek apakah pengguna ada di database dan mengambil password dari database
	var dbPassword string
	err := config.DB.QueryRow("SELECT password FROM users WHERE username=?", username).Scan(&dbPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			responses.ErrorResponse(w, "User tidak ditemukan", http.StatusNotFound)
			return
		}
		responses.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Membandingkan password yang dimasukkan dengan password yang ada di database
	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password)); err != nil {
		responses.ErrorResponse(w, "Password salah", http.StatusUnauthorized)
		return
	}
	// Jika login berhasil, buat token JWT
	token := jwt.New(jwt.SigningMethodHS256)

	// Menentukan klaim (claims) token
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // Token berlaku selama 1 jam

	// Menandatangani token dengan secret key (gantilah dengan secret key yang kuat)
	secretKey := []byte("secretKey") // Ganti dengan secret key yang lebih kuat
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		responses.ErrorResponse(w, "Gagal membuat token JWT", http.StatusInternalServerError)
		return
	}

	// Mengembalikan token dan pesan sukses
	response := map[string]interface{}{"token": tokenString}
	responses.SuccessResponse(w, "berhasil login", response, http.StatusCreated)
}

func Register(w http.ResponseWriter, r *http.Request) {
	// Inisialisasi koneksi ke database
	db := config.InitDB()
	defer db.Close() // Tutup koneksi database setelah selesai

	var user map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		// Mengembalikan respons JSON jika gagal membaca data dari permintaan
		responses.ErrorResponse(w, "Gagal membaca data pengguna dari permintaan", http.StatusBadRequest)
		return
	}

	username, ok := user["username"].(string)
	if !ok {
		responses.ErrorResponse(w, "Data username tidak valid", http.StatusBadRequest)
		return
	}
	email, ok := user["email"].(string)
	if !ok {
		responses.ErrorResponse(w, "Data email tidak valid", http.StatusBadRequest)
		return
	}
	password, ok := user["password"].(string)
	if !ok {
		responses.ErrorResponse(w, "Data password tidak valid", http.StatusBadRequest)
		return
	}
	// Hashing password sebelum disimpan ke database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		responses.ErrorResponse(w, "Gagal melakukan hashing password", http.StatusInternalServerError)
		return
	}

	// Simpan pengguna ke database dengan menggunakan data yang telah Anda validasi
	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, hashedPassword)
	if err != nil {
		// Menangani kesalahan jika gagal menyimpan pengguna ke database
		errorMessage := "Gagal menyimpan pengguna ke database"
		if strings.Contains(err.Error(), "Duplicate entry") {
			errorMessage = "Email sudah digunakan. Silakan gunakan email lain."
		}
		responses.ErrorResponse(w, errorMessage, http.StatusInternalServerError)
		return
	}

	// Membuat objek data pengguna untuk dikirim dalam respons
	userData := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		// Anda dapat menambahkan lebih banyak data pengguna sesuai kebutuhan
	}{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	responses.SuccessResponse(w, "Pengguna telah berhasil dibuat", userData, http.StatusCreated)
}

// mengambil user berdasarkan ID.
func GetUser(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan ID pengguna dari parameter URL
	vars := mux.Vars(r)
	userID := vars["id"]

	if userID == "" {
		// Tangani jika ID pengguna tidak ada
		http.Error(w, "ID pengguna harus diisi", http.StatusBadRequest)
		return
	}

	// Mengambil pengguna dari database berdasarkan ID
	var (
		id       int
		username string
		email    string
		password string
	)

	err := config.DB.QueryRow("SELECT id, username, email, Password FROM users WHERE id=?", userID).Scan(&id, &username, &email, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			responses.ErrorResponse(w, "user tidak ditemukan", http.StatusNotFound)
			return
		}
		responses.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Membuat objek data pengguna untuk dikirim dalam respons
	userData := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		// Anda dapat menambahkan lebih banyak data pengguna sesuai kebutuhan
	}{
		Username: username,
		Email:    email,
		Password: password,
	}

	// Mengembalikan data pengguna sebagai JSON
	w.Header().Set("Content-Type", "application/json")
	responses.SuccessResponse(w, "berhasil ambil user", userData, http.StatusCreated)
}

// UpdateUserHandler memperbarui data pengguna berdasarkan ID.
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan ID pengguna dari parameter URL
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "ID pengguna harus disertakan", http.StatusBadRequest)
		return
	}

	// Mendapatkan data pengguna dari body permintaan
	var updatedUser struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"Password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Memperbarui pengguna di database
	_, err := config.DB.Exec("UPDATE users SET username=?, email=? WHERE id=?", updatedUser.Username, updatedUser.Email, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Pengguna telah diperbarui")
}

// DeleteUserHandler menghapus pengguna berdasarkan ID.
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan ID pengguna dari parameter URL
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "ID pengguna harus disertakan", http.StatusBadRequest)
		return
	}

	// Menghapus pengguna dari database
	_, err := config.DB.Exec("DELETE FROM users WHERE id=?", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Pengguna telah dihapus")
}
