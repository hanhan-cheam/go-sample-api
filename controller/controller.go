package controller

import (
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
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

// func (h *Handler) Login(c echo.Context) error {

// 	current_user := new(model.User)
// 	_ = c.Bind(current_user)
// 	user := new(model.User)
// 	fmt.Println(current_user.Username)

// 	uNameCheck := helpers.IsEmpty(current_user.Username)
// 	pwdCheck := helpers.IsEmpty(current_user.Password)
// 	// fmt.Println(uNameCheck)
// 	if uNameCheck || pwdCheck {
// 		return c.JSONPretty(http.StatusOK, model.Response{Msg: "please fill in the require information"}, "\t")
// 	}

// 	if err := h.Db.Where("username = ?", current_user.Username).Find(&user).Error; err != nil {
// 		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error(), Msg: "username does not exists"}, "\t")
// 	}
// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(current_user.Password)); err != nil {
// 		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error(), Msg: "unauthorized password,please try again"}, "\t")
// 	}

// 	claims := jwt.StandardClaims{
// 		Id:        strconv.FormatUint(uint64(user.ID), 10),
// 		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
// 		Audience:  user.Role,
// 		// Username:  user.Username,
// 	}
// 	mytoken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
// 	token, err := mytoken.SignedString([]byte("secretpassword"))

// 	if err != nil {
// 		log.Println("Error Creating JWT token", err)
// 		return c.String(http.StatusOK, "StatusInternalServerError")
// 	}

// 	return c.JSONPretty(http.StatusOK, model.Response{Data: &model.Responsetoken{Token: token, Message: "You were logged in!"}}, "\t")
// }

// func (h *Handler) AllUsers(c echo.Context) error {
// 	var users []model.User
// 	h.Db.Find(&users)
// 	return c.JSONPretty(http.StatusCreated, model.Response{Data: users}, "\t")
// }

// func (h *Handler) User(c echo.Context) error {
// 	// var user []model.User //use this will  be array
// 	user := new(model.User)
// 	id := c.Param("id")
// 	h.Db.Where("id = ?", id).First(&user)
// 	return c.JSONPretty(http.StatusCreated, user, "\t")
// }

// func (h *Handler) Checkaccess(c echo.Context) error {
// 	reqToken := c.Request().Header.Get("Authorization")
// 	splitToken := strings.Split(reqToken, "Bearer")
// 	reqToken = strings.TrimSpace(splitToken[1])
// 	fmt.Println(reqToken)
// 	token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
// 		return jwtKey, nil
// 	})
// 	if err != nil {
// 		fmt.Println("Cannot get token")
// 		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error()}, "\t")
// 	}
// 	claims := token.Claims.(jwt.MapClaims)
// 	return c.JSONPretty(http.StatusOK, model.Response{Data: claims}, "\t")

// }

// func (h *Handler) Uploadfilepath(c echo.Context) error {
// 	id := c.Param("id")
// 	return c.File(UPLOAD_PATH + "/" + id)
// }

// func (h *Handler) Uploaduserfilepath(c echo.Context) error {
// 	id := c.Param("id")
// 	return c.File(UserimageUPLOAD_PATH + "/" + id)

// }

// func (h *Handler) Userimageupload(c echo.Context) error {
// 	user := new(model.User)
// 	// _ = c.Bind(user)
// 	login_id := c.FormValue("profile_user_id")
// 	var uploadedFileName string
// 	h.Db.Where("id = ?", login_id).Find(&user)

// 	fmt.Println("user id is ", login_id)
// 	stringuserid := strconv.FormatUint(uint64(user.ID), 10)
// 	file, err := c.FormFile("imagefilename")
// 	if err != nil {
// 		return err
// 	}
// 	src, err := file.Open()
// 	if err != nil {
// 		return err
// 	}
// 	defer src.Close()
// 	getextension := strings.Split(file.Filename, ".")
// 	uploadedFileName = stringuserid + "." + getextension[1] //must use leave id not available id cause duplicate //importan need fix
// 	uploadedFilePath := path.Join(UserimageUPLOAD_PATH, uploadedFileName)
// 	dst, err := os.Create(uploadedFilePath)
// 	if err != nil {
// 		return err
// 	}
// 	defer dst.Close()
// 	io.Copy(dst, src)
// 	// if _, err = io.Copy(dst, src); err != nil {
// 	// 	return err
// 	// }
// 	if err := h.Db.Model(&user).Where("id = ?", stringuserid).Update("image", uploadedFileName).Error; err != nil {
// 		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error()}, "\t")
// 	}
// 	return c.JSONPretty(http.StatusOK, model.Response{Data: user}, "\t")
// }

// func (h *Handler) Changepassword(c echo.Context) error {
// 	user := new(model.User)
// 	input_user := new(model.User)
// 	_ = c.Bind(input_user)

// 	new_input_user := new(model.Changepassword)
// 	_ = c.Bind(new_input_user)
// 	h.Db.Where("id = ?", input_user.ID).First(&user)
// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input_user.Password)); err != nil {
// 		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error(), Msg: "current password are wrong,please try again"}, "\t")
// 	}
// 	if new_input_user.Newpassword != new_input_user.Retypepassword {
// 		return c.JSONPretty(http.StatusOK, model.Response{Msg: "new password and retype password are not match"}, "\t")
// 	} else {
// 		fmt.Println(new_input_user.Newpassword)
// 		fmt.Println(new_input_user.Retypepassword)
// 		hash, _ := bcrypt.GenerateFromPassword([]byte(new_input_user.Newpassword), bcrypt.DefaultCost)
// 		if err := h.Db.Model(&user).Where("id = ?", input_user.ID).Update("password", string(hash)).Error; err != nil {
// 			return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error(), Msg: "password change failed"}, "\t")
// 		}
// 	}
// 	return c.JSONPretty(http.StatusOK, model.Response{Data: user}, "\t")
// }

// func (h *Handler) Test(c echo.Context) error {

// 	return c.JSONPretty(http.StatusOK, model.Response{Data: "dsadsadsadsadsa"}, "\t")
// }
