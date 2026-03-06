package bootstrap

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/domain"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/crypto"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB, logger *zap.Logger) {
	logger.Info("Checking database seeds...")

	permissions := []domain.Permission{
		{ID: "iam:tenants:write", Description: "Can onboard new schools"},
		{ID: "iam:roles:write", Description: "Can create custom roles"},
		{ID: "acad:students:read", Description: "Can view student records"},
		{ID: "acad:students:write", Description: "Can manage student admissions"},
		{ID: "acad:fees:write", Description: "Can configure fee heads"},
		{ID: "utility:docs:upload", Description: "Can upload documents to S3"},
	}

	for _, p := range permissions {
		db.FirstOrCreate(&p, domain.Permission{ID: p.ID})
	}

	superAdmin := domain.Role{Name: "SUPER_ADMIN", Description: "Global Admin", IsSystem: true}
	db.Where("name = ? AND tenant_id IS NULL", superAdmin.Name).FirstOrCreate(&superAdmin)

	var allPerms []domain.Permission
	db.Find(&allPerms)
	db.Model(&superAdmin).Association("Permissions").Replace(&allPerms)

	rootEmail := "root@shikshakul.com"
	var count int64
	db.Model(&domain.User{}).Where("primary_identifier = ?", rootEmail).Count(&count)

	if count == 0 {
		bytes := make([]byte, 4)
		rand.Read(bytes)
		plainPassword := hex.EncodeToString(bytes)
		hashedPassword, _ := crypto.HashPassword(plainPassword)

		rootUser := domain.User{
			IdentifierType:    domain.IdentifierEmail,
			PrimaryIdentifier: rootEmail,
			PasswordHash:      hashedPassword,
			FirstName:         "System",
			LastName:          "Administrator",
			Status:            domain.StatusActive,
		}

		db.Create(&rootUser)
		db.Model(&rootUser).Association("Roles").Append(&superAdmin)

		logger.Info("INITIAL SUPER ADMIN CREATED",
			zap.String("email", rootEmail),
			zap.String("password", plainPassword),
		)
	} else {
		logger.Info("Root user already exists, skipping creation.")
	}
}
