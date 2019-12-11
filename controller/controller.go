package controller

import (
	"fmt"
	"leave-order/helpers"
	"leave-order/io"
	"math"
	"log"
	"leave-order/model"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"path"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const UPLOAD_PATH = "./upload"
const UserimageUPLOAD_PATH = "./userupload"
const arraysize = 10 //maximum 10 food in poll,can change if need more*******
var jwtKey = []byte("secretpassword")
type Handler struct {
	Db *gorm.DB
}

type ErrorResponse struct {
	Err string
}

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}

func (h *Handler) CreateUser(c echo.Context) error {
	user := new(model.User)
	_ = c.Bind(user)

	find_user := new(model.User)

	confirm_input_user := new(model.Confirmpassword)
	_ = c.Bind(confirm_input_user)

	check1 := helpers.IsEmpty(user.Username)
	check2 := helpers.IsEmpty(user.Employeeid)
	check3 := helpers.IsEmpty(user.Password)
	check4 := helpers.IsEmpty(user.Position)
	check5 := helpers.IsEmpty(user.Firstname)
	check6 := helpers.IsEmpty(user.Lastname)
	check7 := helpers.IsEmpty(user.Email)

	fmt.Println(confirm_input_user.Password)
	fmt.Println(confirm_input_user.Repeatpassword)

	if check1 || check2 || check3 || check4 || check5 || check6 || check7 {
		return c.JSONPretty(http.StatusOK, model.Response{Error: "plese fill in the require field"}, "\t")
	}else if !validateEmail(user.Email) {
		return c.JSONPretty(http.StatusOK, model.Response{Error: "email format wrong !!! eg:example@gmail.com"}, "\t")
	}else if confirm_input_user.Password != confirm_input_user.Repeatpassword {
		return c.JSONPretty(http.StatusOK, model.Response{Error: "password and confirm password are not match"}, "\t")
	}else if err := h.Db.Where("username = ?", user.Username).First(&find_user).Error; err == nil {
		return c.JSONPretty(http.StatusOK, model.Response{Error:"username already exists"}, "\t") 
		// Error: err.Error(),
	}else {
		hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err := h.Db.Create(&model.User{Username: user.Username, Password: string(hash),
			Employeeid: user.Employeeid,
			Position:   user.Position, Firstname: user.Firstname,
			Lastname: user.Lastname, Email: user.Email, Funds: 0}).Error; err != nil {
			return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error()}, "\t")
		}
		// user.Role = "user"
		user.Status = "success"
		user.Msg = "register success"

		lastinsertuser := new(model.User)
		h.Db.Last(&lastinsertuser)

		fmt.Println("my id is ", lastinsertuser.ID)

		defaultleavetype := [11]string{"Annual Leave", "Medical Leave", "Maternity Leave", "Paternity Leave", "Marriage Leave", "Hospitalisation Leave", "Landon Leave", "Examination Leave", "Unpaid Leave", "Compassionate Leave", "Emergency Leave"}

		defaultleaveday := [11]float64{14, 14, 60, 1, 3, 60, 1, 10, 5, 0, 0}

		// availaleleave := new(model.Availableleave)
		for x := 0; x < len(defaultleavetype); x++ {
			if err := h.Db.Create(&model.Availableleave{Amount: defaultleaveday[x], Leavetype: defaultleavetype[x], UserID: lastinsertuser.ID}).Error; err != nil {
				return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error()}, "\t")
			}
		}
		return c.JSONPretty(http.StatusCreated, model.Response{Data: &user}, "\t")
	}
}

func (h *Handler) Login(c echo.Context) error {

	current_user := new(model.User)
	_ = c.Bind(current_user)
	user := new(model.User)
	fmt.Println(current_user.Username)

	uNameCheck := helpers.IsEmpty(current_user.Username)
	pwdCheck := helpers.IsEmpty(current_user.Password)
	// fmt.Println(uNameCheck)
	if uNameCheck || pwdCheck {
		return c.JSONPretty(http.StatusOK, model.Response{Msg: "please fill in the require information"}, "\t")
	}

	if err := h.Db.Where("username = ?", current_user.Username).Find(&user).Error; err != nil {
		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error(), Msg: "username does not exists"}, "\t")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(current_user.Password)); err != nil {
		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error(), Msg: "unauthorized password,please try again"}, "\t")
	}



	// 	token := jwt.New(jwt.SigningMethodHS256)
	// token.Claims = &model.Tokens{
	// 	Id:        strconv.FormatUint(uint64(user.ID), 10),
	// 	ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	// 	Role:      user.Role,
	// 	Username:  user.Username,
	// }
	// tokenString, _ := token.SignedString([]byte("secretpassword"))
	// type tempToken struct{ Token string }
	// // return c.JSONPretty(http.StatusOK, model.Response{Data: &tempToken{Token: tokenString}}, "\t")
	//  return c.JSONPretty(http.StatusOK, model.Response{Data: &model.Responsetoken{Token: tokenString, Message: "You were logged in!"}}, "\t")




	claims := jwt.StandardClaims{
		Id:        strconv.FormatUint(uint64(user.ID), 10),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Audience:      user.Role,
		// Username:  user.Username,
	}
	mytoken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err := mytoken.SignedString([]byte("secretpassword"))

	if err != nil {
		log.Println("Error Creating JWT token", err)
		return c.String(http.StatusOK, "StatusInternalServerError")
	}

	return c.JSONPretty(http.StatusOK, model.Response{Data: &model.Responsetoken{Token: token, Message: "You were logged in!"}}, "\t")
}

