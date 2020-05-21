package admin

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	guuid "github.com/google/uuid"

	"gitlab.com/IIIS/backend/facade/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/spf13/viper"
	"gitlab.com/IIIS/backend/facade/auth/jwt"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// The list of error types returned from account resource.
var (
//ErrRoomValidation = errors.New("room validation error")
)

// RoomStore defines database operations for account management.
type FileStore interface {
	//List(*database.RoomFilter) ([]models.Room, int, error)
	Create(*models.File) error
	Get(id int) (*models.File, error)
	GetByUUID(uuid string) (*models.File, error)
	Update(*models.File) error
	Delete(*models.File) error
}

// RoomResource implements account management handler.
type FileResource struct {
	Store FileStore
	Acc   AccountStore
}

// NewFileResource creates and returns an account resource.
func NewFileResource(store FileStore, acc AccountStore) *FileResource {
	return &FileResource{
		Store: store,
		Acc:   acc,
	}
}

func (rs *FileResource) router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(rs.accountCtx)
	//r.Get("/", rs.list)
	r.Post("/upload", rs.upload)
	r.Route("/{fileID}", func(r chi.Router) {
		r.Use(rs.fileCtx)
		r.Get("/head", rs.header)
		/*
			r.Put("/", rs.update)
			r.Delete("/", rs.delete)*/
	})
	return r
}
func (rs *FileResource) fileCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "fileID")
		file, err := rs.Store.GetByUUID(id)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), ctxFile, file)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (rs *FileResource) accountCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := jwt.ClaimsFromCtx(r.Context())
		fmt.Println(claims.ID)
		account, err := rs.Acc.GetByUUID(claims.ID)
		if err != nil {
			// account deleted while access token still valid
			render.Render(w, r, ErrUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ctxAccount, account)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (rs *FileResource) header(w http.ResponseWriter, r *http.Request) {
	file := r.Context().Value(ctxFile).(*models.File)

	fmt.Println("Mime")
	fmt.Println(file.Mime)

	if !file.IsTableFile() {
		render.Respond(w, r, ErrUnsupportedMediaType)
		return
	}

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
	row := rows[0]
	var header []string

	//for _, row := range rows {
	for _, colCell := range row {
		if len(colCell) > 0 {
			header = append(header, colCell)
			fmt.Print(colCell, "\t")
		}
	}
	fmt.Println()
	//}

	render.Respond(w, r, header)
}

func (rs *FileResource) upload(w http.ResponseWriter, r *http.Request) {
	/*f, err := database.NewRoomFilter(r.URL.Query())
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	al, count, err := rs.Store.List(f)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	render.Respond(w, r, newRoomListResponse(al, count))*/
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)
	fmt.Printf("MIME Header: %+v\n", handler.Header["Content-Type"][0])

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	mime := handler.Header["Content-Type"][0]

	uuid := guuid.New().String()

	acc := r.Context().Value(ctxAccount).(*models.Account)

	fileDB := &models.File{
		Name:      handler.Filename,
		Size:      handler.Size,
		Mime:      mime,
		UUID:      uuid,
		AccountID: acc.ID,
	}

	if err := rs.Store.Create(fileDB); err != nil {
		fmt.Println(err)
		return
	}

	fileData, err := rs.Store.GetByUUID(uuid)

	fmt.Println(fileData)

	fileData.AccountData = (*models.AccountPure)(fileData.Account)

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	filePathBase := viper.GetString("file_upload_path")
	//fn := []byte(handler.Filename + time.Now().Format("2006-01-02 15:04:05"))
	//hash := sha.Sum256(fn)
	filePath := filePathBase + "/" + uuid

	/*err :=*/
	ioutil.WriteFile(filePath, fileBytes, 0644)

	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!

	render.Respond(w, r, fileData)
}

/*


type roomRequest struct {
	*models.Room
}

func (d *roomRequest) Bind(r *http.Request) error {
	return nil
}

type roomResponse struct {
	*models.Room
}

func newRoomResponse(r *models.Room) *roomResponse {
	resp := &roomResponse{Room: r}
	return resp
}

type roomListResponse struct {
	Rooms []models.Room `json:"rooms"`
	Count int           `json:"count"`
}

func newRoomListResponse(a []models.Room, count int) *roomListResponse {
	resp := &roomListResponse{
		Rooms: a,
		Count: count,
	}
	return resp
}

func (rs *RoomResource) list(w http.ResponseWriter, r *http.Request) {
	f, err := database.NewRoomFilter(r.URL.Query())
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	al, count, err := rs.Store.List(f)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	render.Respond(w, r, newRoomListResponse(al, count))
}

func (rs *RoomResource) get(w http.ResponseWriter, r *http.Request) {
	rm := r.Context().Value(ctxRoom).(*models.Room)
	render.Respond(w, r, newRoomResponse(rm))
}

func (rs *RoomResource) delete(w http.ResponseWriter, r *http.Request) {
	rm := r.Context().Value(ctxRoom).(*models.Room)
	if err := rs.Store.Delete(rm); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Respond(w, r, http.NoBody)
}

func (rs *RoomResource) create(w http.ResponseWriter, r *http.Request) {
	data := &roomRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := rs.Store.Create(data.Room); err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrRoomValidation, err.(validation.Errors)))
			return
		}
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Respond(w, r, newRoomResponse(data.Room))
}

func (rs *RoomResource) update(w http.ResponseWriter, r *http.Request) {
	rm := r.Context().Value(ctxRoom).(*models.Room)
	data := &roomRequest{Room: rm}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := rs.Store.Update(rm); err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrRoomValidation, err.(validation.Errors)))
			return
		}
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Respond(w, r, newRoomResponse(rm))
}
*/
