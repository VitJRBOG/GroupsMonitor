package data_manager

import (
	"errors"
	"github.com/VitJRBOG/GroupsMonitor/db"
)

type AccessToken db.AccessToken

func (a *AccessToken) SetName(name string) error {
	nameIsUnique := a.uniquenessCheck(name)
	if nameIsUnique {
		a.Name = name
	} else {
		err := errors.New("access token with this name already exists")
		return err
	}
	return nil
}

func (a *AccessToken) uniquenessCheck(name string) bool {
	accessTokens := db.SelectAccessTokens()
	nameIsUnique := true
	for _, accessToken := range accessTokens {
		if a.Name != name {
			if accessToken.Name == name {
				nameIsUnique = false
				break
			}
		}
	}
	return nameIsUnique
}

func (a *AccessToken) SetValue(value string) {
	a.Value = value
}

func (a *AccessToken) SelectFromDB(name string) {
	var accessToken db.AccessToken
	accessToken.SelectByName(name)

	a.ID = accessToken.ID
	a.Name = accessToken.Name
	a.Value = accessToken.Value
}

func (a *AccessToken) SaveToDB() {
	var accessToken db.AccessToken
	accessToken.ID = a.ID
	accessToken.Name = a.Name
	accessToken.Value = a.Value

	accessToken.InsertIntoDB()
}

func (a *AccessToken) UpdateInDB() {
	var accessToken db.AccessToken
	accessToken.ID = a.ID
	accessToken.Name = a.Name
	accessToken.Value = a.Value

	accessToken.UpdateInDB()
}