func (h *Handler) AllUsers(c echo.Context) error {
	var users []model.User
	h.Db.Find(&users)
	return c.JSONPretty(http.StatusCreated, model.Response{Data: users}, "\t")
}

func (h *Handler) User(c echo.Context) error {
	// var user []model.User //use this will  be array
	user := new(model.User)
	id := c.Param("id")
	h.Db.Where("id = ?", id).First(&user)
	return c.JSONPretty(http.StatusCreated, user, "\t")
}

func (h *Handler) Checkaccess(c echo.Context) error {
	reqToken := c.Request().Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	reqToken = strings.TrimSpace(splitToken[1])
	fmt.Println(reqToken)
	token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		fmt.Println("Cannot get token")
		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error()}, "\t")
	}
	claims := token.Claims.(jwt.MapClaims)
	return c.JSONPretty(http.StatusOK, model.Response{Data: claims}, "\t")

}

func (h *Handler) Addfunds(c echo.Context) error {
	current_user := new(model.User)
	_ = c.Bind(current_user)

	user := new(model.User)

	if current_user.Funds == 0 {
		return c.JSONPretty(http.StatusOK, model.Response{Msg: "Please enter an amount to top up"}, "\t")
	}

	fmt.Println("before topup value", user.Funds)

	transaction := new(model.Transaction)
	transaction.Description = "Top up"
	transaction.Type = "In"
	if err := h.Db.Create(&model.Transaction{Description: transaction.Description,
		Type: transaction.Type, Amount: current_user.Funds,
		UserID: current_user.ID}).Error; err != nil {
		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error()}, "\t")
	}

	if err := h.Db.Where("id = ?", current_user.ID).Find(&user).Error; err != nil {
		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error(), Msg: "Top up failed"}, "\t")
	}

	user.Funds = user.Funds + current_user.Funds
	fmt.Println("input value", current_user.Funds)
	fmt.Println("database value", user.Funds)

	user.Funds = math.Ceil(user.Funds*100)/100

	h.Db.Model(&current_user).Where("id = ?", current_user.ID).Update("funds", user.Funds)
	return c.JSONPretty(http.StatusCreated, model.Response{Msg: "Top up Success"}, "\t")
}


func (h *Handler) AllTransactions(c echo.Context) error {
	var transactions []model.Transaction
	h.Db.Find(&transactions)
	return c.JSONPretty(http.StatusCreated, model.Response{Data: transactions}, "\t")
}

func (h *Handler) IndividualTransactions(c echo.Context) error {
	current_user := new(model.User)
	_ = c.Bind(current_user)
	user_id := current_user.ID
	fmt.Println(user_id)
	var transactions []model.Transaction
	h.Db.Where("user_id = ?", current_user.ID).Find(&transactions)
	return c.JSONPretty(http.StatusOK, model.Response{Data: transactions}, "\t")
}

func (h *Handler) Availableleaves(c echo.Context) error {
	var availaleleaves []model.Availableleave
	h.Db.Find(&availaleleaves)
	return c.JSONPretty(http.StatusCreated, model.Response{Data: availaleleaves}, "\t")
}

func (h *Handler) Availableleavesindividual(c echo.Context) error {
	current_user := new(model.User)
	_ = c.Bind(current_user)
	userid := current_user.ID
	// json need pass id
	fmt.Println(userid)
	var availableleaves []model.Availableleave
	h.Db.Where("user_id = ?", current_user.ID).Find(&availableleaves)
	return c.JSONPretty(http.StatusOK, model.Response{Data: availableleaves}, "\t")
}

func (h *Handler) Addleave(c echo.Context) error {
	availaleleave := new(model.Availableleave)
	_ = c.Bind(availaleleave)
	if err := h.Db.Create(&model.Availableleave{Amount: availaleleave.Amount,
		Leavetype: availaleleave.Leavetype, UserID: availaleleave.UserID}).Error; err != nil {
		return c.JSONPretty(0, model.Response{Error: err.Error()}, "\t")
	}
	return c.JSONPretty(http.StatusCreated, model.Response{Data: &availaleleave}, "\t")
}

