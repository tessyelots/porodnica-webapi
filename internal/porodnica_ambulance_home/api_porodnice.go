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
	"github.com/gin-gonic/gin"
)

type PorodniceAPI interface {


    // CreatePorodnica Post /api/porodnica
    // Saves new porodnica definition 
     CreatePorodnica(c *gin.Context)

    // DeletePorodnica Delete /api/porodnica/:porodnicaId
    // Deletes specific porodnica 
     DeletePorodnica(c *gin.Context)

}