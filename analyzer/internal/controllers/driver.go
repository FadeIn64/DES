package controllers

import (
	"DAS/internal/repositories"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strconv"
)

type DriverController struct {
	DriverRepo *repositories.DriverRepository
}

func NewDriverController(driverRepo *repositories.DriverRepository) Controller {
	return &DriverController{DriverRepo: driverRepo}
}

func (d DriverController) Register(r *gin.Engine) {
	r.GET("/drivers/:id", d.GetDriver)
	r.GET("/drivers", d.GetDriversByTeam)
}

func (d DriverController) GetDriver(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	driver, err := d.DriverRepo.GetDriverByNumber(c, int(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, driver)
}

func (d DriverController) GetDriversByTeam(c *gin.Context) {

	teamIDParam := c.Query("team")
	if teamIDParam == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("team is required"))
		return
	}

	teamID, err := strconv.ParseInt(teamIDParam, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	drivers, err := d.DriverRepo.GetDriversByTeam(c, int(teamID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, drivers)
}
