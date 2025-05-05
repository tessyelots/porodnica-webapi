package porodnica_ambulance_home

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type implPorodnicaWaitingListAPI struct {
}

func NewPorodnicaWaitingListApi() PorodnicaWaitingListAPI {
	return &implPorodnicaWaitingListAPI{}
}

func (o implPorodnicaWaitingListAPI) CreateWaitingListEntry(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implPorodnicaWaitingListAPI) DeleteWaitingListEntry(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implPorodnicaWaitingListAPI) GetWaitingListEntries(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implPorodnicaWaitingListAPI) GetWaitingListEntry(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implPorodnicaWaitingListAPI) UpdateWaitingListEntry(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
