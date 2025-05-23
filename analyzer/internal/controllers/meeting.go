package controllers

import (
	"DAS/internal/repositories"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strconv"
)

type MeetingController struct {
	MeetingRepo *repositories.MeetingRepository
}

func NewMeetingController(meetingRepo *repositories.MeetingRepository) Controller {
	return &MeetingController{MeetingRepo: meetingRepo}
}

func (m MeetingController) Register(r *gin.Engine) {
	r.GET("/meetings", m.GetMeetings)
	r.GET("/meetings/:id", m.GetMeeting)
}

func (m MeetingController) GetMeetings(c *gin.Context) {
	meetings, err := m.MeetingRepo.GetAll(c)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
		}
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, meetings)
}

func (m MeetingController) GetMeeting(c *gin.Context) {
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

	meeting, err := m.MeetingRepo.GetMeeting(c, int(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, meeting)
}
