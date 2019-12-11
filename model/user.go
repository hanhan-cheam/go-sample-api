package model

import (
	"github.com/jinzhu/gorm"
	"github.com/dgrijalva/jwt-go"
)

// "github.com/jinzhu/gorm"

//User struct declaration
type User struct {
	gorm.Model
	Username   string  `json:"username"`
	Employeeid string  `json:"employeeid"`
	Password   string  `json:"password"`
	Role       string  `gorm:"default:'user'" json:"role"`
	Position   string  `json:"position"`
	Firstname  string  `json:"firstname"`
	Lastname   string  `json:"lastname"`
	Email      string  `json:"email"`
	Funds      float64 `gorm:"default:0" json:"funds"`
	Status     string  `json:"status"`
	Msg        string  `json:"msg"`
	Image      string  `gorm:"default:'default.png'" json:"image"`

	// Availablefunds int `json:"availablefunds"`
}

type Response struct {
	Data   interface{} `json:"data"`
	Data2 interface{} `json:"data2"`
	Data3 interface{} `json:"data3"`
	Data4 interface{} `json:"data4"`
	Data5 interface{} `json:"data5"`
	Error  interface{} `json:"error"`
	Msg    string      `json:"msg"`
	Status int         `json:"status"`
	// Pagination postgres.Jsonb `json:"pagination"`
}

func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&User{}, &Food{}, &Transaction{}, &Availableleave{}, &Leave{}, &Poll{}, &Polldetail{},&Order{}, &Orderdetail{})
	return db
}

type Transaction struct {
	gorm.Model
	Description   string  `json:"description"`
	Type          string  `json:"type"`
	Amount        float64 `gorm:"default:0" json:"amount"`
	UserID        uint    `gorm:"default:null" json:"userid"`
	OrderdetailID uint    `gorm:"default:null" json:"orderdetailid"`
	Users         []User
	// Availablefunds int `json:"availablefunds"`
}

type Availableleave struct {
	gorm.Model
	Amount    float64 `gorm:"default:0" json:"amount"`
	Leavetype string  `json:"leavetype"`
	UserID    uint    `gorm:"default:null" json:"userid"`
	// Leaves []Leave
	User   User
	// Availablefunds int `json:"availablefunds"`
}

type Leave struct {
	gorm.Model
	Fromdate         string  `json:"fromdate"`
	Todate           string  `gorm:"default:null" json:"todate"`
	Days             float64 `json:"days"` //apply how many days
	Reason           string  `json:"reason"`
	Attachment       string  `gorm:"default:null" json:"attachment"`
	Status           string  `gorm:"default:'pending'" json:"status"`
	Approvedate      string  `gorm:"default:null" json:"approvedate"`
	// Availableleaveid uint    `gorm:"default:null foreignkey:availableleaveid"  json:"availableleaveid"`
	// AvailableleaveID uint    `gorm:"default:null column:availableleaveid"  json:"availableleaveid"`  //correct
	AvailableleaveID uint    `gorm:"default:null"  json:"availableleaveid"`
	Availableleave         Availableleave
}



type Leavejoin struct {
	gorm.Model
	Fromdate         string  `json:"fromdate"`
	Todate           string  `gorm:"default:null" json:"todate"`
	Days             float64 `json:"days"` //apply how many days
	Reason           string  `json:"reason"`
	Attachment       string  `gorm:"default:null" json:"attachment"`
	Status           string  `gorm:"default:'pending'" json:"status"`
	Approvedate      string  `gorm:"default:null" json:"approvedate"`
	// Availableleaveid uint    `gorm:"default:null foreignkey:availableleaveid"  json:"availableleaveid"`
	// AvailableleaveID uint    `gorm:"default:null column:availableleaveid"  json:"availableleaveid"`  //correct
	AvailableleaveID uint    `gorm:"default:null"  json:"availableleaveid"`

	//Availableleave
	Amount    float64 `gorm:"default:0" json:"amount"`
	Leavetype string  `json:"leavetype"`
	UserID    uint    `gorm:"default:null" json:"userid"`
}








type Food struct {
	gorm.Model
	Name  string  `json:"name"`
	Price float64 `gorm:"default:0" json:"price"`
}

type Poll struct {
	gorm.Model
	Closeby      string  `json:"closeby"`
	Date         string  `json:"date"`
	Status       bool    `gorm:"default:true" json:"status"`
	Paymentprice float64 `gorm:"default:null" json:"paymentprice"`
	// OrderID      uint    `gorm:"default:null" json:"orderid"`
	Size     float64 `json:"size"`
	Polldetail  []Polldetail  //trying ??

	// Userid           uint    `gorm:"default:null" json:"userid"`
}

type Polldetail struct {
	gorm.Model
	FoodID uint `gorm:"default:null" json:"foodid"`
	PollID uint `gorm:"default:null" json:"poolid"`
	Poll   Poll
	Food   Food
}


type Selectedquantity struct{
	Pollsize string `json:"pollsize"`
	Quantity string `json:"quantity"`
}

type Order struct{
	gorm.Model
	Orderdate string `json:"orderdate"`
	Quantity int `json:"quantity"`
	Totalprice float64 `json:"totalprice"`
	PollID uint `gorm:"default:null" json:"pollid"`
	UserID uint `gorm:"default:null" json:"userid"`
	Orderdetail []Orderdetail
	User User
	Poll Poll

}

type Orderdetail struct{
	gorm.Model
	Price float64 `json:"price"`
	Quantity int `json:"quantity"`
	OrderID uint `gorm:"default:null" json:"orderid"`
	FoodID uint  `gorm:"default:null" json:"foodid"`
	Order Order
	Food Food
}


type Orderdetailjoin struct{
	gorm.Model
	Orderdate string `json:"orderdate"`
	Quantity int `json:"quantity"`
	Totalprice float64 `json:"totalprice"`
	PollID uint `gorm:"default:null" json:"pollid"`
	UserID uint `gorm:"default:null" json:"userid"`
	Price float64 `json:"price"`
	// Quantity int `json:"quantity"`
	OrderID uint `gorm:"default:null" json:"orderid"`
	FoodID uint  `gorm:"default:null" json:"foodid"`
	Firstname  string  `json:"firstname"`
	Lastname   string  `json:"lastname"`
	
}



type Tokens struct {
	Role      string `json:"role,omitempty"`
	Username  string `json:"username,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	Id        string `json:"jti,omitempty"`
	Is_revoked bool   `json:"is_revoked"`
	jwt.StandardClaims
	// jwt.StandardClaims
}




type Sumresponse struct{

	Quantity float64 `json:"quantity"`
	FoodID uint `json:"foodid"`
	Size float64 `json:"size"`
	// Idsize float64 `json:"idsize"`
	Overallpayment float64 `json:"overallpayment"`
}

type Idsizeresponse struct{

	Idsize float64 `json:"idsize"`
	
}

type Changepassword struct{
	Currentpassword   string  `json:"currentpassword"`
	Newpassword   string  `json:"newpassword"`
	Retypepassword   string  `json:"retypepassword"`
}


type Confirmpassword struct{
	Password   string  `json:"password"`
	Repeatpassword   string  `json:"repeatpassword"`
}