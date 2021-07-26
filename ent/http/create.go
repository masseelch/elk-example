// Code generated by entc, DO NOT EDIT.

package http

import (
	"elk-example/ent"
	"elk-example/ent/group"
	"elk-example/ent/pet"
	"elk-example/ent/user"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/liip/sheriff"
	"github.com/masseelch/render"
	"go.uber.org/zap"
)

// Payload of a ent.Group create request.
type GroupCreateRequest struct {
}

// Create creates a new ent.Group and stores it in the database.
func (h GroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Create"))
	// Get the post data.
	var d GroupCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		l.Error("error decoding json", zap.Error(err))
		render.BadRequest(w, r, "invalid json string")
		return
	}
	// Validate the data.
	if err := h.validator.Struct(d); err != nil {
		if err, ok := err.(*validator.InvalidValidationError); ok {
			l.Error("error validating request data", zap.Error(err))
			render.InternalServerError(w, r, nil)
			return
		}
		l.Info("validation failed", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}
	// Save the data.
	b := h.client.Group.Create()
	// TODO: what about slice fields that have custom marshallers?
	// Store in database.
	e, err := b.Save(r.Context())
	if err != nil {
		l.Error("error saving group", zap.Error(err))
		render.InternalServerError(w, r, nil)
		return
	}
	// Reload entry.
	q := h.client.Group.Query().Where(group.ID(e.ID))
	e, err = q.Only(r.Context())
	if err != nil {
		switch err.(type) {
		case *ent.NotFoundError:
			msg := h.stripEntError(err)
			l.Info(msg, zap.Int("id", e.ID), zap.Error(err))
			render.NotFound(w, r, msg)
		default:
			l.Error("error fetching group from db", zap.Int("id", e.ID), zap.Error(err))
			render.InternalServerError(w, r, nil)
		}
		return
	}
	j, err := sheriff.Marshal(&sheriff.Options{
		IncludeEmptyTag: true,
		Groups:          []string{"group"},
	}, e)
	if err != nil {
		l.Error("serialization error", zap.Int("id", e.ID), zap.Error(err))
		render.InternalServerError(w, r, nil)
		return
	}
	l.Info("group rendered", zap.Int("id", e.ID))
	render.OK(w, r, j)
}

// Payload of a ent.Pet create request.
type PetCreateRequest struct {
	Name  *string `json:"name"`
	Age   *int    `json:"age" validate:"required,gt=0"`
	Owner *int    `json:"owner" validate:"required"`
}

// Create creates a new ent.Pet and stores it in the database.
func (h PetHandler) Create(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Create"))
	// Get the post data.
	var d PetCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		l.Error("error decoding json", zap.Error(err))
		render.BadRequest(w, r, "invalid json string")
		return
	}
	// Validate the data.
	if err := h.validator.Struct(d); err != nil {
		if err, ok := err.(*validator.InvalidValidationError); ok {
			l.Error("error validating request data", zap.Error(err))
			render.InternalServerError(w, r, nil)
			return
		}
		l.Info("validation failed", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}
	// Save the data.
	b := h.client.Pet.Create()
	// TODO: what about slice fields that have custom marshallers?
	if d.Name != nil {
		b.SetName(*d.Name)
	}
	if d.Age != nil {
		b.SetAge(*d.Age)
	}
	if d.Owner != nil {
		b.SetOwnerID(*d.Owner)

	}
	// Store in database.
	e, err := b.Save(r.Context())
	if err != nil {
		l.Error("error saving pet", zap.Error(err))
		render.InternalServerError(w, r, nil)
		return
	}
	// Reload entry.
	q := h.client.Pet.Query().Where(pet.ID(e.ID))
	e, err = q.Only(r.Context())
	if err != nil {
		switch err.(type) {
		case *ent.NotFoundError:
			msg := h.stripEntError(err)
			l.Info(msg, zap.Int("id", e.ID), zap.Error(err))
			render.NotFound(w, r, msg)
		default:
			l.Error("error fetching pet from db", zap.Int("id", e.ID), zap.Error(err))
			render.InternalServerError(w, r, nil)
		}
		return
	}
	j, err := sheriff.Marshal(&sheriff.Options{
		IncludeEmptyTag: true,
		Groups:          []string{"pet"},
	}, e)
	if err != nil {
		l.Error("serialization error", zap.Int("id", e.ID), zap.Error(err))
		render.InternalServerError(w, r, nil)
		return
	}
	l.Info("pet rendered", zap.Int("id", e.ID))
	render.OK(w, r, j)
}

// Payload of a ent.User create request.
type UserCreateRequest struct {
	Name *string `json:"name"`
	Age  *int    `json:"age"`
	Pets []int   `json:"pets"`
}

// Create creates a new ent.User and stores it in the database.
func (h UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Create"))
	// Get the post data.
	var d UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		l.Error("error decoding json", zap.Error(err))
		render.BadRequest(w, r, "invalid json string")
		return
	}
	// Validate the data.
	if err := h.validator.Struct(d); err != nil {
		if err, ok := err.(*validator.InvalidValidationError); ok {
			l.Error("error validating request data", zap.Error(err))
			render.InternalServerError(w, r, nil)
			return
		}
		l.Info("validation failed", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}
	// Save the data.
	b := h.client.User.Create()
	// TODO: what about slice fields that have custom marshallers?
	if d.Name != nil {
		b.SetName(*d.Name)
	}
	if d.Age != nil {
		b.SetAge(*d.Age)
	}
	if d.Pets != nil {
		b.AddPetIDs(d.Pets...)
	}
	// Store in database.
	e, err := b.Save(r.Context())
	if err != nil {
		l.Error("error saving user", zap.Error(err))
		render.InternalServerError(w, r, nil)
		return
	}
	// Reload entry.
	q := h.client.User.Query().Where(user.ID(e.ID))
	e, err = q.Only(r.Context())
	if err != nil {
		switch err.(type) {
		case *ent.NotFoundError:
			msg := h.stripEntError(err)
			l.Info(msg, zap.Int("id", e.ID), zap.Error(err))
			render.NotFound(w, r, msg)
		default:
			l.Error("error fetching user from db", zap.Int("id", e.ID), zap.Error(err))
			render.InternalServerError(w, r, nil)
		}
		return
	}
	j, err := sheriff.Marshal(&sheriff.Options{
		IncludeEmptyTag: true,
		Groups:          []string{"user"},
	}, e)
	if err != nil {
		l.Error("serialization error", zap.Int("id", e.ID), zap.Error(err))
		render.InternalServerError(w, r, nil)
		return
	}
	l.Info("user rendered", zap.Int("id", e.ID))
	render.OK(w, r, j)
}