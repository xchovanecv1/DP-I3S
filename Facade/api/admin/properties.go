package admin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gitlab.com/IIIS/backend/facade/models"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gitlab.com/IIIS/backend/facade/database"
)

// The list of error types returned from account resource.
var (
	ErrPropertyValidation = errors.New("property validation error")
)

// PropertyStore defines database operations for account management.
type PropertyStore interface {
	List(*database.PropertyFilter) ([]models.Property, int, error)
	Create(*models.Property) error
	CreateMultiple([]*models.Property) error
	Get(id int) (*models.Property, error)
	GetByUUID(uuid string) (*models.Property, error)
	Update(*models.Property) error
	Delete(*models.Property) error
}

// PropertyResource implements account management handler.
type PropertyResource struct {
	Store      PropertyStore
	Transition *database.PropertyTransitStore
	Accounts   *database.AdmAccountStore
	Files      *FileResource
	Rooms      *RoomResource
	Cards      *CardResource
	CardLogs   *database.CardLogStore
	Imports    *ImportResource
}

// NewPropertyResource creates and returns an account resource.
func NewPropertyResource(store PropertyStore, trans *database.PropertyTransitStore, ac *database.AdmAccountStore, file *FileResource, room *RoomResource, card *CardResource, cl *database.CardLogStore, imp *ImportResource) *PropertyResource {
	return &PropertyResource{
		Store:      store,
		Transition: trans,
		Accounts:   ac,
		Files:      file,
		Rooms:      room,
		Cards:      card,
		CardLogs:   cl,
		Imports:    imp,
	}
}

func (rs *PropertyResource) router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", rs.list)
	r.Get("/pending", rs.possesPending)
	r.Post("/pending", rs.postPending)
	r.Post("/", rs.create)
	r.Post("/commit", rs.commit)
	r.Route("/{propertyID}", func(r chi.Router) {
		r.Use(rs.propertyCtx)
		r.Get("/", rs.get)
		r.Put("/", rs.update)
		r.Delete("/", rs.delete)
		r.Post("/claim", rs.claim)
		r.Route("/cards", func(r chi.Router) {
			r.Get("/", rs.getCards)
			r.Delete("/", rs.deleteCardLink)
			//r.Post("/", rs.addCardToProp)
			r.Post("/log", rs.addCardFromLog)
		})
		r.Route("/location", func(r chi.Router) {
			r.Get("/", rs.getLocation)
		})
	})
	r.Route("/import", func(r chi.Router) {
		r.Route("/{fileID}", func(r chi.Router) {
			r.Use(rs.fileCtx)
			r.Post("/", rs.importFile)
			/*
				r.Put("/", rs.update)
				r.Delete("/", rs.delete)*/
		})
		r.Get("/headers", rs.importHeader)
	})
	return r
}

// Import

func (rs *PropertyResource) fileCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "fileID")
		file, err := rs.Files.Store.GetByUUID(id)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), ctxFile, file)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type ImportItem struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	FileKey   int    `json:"fileKey"`
	Code      string `json:"code"`
	Operation string `json:"op_code"`
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func (rs *PropertyResource) importHeader(w http.ResponseWriter, r *http.Request) {
	hd := models.Property{}
	t := reflect.TypeOf(hd)

	// Get the type and kind of our user variable
	//fmt.Println("Type:", t.Name())
	//fmt.Println("Kind:", t.Kind())

	importHeader := []ImportItem{}

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := t.Field(i)

		// Get the field tag value
		tag := field.Tag.Get("import")

		params := strings.Split(tag, ",")

		importable := contains(params, "true")
		if importable {
			jsonTag := field.Tag.Get("json")
			jparams := strings.Split(jsonTag, ",")
			jsonName := jparams[0]
			importHeader = append(importHeader, ImportItem{
				Name: field.Name,
				Type: field.Type.String(),
				Code: jsonName,
			})

			//fmt.Printf("%d. %v (%v) %v\n", i+1, field.Name, field.Type.Name(), jsonName)
		}

	}
	render.Respond(w, r, importHeader)
}

