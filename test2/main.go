package main

import (
	"fmt"

	surrealdb "github.com/surrealdb/surrealdb.go/v2"
	"github.com/surrealdb/surrealdb.go/v2/pkg/models"
)

// UGC
type UGC struct {
	ID        *models.RecordID             `json:"id"`
	Name      string                       `json:"name"`
	State     *models.RecordID             `json:"state"`
	Type      string                       `json:"type"`
	Number    string                       `json:"number"`
	Area      float32                      `json:"area,omitempty"`
	Centre    *models.GeometryPoint        `json:"centre"`
	Geometry  *models.GeometryMultiPolygon `json:"geometry"`
	CWA       []*models.RecordID           `json:"cwa"`
	IsMarine  bool                         `json:"is_marine"`
	IsFire    bool                         `json:"is_fire"`
	ValidFrom *models.CustomDateTime       `json:"valid_from"`
	ValidTo   *models.CustomDateTime       `json:"valid_to"`
}

// UGC ID
type UGCID struct {
	ID    string          `json:"id"`
	State models.RecordID `json:"state"`
	Type  string          `json:"type"`
}

type VTECEvent struct {
	ID         models.RecordID
	CreatedAt  models.CustomDateTime `json:"created_at"`
	UpdatedAt  models.CustomDateTime `json:"updated_at"`
	Start      models.CustomDateTime `json:"start"`
	Expires    models.CustomDateTime `json:"expires"`
	End        models.CustomDateTime `json:"end"`
	EndInitial models.CustomDateTime `json:"end_initial"`
	Title      string                `json:"title"`
}

type VTECEventID struct {
	EventNumber  int                   `json:"event_number"`
	Phenomena    models.RecordID       `json:"phenomena"`
	WFO          models.RecordID       `json:"wfo"`
	Significance models.RecordID       `json:"significance"`
	Issued       models.CustomDateTime `json:"issued"`
}

func main() {
	// Connect to SurrealDB
	db, err := surrealdb.New("ws://localhost:8000/rpc")
	if err != nil {
		panic(err)
	}

	if err = db.Use("WITS", "AWIPS"); err != nil {
		panic(err)
	}

	authData := surrealdb.Auth{
		Username: "root",
		Password: "root",
	}
	_, err = db.SignIn(&authData)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to SurrealDB")

	res, err := surrealdb.Query[[]UGC](db, "SELECT id FROM ugc", map[string]interface{}{})
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
