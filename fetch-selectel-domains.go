package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/viper"
)

func init() {
	// проверка конфигурационного файла
	viper.SetConfigFile(filepath.Base(os.Args[0]) + ".yaml")
	if err := viper.ReadInConfig(); err != nil {
		exitWithMsg(err.Error())
	}
	// check requires parameters
	reqPars := [...]string{
		"APIURL",
		"APItoken",
		"TargetPath",
	}
	for i := range reqPars {
		if !viper.IsSet(reqPars[i]) {
			exitWithMsg(fmt.Sprintf("'%s' not found in config", reqPars[i]))
		}
		log.Printf("INFO: Use '%s' = %s", reqPars[i], viper.GetString(reqPars[i]))
	}
	if _, err := os.Stat(viper.GetString("TargetPath")); err != nil {
		exitWithMsg(fmt.Sprintf("Path '%s' not found",
			viper.GetString("TargetPath")))
	}
}

func main() {
	log.Println("INFO: Start Fetch")
	listZones, err := getZonesList()
	if err != nil {
		exitWithMsg(err.Error())
	}
	for _, z := range listZones {
		fmt.Println(strconv.Itoa(z.ID) + "\t" + z.Name)

	}
	log.Println("INFO: Stop Successfull")
}

func getZonesList() (z []selectelZone, err error) {
	client := &http.Client{}
	r, _ := http.NewRequest("GET", viper.GetString("APIURL"), nil)
	r.Header.Add("X-Token", viper.GetString("APItoken"))
	resp, err := client.Do(r)
	if err != nil {
		return z, err
	}
	defer resp.Body.Close()
	JSONdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return z, err
	}
	if resp.StatusCode != 200 {
		return z, fmt.Errorf("Error in response: %s\n", JSONdata)
	}
	//fmt.Println(string(JSONdata))
	err = json.Unmarshal(JSONdata, &z)
	if err != nil {
		return z, err
	}
	return z, nil
}

func exitWithMsg(msg string) {
	log.Printf("Exit with fatal error: %s\n", msg)
	os.Exit(1)
}

type selectelZone struct {
	ID          int    // Идентификатор домена,
	Name        string // Имя домена,
	Tags        []int  // Список тэгов домена,
	Create_date int64  // (UNIX Timestamp): Дата создания,
	Change_date int64  //(UNIX Timestamp): Дата изменения,
	User_id     int    // Идентификатор пользователя,
}
