// Package main mendemonstrasikan CRUD DENGAN menggunakan interface.
// Ini adalah pendekatan yang lebih fleksibel dan mengikuti prinsip Dependency Inversion.
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
// INTERFACE DEFINITION
// =============================================================================

// UserRepository adalah INTERFACE yang mendefinisikan kontrak untuk operasi CRUD User.
// Interface ini adalah "kontrak" yang harus dipenuhi oleh SEMUA implementasi storage.
// Dengan interface ini, kita bisa mengganti storage kapanpun tanpa ubah business logic!
type UserRepository interface {
	Create(name, email string) User
	GetByID(id int) (User, error)
	GetAll() []User
	Update(id int, name, email string) (User, error)
	Delete(id int) error
}

// =============================================================================
// IMPLEMENTASI 1: Memory Storage
// =============================================================================

// MemoryUserStore adalah implementasi UserRepository menggunakan memory (slice).
// Cocok untuk development, testing, atau aplikasi sederhana.
type MemoryUserStore struct {
	users  []User
	nextID int
}

// NewMemoryUserStore membuat instance baru dari MemoryUserStore
func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{
		users:  make([]User, 0),
		nextID: 1,
	}
}

// Implementasi semua method dari UserRepository interface

func (s *MemoryUserStore) Create(name, email string) User {
	user := User{
		ID:    s.nextID,
		Name:  name,
		Email: email,
	}
	s.users = append(s.users, user)
	s.nextID++
	fmt.Printf("   [MemoryStore] User '%s' disimpan di memory\n", name)
	return user
}

func (s *MemoryUserStore) GetByID(id int) (User, error) {
	for _, user := range s.users {
		if user.ID == id {
			return user, nil
		}
	}
	return User{}, errors.New("user tidak ditemukan")
}

func (s *MemoryUserStore) GetAll() []User {
	return s.users
}

func (s *MemoryUserStore) Update(id int, name, email string) (User, error) {
	for i, user := range s.users {
		if user.ID == id {
			s.users[i].Name = name
			s.users[i].Email = email
			fmt.Printf("   [MemoryStore] User ID %d diupdate di memory\n", id)
			return s.users[i], nil
		}
	}
	return User{}, errors.New("user tidak ditemukan")
}

func (s *MemoryUserStore) Delete(id int) error {
	for i, user := range s.users {
		if user.ID == id {
			s.users = append(s.users[:i], s.users[i+1:]...)
			fmt.Printf("   [MemoryStore] User ID %d dihapus dari memory\n", id)
			return nil
		}
	}
	return errors.New("user tidak ditemukan")
}

// =============================================================================
// IMPLEMENTASI 2: Mock Storage (untuk Testing)
// =============================================================================

// MockUserStore adalah implementasi UserRepository untuk testing.
// Bisa dikonfigurasi untuk return nilai tertentu atau error.
type MockUserStore struct {
	MockUsers     []User
	ShouldFail    bool
	FailMessage   string
	CreateCalled  int
	GetByIDCalled int
	UpdateCalled  int
	DeleteCalled  int
}

// NewMockUserStore membuat instance baru MockUserStore
func NewMockUserStore() *MockUserStore {
	return &MockUserStore{
		MockUsers: make([]User, 0),
	}
}

func (m *MockUserStore) Create(name, email string) User {
	m.CreateCalled++
	user := User{ID: len(m.MockUsers) + 1, Name: name, Email: email}
	m.MockUsers = append(m.MockUsers, user)
	fmt.Printf("   [MockStore] CREATE dipanggil (call ke-%d)\n", m.CreateCalled)
	return user
}

func (m *MockUserStore) GetByID(id int) (User, error) {
	m.GetByIDCalled++
	fmt.Printf("   [MockStore] GET BY ID dipanggil (call ke-%d)\n", m.GetByIDCalled)
	if m.ShouldFail {
		return User{}, errors.New(m.FailMessage)
	}
	for _, user := range m.MockUsers {
		if user.ID == id {
			return user, nil
		}
	}
	return User{}, errors.New("user tidak ditemukan")
}

func (m *MockUserStore) GetAll() []User {
	return m.MockUsers
}

func (m *MockUserStore) Update(id int, name, email string) (User, error) {
	m.UpdateCalled++
	fmt.Printf("   [MockStore] UPDATE dipanggil (call ke-%d)\n", m.UpdateCalled)
	if m.ShouldFail {
		return User{}, errors.New(m.FailMessage)
	}
	for i, user := range m.MockUsers {
		if user.ID == id {
			m.MockUsers[i].Name = name
			m.MockUsers[i].Email = email
			return m.MockUsers[i], nil
		}
	}
	return User{}, errors.New("user tidak ditemukan")
}

func (m *MockUserStore) Delete(id int) error {
	m.DeleteCalled++
	fmt.Printf("   [MockStore] DELETE dipanggil (call ke-%d)\n", m.DeleteCalled)
	if m.ShouldFail {
		return errors.New(m.FailMessage)
	}
	for i, user := range m.MockUsers {
		if user.ID == id {
			m.MockUsers = append(m.MockUsers[:i], m.MockUsers[i+1:]...)
			return nil
		}
	}
	return errors.New("user tidak ditemukan")
}

