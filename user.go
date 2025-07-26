package go_gorm

import (
	"time"
	"gorm.io/gorm"
)


type Sample struct {
	ID   string `gorm:"type:varchar(100);primaryKey" json:"id"`
	Name string `gorm:"type:varchar(100);not null" json:"name"`
}


type User struct {
	ID         string     `gorm:"type:varchar(100);primaryKey" json:"id"`
	Password   string     `gorm:"type:varchar(100);not null" json:"password"`
	FirstName  string     `gorm:"type:varchar(100);not null;column:first_name" json:"first_name"`
	MiddleName string    `gorm:"type:varchar(100);column:middle_name" json:"middle_name,omitempty"`
	LastName   string    `gorm:"type:varchar(100);column:last_name" json:"last_name,omitempty"`
	CreatedAt  time.Time  `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
	
	// Relationships
	UserLogs []UserLog  `gorm:"foreignKey:UserID" json:"user_logs,omitempty"`
	Todos    []Todo     `gorm:"foreignKey:UserID" json:"todos,omitempty"`
	Wallets  []Wallet   `gorm:"foreignKey:UserID" json:"wallets,omitempty"`
	Addresses []Address `gorm:"foreignKey:UserID" json:"addresses,omitempty"`
	Products []Product  `gorm:"many2many:user_like_product;foreignKey:ID;joinForeignKey:UserID;References:ID;joinReferences:ProductID" json:"liked_products,omitempty"`
}


type UserLog struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    string `gorm:"type:varchar(100);not null;column:user_id" json:"user_id"`
	Action    string `gorm:"type:varchar(100);not null" json:"action"`
	CreatedAt int64  `gorm:"not null;column:created_at" json:"created_at"`
	UpdatedAt int64  `gorm:"not null;column:updated_at" json:"updated_at"`
	
	// Relationships
	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}


type Todo struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      string         `gorm:"type:varchar(100);not null;column:user_id" json:"user_id"`
	Title       string         `gorm:"type:varchar(100);not null" json:"title"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at,omitempty"`
	
	
	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}


type Wallet struct {
	ID        string    `gorm:"type:varchar(100);primaryKey" json:"id"`
	UserID    string    `gorm:"type:varchar(100);not null;column:user_id" json:"user_id"`
	Balance   int64     `gorm:"not null" json:"balance"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
	
	
	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}


type Address struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    string    `gorm:"type:varchar(100);not null;column:user_id" json:"user_id"`
	Address   string    `gorm:"type:varchar(100);not null" json:"address"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
	
	// Relationships
	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}


type Product struct {
	ID        string    `gorm:"type:varchar(100);primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Price     int64     `gorm:"not null" json:"price"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
	
	// Relationships
	Users []User `gorm:"many2many:user_like_product;foreignKey:ID;joinForeignKey:ProductID;References:ID;joinReferences:UserID" json:"liked_by_users,omitempty"`
}


type UserLikeProduct struct {
	UserID    string `gorm:"type:varchar(100);primaryKey;column:user_id" json:"user_id"`
	ProductID string `gorm:"type:varchar(100);primaryKey;column:product_id" json:"product_id"`
	
	// Relationships
	User    User    `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Product Product `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
}


func (u *Sample) TableName() string {
	return "sample"
}

func (u *User) TableName() string {
	return "users"
}

func (u *UserLog) TableName() string {
	return "user_logs"
}

func (u *Todo) TableName() string {
	return "todos"
}

func (u *Wallet) TableName() string {
	return "wallets"
}

func (u *Address) TableName() string {
	return "addresses"
}

func (u *Product) TableName() string {
	return "products"
}

func (u *UserLikeProduct) TableName() string {
	return "user_like_product"
}

