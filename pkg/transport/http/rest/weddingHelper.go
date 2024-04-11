package rest

import (
	"github.com/babafemi99/WR/internal/util"
	"github.com/babafemi99/WR/internal/values"
	"github.com/babafemi99/WR/pkg/model"
)

func (a *API) VerifyWeddingId(id string) (any, string, string, error) {

	// check if wedding is today

	// check if the wedding link has been toggled

	// return error

	return nil, values.Success, "verification successful", nil
}

func (a *API) ToggleWedding(req model.ToggleWeddingReq) (string, string, error) {

	// verify body

	// persist to database

	// return values

	return values.Success, "toggled link successfully", nil
}

func (a *API) ToggleWeddingOff(weddingId string) (string, string, error) {

	// verify body

	// check if this wedding is toggled

	// if true set as not offline

	// return values

	return values.Success, "toggled link successfully", nil
}

func (a *API) VerifyWeddingIdForToday(id string) (any, string, string, error) {

	// check if wedding is today

	// check if the wedding link not been toggled

	// return error

	return nil, values.Success, "verification successful", nil
}

func (a *API) DoAddMember(member model.Member) (string, string, error) {
	// verify data
	// check if member has been added
	//generate code
	member.MemberCode = util.RandomString(6, values.Alphabet)
	// add member
	// send email to member added

	return values.Success, "member added successfully", nil
}

func (a *API) DoPersistWedding() (string, string, error) {
	return values.Success, "member added successfully", nil
}
