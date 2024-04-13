package password

import (
	"github.com/mistandok/auth/internal/config"
	"golang.org/x/crypto/bcrypt"
)

// Manager мэнеджер для работы с паролями.
type Manager struct {
	passwordConfig *config.PasswordConfig
}

// NewManager ..
func NewManager(passwordConfig *config.PasswordConfig) *Manager {
	return &Manager{passwordConfig: passwordConfig}
}

// HashPassword хэширует пассворд с учетом соли.
func (m *Manager) HashPassword(password string) (string, error) {
	passBytes, err := bcrypt.GenerateFromPassword([]byte(m.passWithSalt(password)), bcrypt.MinCost)
	return string(passBytes), err
}

// CheckPasswordHash проверяет пароль с учетом соли.
func (m *Manager) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(m.passWithSalt(password)))
	return err == nil
}

func (m *Manager) passWithSalt(password string) string {
	return password + m.passwordConfig.PasswordSalt
}