func (h *Handler) Upload(c echo.Context) error {

	filestatus := c.FormValue("filestatus") // check if there is file upload eg:if 1 have uplaod , 0 no upload
	availableannualleave := c.FormValue("availableannualleave")
	applyingdays := c.FormValue("applyingdays") 
	leavetypeholder := c.FormValue("leavetypeholder")
	availableday := c.FormValue("availableday")
	datefrom := c.FormValue("datefrom")
	dateto := c.FormValue("dateto")
	dataid := c.FormValue("dataid")
	reasonText := c.FormValue("reasonText")
	check1 := helpers.IsEmpty(applyingdays)
	check2 := helpers.IsEmpty(leavetypeholder)
	check3 := helpers.IsEmpty(availableannualleave)
	check4 := helpers.IsEmpty(availableday)
	check5 := helpers.IsEmpty(reasonText)
	check6 := helpers.IsEmpty(datefrom)
	// check7 := helpers.IsEmpty(dateto)
	if check1 || check2 || check3 || check4 || check5 || check6 {
		return c.JSONPretty(http.StatusOK, model.Response{Status: 0, Msg: "please fill in the require information"}, "\t")
	}
	if applyingdays == "0" {
		return c.JSONPretty(http.StatusOK, model.Response{Status: 0, Msg: "please pick a day to apply before submit"}, "\t")
	}
	//test
	floatapplyingdays, _ := strconv.ParseFloat(applyingdays, 64)
	floatavailableday, _ := strconv.ParseFloat(availableday, 64)
		
    u64, err := strconv.ParseUint(dataid, 10, 32)
    if err != nil {
        fmt.Println(err)
    }
    uintdataid := uint(u64)
    // fmt.Println("i am here",uintdataid)

	if floatapplyingdays > floatavailableday {
		return c.JSONPretty(http.StatusOK, model.Response{Status: 0, Msg: "not enough leave"}, "\t")
	}

	//delete later
	fmt.Println("see below" + dataid)
	fmt.Println(uintdataid)
	fmt.Println("applying day " + applyingdays)
	fmt.Println(datefrom)
	fmt.Println(dateto)
	fmt.Println(leavetypeholder)
	fmt.Println(availableannualleave)
	fmt.Println(availableday)
	fmt.Println(reasonText)
	fmt.Println(filestatus)
	//delete later

	var uploadedFileName string

	if err := h.Db.Create(&model.Leave{Fromdate: datefrom,Todate: dateto ,Days: floatapplyingdays ,Reason:reasonText ,Attachment: uploadedFileName,AvailableleaveID:uintdataid}).Error; err != nil {
		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error()}, "\t")
	}


	lastapplyleave := new(model.Leave)
	h.Db.Last(&lastapplyleave)

	fmt.Println("my last apply leave id is ", lastapplyleave.ID)
	stringlastapplyleaveid := strconv.FormatUint(uint64(lastapplyleave.ID), 10)

	if filestatus == "1" {
		file, err := c.FormFile("filename")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()
		getextension := strings.Split(file.Filename, ".")
		uploadedFileName = stringlastapplyleaveid + "." + getextension[1]   //must use leave id not available id cause duplicate //importan need fix
		uploadedFilePath := path.Join(UPLOAD_PATH, uploadedFileName)
		dst, err := os.Create(uploadedFilePath)
		if err != nil {
			return err
		}
		defer dst.Close()
		io.Copy(dst, src)
		// if _, err = io.Copy(dst, src); err != nil {
		// 	return err
		// }
		leave := new(model.Leave)
		if err := h.Db.Model(&leave).Where("id = ?", stringlastapplyleaveid).Update("attachment", uploadedFileName).Error; err != nil {
			return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error()}, "\t")
		}	
	}
	currentleave := new(model.Availableleave)
	h.Db.Where("id = ?", uintdataid).Find(&currentleave)


	leaveafterdeduct := 0.0
	if leavetypeholder == "Unpaid Leave"{
		leaveafterdeduct =	floatavailableday  //unpaid leave always 5 available
	}else{
		leaveafterdeduct =	floatavailableday-floatapplyingdays
	}




	// leaveafterdeduct :=	floatavailableday-floatapplyingdays


	fmt.Println(leaveafterdeduct)
	if err := h.Db.Model(&currentleave).Where("id = ?", uintdataid).Update("amount", leaveafterdeduct).Error; err != nil {
		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error()}, "\t")
	} 
	return c.JSONPretty(http.StatusOK, model.Response{Status: 1, Msg: "Leave apply success,please wait for update "}, "\t")
}

func (h *Handler) Allpendingleaves(c echo.Context) error {
	var pendings []model.Leave  //can get many
	if err := 	h.Db.Preload("Availableleave").Preload("Availableleave.User").Where("status=?","pending").Find(&pendings).Error; err != nil {
		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error(), Msg: "data does not exists"}, "\t")
	}
	return c.JSONPretty(http.StatusCreated, model.Response{Data: &pendings}, "\t")
}

func (h *Handler) Allnotpendingleaves(c echo.Context) error {
	var leave []model.Leave  //can get many
	if err := 	h.Db.Preload("Availableleave").Preload("Availableleave.User").Where("status!=?","pending").Find(&leave).Error; err != nil {
		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error(), Msg: "data does not exists"}, "\t")
	}
	return c.JSONPretty(http.StatusCreated, model.Response{Data: &leave}, "\t")
}


func (h *Handler) Approveleave(c echo.Context) error {	
	res := make(map[string]string)
	res["status"] = "success" ;
	res["msg"] = "approved success" ;
	leave := new(model.Leave)
	id := c.Param("id")
	if err := 	h.Db.Model(&leave).Where("id = ?", id).Update("status", "approved").Error; err != nil {
		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error(), Msg: "data does not exists"}, "\t")
	}
	return c.JSONPretty(http.StatusCreated, res, "\t")
}


