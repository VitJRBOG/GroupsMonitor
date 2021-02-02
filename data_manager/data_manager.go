package data_manager

import (
	"errors"
	"github.com/VitJRBOG/GroupsMonitor/db"
	"strconv"
	"strings"
)

type AccessToken db.AccessToken

func (a *AccessToken) SetName(name string) error {
	err := stringLengthCheck(name)
	if err != nil {
		return err
	}
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

func (a *AccessToken) SetValue(value string) error {
	err := stringLengthCheck(value)
	if err != nil {
		return err
	}
	a.Value = value
	return nil
}

func (a *AccessToken) SelectFromDB(name string) error {
	var accessToken db.AccessToken
	accessToken.SelectByName(name)

	a.ID = accessToken.ID
	a.Name = accessToken.Name
	a.Value = accessToken.Value

	err := a.checkExistence()
	if err != nil {
		return err
	}
	return nil
}

func (a *AccessToken) checkExistence() error {
	if a.ID == 0 && a.Name == "" && a.Value == "" {
		err := errors.New("no such access token found")
		return err
	}
	return nil
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
	err := stringLengthCheck(name)
	if err != nil {
		return err
	}
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
	err := stringLengthCheck(strVkID)
	if err != nil {
		return err
	}
	err = o.checkZeroInTheBeginning(strVkID)
	if err != nil {
		return err
	}
	vkID, err := o.integerCheck(strVkID)
	if err != nil {
		return err
	} else {
		o.VkID = vkID
	}
	return nil
}

func (o *Operator) checkZeroInTheBeginning(strVkID string) error {
	if strings.Split(strVkID, "")[0] == "0" {
		err := errors.New("vk id starts with zero")
		return err
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

func (o *Operator) SelectFromDB(name string) error {
	var operator db.Operator
	operator.SelectByName(name)

	o.ID = operator.ID
	o.Name = operator.Name
	o.VkID = operator.VkID

	err := o.checkExistence()
	if err != nil {
		return err
	}
	return nil
}

func (o *Operator) checkExistence() error {
	if o.ID == 0 && o.Name == "" && o.VkID == 0 {
		err := errors.New("no such operator found")
		return err
	}
	return nil
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

type Ward db.Ward

func (w *Ward) SetName(name string) error {
	err := stringLengthCheck(name)
	if err != nil {
		return err
	}
	nameIsUnique := w.uniquenessCheck(name)
	if nameIsUnique {
		w.Name = name
	} else {
		err := errors.New("ward with this name already exists")
		return err
	}
	return nil
}

func (w *Ward) uniquenessCheck(name string) bool {
	wards := db.SelectWards()
	nameIsUnique := true
	for _, ward := range wards {
		if w.Name != name {
			if ward.Name == name {
				nameIsUnique = false
				break
			}
		}
	}
	return nameIsUnique
}

func (w *Ward) SetVkID(strVkID string) error {
	err := stringLengthCheck(strVkID)
	if err != nil {
		return err
	}
	err = w.checkZeroInTheBeginning(strVkID)
	if err != nil {
		return err
	}
	vkID, err := w.integerCheck(strVkID)
	if err != nil {
		return err
	}
	err = w.checkNegativeNumber(vkID)
	if err != nil {
		return err
	}
	w.VkID = vkID
	return nil
}

func (w *Ward) checkZeroInTheBeginning(strVkID string) error {
	if strings.Split(strVkID, "")[0] == "0" {
		err := errors.New("vk id starts with zero")
		return err
	}
	return nil
}

func (w *Ward) integerCheck(strVkID string) (int, error) {
	vkID, err := strconv.Atoi(strVkID)
	if err != nil {
		return 0, err
	}
	return vkID, nil
}

func (w *Ward) checkNegativeNumber(vkID int) error {
	if vkID > 0 {
		err := errors.New("vk group id positive")
		return err
	}
	return nil
}

func (w *Ward) SetAccessToken(accessTokenName string) error {
	err := stringLengthCheck(accessTokenName)
	if err != nil {
		return err
	}

	var accessToken db.AccessToken
	accessToken.SelectByName(accessTokenName)

	var a AccessToken
	a.ID = accessToken.ID
	a.Name = accessToken.Name
	a.Value = accessToken.Value
	err = a.checkExistence()
	if err != nil {
		return err
	}

	w.GetAccessTokenID = a.ID

	return nil
}

func (w *Ward) SelectFromDB(name string) error {
	var ward db.Ward
	ward.SelectByName(name)

	w.ID = ward.ID
	w.Name = ward.Name
	w.VkID = ward.VkID
	w.IsOwned = ward.IsOwned
	w.LastTS = ward.LastTS
	w.GetAccessTokenID = ward.GetAccessTokenID

	err := w.checkExistence()
	if err != nil {
		return err
	}
	return nil
}

func (w *Ward) checkExistence() error {
	if w.ID == 0 && w.Name == "" && w.VkID == 0 && w.IsOwned == 0 && w.LastTS == 0 && w.GetAccessTokenID == 0 {
		err := errors.New("no such ward found")
		return err
	}
	return nil
}

func (w *Ward) SaveToDB() {
	var ward db.Ward
	ward.Name = w.Name
	ward.VkID = w.VkID
	ward.IsOwned = 1
	ward.LastTS = 1
	ward.GetAccessTokenID = w.GetAccessTokenID

	ward.InsertIntoDB()
}

func (w *Ward) UpdateInDB() {
	var ward db.Ward
	ward.ID = w.ID
	ward.Name = w.Name
	ward.VkID = w.VkID
	ward.IsOwned = w.IsOwned
	ward.LastTS = w.LastTS
	ward.GetAccessTokenID = w.GetAccessTokenID

	ward.UpdateDB()
}

type Observer db.Observer

func (o *Observer) SetName(name string) {
	o.Name = name
}

func (o *Observer) SetWardID(wardID int) {
	o.WardID = wardID
}

func (o *Observer) SetOperator(operatorName string) error {
	err := stringLengthCheck(operatorName)
	if err != nil {
		return err
	}

	var operator db.Operator
	operator.SelectByName(operatorName)

	var op Operator
	op.ID = operator.ID
	op.Name = operator.Name
	op.VkID = operator.VkID
	err = op.checkExistence()
	if err != nil {
		return err
	}

	o.OperatorID = op.ID

	return nil
}

func (o *Observer) SetAccessToken(accessTokenName string) error {
	err := stringLengthCheck(accessTokenName)
	if err != nil {
		return err
	}

	var accessToken db.AccessToken
	accessToken.SelectByName(accessTokenName)

	var a AccessToken
	a.ID = accessToken.ID
	a.Name = accessToken.Name
	a.Value = accessToken.Value
	err = a.checkExistence()
	if err != nil {
		return err
	}

	o.SendAccessTokenID = a.ID

	return nil
}

func (o *Observer) SetAdditionalParams(wallPostType string) {
	o.AdditionalParams.WallPost.PostType = wallPostType
}

func (o *Observer) SaveToDB() {
	var observer db.Observer
	observer.Name = o.Name
	observer.WardID = o.WardID
	observer.OperatorID = o.OperatorID
	observer.SendAccessTokenID = o.SendAccessTokenID
	observer.AdditionalParams = o.AdditionalParams

	observer.InsertIntoDB()
}

func stringLengthCheck(s string) error {
	if len(s) == 0 {
		err := errors.New("string length is zero")
		return err
	}
	return nil
}
