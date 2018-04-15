package handlers

import (
	"database/sql"
	"net/http"

	"github.com/takama/backer"
	"github.com/takama/backer/datastore"
	"github.com/takama/backer/player"
	"github.com/takama/bit"
)

// PlayerDetails shows player information
func (h *Handler) PlayerDetails(c bit.Control) {
	id, ok := decodeString(":id", c)
	if !ok || id == "" {
		badRequest(c)
		return
	}
	p, err := player.Find(id, h.db)
	if err != nil {
		if err == sql.ErrNoRows || err == datastore.ErrRecordNotFound {
			notFound(c)
		} else {
			serviceError(err, c)
		}
		return
	}
	c.Code(http.StatusOK)
	c.Body(p.Player)
}

// PlayerFund funds existing player or create new one
func (h *Handler) PlayerFund(c bit.Control) {
	id, ok := decodeString(":id", c)
	if !ok || id == "" {
		badRequest(c)
		return
	}
	p, err := player.Find(id, h.db)
	if err != nil {
		if err != sql.ErrNoRows && err != datastore.ErrRecordNotFound {
			serviceError(err, c)
			return
		}
		p, err = player.New(id, h.db)
		if err != nil {
			serviceError(err, c)
			return
		}
	}
	m := new(struct {
		backer.Points `json:"points"`
	})
	if !decodeRecord(m, c) {
		return
	}
	err = p.Fund(m.Points)
	if err != nil {
		serviceError(err, c)
		return
	}
	c.Code(http.StatusOK)
	c.Body(p.Player)
}

// PlayerTake takes points from existing player
func (h *Handler) PlayerTake(c bit.Control) {
	id, ok := decodeString(":id", c)
	if !ok || id == "" {
		badRequest(c)
		return
	}
	p, err := player.Find(id, h.db)
	if err != nil {
		if err == sql.ErrNoRows || err == datastore.ErrRecordNotFound {
			notFound(c)
		} else {
			serviceError(err, c)
		}
		return
	}
	m := new(struct {
		backer.Points `json:"points"`
	})
	if !decodeRecord(m, c) {
		return
	}
	err = p.Take(m.Points)
	if err != nil {
		if err == player.ErrInsufficientPoints {
			notAcceptable(err, c)
			return
		}
		serviceError(err, c)
		return
	}
	c.Code(http.StatusOK)
	c.Body(p.Player)
}
