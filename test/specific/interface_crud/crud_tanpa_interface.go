// Package main mendemonstrasikan CRUD TANPA menggunakan interface.
// Ini adalah pendekatan langsung yang lebih sederhana tapi kurang fleksibel.
package main

import (
	"errors"
	"fmt"
)

// =============================================================================
// MODEL
// =============================================================================

// User adalah model data untuk pengguna
type User struct {
	ID    int
	Name  string
	Email string
}

// =============================================================================
// STORAGE (Implementasi Langsung Tanpa Interface)
// =============================================================================

// UserStore adalah penyimpanan data user menggunakan slice (in-memory).
// PERHATIKAN: Tidak ada interface yang didefinisikan - langsung implementasi konkret.
type UserStore struct {
	users  []User
	nextID int
}

// NewUserStore membuat instance baru dari UserStore
func NewUserStore() *UserStore {
	return &UserStore{
		users:  make([]User, 0),
		nextID: 1,
	}
}

// Create menambahkan user baru ke storage
func (s *UserStore) Create(name, email string) User {
	user := User{
		ID:    s.nextID,
		Name:  name,
		Email: email,
	}
	s.users = append(s.users, user)
	s.nextID++
	return user
}

// GetByID mengambil user berdasarkan ID
func (s *UserStore) GetByID(id int) (User, error) {
	for _, user := range s.users {
		if user.ID == id {
			return user, nil
		}
	}
	return User{}, errors.New("user tidak ditemukan")
}

// GetAll mengambil semua user
func (s *UserStore) GetAll() []User {
	return s.users
}

// Update mengubah data user yang sudah ada
func (s *UserStore) Update(id int, name, email string) (User, error) {
	for i, user := range s.users {
		if user.ID == id {
			s.users[i].Name = name
			s.users[i].Email = email
			return s.users[i], nil
		}
	}
	return User{}, errors.New("user tidak ditemukan")
}

// Delete menghapus user berdasarkan ID
func (s *UserStore) Delete(id int) error {
	for i, user := range s.users {
		if user.ID == id {
			// Hapus elemen dari slice
			s.users = append(s.users[:i], s.users[i+1:]...)
			return nil
		}
	}
	return errors.New("user tidak ditemukan")
}

// =============================================================================
// SERVICE LAYER
// =============================================================================

// UserService adalah layer bisnis logic yang menggunakan UserStore LANGSUNG.
// MASALAH: Service ini tightly coupled dengan UserStore konkret.
// Tidak bisa diganti ke database lain tanpa mengubah kode ini.
type UserService struct {
	store *UserStore // <-- Dependency langsung ke implementasi konkret!
}

// NewUserService membuat instance baru UserService
func NewUserService(store *UserStore) *UserService {
	return &UserService{store: store}
}

// RegisterUser mendaftarkan user baru dengan validasi
func (s *UserService) RegisterUser(name, email string) (User, error) {
	// Validasi sederhana
	if name == "" || email == "" {
		return User{}, errors.New("nama dan email harus diisi")
	}

	user := s.store.Create(name, email)
	fmt.Printf("✅ User '%s' berhasil didaftarkan dengan ID %d\n", name, user.ID)
	return user, nil
}

// GetUserProfile mengambil profil user
func (s *UserService) GetUserProfile(id int) (User, error) {
	return s.store.GetByID(id)
}

// UpdateUserProfile mengupdate profil user
func (s *UserService) UpdateUserProfile(id int, name, email string) (User, error) {
	return s.store.Update(id, name, email)
}

// DeleteUser menghapus user
func (s *UserService) DeleteUser(id int) error {
	return s.store.Delete(id)
}

// =============================================================================
// MAIN - Demo Penggunaan
// =============================================================================

func main() {
	fmt.Println("=" + "=========================")
	fmt.Println("CRUD TANPA INTERFACE")
	fmt.Println("==========================")
	fmt.Println()

	// Inisialisasi storage dan service
	store := NewUserStore()
	service := NewUserService(store)

	// CREATE - Tambah beberapa user
	fmt.Println("📝 CREATE: Menambahkan user...")
	user1, _ := service.RegisterUser("Budi", "budi@email.com")
	user2, _ := service.RegisterUser("Ani", "ani@email.com")
	user3, _ := service.RegisterUser("Citra", "citra@email.com")
	fmt.Println()

	// READ - Baca semua user
	fmt.Println("📖 READ: Daftar semua user...")
	for _, user := range store.GetAll() {
		fmt.Printf("   - ID: %d, Nama: %s, Email: %s\n", user.ID, user.Name, user.Email)
	}
	fmt.Println()

	// READ - Baca user spesifik
	fmt.Println("📖 READ: Mengambil user dengan ID 2...")
	user, err := service.GetUserProfile(2)
	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Printf("   ✅ Ditemukan: %s (%s)\n", user.Name, user.Email)
	}
	fmt.Println()

	// UPDATE - Update user
	fmt.Println("✏️ UPDATE: Mengubah nama user ID 1 menjadi 'Budi Santoso'...")
	updatedUser, err := service.UpdateUserProfile(user1.ID, "Budi Santoso", "budi.santoso@email.com")
	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Printf("   ✅ Berhasil diupdate: %s (%s)\n", updatedUser.Name, updatedUser.Email)
	}
	fmt.Println()

	// DELETE - Hapus user
	fmt.Println("🗑️ DELETE: Menghapus user ID 3...")
	err = service.DeleteUser(user3.ID)
	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Printf("   ✅ User berhasil dihapus\n")
	}
	fmt.Println()

	// Tampilkan hasil akhir
	fmt.Println("📋 HASIL AKHIR:")
	for _, user := range store.GetAll() {
		fmt.Printf("   - ID: %d, Nama: %s, Email: %s\n", user.ID, user.Name, user.Email)
	}
	fmt.Println()

	// MASALAH dengan pendekatan ini:
	fmt.Println("⚠️ MASALAH TANPA INTERFACE:")
	fmt.Println("   1. UserService tightly coupled dengan UserStore")
	fmt.Println("   2. Tidak bisa mengganti storage (misal ke database) tanpa ubah kode service")
	fmt.Println("   3. Sulit melakukan unit testing karena tidak bisa di-mock")
	fmt.Println("   4. Tidak fleksibel untuk extend atau ganti implementasi")

	_ = user2 // suppress unused variable warning
}
