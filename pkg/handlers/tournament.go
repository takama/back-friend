package handlers

import (
	"database/sql"
	"net/http"

	"github.com/takama/backer"
	"github.com/takama/backer/datastore"
	"github.com/takama/backer/player"
	"github.com/takama/backer/tournament"
	"github.com/takama/bit"
)

// TournamentDetails shows tournament information
func (h *Handler) TournamentDetails(c bit.Control) {
	id, ok := decodeNumber(":id", c)
	if !ok || id == 0 {
		badRequest(c)
		return
	}
	t, err := tournament.Find(id, h.db)
	if err != nil {
		switch err {
		case sql.ErrNoRows, datastore.ErrRecordNotFound:
			notFound(c)
		default:
			serviceError(err, c)
		}
		return
	}
	c.Code(http.StatusOK)
	c.Body(t.Tournament)
}

// TournamentAnnounce announce existing tournament with specified deposit
func (h *Handler) TournamentAnnounce(c bit.Control) {
	id, ok := decodeNumber(":id", c)
	if !ok || id == 0 {
		badRequest(c)
		return
	}
	t, err := tournament.Find(id, h.db)
	if err != nil {
		if err != sql.ErrNoRows && err != datastore.ErrRecordNotFound {
			serviceError(err, c)
			return
		}
		t, err = tournament.New(id, h.db)
		if err != nil {
			serviceError(err, c)
			return
		}
	}
	record := new(struct {
		backer.Points `json:"deposit"`
	})
	if !decodeRecord(record, c) {
		return
	}
	err = t.Announce(record.Points)
	if err != nil {
		switch err {
		case sql.ErrNoRows, datastore.ErrRecordNotFound:
			notFound(c)
		case tournament.ErrAllreadyFinished,
			tournament.ErrPlayersAlreadyJoined:
			notAcceptable(err, c)
		default:
			serviceError(err, c)
		}
		return
	}
	c.Code(http.StatusOK)
	c.Body(t.Tournament)
}

// TournamentJoin joins players and backers if necessary
func (h *Handler) TournamentJoin(c bit.Control) {
	id, ok := decodeNumber(":id", c)
	if !ok || id == 0 {
		badRequest(c)
		return
	}
	t, err := tournament.Find(id, h.db)
	if err != nil {
		switch err {
		case sql.ErrNoRows, datastore.ErrRecordNotFound:
			notFound(c)
		default:
			serviceError(err, c)
		}
		return
	}
	record := new(struct {
		Player  string   `json:"player"`
		Backers []string `json:"backers"`
	})
	if !decodeRecord(record, c) {
		return
	}
	bidder, err := player.Find(record.Player, h.db)
	if err != nil {
		switch err {
		case sql.ErrNoRows, datastore.ErrRecordNotFound:
			notFound(c)
		default:
			serviceError(err, c)
		}
		return
	}
	bidders := make([]backer.Player, 0)
	bidders = append(bidders, bidder)
	if len(record.Backers) > 0 {
		for _, backerID := range record.Backers {
			backer, err := player.Find(backerID, h.db)
			if err != nil {
				switch err {
				case sql.ErrNoRows, datastore.ErrRecordNotFound:
					notFound(c)
				default:
					serviceError(err, c)
				}
				return
			}
			bidders = append(bidders, backer)
		}
	}
	err = t.Join(bidders...)
	if err != nil {
		switch err {
		case sql.ErrNoRows, datastore.ErrRecordNotFound:
			notFound(c)
		case tournament.ErrAllreadyFinished,
			tournament.ErrCouldNotJoinTwice,
			player.ErrInsufficientPoints:
			notAcceptable(err, c)
		default:
			serviceError(err, c)
		}
		return
	}
	c.Code(http.StatusOK)
	c.Body(t.Tournament)
}

// TournamentResult finishes the tournament with winners result
func (h *Handler) TournamentResult(c bit.Control) {
	id, ok := decodeNumber(":id", c)
	if !ok || id == 0 {
		badRequest(c)
		return
	}
	t, err := tournament.Find(id, h.db)
	if err != nil {
		switch err {
		case sql.ErrNoRows, datastore.ErrRecordNotFound:
			notFound(c)
		default:
			serviceError(err, c)
		}
		return
	}
	record := new(struct {
		Winners []struct {
			Player string        `json:"player"`
			Prize  backer.Points `json:"prize"`
		} `json:"winners"`
	})
	if !decodeRecord(record, c) {
		return
	}
	winners := make(map[backer.Player]backer.Points, len(record.Winners))
	for _, winner := range record.Winners {
		participant, err := player.Find(winner.Player, h.db)
		if err != nil {
			switch err {
			case sql.ErrNoRows, datastore.ErrRecordNotFound:
				notFound(c)
			default:
				serviceError(err, c)
			}
			return
		}
		winners[participant] = winner.Prize
	}
	err = t.Result(winners)
	if err != nil {
		switch err {
		case sql.ErrNoRows, datastore.ErrRecordNotFound:
			notFound(c)
		case tournament.ErrAllreadyFinished,
			tournament.ErrWinnerIsNotMember:
			notAcceptable(err, c)
		default:
			serviceError(err, c)
		}
		return
	}
	c.Code(http.StatusOK)
	c.Body(t.Tournament)
}
