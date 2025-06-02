package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TheRangiCrew/WITS/services/api/archive/db"
	"github.com/TheRangiCrew/WITS/services/api/archive/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type VTECEvent struct {
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Issued         time.Time       `json:"issued"`
	Starts         time.Time       `json:"starts"`
	Expires        time.Time       `json:"expires"`
	Ends           time.Time       `json:"ends"`
	EndInitial     time.Time       `json:"endInitial"`
	Class          string          `json:"class"`
	Phenomena      string          `json:"phenomena"`
	WFO            string          `json:"wfo"`
	Significance   string          `json:"significance"`
	EventNumber    int             `json:"eventNumber"`
	Year           int             `json:"year"`
	Title          string          `json:"title"`
	IsEmergency    bool            `json:"isEmergency"`
	IsPDS          bool            `json:"isPDS"`
	InitialPolygon json.RawMessage `json:"initialPolygon,omitempty"`
}

func RegisterVTECRoutes(router *gin.Engine) {
	// VTEC routes can be registered here
	router.GET("/vtec", func(c *gin.Context) {
		var params struct {
			WFO          string  `form:"wfo" binding:"required,min=3,max=4"`
			Phenomena    *string `form:"phenomena"`
			Significance *string `form:"significance"`
			EventNumber  *int    `form:"eventNumber"`
			Year         int     `form:"year" binding:"required"`
		}
		// Trying binding query parameters
		if err := c.ShouldBindQuery(&params); err != nil {
			// Provide a user-friendly error message
			var msg string
			if ve, ok := err.(validator.ValidationErrors); ok {
				for _, fe := range ve {
					msg += fe.Field() + ": " + fe.Tag() + ". "
				}
			} else {
				msg = "Invalid or missing parameters. Please check your input."
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		// Validate the event number
		if params.EventNumber != nil && *params.EventNumber < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Event number must be greater than 0."})
			return
		}

		// Validate the year
		if params.Year < 1960 || params.Year > time.Now().Year() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Year must be between 1960 and the current year."})
			return
		}

		if len(params.WFO) == 3 {
			rows, err := db.DB.Query(context.Background(), "SELECT icao FROM postgis.offices WHERE id = $1", params.WFO)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "WFO query failed.\n" + err.Error()})
				return
			}

			if !rows.Next() {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid WFO code. Could not find WFO with code: " + params.WFO})
				rows.Close()
				return
			} else {
				var icao string
				err = rows.Scan(&icao)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning WFO code.\n" + err.Error()})
					rows.Close()
					return
				}
				params.WFO = icao // Use the ICAO code for the query
			}
		}

		// Build the query
		query := "SELECT * FROM vtec.events WHERE wfo = $1 AND year = $2"
		args := []interface{}{params.WFO, params.Year}
		paramIdx := 3 // Next placeholder index
		// Optional parameters
		if params.Phenomena != nil {
			query += " AND phenomena = $" + fmt.Sprint(paramIdx)
			args = append(args, *params.Phenomena)
			paramIdx++
		}
		if params.Significance != nil {
			query += " AND significance = $" + fmt.Sprint(paramIdx)
			args = append(args, *params.Significance)
			paramIdx++
		}
		if params.EventNumber != nil {
			query += " AND event_number = $" + fmt.Sprint(paramIdx)
			args = append(args, *params.EventNumber)
			paramIdx++
		}

		// Execute the query
		rows, err := db.DB.Query(context.Background(), query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Event query failed.\n" + err.Error()})
			return
		}
		defer rows.Close()

		events := []VTECEvent{}
		for rows.Next() {
			var result models.VTECEvent
			err = rows.Scan(
				&result.ID,
				&result.CreatedAt,
				&result.UpdatedAt,
				&result.Issued,
				&result.Starts,
				&result.Expires,
				&result.Ends,
				&result.EndInitial,
				&result.Class,
				&result.Phenomena,
				&result.WFO,
				&result.Significance,
				&result.EventNumber,
				&result.Year,
				&result.Title,
				&result.IsEmergency,
				&result.IsPDS,
				&result.PolygonStart)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning database row.\n" + err.Error()})
				return
			}

			var polygonGeoJSON json.RawMessage = nil
			if result.PolygonStart != nil {
				geoJSONString := result.PolygonStart.ToGeoJSON(1)
				polygonGeoJSON = json.RawMessage(geoJSONString)
			}

			event := VTECEvent{
				CreatedAt:      result.CreatedAt,
				UpdatedAt:      result.UpdatedAt,
				Issued:         result.Issued,
				Starts:         result.Starts,
				Expires:        result.Expires,
				Ends:           result.Ends,
				EndInitial:     result.EndInitial,
				Class:          result.Class,
				Phenomena:      result.Phenomena,
				WFO:            result.WFO,
				Significance:   result.Significance,
				EventNumber:    result.EventNumber,
				Year:           result.Year,
				Title:          result.Title,
				IsEmergency:    result.IsEmergency,
				IsPDS:          result.IsPDS,
				InitialPolygon: polygonGeoJSON,
			}

			events = append(events, event)
		}

		if rows.Err() != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading from database.\n" + rows.Err().Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"events": events,
			"count":  len(events),
		})
	})
}