//https://www.myonlinetraininghub.com/excel-date-and-time
//https://stackoverflow.com/questions/13850605/convert-excel-date-serial-number-to-regular-date
//ParseExcelDate convert serial date format used in xlsx files to time.Time
func ParseExcelDate(val string) (*time.Time, error) {
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return nil, err
	}
	ipart := int(f)
	decpart := f - float64(ipart)

	hours := int(decpart / 0.04166666666)
	decpart = decpart - float64(hours)*0.04166666666

	minutes := int(decpart / 0.00069444)

	decpart = decpart - float64(minutes)*0.00069444
	seconds := int(decpart / 0.0000115740740740741)

	var buf = time.Date(1899, time.December, 30, hours, minutes, seconds, 0, time.UTC)
	buf = buf.AddDate(0, 0, ipart)
	return &buf, nil
}

type importList []ImportItem

func (d *importList) Bind(r *http.Request) error {
	return nil
}

type importParseResult struct {
	Valid       []*models.Property       `json:"valid"`
	ValidRows   []int                    `json:"valid_rows"`
	InvalidRows []*importValidationError `json:"invalid_rows"`
	Count       int                      `json:"count"`
}

type importValidationError struct {
	Item   int    `json:"item"`
	Reason string `json:"reason"`
}

func (rs *PropertyResource) importFile(w http.ResponseWriter, r *http.Request) {
	file := r.Context().Value(ctxFile).(*models.File)

	xlsHeader := r.URL.Query().Get("header")
	includeHeader := len(xlsHeader) > 0
	fmt.Printf("Got header %\n", includeHeader)
	//includeHeader := true

	fmt.Println("Mime")
	fmt.Println(file.Mime)

	if !file.IsTableFile() {
		render.Respond(w, r, ErrUnsupportedMediaType)
		return
	}

	var data importList
	if err := render.Bind(r, &data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	//

	//for i, d := range data {
	//	fmt.Println(i, d)
	//mapData[d.Name] =
	//}
	filePathBase := viper.GetString("file_upload_path")
	filePath := filePathBase + "/" + file.UUID

	xlsx, err := excelize.OpenFile(filePath)

	if err != nil {
		render.Respond(w, r, ErrUnsupportedMediaType)
		fmt.Println(err)

		return
	}
	fSheetC := xlsx.SheetCount
	fmt.Println(fSheetC)

	// Retrieve name of first sheet
	fSheet := xlsx.GetSheetName(1)
	fmt.Println(fSheet)

	rows := xlsx.GetRows(fSheet)

	if len(rows) < 1 {
		render.Respond(w, r, ErrUnprocessableEntity)
	}

	valid := []*models.Property{}
	var iRows []*importValidationError
	count := 0

	for i, row := range rows {
		if includeHeader && i == 0 {
			continue
		}
		mapData := make(map[string]interface{})
		for idx, colCell := range row {

			for _, d := range data {

				if idx == d.FileKey {
					var val = strings.TrimSpace(colCell)
					if d.Type == "time.Time" {
						date, err := ParseExcelDate(val)
						if err != nil {
							continue
						}
						mapData[d.Name] = date
						val = date.String()
					} else {
						mapData[d.Name] = val
					}

					if d.Operation == "F_RM_ID_CODE" || d.Operation == "FoC_RM_ID_CODE" {
						fmt.Println("Find room by CODE", val)
						room, err := rs.Rooms.Store.GetCodePure(val)
						if err != nil {
							fmt.Printf("[IMPORT][%d - %d] %s", i, idx, err)
						} else {
							mapData[d.Name] = room.ID
							mapData["Room"] = room
							continue
						}
					}

					if d.Operation == "FoC_RM_ID_CODE" {
						fmt.Println("Creating room by CODE", val)
						rm := &models.Room{
							Name: val,
							Code: val,
						}

						err := rm.Validate()

						if err == nil {
							err = rs.Rooms.Store.Create(rm)
							if err != nil {
								fmt.Printf("[IMPORT][%d - %d] Failed to create room: %s", i, idx, err)
							} else {
								mapData[d.Name] = rm.ID
								mapData["Room"] = rm
							}
						} else {
							fmt.Printf("[IMPORT][%d - %d] Failed to create room: %s", i, idx, err)
						}
						continue
					}
				}
			}
		}
		var prop models.Property

		config := &mapstructure.DecoderConfig{
			WeaklyTypedInput: true,
			Result:           &prop,
		}

		decoder, _ := mapstructure.NewDecoder(config)

		decoder.Decode(mapData)
		vld := prop.Validate()
		if vld == nil {
			valid = append(valid, &prop)
		} else {
			inv := importValidationError{
				Reason: vld.Error(),
				Item:   i,
			}
			iRows = append(iRows, &inv)
		}
		count++
	}

	res := importParseResult{
		Valid:       valid,
		InvalidRows: iRows,
		Count:       count,
	}
	//https://www.thepolyglotdeveloper.com/2017/04/decode-map-values-native-golang-structures/
	render.Respond(w, r, res)
}

//////////////

func (rs *PropertyResource) propertyCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "propertyID")
		property, err := rs.Store.GetByUUID(id)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), ctxProperty, property)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type propertyRequest struct {
	*models.Property
}

