package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dujiao-next/internal/cache"
	"github.com/dujiao-next/internal/constants"
	"github.com/dujiao-next/internal/models"

	"golang.org/x/crypto/bcrypt"
)

// LoginWithTelegramInput Telegram 登录输入
type LoginWithTelegramInput struct {
	Payload TelegramLoginPayload
	Context context.Context
}

// BindTelegramInput 绑定 Telegram 输入
type BindTelegramInput struct {
	UserID  uint
	Payload TelegramLoginPayload
	Context context.Context
}

// LoginWithTelegram Telegram 登录
func (s *UserAuthService) LoginWithTelegram(input LoginWithTelegramInput) (*models.User, string, time.Time, error) {
	if s.telegramAuthService == nil || s.userOAuthIdentityRepo == nil {
		return nil, "", time.Time{}, ErrTelegramAuthConfigInvalid
	}
	ctx := input.Context
	if ctx == nil {
		ctx = context.Background()
	}
	verified, err := s.telegramAuthService.VerifyLogin(ctx, input.Payload)
	if err != nil {
		return nil, "", time.Time{}, err
	}

	identity, err := s.userOAuthIdentityRepo.GetByProviderUserID(verified.Provider, verified.ProviderUserID)
	if err != nil {
		return nil, "", time.Time{}, err
	}

	var user *models.User
	if identity != nil {
		user, err = s.getActiveUserByID(identity.UserID)
		if err != nil {
			return nil, "", time.Time{}, err
		}
		identityChanged := applyTelegramIdentity(verified, identity)
		if identityChanged {
			identity.UpdatedAt = time.Now()
			if err := s.userOAuthIdentityRepo.Update(identity); err != nil {
				return nil, "", time.Time{}, err
			}
		}
	} else {
		user, err = s.findOrCreateTelegramUser(verified)
		if err != nil {
			return nil, "", time.Time{}, err
		}
		identity = &models.UserOAuthIdentity{
			UserID:         user.ID,
			Provider:       verified.Provider,
			ProviderUserID: verified.ProviderUserID,
			Username:       verified.Username,
			AvatarURL:      verified.AvatarURL,
			AuthAt:         &verified.AuthAt,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		if err := s.userOAuthIdentityRepo.Create(identity); err != nil {
			existing, getErr := s.userOAuthIdentityRepo.GetByProviderUserID(verified.Provider, verified.ProviderUserID)
			if getErr != nil {
				return nil, "", time.Time{}, err
			}
			if existing == nil {
				return nil, "", time.Time{}, err
			}
			identity = existing
			user, err = s.getActiveUserByID(existing.UserID)
			if err != nil {
				return nil, "", time.Time{}, err
			}
		}
	}

	token, expiresAt, err := s.GenerateUserJWT(user, 0)
	if err != nil {
		return nil, "", time.Time{}, err
	}

	now := time.Now()
	user.LastLoginAt = &now
	user.UpdatedAt = now
	if err := s.userRepo.Update(user); err != nil {
		return nil, "", time.Time{}, err
	}
	_ = cache.SetUserAuthState(context.Background(), cache.BuildUserAuthState(user))
	return user, token, expiresAt, nil
}

// BindTelegram 绑定 Telegram
func (s *UserAuthService) BindTelegram(input BindTelegramInput) (*models.UserOAuthIdentity, error) {
	if input.UserID == 0 {
		return nil, ErrNotFound
	}
	if s.telegramAuthService == nil || s.userOAuthIdentityRepo == nil {
		return nil, ErrTelegramAuthConfigInvalid
	}
	ctx := input.Context
	if ctx == nil {
		ctx = context.Background()
	}
	verified, err := s.telegramAuthService.VerifyLogin(ctx, input.Payload)
	if err != nil {
		return nil, err
	}
	if _, err := s.getActiveUserByID(input.UserID); err != nil {
		return nil, err
	}

	occupied, err := s.userOAuthIdentityRepo.GetByProviderUserID(verified.Provider, verified.ProviderUserID)
	if err != nil {
		return nil, err
	}
	if occupied != nil && occupied.UserID != input.UserID {
		return nil, ErrUserOAuthIdentityExists
	}

	current, err := s.userOAuthIdentityRepo.GetByUserProvider(input.UserID, verified.Provider)
	if err != nil {
		return nil, err
	}
	if current != nil && current.ProviderUserID != verified.ProviderUserID {
		return nil, ErrUserOAuthAlreadyBound
	}
	if current == nil {
		current = &models.UserOAuthIdentity{
			UserID:         input.UserID,
			Provider:       verified.Provider,
			ProviderUserID: verified.ProviderUserID,
			Username:       verified.Username,
			AvatarURL:      verified.AvatarURL,
			AuthAt:         &verified.AuthAt,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		if err := s.userOAuthIdentityRepo.Create(current); err != nil {
			return nil, err
		}
		return current, nil
	}

	if applyTelegramIdentity(verified, current) {
		current.UpdatedAt = time.Now()
		if err := s.userOAuthIdentityRepo.Update(current); err != nil {
			return nil, err
		}
	}
	return current, nil
}

// UnbindTelegram 解绑 Telegram
func (s *UserAuthService) UnbindTelegram(userID uint) error {
	if userID == 0 {
		return ErrNotFound
	}
	if s.userOAuthIdentityRepo == nil {
		return ErrTelegramAuthConfigInvalid
	}
	user, err := s.getActiveUserByID(userID)
	if err != nil {
		return err
	}
	mode, err := s.ResolveEmailChangeMode(user)
	if err != nil {
		return err
	}
	if mode == EmailChangeModeBindOnly {
		return ErrTelegramUnbindRequiresEmail
	}
	identity, err := s.userOAuthIdentityRepo.GetByUserProvider(userID, constants.UserOAuthProviderTelegram)
	if err != nil {
		return err
	}
	if identity == nil {
		return ErrUserOAuthNotBound
	}
	return s.userOAuthIdentityRepo.DeleteByID(identity.ID)
}

// GetTelegramBinding 获取 Telegram 绑定
func (s *UserAuthService) GetTelegramBinding(userID uint) (*models.UserOAuthIdentity, error) {
	if userID == 0 {
		return nil, ErrNotFound
	}
	if s.userOAuthIdentityRepo == nil {
		return nil, ErrTelegramAuthConfigInvalid
	}
	return s.userOAuthIdentityRepo.GetByUserProvider(userID, constants.UserOAuthProviderTelegram)
}

func (s *UserAuthService) getActiveUserByID(userID uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNotFound
	}
	if strings.ToLower(strings.TrimSpace(user.Status)) != constants.UserStatusActive {
		return nil, ErrUserDisabled
	}
	return user, nil
}

func (s *UserAuthService) findOrCreateTelegramUser(verified *TelegramIdentityVerified) (*models.User, error) {
	if verified == nil {
		return nil, ErrTelegramAuthPayloadInvalid
	}
	email := buildTelegramPlaceholderEmail(verified.ProviderUserID)
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user != nil {
		if strings.ToLower(strings.TrimSpace(user.Status)) != constants.UserStatusActive {
			return nil, ErrUserDisabled
		}
		return user, nil
	}

	randomSuffix, err := randomNumericCode(16)
	if err != nil {
		return nil, err
	}
	passwordSeed := fmt.Sprintf("tg_%s_%s", verified.ProviderUserID, randomSuffix)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordSeed), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user = &models.User{
		Email:                 email,
		PasswordHash:          string(hashedPassword),
		PasswordSetupRequired: true,
		DisplayName:           resolveTelegramDisplayName(verified),
		Status:                constants.UserStatusActive,
		LastLoginAt:           &now,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func applyTelegramIdentity(verified *TelegramIdentityVerified, identity *models.UserOAuthIdentity) bool {
	if verified == nil || identity == nil {
		return false
	}
	changed := false
	if identity.Provider == "" {
		identity.Provider = verified.Provider
		changed = true
	}
	if identity.ProviderUserID == "" {
		identity.ProviderUserID = verified.ProviderUserID
		changed = true
	}
	if identity.Username != verified.Username {
		identity.Username = verified.Username
		changed = true
	}
	if identity.AvatarURL != verified.AvatarURL {
		identity.AvatarURL = verified.AvatarURL
		changed = true
	}
	if identity.AuthAt == nil || !identity.AuthAt.Equal(verified.AuthAt) {
		authAt := verified.AuthAt
		identity.AuthAt = &authAt
		changed = true
	}
	return changed
}

func buildTelegramPlaceholderEmail(providerUserID string) string {
	normalizedID := strings.TrimSpace(providerUserID)
	if normalizedID == "" {
		normalizedID = "unknown"
	}
	return fmt.Sprintf("%s%s%s", telegramPlaceholderEmailPrefix, normalizedID, telegramPlaceholderEmailDomain)
}

func isTelegramPlaceholderEmail(email string) bool {
	normalized := strings.ToLower(strings.TrimSpace(email))
	if normalized == "" {
		return false
	}
	return strings.HasPrefix(normalized, telegramPlaceholderEmailPrefix) &&
		strings.HasSuffix(normalized, telegramPlaceholderEmailDomain)
}

func resolveTelegramDisplayName(verified *TelegramIdentityVerified) string {
	if verified == nil {
		return "Telegram User"
	}
	fullName := strings.TrimSpace(strings.TrimSpace(verified.FirstName) + " " + strings.TrimSpace(verified.LastName))
	if fullName != "" {
		return fullName
	}
	if strings.TrimSpace(verified.Username) != "" {
		return verified.Username
	}
	if strings.TrimSpace(verified.ProviderUserID) != "" {
		return fmt.Sprintf("telegram_%s", strings.TrimSpace(verified.ProviderUserID))
	}
	return "Telegram User"
}