func (h *Handler) Rejectleave(c echo.Context) error {	
	res := make(map[string]string)
	res["status"] = "success" ;
	res["msg"] = "reject success" ;
	leave := new(model.Leave)
	availableleave := new(model.Availableleave)
	id := c.Param("id")
	applydays := c.Param("applydays")
	availableleaveid := c.Param("availableleaveid")
	if err := 	h.Db.Model(&leave).Where("id = ?", id).Update("status", "rejected").Error; err != nil {
		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error(), Msg: "data does not exists"}, "\t")
	}
	if err := h.Db.Where("id = ?",availableleaveid).Find(&availableleave).Error; err != nil {
		return c.JSONPretty(http.StatusOK, model.Response{Error:err.Error(),Msg:"not found"},"\t")
	}
	floatapplydays, _ := strconv.ParseFloat(applydays, 64)
	availableleave.Amount = availableleave.Amount + floatapplydays  
	fmt.Println("apply day abelow")
	fmt.Println(floatapplydays)
	fmt.Println("available leave amount below")
	fmt.Println(availableleave.Amount)
	if err := 	h.Db.Model(&availableleave).Where("id = ?", availableleaveid).Update("amount", availableleave.Amount).Error; err != nil {
		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error(), Msg: "data does not exists"}, "\t")
	}
	return c.JSONPretty(http.StatusCreated, res, "\t")
}

func (h *Handler) Uploadfilepath(c echo.Context) error {	
	id := c.Param("id")
	return c.File(UPLOAD_PATH +"/"+ id)
}

func (h *Handler) Uploaduserfilepath(c echo.Context) error {	
	id := c.Param("id")
	return c.File(UserimageUPLOAD_PATH +"/"+ id)

}

func (h *Handler) Addfood(c echo.Context) error {
	food := new(model.Food)
	_ = c.Bind(food)
	check1 := helpers.IsEmpty(food.Name)
	fmt.Println(food.Price)
	if check1 {
		return c.JSONPretty(http.StatusOK, model.Response{Error: "please fill in the require field"}, "\t")
	}
	if food.Price <= 0 {
		return c.JSONPretty(http.StatusOK, model.Response{Error: "please fill in the price"}, "\t")
	}
	fmt.Println(food.Name)
	h.Db.Create(&model.Food{Name: food.Name, Price: food.Price})
	return c.JSONPretty(http.StatusCreated, model.Response{Data: &food}, "\t")	
}

func (h *Handler) Findfood(c echo.Context) error {
	var foods []model.Food
	h.Db.Find(&foods)
	return c.JSONPretty(http.StatusCreated, model.Response{Data: foods}, "\t")
}

func (h *Handler) Updatefood(c echo.Context) error {
	food := new(model.Food)
	_ = c.Bind(food)
	fmt.Println("dasdsaas")
	if err := h.Db.Model(food).Where("ID = ?", food.ID).Update("price", food.Price).Error; err != nil {
		panic(err)
	}
	return c.JSONPretty(http.StatusCreated, model.Response{Data: &food}, "\t")
}

func (h *Handler) Deletefood(c echo.Context) error {
	food := new(model.Food)
	_ = c.Bind(food)
	fmt.Println(food.ID)
	h.Db.Where("ID =?", &food.ID).Find(&food)
	h.Db.Delete(&food)
	return c.JSONPretty(http.StatusCreated, model.Response{Data: &food}, "\t")
}

func (h *Handler) Foodscheckbox(c echo.Context) error {
	var food []model.Food
	checkedID := c.Param("checkedid")
	convertcheckedID := strings.Split(checkedID, ",")
	h.Db.Raw("SELECT * FROM foods WHERE id IN (?)", convertcheckedID).Scan(&food)
	return c.JSONPretty(http.StatusCreated, model.Response{Data: food}, "\t")
}

func (h *Handler) Addpoll(c echo.Context) error {
	poll := new(model.Poll)
	_ = c.Bind(poll)
	h.Db.Create(&model.Poll{Closeby: poll.Closeby, Date: poll.Date, Status: poll.Status, Paymentprice: poll.Paymentprice, Size: poll.Size})
	lastpoll := new(model.Poll)
	h.Db.Last(&lastpoll)
	fmt.Println("my id is ", lastpoll.ID)
	fmt.Println(c.FormValue("FoodID"))
	foods := strings.Split(c.FormValue("FoodID"), ",")
	for _, food := range foods {
		fmt.Println(food)
		h.Db.Create(&model.Polldetail{FoodID: toInt(food), PollID: lastpoll.ID})
	}
	return c.JSONPretty(http.StatusCreated, model.Response{Data: &poll}, "\t")
}

func toInt(value string) uint {
	integer, _ := strconv.ParseUint(value, 10, 32)
	return uint(integer)
}

