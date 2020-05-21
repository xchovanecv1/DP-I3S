// Package gateway ties together gateway resources and handlers.
package gateway

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/IIIS/backend/facade/database"
	"gitlab.com/IIIS/backend/facade/logging"
	"gitlab.com/IIIS/backend/facade/models"
)

// API provides application resources and handlers.
type API struct {
	DoorAccess      *database.DoorAccessStore
	Card            *database.CardStore
	CardLog         *database.CardLogStore
	Gateway         *database.GatewayStore
	Account         *database.AccountStore
	Room            *database.RoomStore
	PropertyTransit *database.PropertyTransitStore
}

type acccessRequest struct {
	MAC  string
	Card string
	Mid  uint64
}

// NewAPI configures and returns application API.
func NewAPI(db *pg.DB) (*API, error) {
	/*	accountStore := database.NewAccountStore(db)
		account := NewAccountResource(accountStore)

		profileStore := database.NewProfileStore(db)
		profile := NewProfileResource(profileStore)
	*/
	doorAccessStore := database.NewDoorAccessStore(db)
	cardStore := database.NewCardStore(db)
	cardLogStore := database.NewCardLogStore(db)
	gwStore := database.NewGatewayStore(db)
	accountStore := database.NewAccountStore(db)
	roomStore := database.NewRoomStore(db)
	propTransitStore := database.NewPropertyTransitStore(db)

	api := &API{
		DoorAccess:      doorAccessStore,
		Card:            cardStore,
		CardLog:         cardLogStore,
		Gateway:         gwStore,
		Account:         accountStore,
		Room:            roomStore,
		PropertyTransit: propTransitStore,
		//		Profile: profile,
	}

	return api, nil
}

// Router provides application routes.
func (a *API) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Post("/access", a.handleAccess)
	r.Post("/props", a.handleProps)
	r.Post("/hb", a.hearthBeat)

	return r
}

