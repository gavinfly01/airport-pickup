package mysqlrepo

import (
	userentity "github.com/gavin/airport-pickup/internal/domain/user/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

func newTestDBUser() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&Passenger{}, &Driver{})
	return db
}

func TestPassengerRepository_SaveAndGet(t *testing.T) {
	db := newTestDBUser()
	repo := NewPassengerRepository(db)
	p := &userentity.Passenger{ID: "p1", Name: "Alice", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	err := repo.Save(p)
	assert.NoError(t, err)
	got, err := repo.GetByID("p1")
	assert.NoError(t, err)
	assert.Equal(t, "p1", got.ID)
	assert.Equal(t, "Alice", got.Name)
}

func TestPassengerRepository_Save_Invalid(t *testing.T) {
	db := newTestDBUser()
	repo := NewPassengerRepository(db)
	err := repo.Save(nil)
	assert.Error(t, err)
	err = repo.Save(&userentity.Passenger{})
	assert.Error(t, err)
}

func TestPassengerRepository_GetByID_NotFound(t *testing.T) {
	db := newTestDBUser()
	repo := NewPassengerRepository(db)
	_, err := repo.GetByID("not_exist")
	assert.Error(t, err)
}

func TestDriverRepository_SaveAndGet(t *testing.T) {
	db := newTestDBUser()
	repo := NewDriverRepository(db)
	d := &userentity.Driver{ID: "d1", Name: "Bob", Rating: 4.8, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	err := repo.Save(d)
	assert.NoError(t, err)
	// 这里只测试保存，GetByID 可按需补充
}

func TestDriverRepository_Save_Invalid(t *testing.T) {
	db := newTestDBUser()
	repo := NewDriverRepository(db)
	err := repo.Save(nil)
	assert.Error(t, err)
	err = repo.Save(&userentity.Driver{})
	assert.Error(t, err)
}