func (h *Handler) Findpolldetail(c echo.Context) error {

	var polldetails []model.Polldetail
	h.Db.Find(&polldetails)
	fmt.Println(polldetails)

	h.Db.Preload("Poll").Preload("Food").Find(&polldetails)

	return c.JSONPretty(http.StatusOK, model.Response{Data: &polldetails}, "\t")
}
func (h *Handler) Lastpolldetail(c echo.Context) error {
	poll := new(model.Poll)
	fmt.Println(poll)
	h.Db.Where("status = ?", "true").Find(&poll)
	return c.JSONPretty(http.StatusOK, model.Response{Data: &poll}, "\t")
}

func (h *Handler) Fooddisplay(c echo.Context) error {
	var polls []model.Poll
	h.Db.Preload("Polldetail").Where("status=?",true).Preload("Polldetail.Food").Find(&polls)
	return c.JSONPretty(http.StatusOK, model.Response{Data: &polls}, "\t")

}

var selectedfoodIDholder [arraysize]int64 //need clear
var finalizeselectedfoodID [arraysize]int64 //need clear

var quantityholder [arraysize]int64  //need clear
var finalizeselectedfoodquantity [arraysize]int64 //need clear

func (h *Handler) Calculateselectedquantityprice(c echo.Context) error {
	fmt.Println("-----------------")
	fmt.Println("-----------------")
	fmt.Println("-----------------")
	var x int64
	res := make(map[string]float64)
	// resInt := make(map[string]float64)
	var totalorderquantity int64 = 0

	currentinputquantityid := c.FormValue("currentinputquantityid")
	intcurrentinputquantityid,_ := strconv.ParseInt(currentinputquantityid, 10, 64)

	pollsize := c.FormValue("pollsize")
	intpollsize,_ := strconv.ParseInt(pollsize, 10, 64)
	
	currentinputfoodprice := c.FormValue("currentinputfoodprice")
	quantity := c.FormValue("qty")
	floatcurrentinputfoodprice,_ := strconv.ParseFloat(currentinputfoodprice, 64)
	intquantity,_ := strconv.ParseInt(quantity,10, 64)

	quantityholder[intcurrentinputquantityid] = intquantity
	for x = 0;x<intpollsize;x++{
		totalorderquantity = totalorderquantity+quantityholder[x]
	}

	fmt.Println("quantity ",intquantity)
	foodid := c.FormValue("foodid")
	intfoodid,_ := strconv.ParseInt(foodid,10,64)
	selectedfoodIDholder[intcurrentinputquantityid] = intfoodid

	checkorderdifferentfoodsize := 0
	for x = 0;x<arraysize;x++{

		if intquantity == 0{
			selectedfoodIDholder[intcurrentinputquantityid] = 0
		}

		if selectedfoodIDholder[x] == 0{
			fmt.Println(x,"isempty")

		}else
		{
			fmt.Println(x,"notempty")
			checkorderdifferentfoodsize = checkorderdifferentfoodsize+1 //check different food order ,eg :order chicken rice,mee goreng,size = 2		
		}
	}
	//can delete
	fmt.Println("check order different foodsize ",checkorderdifferentfoodsize)
	i := 0
	q := 0
	
	// fmt.Println("checkcontent below")
	for x = 0;x<arraysize;x++{
		
		fmt.Println("holdervalue",selectedfoodIDholder[x])


		if selectedfoodIDholder[x] > 0{

			finalizeselectedfoodID[i] = selectedfoodIDholder[x]
			i++
		}else{
			finalizeselectedfoodID[i] = 0
		}


		if quantityholder[x] > 0{
			finalizeselectedfoodquantity[q] = quantityholder[x]
			q++
		}else{
			finalizeselectedfoodquantity[q] = 0
		}
	}

	for x = 0;x<arraysize;x++{
		fmt.Println("finalize foodid",finalizeselectedfoodID[x])
	}

	fmt.Println("-----------------")

	for x = 0;x<arraysize;x++{
		fmt.Println("finalize quantity",finalizeselectedfoodquantity[x])
	}

	totalselectedquantityprice := floatcurrentinputfoodprice*float64(intquantity)
	totalselectedquantityprice = math.Ceil(totalselectedquantityprice*100)/100
	res["totalorderquantity"] = float64(totalorderquantity) ;
	res["totalselectedquantityprice"] = totalselectedquantityprice ;
	res["checkorderdifferentfoodsize"] = float64(checkorderdifferentfoodsize) ;
	fmt.Println("-----------------")
	return c.JSONPretty(http.StatusCreated, res, "\t")
}

var amountholder [arraysize]float64
// var finalizeamount [arraysize]float64

func (h *Handler) Calculatetotalselectedquantityprice(c echo.Context) error {
	var x int64
	var totalorderprice float64 = 0
	res := make(map[string]float64)
	selectedtotalprice := c.FormValue("selectedtotalprice")
	floatselectedtotalprice,_ := strconv.ParseFloat(selectedtotalprice, 64)
	currentinputquantityid := c.FormValue("currentinputquantityid")
	intcurrentinputquantityid,_ := strconv.ParseInt(currentinputquantityid, 10, 64)
	pollsize := c.FormValue("pollsize")
	intpollsize,_ := strconv.ParseInt(pollsize, 10, 64)
	amountholder[intcurrentinputquantityid] = floatselectedtotalprice
	// a := 0
	for x = 0;x<intpollsize;x++{
		totalorderprice = totalorderprice+amountholder[x]

		fmt.Println("amount holder : ",amountholder[x])
	}

	res["totalorderprice"] = math.Ceil(totalorderprice*100)/100    
	return c.JSONPretty(http.StatusCreated, res, "\t")
}

