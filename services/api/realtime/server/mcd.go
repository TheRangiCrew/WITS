package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/twpayne/go-geos"
)

type MCDDB struct {
	ID               int        `json:"id"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	Product          string     `json:"product"`
	Issued           time.Time  `json:"issued"`
	Expires          time.Time  `json:"expires"`
	Year             int        `json:"year"`
	Concerning       string     `json:"concerning"`
	Geometry         *geos.Geom `json:"geometry"`
	WatchProbability int        `json:"watchProbability"`
	MostProbTornado  string     `json:"mostProbTornado"`
	MostProbGust     string     `json:"mostProbGust"`
	MostProbHail     string     `json:"mostProbHail"`
	Text             string     `json:"text"`
}

type MCD struct {
	ID               int             `json:"id"`
	CreatedAt        time.Time       `json:"createdAt"`
	UpdatedAt        time.Time       `json:"updatedAt"`
	Product          string          `json:"product"`
	Issued           time.Time       `json:"issued"`
	Expires          time.Time       `json:"expires"`
	Year             int             `json:"year"`
	Concerning       string          `json:"concerning"`
	Geometry         json.RawMessage `json:"geometry"`
	WatchProbability int             `json:"watchProbability"`
	MostProbTornado  string          `json:"mostProbTornado"`
	MostProbGust     string          `json:"mostProbGust"`
	MostProbHail     string          `json:"mostProbHail"`
	Text             string          `json:"text"`
}

func (server *Server) InitialiseMCD() error {

	rows, err := server.DB.Query(context.Background(), "SELECT m.*, p.data AS text FROM mcd.mcd m JOIN awips.products p ON m.product = p.product_id")
	if err != nil {
		return err
	}
	defer rows.Close()

	now := time.Now()

	for rows.Next() {
		var mcddb MCDDB
		err := rows.Scan(
			&mcddb.ID,
			&mcddb.CreatedAt,
			&mcddb.UpdatedAt,
			&mcddb.Product,
			&mcddb.Issued,
			&mcddb.Expires,
			&mcddb.Year,
			&mcddb.Concerning,
			&mcddb.Geometry,
			&mcddb.WatchProbability,
			&mcddb.MostProbTornado,
			&mcddb.MostProbGust,
			&mcddb.MostProbHail,
			&mcddb.Text)
		if err != nil {
			slog.Error("Failed to scan MCD data", "error", err)
			continue
		}
		mcd := MCD{
			ID:               mcddb.ID,
			CreatedAt:        mcddb.CreatedAt,
			UpdatedAt:        mcddb.UpdatedAt,
			Product:          mcddb.Product,
			Issued:           mcddb.Issued,
			Expires:          mcddb.Expires,
			Year:             mcddb.Year,
			Concerning:       mcddb.Concerning,
			Geometry:         json.RawMessage(mcddb.Geometry.ToGeoJSON(1)),
			WatchProbability: mcddb.WatchProbability,
			MostProbTornado:  mcddb.MostProbTornado,
			MostProbGust:     mcddb.MostProbGust,
			MostProbHail:     mcddb.MostProbHail,
			Text:             mcddb.Text,
		}

		bytes, err := json.Marshal(mcd)
		if err != nil {
			slog.Error("Failed to marshal MCD data to JSON", "error", err)
			continue
		}

		stat := server.RDB.Set(context.Background(), "mcd:"+strconv.Itoa(mcd.ID), bytes, mcd.Expires.Sub(now))
		if stat.Err() != nil {
			slog.Error("Failed to set MCD data in Redis", "error", stat.Err())
			continue
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	server.Router.GET("/mcd/current", func(c *gin.Context) {
		mcds, err := server.getAllMCD()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve MCD data: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, mcds)
	})

	return nil
}

func (server *Server) getAllMCD() ([]MCD, error) {
	ctx := context.Background()
	iter := server.RDB.Scan(ctx, 0, "mcd:*", 0).Iterator()
	mcds := []MCD{}
	for iter.Next(ctx) {
		res := server.RDB.Get(ctx, iter.Val())
		var mcd MCD
		err := json.Unmarshal([]byte(res.Val()), &mcd)
		if err != nil {
			slog.Error("Failed to unmarshal MCD data", "error", err, "key", iter.Val())
			return nil, err
		}
		mcds = append(mcds, mcd)
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	return mcds, nil
}