// =============================================================================
// SERVICE LAYER (Menggunakan Interface!)
// =============================================================================

// UserService adalah layer bisnis logic yang menggunakan UserRepository INTERFACE.
// KEUNTUNGAN: Service ini loosely coupled - bisa bekerja dengan implementasi apapun!
type UserService struct {
	repo UserRepository // <-- Dependency ke INTERFACE, bukan implementasi konkret!
}

// NewUserService membuat instance baru UserService.
// Parameter menerima UserRepository (interface), bukan konkret struct.
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// RegisterUser mendaftarkan user baru dengan validasi
func (s *UserService) RegisterUser(name, email string) (User, error) {
	if name == "" || email == "" {
		return User{}, errors.New("nama dan email harus diisi")
	}
	user := s.repo.Create(name, email)
	fmt.Printf("   ✅ User '%s' berhasil didaftarkan dengan ID %d\n", name, user.ID)
	return user, nil
}

// GetUserProfile mengambil profil user
func (s *UserService) GetUserProfile(id int) (User, error) {
	return s.repo.GetByID(id)
}

// UpdateUserProfile mengupdate profil user
func (s *UserService) UpdateUserProfile(id int, name, email string) (User, error) {
	return s.repo.Update(id, name, email)
}

// DeleteUser menghapus user
func (s *UserService) DeleteUser(id int) error {
	return s.repo.Delete(id)
}

// =============================================================================
// MAIN - Demo Penggunaan
// =============================================================================

func main() {
	fmt.Println("===========================")
	fmt.Println("CRUD DENGAN INTERFACE")
	fmt.Println("===========================")
	fmt.Println()

	// =========================================================================
	// DEMO 1: Menggunakan Memory Store
	// =========================================================================
	fmt.Println("🔹 DEMO 1: Menggunakan MemoryUserStore")
	fmt.Println("-" + "--------------------------------------")

	memoryStore := NewMemoryUserStore()
	service := NewUserService(memoryStore) // Inject MemoryUserStore

	fmt.Println("\n📝 CREATE:")
	service.RegisterUser("Budi", "budi@email.com")
	service.RegisterUser("Ani", "ani@email.com")

	fmt.Println("\n📖 READ ALL:")
	for _, user := range memoryStore.GetAll() {
		fmt.Printf("   - ID: %d, Nama: %s, Email: %s\n", user.ID, user.Name, user.Email)
	}

	fmt.Println("\n✏️ UPDATE:")
	service.UpdateUserProfile(1, "Budi Santoso", "budi.s@email.com")

	fmt.Println("\n🗑️ DELETE:")
	service.DeleteUser(2)

	fmt.Println("\n📋 HASIL AKHIR (Memory Store):")
	for _, user := range memoryStore.GetAll() {
		fmt.Printf("   - ID: %d, Nama: %s, Email: %s\n", user.ID, user.Name, user.Email)
	}

	// =========================================================================
	// DEMO 2: Menggunakan Mock Store (untuk Testing)
	// =========================================================================
	fmt.Println("\n")
	fmt.Println("🔹 DEMO 2: Menggunakan MockUserStore (untuk Testing)")
	fmt.Println("-" + "-------------------------------------------------")

	mockStore := NewMockUserStore()
	serviceWithMock := NewUserService(mockStore) // Inject MockUserStore - KODE SERVICE SAMA!

	fmt.Println("\n📝 CREATE (dengan Mock):")
	serviceWithMock.RegisterUser("Test User", "test@email.com")

	fmt.Println("\n📖 READ (dengan Mock):")
	_, _ = serviceWithMock.GetUserProfile(1)

	fmt.Println("\n📊 Statistik Mock:")
	fmt.Printf("   - Create dipanggil: %d kali\n", mockStore.CreateCalled)
	fmt.Printf("   - GetByID dipanggil: %d kali\n", mockStore.GetByIDCalled)
	fmt.Printf("   - Update dipanggil: %d kali\n", mockStore.UpdateCalled)
	fmt.Printf("   - Delete dipanggil: %d kali\n", mockStore.DeleteCalled)

	// =========================================================================
	// KEUNTUNGAN MENGGUNAKAN INTERFACE
	// =========================================================================
	fmt.Println("\n")
	fmt.Println("✅ KEUNTUNGAN DENGAN INTERFACE:")
	fmt.Println("   1. UserService TIDAK PERLU DIUBAH ketika ganti storage")
	fmt.Println("   2. Bisa inject MemoryStore, MockStore, atau DatabaseStore")
	fmt.Println("   3. Mudah melakukan unit testing dengan mock")
	fmt.Println("   4. Mengikuti prinsip Dependency Inversion (SOLID)")
	fmt.Println("   5. Loose coupling = kode lebih maintainable")
	fmt.Println()
	fmt.Println("💡 CATATAN:")
	fmt.Println("   - Perhatikan bahwa UserService menerima interface UserRepository")
	fmt.Println("   - Baik MemoryUserStore maupun MockUserStore bisa digunakan")
	fmt.Println("   - Di production, bisa dibuat DatabaseUserStore dengan interface yang sama")
}