func (h *Handler) Clearamountholder(c echo.Context) error {
	for x:=0;x<arraysize;x++{
		amountholder[x] = 0
		quantityholder[x] = 0
		selectedfoodIDholder[x] = 0
		finalizeselectedfoodID[x] = 0
		finalizeselectedfoodquantity[x] = 0
		orderdetailIDholder[x] = 0
		orderIDholder[x] = 0 //????
	}
	return c.JSONPretty(http.StatusCreated, model.Response{Data: selectedfoodIDholder}, "\t")
}


var orderdetailIDholder [arraysize]uint

func (h *Handler) Createorder(c echo.Context) error {
	order := new(model.Order)
	_ = c.Bind(order)
	if err := h.Db.Create(&model.Order{Orderdate: order.Orderdate, Quantity: order.Quantity,
		Totalprice: order.Totalprice,PollID: order.PollID, UserID: order.UserID }).Error; err != nil {
		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error()}, "\t")
	}

	lastcreateorder := new(model.Order)
	h.Db.Last(&lastcreateorder)
	fmt.Println("my id is ", lastcreateorder.ID)

	differentfoodsize := c.FormValue("differentfoodsize")
	intdifferentfoodsize,_ := strconv.ParseInt(differentfoodsize,10,64)
	
	for x:=0;x< int(intdifferentfoodsize);x++ {
		var totalselecteditemprice float64 = 0
		foodsearch := new(model.Food)
		fmt.Println("foodid is",uint(finalizeselectedfoodID[x]))
		h.Db.Where("id = ?", uint(finalizeselectedfoodID[x]) ).Find(&foodsearch)
		// fmt.Println("food name is",foodsearch.Price)
		totalselecteditemprice = float64(finalizeselectedfoodquantity[x]) * foodsearch.Price
		totalselecteditemprice = math.Ceil(totalselecteditemprice*100)/100
		if err := h.Db.Create(&model.Orderdetail{Price: totalselecteditemprice , Quantity: int(finalizeselectedfoodquantity[x]),
			FoodID: uint(finalizeselectedfoodID[x]),OrderID: lastcreateorder.ID}).Error; err != nil {
			return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error()}, "\t")
		}
		lastcreateorderdetail := new(model.Orderdetail)
		h.Db.Last(&lastcreateorderdetail)
		orderdetailIDholder[x] = lastcreateorderdetail.ID
		fmt.Println("orderdetailIDholder ", orderdetailIDholder[x])
	}
	fmt.Println("foodsize is",intdifferentfoodsize)

	//update funds after order food
	current_funds := c.FormValue("current_funds")
	fmt.Println("current_funds",current_funds)
	userid := c.FormValue("userid")
	totalprice := c.FormValue("totalprice")
	fmt.Println("userid",userid)

	user := new(model.User)

	if err := h.Db.Where("id = ?", userid).Find(&user).Error; err != nil {
		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error(), Msg: "cannot find record"}, "\t")
	}

	floatcurrent_funds,_ := strconv.ParseFloat(current_funds, 64)
	floattotalprice,_ := strconv.ParseFloat(totalprice, 64)

	floatcurrent_funds = floatcurrent_funds - floattotalprice
	// fmt.Println("user.Funds",user.Funds)
	h.Db.Model(&user).Where("id = ?", userid).Update("funds", floatcurrent_funds)
	//update funds after order food
	for x:=0;x< int(intdifferentfoodsize);x++ {

		foodsearch := new(model.Food)
		h.Db.Where("id = ?", uint(finalizeselectedfoodID[x]) ).Find(&foodsearch)


		u64, err := strconv.ParseUint(userid, 10, 32)
	    if err != nil {
	        fmt.Println(err)
	    }
	    uintuserid := uint(u64)

		transaction := new(model.Transaction)
		transaction.Description = foodsearch.Name
		transaction.Type = "Out"
		if err := h.Db.Create(&model.Transaction{Description: transaction.Description,
			Type: transaction.Type, Amount: amountholder[x],
			UserID: uintuserid,OrderdetailID:orderdetailIDholder[x]}).Error; err != nil {
			return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error()}, "\t")
		}
	}
	//transaction history
	res := make(map[string]float64)
	res["current_funds"] =  math.Ceil(floatcurrent_funds*100)/100;
	return c.JSONPretty(http.StatusOK, model.Response{Data: res}, "\t")
}


//not using
type NResult struct{
    N int64 //or int ,or some else
}

func (h *Handler) SumSame() int64 {
    var n  NResult
    h.Db.Table("orderdetails").Select("sum(quantity) as n").Scan(&n)
    return n.N
}
//not using

var totalquantityofAnitem int = 0

var currentpollfoodquantityholder [arraysize]int

