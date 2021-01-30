package data_manager

import (
	"errors"
	"github.com/VitJRBOG/GroupsMonitor/db"
	"strconv"
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

type Operator db.Operator

func (o *Operator) SetName(name string) error {
	nameIsUnique := o.uniquenessCheck(name)
	if nameIsUnique {
		o.Name = name
	} else {
		err := errors.New("operator with this name already exists")
		return err
	}
	return nil
}

func (o *Operator) uniquenessCheck(name string) bool {
	operators := db.SelectOperators()
	nameIsUnique := true
	for _, operator := range operators {
		if o.Name != name {
			if operator.Name == name {
				nameIsUnique = false
				break
			}
		}
	}
	return nameIsUnique
}

func (o *Operator) SetVkID(strVkID string) error {
	vkID, err := o.integerCheck(strVkID)
	if err != nil {
		return err
	} else {
		o.VkID = vkID
	}
	return nil
}

func (o *Operator) integerCheck(strVkID string) (int, error) {
	vkID, err := strconv.Atoi(strVkID)
	if err != nil {
		return 0, err
	}
	return vkID, nil
}

func (o *Operator) SelectFromDB(name string) {
	var operator db.Operator
	operator.SelectByName(name)

	o.ID = operator.ID
	o.Name = operator.Name
	o.VkID = operator.VkID
}

func (o *Operator) SaveToDB() {
	var operator db.Operator
	operator.Name = o.Name
	operator.VkID = o.VkID

	operator.InsertIntoDB()
}

func (o *Operator) UpdateIdDB() {
	var operator db.Operator
	operator.ID = o.ID
	operator.Name = o.Name
	operator.VkID = o.VkID

	operator.UpdateInDB()
}
