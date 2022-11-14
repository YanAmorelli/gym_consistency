package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgconn"
	"github.com/labstack/echo"
	"github.com/yanamorelli/gym_consistency/models"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

var duplicateEntryError = &pgconn.PgError{Code: "23505"}

func (h *Handler) WentGym(c echo.Context) error {
	var ok = new(models.Ok)
	if err := c.Bind(ok); err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := h.DB.Table("gym_consistency").Create(ok).Error; err != nil {
		log.Println(err.Error())

		if err != nil && errors.As(err, &duplicateEntryError) {
			return c.JSON(http.StatusInternalServerError, models.JsonObj{
				"error":        err.Error(),
				"errorMessage": "Você já registrou seu treino de hoje! :)",
			})
		}

		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error": err.Error(),
		})

	}

	c.JSON(http.StatusOK, ok)
	return nil
}

func (h *Handler) GetDate(c echo.Context) error {
	dates := c.Param("date")
	var ok = new(models.Ok)
	if err := h.DB.Table("gym_consistency").Where("date_gym= ?", dates).First(ok).Error; err != nil {
		log.Println(err.Error())

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusInternalServerError, models.JsonObj{
				"error":        err.Error(),
				"errorMessage": "Não foi registrado se foi para a academia ou não :(!",
			})
		}

		return c.JSON(http.StatusInternalServerError, models.JsonObj{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, models.JsonObj{
		"Went": ok.Ok,
	})

	return nil
}

func (h *Handler) StatsOfMonth(c echo.Context) error {
	// comparason dates not working
	var stats = new(models.Stats)
	firstDay, lastDay := getFirstAndLastDayOfMonth()
	log.Println("Datas", firstDay, lastDay)

	h.DB.Raw("select count(*) from gym_consistency where ok = true and date_gym between ? and ?", firstDay, lastDay).Scan(&stats.PresentDays)
	h.DB.Raw("select count(*) from gym_consistency where ok = false and date_gym between ? and ?", firstDay, lastDay).Scan(&stats.MissedDays)

	c.JSON(http.StatusOK, stats)
	return nil
}

func getFirstAndLastDayOfMonth() (string, string) {
	year := time.Now().Year()
	month := time.Now().Month()
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	lastDay := time.Date(year, month+1, 1, 0, 0, 0, -1, time.Local).Format("2006-01-02")
	return firstDay, lastDay
}