var orderIDholder [100]int

func (h *Handler) Ordersummarydisplay(c echo.Context) error {
	var orderdetailfirst_last_holder [2]int
	var orderdetail []model.Orderdetail 
	var count_result int64
	var count_orderdetail int64
	poll := new(model.Poll)

	h.Db.Where("status=?",true).Find(&poll)
	var order []model.Order 
	h.Db.Where("poll_id=?",poll.ID).Find(&order).Count(&count_result)

    for x:=0;x<int(count_result);x++{
    	orderIDholder[x] =  int(order[x].ID)
	}

	for x:=0;x<int(count_result);x++{
		fmt.Println("orderIDholder has ",orderIDholder[x])
	}

	stringOrderIdInsideOrderdetail := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(orderIDholder)), ","), "[]")
   	fmt.Println("stringOrderIdInsideOrderdetail",stringOrderIdInsideOrderdetail)
   	// h.Db.Where("poll_id=?",poll.ID).Find(&order).Count(&count_result) //count order detail to find the order detail with id in orderIDholder[?]
   	h.Db.Where("order_id IN (?)", strings.Split(stringOrderIdInsideOrderdetail, ",")).Find(&orderdetail).Count(&count_orderdetail)  //fix 11,12,13
	fmt.Println("count_orderdetail",count_orderdetail)
	for x:=0;x<int(count_orderdetail);x++ {

		if x==0 {
			orderdetailfirst_last_holder[0] = int(orderdetail[0].ID)
		}

		if x==(int(count_orderdetail)-1) {
			orderdetailfirst_last_holder[1] = int(orderdetail[count_orderdetail-1].ID)
		}

		h.Db.Where("order_id IN (?)", strings.Split(stringOrderIdInsideOrderdetail, ",")).Find(&orderdetail)
	}

	for x:=0;x<2;x++ { 
		fmt.Println("orderdetailfirst_last_holder is",orderdetailfirst_last_holder[x])
	}

	var sumresponse []model.Sumresponse 
	// SELECT SUM(quantity) FROM orderdetails GROUP BY food_id;
	h.Db.Raw("SELECT food_id,SUM(quantity) as quantity FROM orderdetails WHERE orderdetails.id BETWEEN ? AND ? GROUP BY food_id ORDER BY food_id ASC",orderdetailfirst_last_holder[0],orderdetailfirst_last_holder[1]).Scan(&sumresponse)
	
	var sizeresponse []model.Sumresponse 
	h.Db.Raw("SELECT COUNT(DISTINCT food_id) as size FROM public.orders INNER JOIN public.orderdetails ON orderdetails.order_id = orders.id WHERE orders.poll_id = ?",poll.ID).Scan(&sizeresponse)

	//sum overall payment price
	var overallpriceresponse []model.Sumresponse 
	h.Db.Raw("SELECT SUM(price) as overallpayment FROM orderdetails WHERE orderdetails.id BETWEEN ? AND ?",orderdetailfirst_last_holder[0],orderdetailfirst_last_holder[1]).Scan(&overallpriceresponse)
	var order_detail_join []model.Orderdetailjoin 
	// h.Db.Raw("SELECT * FROM public.orderdetails INNER JOIN public.orders ON orders.ID = orderdetails.order_id INNER JOIN public.users ON users.ID = orders.user_id where orders.poll_id = ?",poll.ID).Scan(&order_detail_join)

	h.Db.Raw("SELECT users.id,users.firstname,users.lastname,orders.poll_id,orderdetails.food_id,sum(orderdetails.quantity) as quantity FROM public.orderdetails INNER JOIN public.orders ON orders.ID = orderdetails.order_id INNER JOIN public.users ON users.id = orders.user_id where orders.poll_id = ? GROUP BY orders.poll_id,food_id,users.id",poll.ID).Scan(&order_detail_join)
	var idsizeresponse []model.Idsizeresponse 
	h.Db.Raw("SELECT COUNT(*) as idsize FROM public.orderdetails INNER JOIN public.orders ON orders.ID = orderdetails.order_id INNER JOIN public.users ON users.id = orders.user_id where orders.poll_id = ? GROUP BY orders.poll_id,food_id,users.id",poll.ID).Scan(&idsizeresponse)

	// h.Db.Raw("SELECT COUNT(*) as idsize FROM public.orderdetails INNER JOIN public.orders ON orders.ID = orderdetails.order_id where orders.poll_id = ?",poll.ID).Scan(&idsizeresponse)

	//refer keep this
	// h.Db.Table("orderdetails").
	// Where("id BETWEEN ? AND ?", 24, 30).
	// Group("food_id").
	// Pluck("sum(quantity)", &sumresponse)
	// SELECT SUM(quantity) FROM public.orderdetails  WHERE id BETWEEN 24 AND 30 GROUP BY food_id
	//refer keep this
	fmt.Println("pollsize is",poll.Size)
 	for x := 0;x<arraysize;x++{
		fmt.Println("currentpollfoodquantityholder",currentpollfoodquantityholder[x])
	}

	fmt.Println("count_result",count_result)
	fmt.Println("pollid",poll.ID)
	fmt.Println("---------------------------------")
	fmt.Println("number of row is",count_result)
	return c.JSONPretty(http.StatusOK, model.Response{Data: sumresponse,Data2:overallpriceresponse,Data3:sizeresponse,Data4:order_detail_join,Data5:idsizeresponse}, "\t")
	// return c.JSONPretty(http.StatusOK, model.Response{Data: order_detail_join}, "\t")
}