func (d *propertyRequest) Bind(r *http.Request) error {
	return nil
}

type propertyResponse struct {
	*models.Property
}

func newPropertyResponse(r *models.Property) *propertyResponse {
	resp := &propertyResponse{Property: r}
	return resp
}

type propertyListResponse struct {
	Properties []models.Property `json:"properties"`
	Count      int               `json:"count"`
}

func newPropertyListResponse(a []models.Property, count int) *propertyListResponse {
	resp := &propertyListResponse{
		Properties: a,
		Count:      count,
	}
	return resp
}

func (rs *PropertyResource) list(w http.ResponseWriter, r *http.Request) {
	f, err := database.NewPropertyFilter(r.URL.Query())
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	al, count, err := rs.Store.List(f)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	render.Respond(w, r, newPropertyListResponse(al, count))
}

func (rs *PropertyResource) get(w http.ResponseWriter, r *http.Request) {
	rm := r.Context().Value(ctxProperty).(*models.Property)
	render.Respond(w, r, newPropertyResponse(rm))
}

func (rs *PropertyResource) getCards(w http.ResponseWriter, r *http.Request) {
	rm := r.Context().Value(ctxProperty).(*models.Property)

	cards, count, err := rs.Cards.Store.ListLink("Property", rm.ID)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, NewCardListResponse(cards, count))
}

type cardClaimRequest struct {
	Reason       string `json:"reason"`
	TransitionID int    `json:"trans_id"`
}

func (d *cardClaimRequest) Bind(r *http.Request) error {
	return nil
}

