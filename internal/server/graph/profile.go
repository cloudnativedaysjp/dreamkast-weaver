package graph

import (
	"errors"

	"dreamkast-weaver/internal/application"
	"dreamkast-weaver/internal/domain/value"
	"dreamkast-weaver/internal/server/graph/model"
)

func newProfile(confName model.ConfName, profileID int) (application.Profile, error) {
	var e, err error
	profile := application.Profile{}
	profile.ConfName, e = value.NewConfName(value.ConferenceKind(confName))
	err = errors.Join(err, e)
	profile.ID, e = value.NewProfileID(int32(profileID))
	err = errors.Join(err, e)

	return profile, err
}
