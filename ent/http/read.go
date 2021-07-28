// Code generated by entc, DO NOT EDIT.

package http

import (
	"elk-example/ent"
	"elk-example/ent/group"
	"elk-example/ent/pet"
	"elk-example/ent/user"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/liip/sheriff"
	"github.com/masseelch/render"
	"go.uber.org/zap"
)

// Read fetches the ent.Group identified by a given url-parameter from the
// database and renders it to the client.
func (h *GroupHandler) Read(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Read"))
	// ID is URL parameter.
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		l.Error("error getting id from url parameter", zap.String("id", chi.URLParam(r, "id")), zap.Error(err))
		render.BadRequest(w, r, "id must be an integer greater zero")
		return
	}
	// Create the query to fetch the Group
	q := h.client.Group.Query().Where(group.ID(id))
	e, err := q.Only(r.Context())
	if err != nil {
		switch {
		case ent.IsNotFound(err):
			msg := stripEntError(err)
			l.Info(msg, zap.Int("id", id), zap.Error(err))
			render.NotFound(w, r, msg)
		case ent.IsNotSingular(err):
			msg := stripEntError(err)
			l.Error(msg, zap.Int("id", id), zap.Error(err))
			render.BadRequest(w, r, msg)
		default:
			l.Error("error fetching group from db", zap.Int("id", id), zap.Error(err))
			render.InternalServerError(w, r, nil)
		}
		return
	}
	d, err := sheriff.Marshal(&sheriff.Options{
		IncludeEmptyTag: true,
		Groups:          []string{"group"},
	}, e)
	if err != nil {
		l.Error("serialization error", zap.Int("id", id), zap.Error(err))
		render.InternalServerError(w, r, nil)
		return
	}
	l.Info("group rendered", zap.Int("id", id))
	render.OK(w, r, d)
}

// Read fetches the ent.Pet identified by a given url-parameter from the
// database and renders it to the client.
func (h *PetHandler) Read(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Read"))
	// ID is URL parameter.
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		l.Error("error getting id from url parameter", zap.String("id", chi.URLParam(r, "id")), zap.Error(err))
		render.BadRequest(w, r, "id must be an integer greater zero")
		return
	}
	// Create the query to fetch the Pet
	q := h.client.Pet.Query().Where(pet.ID(id))
	e, err := q.Only(r.Context())
	if err != nil {
		switch {
		case ent.IsNotFound(err):
			msg := stripEntError(err)
			l.Info(msg, zap.Int("id", id), zap.Error(err))
			render.NotFound(w, r, msg)
		case ent.IsNotSingular(err):
			msg := stripEntError(err)
			l.Error(msg, zap.Int("id", id), zap.Error(err))
			render.BadRequest(w, r, msg)
		default:
			l.Error("error fetching pet from db", zap.Int("id", id), zap.Error(err))
			render.InternalServerError(w, r, nil)
		}
		return
	}
	d, err := sheriff.Marshal(&sheriff.Options{
		IncludeEmptyTag: true,
		Groups:          []string{"pet"},
	}, e)
	if err != nil {
		l.Error("serialization error", zap.Int("id", id), zap.Error(err))
		render.InternalServerError(w, r, nil)
		return
	}
	l.Info("pet rendered", zap.Int("id", id))
	render.OK(w, r, d)
}

// Read fetches the ent.User identified by a given url-parameter from the
// database and renders it to the client.
func (h *UserHandler) Read(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Read"))
	// ID is URL parameter.
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		l.Error("error getting id from url parameter", zap.String("id", chi.URLParam(r, "id")), zap.Error(err))
		render.BadRequest(w, r, "id must be an integer greater zero")
		return
	}
	// Create the query to fetch the User
	q := h.client.User.Query().Where(user.ID(id))
	// Eager load edges that are required on read operation.
	q.WithPets()
	e, err := q.Only(r.Context())
	if err != nil {
		switch {
		case ent.IsNotFound(err):
			msg := stripEntError(err)
			l.Info(msg, zap.Int("id", id), zap.Error(err))
			render.NotFound(w, r, msg)
		case ent.IsNotSingular(err):
			msg := stripEntError(err)
			l.Error(msg, zap.Int("id", id), zap.Error(err))
			render.BadRequest(w, r, msg)
		default:
			l.Error("error fetching user from db", zap.Int("id", id), zap.Error(err))
			render.InternalServerError(w, r, nil)
		}
		return
	}
	d, err := sheriff.Marshal(&sheriff.Options{
		IncludeEmptyTag: true,
		Groups:          []string{"user"},
	}, e)
	if err != nil {
		l.Error("serialization error", zap.Int("id", id), zap.Error(err))
		render.InternalServerError(w, r, nil)
		return
	}
	l.Info("user rendered", zap.Int("id", id))
	render.OK(w, r, d)
}
