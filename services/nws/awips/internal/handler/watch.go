package handler

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/db"
	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/handler/util"
	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/logger"
	"github.com/TheRangiCrew/go-nws/pkg/awips"
	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

type Watches struct {
	watches []*Watch
	lock    sync.Mutex
}

var watchesInstance *Watches
var once sync.Once

func watches() *Watches {
	once.Do(func() {
		watchesInstance = &Watches{
			watches: []*Watch{},
		}
	})
	return watchesInstance
}

type Watch struct {
	Number  int
	Type    string
	Issued  *time.Time
	Expires *time.Time
	SAW     *models.RecordID
	SEL     *models.RecordID
	WOU     *models.RecordID
	WWP     *db.WatchWWP
	IsPDS   bool
	Polygon *models.GeometryPolygon
}

func (handler *Handler) watch(product *awips.TextProduct, productID *models.RecordID) {

	watchRegexp := regexp.MustCompile("(((SEVERE THUNDERSTORM|TORNADO) WATCH [0-9]{1,4})|(WW [0-9]{1,4} (SEVERE TSTM|TORNADO))|((Severe Thunderstorm|Tornado) Watch Number [0-9]{1,4})|(W(T|S) [0-9]{1,4}))")
	str := watchRegexp.FindString(product.Text)
	numRegexp := regexp.MustCompile("[0-9]{1,4}")
	numstr := numRegexp.FindString(str)

	watchNum, err := strconv.Atoi(numstr)
	if err != nil {
		handler.Logger.Error("error getting watch number: " + err.Error())
		return
	}

	fmt.Println(watchNum)

	watch := watches().getWatch(watchNum)

	if watch == nil {
		handler.Logger.Info("Creating new watch " + strconv.Itoa(watchNum))
		watch = watches().newWatch(watchNum)
	}

	prod := product.AWIPS.Product

	fmt.Println(prod)

	switch prod {
	case "SAW":
		err = watch.saw(product, productID, handler.Logger)
	case "SEL":
		err = watch.sel(product, productID, handler.Logger)
	case "WOU":
		update, err := watch.wou(product, productID, handler)
		if err != nil {
			handler.Logger.Error(err.Error())
			return
		}
		if update {
			watches().removeWatch(watch.Number)
			return
		}
	case "WWP":
	}

	if err != nil {
		handler.Logger.Error(err.Error())
		return
	}

	fmt.Println(watch.isComplete())
}

func (watch *Watch) saw(product *awips.TextProduct, productID *models.RecordID, logger *logger.Logger) error {

	watch.SAW = productID

	if watch.Type == "" {
		typeRegexp := regexp.MustCompile("(SEVERE TSTM|TORNADO)")
		typeStr := typeRegexp.FindString(product.Text)
		if typeStr == "SEVERE TSTM" {
			watch.Type = "SV"
		} else if typeStr == "TORNADO" {
			watch.Type = "TO"
		}
	}

	if len(product.Segments) == 0 {
		return fmt.Errorf("SAW has 0 segments")
	}

	if len(product.Segments) != 1 {
		logger.Warn(fmt.Sprintf("SAW contains %d segments", len(product.Segments)))
	}

	latlon := product.Segments[0].LatLon
	p := util.PolygonFromAwips(*latlon.Polygon)
	watch.Polygon = &p

	return nil
}

func (watch *Watch) sel(product *awips.TextProduct, productID *models.RecordID, logger *logger.Logger) error {
	watch.SEL = productID

	if watch.Type == "" {
		typeRegexp := regexp.MustCompile("^Severe Thunderstorm|Tornado")
		typeStr := typeRegexp.FindString(product.Text)
		if typeStr == "Severe Thunderstorm" {
			watch.Type = "SV"
		} else if typeStr == "Tornado" {
			watch.Type = "TO"
		}
	}

	if len(product.Segments) == 0 {
		return fmt.Errorf("SEL has 0 segments")
	}

	if len(product.Segments) != 1 {
		logger.Warn(fmt.Sprintf("SEL contains %d segments", len(product.Segments)))
	}

	issued := product.Issued.UTC()
	watch.Issued = &issued
	watch.Expires = &product.Segments[0].Expires

	watch.IsPDS = product.Segments[0].IsPDS()

	return nil
}

func (watch *Watch) wou(product *awips.TextProduct, productID *models.RecordID, handler *Handler) (bool, error) {
	watch.WOU = productID

	if watch.Type == "" {
		typeRegexp := regexp.MustCompile("W(S|T)")
		typeStr := typeRegexp.FindString(product.Text)
		if typeStr == "WS" {
			watch.Type = "SV"
		} else if typeStr == "WT" {
			watch.Type = "TO"
		}
	}

	issued := product.Issued.UTC()
	watch.Issued = &issued

	if len(product.Segments) == 0 {
		return false, fmt.Errorf("WOU has 0 segments")
	}

	watch.Expires = &product.Segments[0].Expires

	var err error
	var current *db.Watch

	effectRegexp := regexp.MustCompile("((REMAINS|IS NO LONGER) IN EFFECT)")
	effectString := effectRegexp.FindString(product.Text)
	if effectString != "" {
		id := models.NewRecordID("severe_watch", db.WatchID{
			Number:    watch.Number,
			Phenomena: watch.Type,
			Year:      issued.Year(),
		})
		current, err = surrealdb.Select[db.Watch](handler.DB, id)
		if err != nil {
			return false, fmt.Errorf("error getting current watch %s: %s", id.String(), err.Error())
		}
	}

	if current.ID != nil {
		if current.Expires.Before(*watch.Expires) {
			current.Expires.Time = *watch.Expires
		}
		_, err = surrealdb.Merge[db.Watch, models.RecordID](handler.DB, *current.ID, current)
		return true, err
	}

	return false, nil
}

func (watch *Watch) wwp(product *awips.TextProduct, productID *models.RecordID, handler *Handler) error {
	watch.WWP = &db.WatchWWP{
		Text: productID,
	}

	stormMotionRegexp := regexp
}

func (watch *Watch) isComplete() bool {
	return watch.Issued != nil && watch.SAW != nil && watch.SEL != nil && watch.WOU != nil && watch.WWP != nil && watch.Type != "" && watch.Polygon != nil
}

func (watches *Watches) newWatch(number int) *Watch {
	watches.lock.Lock()
	defer watches.lock.Unlock()
	watch := &Watch{
		Number: number,
	}
	watches.watches = append(watches.watches, watch)
	return watch
}

func (watches *Watches) getWatch(number int) *Watch {
	watches.lock.Lock()
	defer watches.lock.Unlock()
	for _, w := range watches.watches {
		if w.Number == number {
			return w
		}
	}
	return nil
}

func (watches *Watches) removeWatch(number int) *Watch {
	watches.lock.Lock()
	defer watches.lock.Unlock()
	newWatches := []*Watch{}
	for _, w := range watches.watches {
		if w.Number != number {
			newWatches = append(newWatches, w)
		}
	}
	watches.watches = newWatches
	return nil
}
