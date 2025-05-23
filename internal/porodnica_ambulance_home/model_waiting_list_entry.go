/*
 * Waiting List Api
 *
 * Porodnica Waiting List management for Web-In-Cloud system
 *
 * API version: 1.0.0
 * Contact: xsmutny@stuba.sk
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package porodnica_ambulance_home

import (
	"time"
)

type WaitingListEntry struct {

	// Unique id of the entry in this waiting list
	Id string `json:"id"`

	// Name of patient in waiting list
	Name string `json:"name,omitempty"`

	// Unique identifier of the patient known to Web-In-Cloud system
	PatientId string `json:"patientId"`

	// Timestamp since when the patient entered the porodnica waiting list
	WaitingSince time.Time `json:"waitingSince"`

	// Estimated time of porod. Ignored on post.
	EstimatedLaborDate time.Time `json:"estimatedLaborDate"`

	GaveBirth bool `json:"gaveBirth,omitempty"`
}
