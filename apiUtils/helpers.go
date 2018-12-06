package apiUtils

import (
	"api/apiSructs"
	"encoding/json"
	"errors"
	"fmt"
)

var CheckCredentials = func(login, password string) (structs.Credential, error){
	//byt := []byte(`{"role":"admin","name":"Вася Пупкинович"}`)
	byt := []byte(`{"role":"admin","name":"Вася Пупкинович", "email":"nvs@nvs.ru"}`)
	res := structs.Credential{}
	json.Unmarshal(byt, &res)
	var err error
	err = nil

	if len(res.Email) <= 0 {
		err = errors.New(fmt.Sprintf("BingoBoom auth server doesn't send required field 'email' for user %d. Sorry. Can't authenticate.", login))
	}

	//err = "login and password are incorrect";

	return res, err
}
