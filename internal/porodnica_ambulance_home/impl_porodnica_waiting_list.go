package porodnica_ambulance_home

import (
	"net/http"
	"time"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type implPorodnicaWaitingListAPI struct {
}

func NewPorodnicaWaitingListApi() PorodnicaWaitingListAPI {
	return &implPorodnicaWaitingListAPI{}
}

func (o implPorodnicaWaitingListAPI) CreateWaitingListEntry(c *gin.Context) {
	updatePorodnicaFunc(c, func(c *gin.Context, porodnica *Porodnica) (*Porodnica, interface{}, int) {
		var entry WaitingListEntry

		if err := c.ShouldBindJSON(&entry); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		if entry.PatientId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Patient ID is required",
			}, http.StatusBadRequest
		}

		if entry.Id == "" || entry.Id == "@new" {
			entry.Id = uuid.NewString()
		}

		conflictIndx := slices.IndexFunc(porodnica.WaitingList, func(waiting WaitingListEntry) bool {
			return entry.Id == waiting.Id || entry.PatientId == waiting.PatientId
		})

		if conflictIndx >= 0 {
			return nil, gin.H{
				"status":  http.StatusConflict,
				"message": "Entry already exists",
			}, http.StatusConflict
		}

		porodnica.WaitingList = append(porodnica.WaitingList, entry)
		// entry was copied by value return reconciled value from the list
		entryIndx := slices.IndexFunc(porodnica.WaitingList, func(waiting WaitingListEntry) bool {
			return entry.Id == waiting.Id
		})
		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to save entry",
			}, http.StatusInternalServerError
		}
		return porodnica, porodnica.WaitingList[entryIndx], http.StatusOK
	})
}

func (o implPorodnicaWaitingListAPI) DeleteWaitingListEntry(c *gin.Context) {
	updatePorodnicaFunc(c, func(c *gin.Context, porodnica *Porodnica) (*Porodnica, interface{}, int) {
		entryId := c.Param("entryId")

		if entryId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Entry ID is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(porodnica.WaitingList, func(waiting WaitingListEntry) bool {
			return entryId == waiting.Id
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		porodnica.WaitingList = append(porodnica.WaitingList[:entryIndx], porodnica.WaitingList[entryIndx+1:]...)
		return porodnica, nil, http.StatusNoContent
	})
}

func (o implPorodnicaWaitingListAPI) GetWaitingListEntries(c *gin.Context) {
	updatePorodnicaFunc(c, func(c *gin.Context, porodnica *Porodnica) (*Porodnica, interface{}, int) {
		result := porodnica.WaitingList
		if result == nil {
			result = []WaitingListEntry{}
		}
		// return nil ambulance - no need to update it in db
		return nil, result, http.StatusOK
	})
}

func (o implPorodnicaWaitingListAPI) GetWaitingListEntry(c *gin.Context) {
	updatePorodnicaFunc(c, func(c *gin.Context, porodnica *Porodnica) (*Porodnica, interface{}, int) {
		entryId := c.Param("entryId")

		if entryId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Entry ID is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(porodnica.WaitingList, func(waiting WaitingListEntry) bool {
			return entryId == waiting.Id
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		// return nil porodnica - no need to update it in db
		return nil, porodnica.WaitingList[entryIndx], http.StatusOK
	})
}

func (o implPorodnicaWaitingListAPI) UpdateWaitingListEntry(c *gin.Context) {
	updatePorodnicaFunc(c, func(c *gin.Context, porodnica *Porodnica) (*Porodnica, interface{}, int) {
		var entry WaitingListEntry

		if err := c.ShouldBindJSON(&entry); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		entryId := c.Param("entryId")

		if entryId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Entry ID is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(porodnica.WaitingList, func(waiting WaitingListEntry) bool {
			return entryId == waiting.Id
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		if entry.PatientId != "" {
			porodnica.WaitingList[entryIndx].PatientId = entry.PatientId
		}

		if entry.Id != "" {
			porodnica.WaitingList[entryIndx].Id = entry.Id
		}

		if entry.Name != "" {
			porodnica.WaitingList[entryIndx].Name = entry.Name
		}

		if entry.WaitingSince.After(time.Time{}) {
			porodnica.WaitingList[entryIndx].WaitingSince = entry.WaitingSince
		}

		if entry.EstimatedLaborDate.After(entry.WaitingSince) {
			porodnica.WaitingList[entryIndx].EstimatedLaborDate = entry.EstimatedLaborDate
		}

		if entry.GaveBirth != porodnica.WaitingList[entryIndx].GaveBirth {
			porodnica.WaitingList[entryIndx].GaveBirth = entry.GaveBirth
		}

		return porodnica, porodnica.WaitingList[entryIndx], http.StatusOK
	})
}
