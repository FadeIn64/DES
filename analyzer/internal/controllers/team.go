package controllers

import (
	"DAS/internal/repositories"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strconv"
)

type TeamController struct {
	TeamRepo *repositories.TeamRepository
}

func NewTeamController(teamRepo *repositories.TeamRepository) Controller {
	return &TeamController{teamRepo}
}

func (t TeamController) Register(r *gin.Engine) {
	r.GET("/teams", t.GetTeams)
	r.GET("/teams/:id", t.GetTeam)
}

func (t TeamController) GetTeams(c *gin.Context) {
	teams, err := t.TeamRepo.GetAll(c)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
		}
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, teams)
}

func (t TeamController) GetTeam(c *gin.Context) {
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

	team, err := t.TeamRepo.GetByID(c, int(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, team)
}
