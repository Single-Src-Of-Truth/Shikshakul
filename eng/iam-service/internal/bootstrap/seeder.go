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
		{ID: "iam:tenants:create", Description: "Can onboard new schools"},
		{ID: "iam:tenants:read", Description: "Can view school details"},
		{ID: "iam:tenants:update", Description: "Can modify school details"},

		{ID: "iam:users:execute", Description: "Can change staff status (suspend/activate)"},
		{ID: "iam:users:delete", Description: "Can remove staff accounts"},

		{ID: "iam:roles:read", Description: "Can view custom roles"},
		{ID: "iam:roles:create", Description: "Can create custom roles"},
		{ID: "iam:roles:update", Description: "Can modify custom roles"},
		{ID: "iam:roles:delete", Description: "Can delete custom roles"},

		{ID: "iam:permissions:read", Description: "Can view available permissions"},
		{ID: "iam:permissions:create", Description: "Can bulk seed new permissions"},

		{ID: "iam:invites:execute", Description: "Can generate and extend invitations"},
		{ID: "iam:invites:delete", Description: "Can cancel invitations"},

		{ID: "acad:students:read", Description: "Can view student records"},
		{ID: "acad:students:create", Description: "Can manage student admissions"},
		{ID: "utility:files:create", Description: "Can upload documents"},
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
			FirstName:         "Golu",
			LastName:          "Billauta",
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
