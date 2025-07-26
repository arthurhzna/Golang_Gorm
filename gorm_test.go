package go_gorm

import (
	// "context"
	"fmt"
	"os"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"strconv"
	"testing"
	"time"
	"github.com/joho/godotenv"
)

func OpenConnection() *gorm.DB {
    err := godotenv.Load()
    if err != nil {
        panic("Error loading .env file: " + err.Error())
    }
	
    databaseURL := os.Getenv("DATABASE_URL")
    if databaseURL == "" {
        panic("DATABASE_URL environment variable is required. Please create .env file with DATABASE_URL")
    }
    
    dialect := postgres.Open(databaseURL)
    
    db, err := gorm.Open(dialect, &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        panic("failed to connect database")
    }

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	return db
}

var db = OpenConnection()

func TestOpenConnection(t *testing.T) {
	assert.NotNil(t, db)
}

func TestCreateUser(t *testing.T) {
	user := User{
		ID:       "1",
		Password: "akubang",
		FirstName: "Arthur",
		MiddleName: "Middle",
		LastName:   "Hozanna",
	}
	response := db.Create(&user)
	assert.Nil(t, response.Error)
	assert.Equal(t, int64(1), response.RowsAffected)
}

func TestBatchInsert(t *testing.T) {
	var users []User
	for i := 2; i < 10; i++ {
		users = append(users, User{
			ID:        strconv.Itoa(i),
			Password:  "rahasia",
			FirstName: "User " + strconv.Itoa(i),  
			
		})
	}

	result := db.Create(&users)
	assert.Nil(t, result.Error)
	assert.Equal(t, 8, int(result.RowsAffected))
}

func TestTransactionSuccess(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&User{ID: "10", Password: "rahasia", FirstName: "User 10"}).Error
		if err != nil {
			return err
		}

		err = tx.Create(&User{ID: "11", Password: "rahasia", FirstName: "User 11"}).Error
		if err != nil {
			return err
		}

		err = tx.Create(&User{ID: "12", Password: "rahasia", FirstName: "User 12"}).Error
		if err != nil {
			return err
		}

		return nil
	})

	assert.Nil(t, err)
}
func TestTransactionRollback(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&User{ID: "13", Password: "rahasia", FirstName: "User 13"}).Error
		if err != nil {
			return err
		}

		err = tx.Create(&User{ID: "11", Password: "rahasia", FirstName: "User 11"}).Error
		if err != nil {
			return err
		}

		return nil
	})

	assert.NotNil(t, err)
}

func TestManualTransactionSuccess(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	err := tx.Create(&User{ID: "13", Password: "rahasia", FirstName: "User 13"}).Error
	assert.Nil(t, err)

	err = tx.Create(&User{ID: "14", Password: "rahasia", FirstName: "User 14"}).Error
	assert.Nil(t, err)

	if err == nil {
		tx.Commit()
	}
}

func TestManualTransactionError(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	err := tx.Create(&User{ID: "15", Password: "rahasia", FirstName: "User 15"}).Error
	assert.Nil(t, err)

	err = tx.Create(&User{ID: "14", Password: "rahasia", FirstName: "User 14"}).Error
	assert.Nil(t, err)

	if err == nil {
		tx.Commit()
	}
}


func TestQuerySingleObject(t *testing.T) {
	user := User{}
	err := db.First(&user).Error
	assert.Nil(t, err)
	assert.Equal(t, "1", user.ID)

	user = User{}
	err = db.Last(&user).Error
	assert.Nil(t, err)
	assert.Equal(t, "9", user.ID)
}

func TestQuerySingleObjectInlineCondition(t *testing.T) {
	user := User{}
	err := db.Take(&user, "id = ?", "5").Error
	fmt.Println(user)
	assert.Nil(t, err)
	assert.Equal(t, "5", user.ID)
	assert.Equal(t, "User 5", user.FirstName)
}

func TestQueryAllObjects(t *testing.T) {
	var users []User
	err := db.Find(&users, "id in ?", []string{"1", "2", "3", "4"}).Error
	fmt.Println(users)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(users))
}

func TestQueryCondition(t *testing.T) {
	var users []User
	err := db.Where("first_name like ?", "%User%").Where("password = ?", "rahasia").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 13, len(users))
}

func stringPtr(s string) *string {
    return &s
}


func TestInsertUserWithTodos(t *testing.T) {
    user := User{
        ID:        "user_456",
        Password:  "password456", 
        FirstName: "Jane",
        Todos: []Todo{
            {
                Title:       "Task 1",
                Description: stringPtr("First task"),
            },
            {
                Title:       "Task 2", 
                Description: stringPtr("Second task"),
            },
        },
    }
    
    result := db.Create(&user)
    assert.Nil(t, result.Error)
}

func TestQueryUserWithTodos(t *testing.T) {
    var user User
    err := db.Preload("Todos").Where("id = ?", "user_456").First(&user).Error
    
    assert.Nil(t, err)
    assert.Equal(t, "user_456", user.ID)
    assert.Equal(t, "Jane", user.FirstName)
    assert.Equal(t, 2, len(user.Todos)) 
    
    t.Logf("User: %s has %d todos", user.FirstName, len(user.Todos))
    for i, todo := range user.Todos {
        t.Logf("Todo %d: %s - %s", i+1, todo.Title, *todo.Description)
    }
}

func TestAutoCreateUpdate(t *testing.T) {
	user := User{
		ID:       "20",
		Password: "rahasia",
		FirstName: "User 20",
		MiddleName: "Middle",
		LastName: "Hozanna",
		Wallets: []Wallet{
			{
				ID: "20",
				UserID: "20",
				Balance: 1000000,
			},
		},
	}

	err := db.Create(&user).Error
	assert.Nil(t, err)
}

func TestSkipAutoCreateUpdate(t *testing.T) {
	user := User{
		ID:       "21",
		Password: "rahasia",
		FirstName: "User 21",
		MiddleName: "Middle",
		LastName: "Hozanna",
		Wallets: []Wallet{
			{ 
				ID:      "21",
				UserID:  "21",
				Balance: 1000000,
			},
		},
	}

	err := db.Omit(clause.Associations).Create(&user).Error
	assert.Nil(t, err)
}


func TestUserAndAddresses(t *testing.T) {
	user := User{
		ID:       "2",
		Password: "rahasia",
		FirstName: "User 50",
		MiddleName: "Middle",
		LastName: "Hozanna",
		Wallets: []Wallet{
			{
				ID:      "2",
				UserID:  "2",
				Balance: 1000000,
			},
		},
		Addresses: []Address{
			{
				UserID:  "2",
				Address: "Jalan A",
			},
			{
				UserID:  "2",
				Address: "Jalan B",
			},
		},
	}

	err := db.Save(&user).Error
	assert.Nil(t, err)
}

func TestPreloadJoinOneToMany(t *testing.T) {
	var users []User
	err := db.Model(&User{}).Preload("Addresses").Preload("Wallets").Find(&users).Error
	assert.Nil(t, err)
}