func (a *API) handleAccess(w http.ResponseWriter, r *http.Request) {

	data := &acccessRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	acs := &models.DoorAccess{}

	card, err := a.Card.GetByCode(data.Card)

	if err != nil && card.ID != 0 {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	var activeUser int

	//Account with card code exists

	cardLog := &models.CardLog{}

	if card.ID != 0 {
		acs.CardID = card.ID
		/*
			acc, _ := a.Account.GetAccountByCard(card.ID);
			if err != nil {
				render.Render(w, r, ErrInvalidRequest(err))
				return
			}
			log.Println("Acc: %+v\n", acc)
		*/
		if card.Link == "Account" && card.LinkID != 0 {
			activeUser = card.LinkID
		}
	} else { // Add access as new request
		acs.CardCode = data.Card
		cardLog.Code = data.Card
		cardLog.Type = "MIFARE 1K"
	}

	fmt.Printf("User: %v \n", activeUser)

	gw, err := a.Gateway.GetByAddress(data.MAC)

	if err != nil && gw.ID != 0 {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if card.ID == 0 {

		if gw.ID != 0 {
			cardLog.GatewayID = gw.ID
		}
		perfLog := viper.GetBool("gateway_card_loging")
		if perfLog && len(cardLog.Code) > 0 {
			a.CardLog.AddLog(cardLog)
		}

		render.Render(w, r, ErrBadRequest)
		return
	}
	//Account with card code exists
	if gw.ID != 0 {
		acs.GatewayID = gw.ID

		// Gateway got room assigned
		if gw.RoomID != 0 {
			rm, err := a.Room.Get(gw.RoomID)

			if err != nil && gw.ID != 0 {
				render.Render(w, r, ErrInvalidRequest(err))
				return
			}

			// Routine for clearing active user
			if rm.ActiveUserID != 0 {
				if rm.ActiveUserID == activeUser {
					rm.ActiveUserID = 0
					//Remove user and declare his exit of room
				} else {
					//Set new active user and close previous one
					rm.ActiveUserID = activeUser
				}
				//Determine if user is exiting room, or other user
				// Remove active users access and log a new one
			} else {
				rm.ActiveUserID = activeUser
				// No active user in room, just log
				err := a.PropertyTransit.AuthPriorTransits(gw.RoomID, activeUser)
				fmt.Printf("Prior auth %v", err)
			}

			err = a.Room.UpdateActive(rm)

			if err != nil {
				render.Render(w, r, ErrInvalidRequest(err))
				return
			}
		}
	} else { // Add access as new request
		acs.MAC = data.MAC
	}

	err = a.DoorAccess.AddAccess(acs)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Respond(w, r, data)

	/*
		if err := rs.Store.Create(data.Account); err != nil {
			switch err.(type) {
			case validation.Errors:
				render.Render(w, r, ErrValidation(ErrAccountValidation, err.(validation.Errors)))
				return
			}
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		render.Respond(w, r, newAccountResponse(data.Account))*/
}

type propResponse struct {
	Card          string `json:"card"`
	Transit       string `json:"trans"`
	Authenticated bool   `json:"auth"`
	Mid           uint64 `json:"mid"`
}

func (a *API) handleProps(w http.ResponseWriter, r *http.Request) {

	data := &acccessRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	fmt.Println("Trans req")
	fmt.Println(data)

	ret := &propResponse{
		Card:          data.Card,
		Mid:           data.Mid,
		Authenticated: true,
	}

	pTrans := &models.PropertyTransit{}

	card, err := a.Card.GetByCode(data.Card)

	if err != nil && card.ID != 0 {
		render.Respond(w, r, ret)
		return
	}

	cardLog := &models.CardLog{}
	//Account with card code exists
	if card.ID != 0 {
		if card.Link == "Property" {
			pTrans.CardID = card.ID
			pTrans.PropertyID = card.LinkID
		}
	} else { // Add access as new request
		//acs.CardCode = data.Card
		cardLog.Code = data.Card
		cardLog.Type = "EM4100"
	}

	gw, err := a.Gateway.GetByAddress(data.MAC)

	if err != nil && gw.ID != 0 {
		render.Respond(w, r, ret)
		return
	}

	//Account with card code exists
	if card.ID == 0 {

		if gw.ID != 0 {
			cardLog.GatewayID = gw.ID
		}
		perfLog := viper.GetBool("gateway_card_loging")
		if perfLog && len(cardLog.Code) > 0 {
			a.CardLog.AddLog(cardLog)
		}

		render.Respond(w, r, ret)
		return
	}
	if gw.ID != 0 {
		//acs.GatewayID = gw.ID
		// Gateway got room assigned
		if gw.RoomID != 0 {

			err := a.Room.ValidateActiveUser(gw.RoomID)
			fmt.Println(gw.RoomID)
			if err != nil {
				render.Render(w, r, ErrInvalidRequest(err))
				return
			}

			rm, err := a.Room.GetPure(gw.RoomID)

			if err != nil && gw.ID != 0 {
				render.Render(w, r, ErrInvalidRequest(err))
				return
			}

			pTrans.RoomID = rm.ID

			// Routine for clearing active user
			if rm.ActiveUserID != 0 {
				pTrans.UserID = rm.ActiveUserID
				//Determine if user is exiting room, or other user
				// Remove active users access and log a new one
			} else {
				pTrans.UserID = 0
				// No active user in room, just log
			}

		}
	}

	lastTrans, _ := a.PropertyTransit.GetLatestByProp(pTrans.PropertyID)

	// Default value if no previous transition is EXIT
	if lastTrans.ID == 0 || len(lastTrans.Transition) == 0 {
		pTrans.Transition = "EXIT"
	} else {
		if lastTrans.Transition == "EXIT" || lastTrans.Transition == "POSSES" {
			pTrans.Transition = "ENTER"
		} else {
			pTrans.Transition = "EXIT"
		}
	}

	err = a.PropertyTransit.Create(pTrans)

	fmt.Printf("trans %+v %v", pTrans, err)
	if err != nil {
		render.Respond(w, r, ret)
		return
	}

	ret.Transit = pTrans.Transition
	ret.Authenticated = pTrans.UserID != 0

	fmt.Printf("resi: %+v\n", ret)

	render.Respond(w, r, ret)

}

func (d *acccessRequest) Bind(r *http.Request) error {
	log.Println(d.MAC)

	return validation.ValidateStruct(d,
		validation.Field(&d.MAC, validation.Required, is.MAC),
		validation.Field(&d.Card, validation.Required, is.Digit),
	)
}

func logRequest(r *http.Request) logrus.FieldLogger {
	return logging.GetLogEntry(r)
}

type hearthRequest struct {
	MAC  string
	Tick uint64
}

func (d *hearthRequest) Bind(r *http.Request) error {
	log.Println(d.MAC)

	return validation.ValidateStruct(d,
		validation.Field(&d.MAC, validation.Required, is.MAC),
	)
}

func (a *API) hearthBeat(w http.ResponseWriter, r *http.Request) {

	data := &hearthRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	a.Gateway.SetLastActive(data.MAC)

	fmt.Println("Hearth req")
	fmt.Println(data)

	render.Respond(w, r, data)

}
