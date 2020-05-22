package base

import (
	"net/http"

	"sample-go/model"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type Handler struct {
	Db *gorm.DB
}

func (h *Handler) CreateUser(c echo.Context) error {
	result := new(model.User)
	_ = c.Bind(result)

	if err := h.Db.Create(result).Error; err != nil {
		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error}, "\t")
	}
	return c.JSONPretty(http.StatusCreated, model.Response{Data: result}, "\t")

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
	return c.JSONPretty(http.StatusCreated, model.Response{Data: user}, "\t")
}

func (h *Handler) UpdateUser(c echo.Context) error {
	result := new(model.User)
	_ = c.Bind(result)
	if err := h.Db.Save(result); err != nil {
		return c.JSONPretty(http.StatusOK, model.Response{Data: result}, "\t")
	}
	return c.JSONPretty(http.StatusOK, model.Response{Data: result}, "\t")
}

func (h *Handler) DeleteUser(c echo.Context) error {
	result := new(model.User)
	_ = c.Bind(result)
	if err := h.Db.Where("id = ?", result.ID).Delete(result); err != nil {
		return c.JSONPretty(http.StatusOK, model.Response{Error: err.Error}, "\t")
	}
	return c.JSONPretty(http.StatusOK, model.Response{Data: result}, "\t")
}