func (h *Handler) Closepoll(c echo.Context) error {

	pollid := c.FormValue("pollid") 
	overallprice_holder := c.FormValue("overallprice_holder") 
	poll := new(model.Poll)
	floatoverallprice_holder, _ := strconv.ParseFloat(overallprice_holder, 64)

	fmt.Println("pollid",pollid)
	fmt.Println("overallprice_holder",overallprice_holder)
	// h.Db.Where("status = ?", "true").Find(&poll)

	h.Db.Model(&poll).Where("id = ?", pollid).Updates(map[string]interface{}{"status": false, "paymentprice": floatoverallprice_holder})

	return c.JSONPretty(http.StatusOK, model.Response{Msg: "Poll Close"}, "\t")

}


func (h *Handler) Orderhistory(c echo.Context) error {
	var orderdetails []model.Orderdetail
	h.Db.Preload("Food").Preload("Order").Preload("Order.User").Find(&orderdetails)
	return c.JSONPretty(http.StatusCreated, model.Response{Data: orderdetails}, "\t")
}

func (h *Handler) Transactionhistoryindividual(c echo.Context) error {
	user := new(model.User)
	var transactions []model.Transaction
	_ = c.Bind(user)
	userid := user.ID
	h.Db.Where("user_id = ?", userid).Find(&transactions)
	return c.JSONPretty(http.StatusOK, model.Response{Data: transactions}, "\t")
}

//do here //do here //do here //do here //do here //do here //do here //do here
func (h *Handler) Leavehistoryindividual(c echo.Context) error {

	user := new(model.User)
	_ = c.Bind(user)
	userid := user.ID
	var leavejoin []model.Leavejoin

	h.Db.Raw("SELECT * FROM public.leaves INNER JOIN public.availableleaves ON leaves.availableleave_id = availableleaves.id where availableleaves.user_id = ?",userid).Scan(&leavejoin)
	return c.JSONPretty(http.StatusOK, model.Response{Data: leavejoin}, "\t")
}


func (h *Handler) Userimageupload(c echo.Context) error {
		user := new(model.User)
		// _ = c.Bind(user)
		login_id := c.FormValue("profile_user_id") 
		var uploadedFileName string
		h.Db.Where("id = ?",login_id).Find(&user)

		fmt.Println("user id is ", login_id)
		stringuserid := strconv.FormatUint(uint64(user.ID), 10)
	    file, err := c.FormFile("imagefilename")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()
		getextension := strings.Split(file.Filename, ".")
		uploadedFileName = stringuserid + "." + getextension[1]   //must use leave id not available id cause duplicate //importan need fix
		uploadedFilePath := path.Join(UserimageUPLOAD_PATH, uploadedFileName)
		dst, err := os.Create(uploadedFilePath)
		if err != nil {
			return err
		}
		defer dst.Close()
		io.Copy(dst, src)
		// if _, err = io.Copy(dst, src); err != nil {
		// 	return err
		// }
		if err := h.Db.Model(&user).Where("id = ?", stringuserid).Update("image", uploadedFileName).Error; err != nil {
			return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error()}, "\t")
		}	
		return c.JSONPretty(http.StatusOK, model.Response{Data: user}, "\t")
}



func (h *Handler) Changepassword(c echo.Context) error {
	user := new(model.User)
	input_user := new(model.User)
	_ = c.Bind(input_user)

	new_input_user := new(model.Changepassword)
	_ = c.Bind(new_input_user)
	h.Db.Where("id = ?",input_user.ID).First(&user)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input_user.Password)); err != nil {
		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error(), Msg: "current password are wrong,please try again"}, "\t")
	}
	if new_input_user.Newpassword != new_input_user.Retypepassword {
		return c.JSONPretty(http.StatusOK, model.Response{Msg: "new password and retype password are not match"}, "\t")
	}else{
		fmt.Println(new_input_user.Newpassword)
		fmt.Println(new_input_user.Retypepassword)
		hash, _ := bcrypt.GenerateFromPassword([]byte(new_input_user.Newpassword), bcrypt.DefaultCost)
		if err := h.Db.Model(&user).Where("id = ?", input_user.ID).Update("password", string(hash)).Error; err != nil {
			return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error(), Msg: "password change failed"}, "\t")
		}
	}
	return c.JSONPretty(http.StatusOK, model.Response{Data: user}, "\t")
}
	

func (h *Handler) Leavetypedisplay(c echo.Context) error {

	var availableleaves []model.Availableleave

	h.Db.Where("leavetype = ? OR leavetype = ?", "Annual Leave", "Medical Leave").Order("leavetype").Find(&availableleaves)

	return c.JSONPretty(http.StatusOK, model.Response{Data: availableleaves}, "\t")
}


