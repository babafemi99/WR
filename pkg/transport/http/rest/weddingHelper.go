package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/babafemi99/WR/internal/util"
	"github.com/babafemi99/WR/internal/values"
	"github.com/babafemi99/WR/pkg/model"
	"github.com/jackc/pgx/v5"
	"time"
)

func (a *API) VerifyWeddingId(id string) (*model.Wedding, string, string, error) {
	// get wedding from BD
	wedding, err := a.Deps.Repository.GetWeddingById(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, values.NotFound, "no wedding with this Id", errors.New("resource not found")
		}
		return nil, values.Error, "failed to load wedding", err
	}

	// check if wedding is today
	if !util.IsSameDayWithToday(wedding.WeddingDate) {
		if util.IsBeforeToday(wedding.WeddingDate) {
			return nil, values.Error, fmt.Sprintf("this wedding has expired, it was slated for %s",
				wedding.WeddingDate.String()), errors.New("wedding has expired")
		} else {
			return nil, values.Error, fmt.Sprintf("this wedding is not slated fot today, it's slated for %s",
				wedding.WeddingDate.String()), errors.New("wedding in future")
		}
	}

	// check if the wedding link has been toggled
	if !wedding.IsLive() {
		return nil, values.Unprocessable, "your link is not active yet, come back later", errors.New("inactive link")
	}

	return wedding, values.Success, "verification successful", nil
}

func (a *API) ToggleWedding(req model.ToggleWeddingReq) (string, string, error) {
	// verify body

	// persist to database
	err := a.Deps.Repository.ToggleWeddingLink(req)
	if err != nil {
		return values.Failed, "failed to toggle link", err
	}

	return values.Success, "toggled link successfully", nil
}

func (a *API) ToggleWeddingOff(weddingId string) (string, string, error) {

	// verify body

	err := a.Deps.Repository.ToggleWeddingLinkOff(weddingId)
	if err != nil {
		return values.Failed, "failed to toggle link", err
	}

	return values.Success, "toggled link successfully", nil
}

func (a *API) DoAddMember(member model.Member) (string, string, error) {
	// verify data

	// check if member has been added
	exist, err := a.Deps.Repository.MemberExist(member)
	if err != nil {
		return values.Failed, "failed to find member", err
	}
	if exist {
		return values.Conflict, "this member already exist on your wedding list", errors.New("duplicate resource")
	}
	err = a.Deps.Repository.RunInTx(context.TODO(), func() error {
		//generate code
		member.MemberCode = util.RandomString(6, values.Alphabet)
		// add sending of email here

		return nil
	})
	if err != nil {
		return values.Failed, "failed to send message to member", err
	}
	return values.Success, "member added successfully", nil
}

func (a *API) DoPersistWedding(req model.NewWeddingReq) (*model.Wedding, string, string, error) {
	// verify data

	// generate weeding key
	req.Link = util.GenerateSpecialKey(req.WeddingId)
	req.CreatedAt = time.Now()

	//persist wedding
	err := a.Deps.Repository.PersistWedding(req)
	if err != nil {
		return nil, values.Failed, "failed to add wedding", err
	}

	// wedding link to watch wedding // wedding code to tch wedding // member link to add // remove member

	// send email to the couple of their member link to add members or remove members // todo verification on link maybe wedding ID

	return nil, values.Success, "member added successfully", nil
}
