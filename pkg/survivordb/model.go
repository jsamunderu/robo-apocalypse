// Package classification of Survivors API
//
// Documentation for Survivors API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package survivordb

import "time"

// LastLocation defines the structure for the last location
// swagger:model
type LastLocation struct {
	// the gps longitude
	//
	// required: true
	// max length: 255
	Longitude float64 `json:"longitude"`

	// the gps latitude
	//
	// required: true
	// max length: 255
	Latitude float64 `json:"latitude"`
}

// Resources defines the structure for a resource
// swagger:model
type Resources struct {
	// the water the survivor currently has
	//
	// required: true
	// max length: 255
	Water float64 `json:"water"`

	// the food the survivor currently has
	//
	// required: true
	// max length: 255
	Food string `json:"food"`

	// the medication the survivor currently has
	//
	// required: true
	// max length: 255
	Medication string `json:"medication"`

	// the ammunition the survivor currently has
	//
	// required: true
	// max length: 255
	Ammunition int `json:"ammunition"`
}

// Survivor defines the structure for a survivor
// swagger:model
type Survivor struct {
	// the name for this poduct
	//
	// required: true
	// max length: 128
	Name string `json:"name"`

	// the age for this survivor
	//
	// required: true
	// max length: 3
	Age int `json:"age"`

	// the gender for this survivor
	//
	// required: true
	// max length: 16
	Gender string `json:"gender"`

	// the id number for this survivor
	//
	// required: true
	// max length: 30
	IdNumber string `json:"id"`

	LastLocation
	Resources

	// the status of infection of the survivor
	//
	// required: true
	Infected bool `json:"infected"`
	// the time when survivor information was recorded
	//
	// required: false
	LastUpdateTime time.Time `json:"timestamp"`
}

//
// NOTE: Types defined here are purely for documentation purposes
// these types are not used by any of the handers

// A list of survivors
// swagger:response surivivorsResponse
type survivorsResponseWrapper struct {
	// All current survivors
	// in: body
	Body []Survivor
}

// Data structure representing a single survivor
// swagger:response survivorResponse
type survivorResponseWrapper struct {
	// Newly created survivor
	// in: body
	Body Survivor
}

// No content is returned by this API endpoint
// swagger:response noContentResponse
type noContentResponseWrapper struct {
}

// swagger:parameters createSurvivor
type survivorParamsWrapper struct {
	// Survivor data structure to = Create.
	// Note: the id field is ignored by update and create operations
	// in: body
	// required: true
	Body Survivor
}

// swagger:parameters setInfected
type survivorIDParamsWrapper struct {
	// The id of the survivor for which the operation relates
	// in: body
	// required: true
	Payload struct {
		IdNumber int `json:"id"`
	}
}

// swagger:parameters updateLocation
type survivorLocationParamsWrapper struct {
	// The id of the survivor for which the operation relates
	// in: body
	// required: true
	Payload struct {
		IdNumber int `json:"id"`
		LastLocation
	}
}

// swagger:parameters updateResource
type survivorResourceParamsWrapper struct {
	// The id of the survivor for which the operation relates
	// in: body
	// required: true
	Payload struct {
		IdNumber int `json:"id"`
		Resources
	}
}

// Data structure representing infected survivor stats
// swagger:response statsResponse
type survivorStatsResponseWrapper struct {
	// Newly created survivor
	// in: body
	Stats struct {
		HealthyPercentage  float64 `json:"healthyPercentage"`
		InfectedPercentage float64 `json:"infectedPercentage"`
	}
}