func (rs *PropertyResource) claim(w http.ResponseWriter, r *http.Request) {

	acc := r.Context().Value(ctxAccount).(*models.Account)

	data := &cardClaimRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	tr, err := rs.Transition.Get(data.TransitionID)

	if err != nil || tr.ID == 0 {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if tr.Transition != "EXIT" {
		render.Render(w, r, ErrUnprocessableEntity)
		return
	}
	tr.Transition = "POSSES"
	tr.UserID = acc.ID
	tr.PostAuth = true
	tr.Comment = data.Reason

	tr.ID = 0

	err = rs.Transition.Create(tr)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Respond(w, r, tr)
}

func (rs *PropertyResource) possesPending(w http.ResponseWriter, r *http.Request) {

	acc := r.Context().Value(ctxAccount).(*models.Account)

	tr, err := rs.Transition.GetLastestUserExits(acc.ID)

	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Respond(w, r, tr)
}

type propTrasRequest []models.PropertyTransit

func (d *propTrasRequest) Bind(r *http.Request) error {
	return nil
}

func (rs *PropertyResource) postPending(w http.ResponseWriter, r *http.Request) {
	//prop := r.Context().Value(ctxProperty).(*models.Property)
	acc := r.Context().Value(ctxAccount).(*models.Account)

	var data propTrasRequest
	if err := render.Bind(r, &data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	//var req []models.PropertyTransit
	//req = data
	leng := len(data)
	if leng > 50 {
		leng = 50
	}
	var ins []models.PropertyTransit

	for _, pt := range data {

		pt.CreatedAt = time.Now()
		pt.UserID = acc.ID
		pt.Transition = "POSSES"
		pt.ID = 0

		ins = append(ins, pt)
	}

	err := rs.Transition.CreateMultiple(&ins)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Respond(w, r, ins)
}

type cardLogRequest struct {
	*models.CardLog
}

func (d *cardLogRequest) Bind(r *http.Request) error {
	return nil
}

func (rs *PropertyResource) addCardFromLog(w http.ResponseWriter, r *http.Request) {
	prop := r.Context().Value(ctxProperty).(*models.Property)

	data := &cardLogRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	cl := data.CardLog

	if cl.ID == 0 {
		render.Render(w, r, ErrUnprocessableEntity)
		return
	}

	log, err := rs.CardLogs.Get(cl.ID)

	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	card := &models.Card{
		Link:   "Property",
		LinkID: prop.ID,
	}

	err = rs.Cards.Store.CreateFromLog(card, log)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	err = rs.CardLogs.DeleteByCode(log)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	/*
		if err := rs.Store.Update(rm); err != nil {
			switch err.(type) {
			case validation.Errors:
				render.Render(w, r, ErrValidation(ErrPropertyValidation, err.(validation.Errors)))
				return
			}
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
	*/
	render.Respond(w, r, card)
}

func (rs *PropertyResource) deleteCardLink(w http.ResponseWriter, r *http.Request) {
	prop := r.Context().Value(ctxProperty).(*models.Property)

	data := &cardRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	fmt.Printf("%+v\n", data.Card)

	cl := data.Card

	if len(cl.UUID) == 0 {
		render.Render(w, r, ErrUnprocessableEntity)
		return
	}

	card, err := rs.Cards.Store.GetByUUID(cl.UUID)

	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if card.Link != "Property" || card.LinkID != prop.ID {
		render.Render(w, r, ErrUnprocessableEntity)
		return
	}

	card.Link = ""
	card.LinkID = 0

	err = rs.Cards.Store.Update(card)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Respond(w, r, http.NoBody)
}

type propLocationResponse struct {
	Transition string              `json:"type"`
	Reason     string              `json:"reason"`
	User       *models.AccountPure `json:"user"`
	Room       *models.Room        `json:"room"`
}

func (rs *PropertyResource) getLocation(w http.ResponseWriter, r *http.Request) {
	prop := r.Context().Value(ctxProperty).(*models.Property)

	res := &propLocationResponse{}

	tr, err := rs.Transition.GetLatestByProp(prop.ID)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if tr.UserID != 0 {
		u, err := rs.Accounts.Get(tr.UserID)

		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		res.User = &models.AccountPure{}
		res.User.Name = u.Name
		res.User.UUID = u.UUID
	}

	if tr.RoomID != 0 {
		rm, err := rs.Rooms.Store.Get(tr.RoomID)

		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		res.Room = rm
		res.Room.ActiveUser = nil
	}

	res.Transition = tr.Transition
	res.Reason = tr.Comment

	/*if err := rs.Store.Delete(rm); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}*/
	render.Respond(w, r, res)
}

func (rs *PropertyResource) delete(w http.ResponseWriter, r *http.Request) {
	rm := r.Context().Value(ctxProperty).(*models.Property)
	if err := rs.Store.Delete(rm); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Respond(w, r, http.NoBody)
}

func (rs *PropertyResource) create(w http.ResponseWriter, r *http.Request) {
	data := &propertyRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := rs.Store.Create(data.Property); err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrPropertyValidation, err.(validation.Errors)))
			return
		}
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Respond(w, r, newPropertyResponse(data.Property))
}

type propertyCommitRequest []*models.Property

func (d *propertyCommitRequest) Bind(r *http.Request) error {
	return nil
}

func (rs *PropertyResource) commit(w http.ResponseWriter, r *http.Request) {
	data := &propertyCommitRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := rs.Store.CreateMultiple(*data); err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrPropertyValidation, err.(validation.Errors)))
			return
		}
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	ids := make([]int, len(*data))

	for i, s := range *data {
		ids[i] = s.ID
	}

	if err := rs.Imports.logImport(r, "Properties", ids); err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrPropertyValidation, err.(validation.Errors)))
			return
		}
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Respond(w, r, (data))
}

func (rs *PropertyResource) update(w http.ResponseWriter, r *http.Request) {
	rm := r.Context().Value(ctxProperty).(*models.Property)
	data := &propertyRequest{Property: rm}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := rs.Store.Update(rm); err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrPropertyValidation, err.(validation.Errors)))
			return
		}
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Respond(w, r, newPropertyResponse(rm))
}
