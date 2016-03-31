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
	"strings"

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
		fmt.Println("\r\n" + strconv.Itoa(z.ID) + "\t" + z.Name + "\r\n" + strings.Repeat("-", 36))
		listRecords, err := getRecordsList(z.ID)
		if err != nil {
			log.Printf("WARN: %s", err)
		}
		for _, r := range listRecords {
			fmt.Println(r.ToString())
		}
	}
	log.Println("INFO: Stop Successfull")
}

func getZonesList() (z []selectelZone, err error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", viper.GetString("APIURL"), nil)
	req.Header.Add("X-Token", viper.GetString("APItoken"))
	resp, err := client.Do(req)
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
	err = json.Unmarshal(JSONdata, &z)
	if err != nil {
		return z, err
	}
	return z, nil
}

func getRecordsList(ID int) (r []selectelRecord, err error) {
	apiURL := fmt.Sprintf("%s/%d/records/", viper.GetString("APIURL"), ID)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Add("X-Token", viper.GetString("APItoken"))
	resp, err := client.Do(req)
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()
	JSONdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}
	if resp.StatusCode != 200 {
		return r, fmt.Errorf("Error in response: %s\n", JSONdata)
	}
	//fmt.Println(string(JSONdata))
	err = json.Unmarshal(JSONdata, &r)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (r selectelRecord) ToString() string {
	switch r.Type {
	case "A":
		return fmt.Sprintf("%s.\t\t%d\tIN\tA\t%s",
			r.Name, r.TTL, r.Content)
	case "CNAME":
		return fmt.Sprintf("%s.\t\t%d\tIN\tCNAME\t%s.",
			r.Name, r.TTL, r.Content)
	case "MX":
		return fmt.Sprintf("%s.\t\t%d\tIN\tMX\t%d\t%s.",
			r.Name, r.TTL, r.Priority, r.Content)
	case "NS":
		return fmt.Sprintf("%s.\t\tIN\tNS\t%s.",
			r.Name, r.Content)
	case "SOA":
		return fmt.Sprintf("%s.\t\tIN\tSOA\t%s",
			r.Name, r.Content)
	case "SRV":
		return fmt.Sprintf("%s.\t\t%d\tIN\tSRV\t%d\t%d\t%d\t%s.",
			r.Name, r.TTL, r.Priority, r.Weight, r.Port, r.Target)
	case "TXT":
		return fmt.Sprintf("%s.\t\t%d\tIN\tTXT\t\"%s\"",
			r.Name, r.TTL, r.Content)
	/*	case "SRV":
			return fmt.Sprintf("%-24s %6d %-6s %d %d %d %s",
				r.IdnName, r.TTL, r.Type, r.SRV.Priority, r.SRV.Weight, r.SRV.Port, r.SRV.Target.IdnName)
		case "TXT":
			return fmt.Sprintf("%-24s %6d %-6s \"%s\"",
				r.IdnName, r.TTL, r.Type, r.TXT.String)*/
	default:
		log.Printf("ERROR: unknown record '%+v'", r)
		return ""
	}
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

type selectelRecord struct {
	ID          int      // Идентификатор записи,
	Name        string   // Имя записи,
	Type        string   // Тип записи (SOA, NS, A/AAAA, CNAME, SRV, MX, TXT, SPF),
	TTL         int      // Время жизни,
	Email       string   // e-mail администратора домена (только у SOA),
	Content     string   // Содержимое записи (нет у SRV),
	Weight      int      // Относительный вес для записей с одинаковым приоритетом (только у SRV),
	Port        int      // Порт сервиса (только у SRV),
	Target      string   // Имя хоста сервиса (только у SRV),
	Geo_records []string // Гео-записи,
	Priority    int      // Приоритет записи (только у MX и SRV),
	Create_date int64    // Create date,
	Change_date int64    // Change date,
}
