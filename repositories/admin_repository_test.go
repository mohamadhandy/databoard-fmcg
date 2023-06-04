package repositories

import (
	"klikdaily-databoard/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// assert mockDB
func assertMockDB() (mockgormDB *gorm.DB, mockSqlmock sqlmock.Sqlmock) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	// Create a new GORM database connection using the mock database
	db, _ := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	return db, mock
}

func Test_Repositories_CreateAdmin_ShouldReturnAdmin(t *testing.T) {
	// Arrage
	db, mock := assertMockDB()
	repoAdmin := adminRepository{
		db: db,
	}
	newAdmin := &models.Admin{
		Name:  "John Cena",
		Email: "johncena@gmail.com",
	}
	adminReq := models.AdminRequest{
		Name:  "John Cena",
		Email: "johncena@gmail.com",
	}
	expectedQuery := "INSERT INTO `Admin` (`name`,`email`) VALUES (?,?)"
	mockQueryResult := sqlmock.NewResult(1, 1)

	mock.ExpectExec(expectedQuery).WithArgs(newAdmin.Name, newAdmin.Email).WillReturnResult(mockQueryResult)

	// call the create method on the repository
	resultChannel := repoAdmin.CreateAdmin(adminReq)
	assert.Equal(t, resultChannel, resultChannel)
}
