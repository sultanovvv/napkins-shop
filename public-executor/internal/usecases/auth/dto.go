package auth

// RegisterInUDTO — регистрация по email+паролю.
type RegisterInUDTO struct {
	Email    string
	Password string
}

// LoginInUDTO — логин по email+паролю. GuestToken (если фронт передал
// guest_cart cookie) используется для merge гостевой корзины в пользовательскую.
type LoginInUDTO struct {
	Provider   string // если пусто — ProviderLocal
	Email      string
	Password   string
	GuestToken string
	UserAgent  string
	IP         string
}

// RefreshInUDTO — обмен refresh-секрета на новую пару.
type RefreshInUDTO struct {
	RefreshSecret string
	UserAgent     string
	IP            string
}

// LogoutInUDTO — единичный logout по текущему refresh-секрету.
type LogoutInUDTO struct {
	RefreshSecret string
}

// RequestPasswordResetInUDTO — отправка запроса на восстановление.
// Реальный email пока не шлём (нет транспорта); secret возвращается из
// usecase только для интеграционных тестов — handler его наружу не отдаёт.
type RequestPasswordResetInUDTO struct {
	Email string
}

// ConfirmPasswordResetInUDTO — подтверждение восстановления.
type ConfirmPasswordResetInUDTO struct {
	Secret      string
	NewPassword string
}

// ConfirmEmailVerificationInUDTO — подтверждение email по secret из письма.
type ConfirmEmailVerificationInUDTO struct {
	Secret string
}