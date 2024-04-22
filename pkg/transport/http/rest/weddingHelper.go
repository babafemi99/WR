package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/babafemi99/WR/internal/util"
	"github.com/babafemi99/WR/internal/values"
	"github.com/babafemi99/WR/pkg/model"
	"github.com/jackc/pgx/v5"
	"log"
	"time"
)

func (a *API) JoinWedding(ctx context.Context, id, code string) (*model.WeddingIdRes, string, string, error) {
	log.Println(id, code)
	// check if code and wedding tally
	exist, err := a.Deps.Repository.MemberCodeExist(ctx, code, id)
	if err != nil {
		return nil, values.Failed, "failed to verify code or wedding ID", err
	}
	if !exist {
		return nil, values.NotAuthorised, "you don't have access to this resource", errors.New("authorization error")
	}

	wedding, err := a.Deps.Repository.GetWeddingById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, values.NotFound, "no wedding with this Id", errors.New("resource not found")
		}
		return nil, values.Error, "failed to load wedding", err
	}
	log.Println(wedding, "ww")

	if wedding.Status == "pending" {
		return nil, values.NotAuthorised, "wedding is not yet live", nil
	}
	if wedding.Status == "done" {
		return nil, values.NotAuthorised, "wedding event is over", nil
	}

	return wedding, values.Success, "verification successful", nil
}

func (a *API) ToggleWedding(ctx context.Context, req model.ToggleWeddingReq) (string, string, error) {
	// verify body

	// getId from context
	executor, ok := ctx.Value(values.Executor).(model.Executor)
	if !ok {
		return values.Failed, "system error", errors.New("failed to get executor")
	}

	req.TogglerId = executor.Id
	req.ModifiedAt = time.Now()

	// persist to database
	err := a.Deps.Repository.ToggleWeddingLink(ctx, req)
	if err != nil {
		return values.Failed, "failed to toggle link", err
	}

	return values.Success, "toggled link successfully", nil
}

func (a *API) ToggleWeddingOff(ctx context.Context, weddingId string) (string, string, error) {

	err := a.Deps.Repository.ToggleWeddingLinkOff(ctx, weddingId)
	if err != nil {
		return values.Failed, "failed to toggle link", err
	}

	return values.Success, "toggled link successfully", nil
}

func (a *API) DoAddMember(ctx context.Context, member model.Member) (string, string, error) {
	// verify data

	// check if member has been added
	exist, err := a.Deps.Repository.MemberExist(ctx, member)
	if err != nil {
		return values.Failed, "failed to find member", err
	}
	if exist {
		return values.Conflict, "this member already exist on your wedding list", errors.New("duplicate resource")
	}

	var status, message string
	err = a.Deps.Repository.RunInTx(context.TODO(), func() error {
		//generate code
		member.MemberCode = util.RandomString(5, values.Alphabet)

		// add to DB
		err = a.Deps.Repository.AddMembers(ctx, member)
		if err != nil {
			status = values.Error
			message = "failed to add members"
			return err
		}
		// add sending of email here
		data := struct {
			WeddingCode string
			Passcode    string
		}{
			WeddingCode: member.WeddingId,
			Passcode:    member.MemberCode,
		}
		patterns := []string{"wedding_invite.tmpl"}
		err = a.Deps.IMailer.SendEmail(member.MemberEmail, nil, data, patterns...)
		if err != nil {
			message = "failed to send email"
			status = values.Failed
			return err
		}

		return nil
	})
	if err != nil {
		return status, message, err
	}
	return values.Success, "member added successfully", nil
}

func (a *API) DoPersistWedding(ctx context.Context, req model.NewWeddingReq) (*model.NewWeddingRes, string, string, error) {
	// verify data

	var status, message string
	var res model.NewWeddingRes
	err := a.Deps.Repository.RunInTx(ctx, func() error {
		req.WeddingId = util.GenerateSpecialKey(req.CoupleId)
		req.CreatedAt = time.Now()
		req.GuestLink = fmt.Sprintf("https://wedding-registy.com/%s/%s/members-list", util.EncodeCID(req.CoupleId), req.WeddingId)
		req.Link = fmt.Sprintf("https://wedding-registry.com/%s", req.WeddingId)

		//persist wedding todo handle duplicate entry
		err := a.Deps.Repository.PersistWedding(ctx, req)
		if err != nil {
			status = values.Failed
			message = "failed to add wedding"
			return err
		}
		res.WeddingId = req.WeddingId
		res.Link = req.Link
		res.GuestLink = req.GuestLink

		// send email to the couple with these credentials
		return nil

	})
	if err != nil {
		return nil, status, message, err
	}

	return &res, values.Success, "wedding added successfully", nil
}
