package database

import (
	"MemberLink-Lite/logger"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Migration 数据库迁移记录
type Migration struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Version   string    `gorm:"uniqueIndex;size:50;not null;comment:迁移版本"`
	Name      string    `gorm:"size:255;not null;comment:迁移名称"`
	Applied   bool      `gorm:"default:false;comment:是否已应用"`
	AppliedAt time.Time `gorm:"comment:应用时间"`
	CreatedAt time.Time `gorm:"autoCreateTime;comment:创建时间"`
}

// MigrationFunc 迁移函数类型
type MigrationFunc func(*gorm.DB) error

// MigrationItem 迁移项
type MigrationItem struct {
	Version string
	Name    string
	Up      MigrationFunc
	Down    MigrationFunc
}

// Migrator 数据库迁移器
type Migrator struct {
	db         *gorm.DB
	migrations []MigrationItem
}

// NewMigrator 创建新的迁移器
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: make([]MigrationItem, 0),
	}
}

// AddMigration 添加迁移
func (m *Migrator) AddMigration(version, name string, up, down MigrationFunc) {
	m.migrations = append(m.migrations, MigrationItem{
		Version: version,
		Name:    name,
		Up:      up,
		Down:    down,
	})
}

// Init 初始化迁移表
func (m *Migrator) Init() error {
	if err := m.db.AutoMigrate(&Migration{}); err != nil {
		logger.Error("Failed to create migration table:", err)
		return err
	}
	logger.Info("Migration table initialized")
	return nil
}

// Migrate 执行迁移
func (m *Migrator) Migrate() error {
	logger.Info("Starting database migration...")

	for _, migration := range m.migrations {
		// 检查迁移是否已经应用
		var existingMigration Migration
		result := m.db.Where("version = ?", migration.Version).First(&existingMigration)

		if result.Error == nil && existingMigration.Applied {
			logger.Info(fmt.Sprintf("Migration %s already applied, skipping", migration.Version))
			continue
		}

		logger.Info(fmt.Sprintf("Applying migration %s: %s", migration.Version, migration.Name))

		// 开始事务
		tx := m.db.Begin()
		if tx.Error != nil {
			logger.Error("Failed to begin transaction:", tx.Error)
			return tx.Error
		}

		// 执行迁移
		if err := migration.Up(tx); err != nil {
			tx.Rollback()
			logger.Error(fmt.Sprintf("Failed to apply migration %s:", migration.Version), err)
			return err
		}

		// 记录迁移
		migrationRecord := Migration{
			Version:   migration.Version,
			Name:      migration.Name,
			Applied:   true,
			AppliedAt: time.Now(),
		}

		if result.Error == gorm.ErrRecordNotFound {
			// 创建新记录
			if err := tx.Create(&migrationRecord).Error; err != nil {
				tx.Rollback()
				logger.Error("Failed to create migration record:", err)
				return err
			}
		} else {
			// 更新现有记录
			if err := tx.Model(&existingMigration).Updates(map[string]interface{}{
				"applied":    true,
				"applied_at": time.Now(),
			}).Error; err != nil {
				tx.Rollback()
				logger.Error("Failed to update migration record:", err)
				return err
			}
		}

		// 提交事务
		if err := tx.Commit().Error; err != nil {
			logger.Error("Failed to commit migration transaction:", err)
			return err
		}

		logger.Info(fmt.Sprintf("Migration %s applied successfully", migration.Version))
	}

	logger.Info("Database migration completed")
	return nil
}

// Rollback 回滚迁移
func (m *Migrator) Rollback(version string) error {
	logger.Info(fmt.Sprintf("Rolling back migration %s", version))

	// 查找迁移
	var migration *MigrationItem
	for _, m := range m.migrations {
		if m.Version == version {
			migration = &m
			break
		}
	}

	if migration == nil {
		return fmt.Errorf("migration %s not found", version)
	}

	// 检查迁移是否已应用
	var migrationRecord Migration
	if err := m.db.Where("version = ? AND applied = ?", version, true).First(&migrationRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("migration %s not applied", version)
		}
		return err
	}

	// 开始事务
	tx := m.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 执行回滚
	if err := migration.Down(tx); err != nil {
		tx.Rollback()
		logger.Error(fmt.Sprintf("Failed to rollback migration %s:", version), err)
		return err
	}

	// 更新迁移记录
	if err := tx.Model(&migrationRecord).Update("applied", false).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to update migration record:", err)
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		logger.Error("Failed to commit rollback transaction:", err)
		return err
	}

	logger.Info(fmt.Sprintf("Migration %s rolled back successfully", version))
	return nil
}

// GetAppliedMigrations 获取已应用的迁移列表
func (m *Migrator) GetAppliedMigrations() ([]Migration, error) {
	var migrations []Migration
	if err := m.db.Where("applied = ?", true).Order("applied_at ASC").Find(&migrations).Error; err != nil {
		return nil, err
	}
	return migrations, nil
}

// GetPendingMigrations 获取待应用的迁移列表
func (m *Migrator) GetPendingMigrations() ([]MigrationItem, error) {
	var appliedVersions []string
	if err := m.db.Model(&Migration{}).Where("applied = ?", true).Pluck("version", &appliedVersions).Error; err != nil {
		return nil, err
	}

	appliedMap := make(map[string]bool)
	for _, version := range appliedVersions {
		appliedMap[version] = true
	}

	var pending []MigrationItem
	for _, migration := range m.migrations {
		if !appliedMap[migration.Version] {
			pending = append(pending, migration)
		}
	}

	return pending, nil
}
