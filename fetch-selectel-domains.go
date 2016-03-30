package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

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
	}
	for i := range reqPars {
		if !viper.IsSet(reqPars[i]) {
			exitWithMsg(fmt.Sprintf("'%s' not found in config", reqPars[i]))
		}
		log.Printf("INFO: Use '%s' = %s", reqPars[i], viper.GetString(reqPars[i]))
	}
}

func main() {
	log.Println("INFO: Start Fetch")
	if err := getZonesList(); err != nil {
		exitWithMsg(err.Error())
	}
	log.Println("INFO: Stop Successfull")
}

func getZonesList() error {
	client := &http.Client{}
	r, _ := http.NewRequest("GET", viper.GetString("APIURL"), nil)
	r.Header.Add("X-Token", viper.GetString("APItoken"))
	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	JSONdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error in response: %s\n", JSONdata)
	}
	fmt.Println(string(JSONdata))
	return nil
}

func exitWithMsg(msg string) {
	log.Printf("Exit with fatal error: %s\n", msg)
	os.Exit(1)
}
