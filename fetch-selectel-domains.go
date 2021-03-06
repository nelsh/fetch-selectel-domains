package main

import (
	"encoding/json"
	"flag"
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

var (
	verbose bool
)

func init() {
	flag.BoolVar(&verbose, "v", false, "Verbose logging")
	flag.Parse()
	// check config file
	viper.SetConfigFile(filepath.Join("/etc/", filepath.Base(os.Args[0])+".yaml"))
	if err := viper.ReadInConfig(); err != nil {
		exitWithMsg(err.Error())
	}
	// check requires parameters
	reqPars := [...]string{
		"APIURL",
		"APItoken",
		"TargetPath",
	}
	if verbose {
		for i := range reqPars {
			if !viper.IsSet(reqPars[i]) {
				exitWithMsg(fmt.Sprintf("'%s' not found in config", reqPars[i]))
			}
			if reqPars[i] == "APItoken" {
				log.Printf("INFO: Use '%s' = ****", reqPars[i])
			} else {
				log.Printf("INFO: Use '%s' = %s", reqPars[i], viper.GetString(reqPars[i]))
			}
		}
	}
	viper.SetDefault("DefaultTTL", 86400)
	if verbose {
		log.Printf("INFO: Use 'DefaultTTL' = %d", viper.GetInt("DefaultTTL"))
	}
	if _, err := os.Stat(viper.GetString("TargetPath")); err != nil {
		exitWithMsg(fmt.Sprintf("Path '%s' not found", viper.GetString("TargetPath")))
	}
}

func main() {
	if verbose {
		log.Println("INFO: Start Fetch")
	}
	listZones, err := getZonesList()
	if err != nil {
		exitWithMsg(err.Error())
	}
	if verbose {
		log.Printf("INFO: Found %d zones", len(listZones))
	}
	totalErr := 0
	for _, z := range listZones {
		if verbose {
			log.Printf("INFO: Fetch '%s'", z.Name)
		}
		zone, err := z.ToString()
		if err != nil {
			log.Printf("WARN: %s", err)
			totalErr++
			continue
		}
		zoneFile := filepath.Join(viper.GetString("TargetPath"), z.Name+".dns")
		if err := ioutil.WriteFile(zoneFile, []byte(zone), 0666); err != nil {
			log.Printf("WARN: write file zone '%s' error: %s", z.Name, err)
			totalErr++
		}
	}
	if verbose {
		log.Println("INFO: Stop Successfull")
	}
	fmt.Printf("Fetched %d from %d zones, errors = %d\r\n",
		len(listZones)-totalErr, len(listZones), totalErr)
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

func (z selectelZone) ToString() (string, error) {
	listRecords, err := getRecordsList(z.ID)
	if err != nil {
		return "", err
	}
	groups := make(map[string]string)
	head := fmt.Sprintf("$ORIGIN %s.\r\n$TTL %d\r\n\r\n",
		z.Name, viper.GetInt("DefaultTTL"))
	for _, r := range listRecords {
		switch r.Type {
		case "A":
			groups["A"] += (r.ToString() + "\r\n")
		case "AAAA":
			groups["AAAA"] += (r.ToString() + "\r\n")
		case "CNAME":
			groups["CNAME"] += (r.ToString() + "\r\n")
		case "MX":
			groups["MX"] += (r.ToString() + "\r\n")
		case "NS":
			groups["NS"] += (r.ToString() + "\r\n")
		case "SOA":
			groups["SOA"] += (r.ToString() + "\r\n")
		case "SRV":
			groups["SRV"] += (r.ToString() + "\r\n")
		case "TXT":
			groups["TXT"] += (r.ToString() + "\r\n")
		default:
			log.Printf("ERROR: unknown record '%+v'", r)
		}
	}
	zone := head +
		groups["SOA"] +
		groups["NS"] +
		groups["MX"] +
		groups["A"] +
		groups["AAAA"] +
		groups["CNAME"] +
		groups["TXT"] +
		groups["SRV"]
	return zone, nil
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
	ttlString := ""
	if viper.GetInt("DefaultTTL") != r.TTL {
		ttlString = strconv.Itoa(r.TTL)
	}
	switch r.Type {
	case "A":
		return fmt.Sprintf("%s.\t\t%s\tIN\tA\t%s",
			r.Name, ttlString, r.Content)
	case "AAAA":
		return fmt.Sprintf("%s.\t\t%s\tIN\tAAAA\t%s",
			r.Name, ttlString, r.Content)
	case "CNAME":
		return fmt.Sprintf("%s.\t\t%s\tIN\tCNAME\t%s.",
			r.Name, ttlString, r.Content)
	case "MX":
		return fmt.Sprintf("%s.\t\t%s\tIN\tMX\t%d\t%s.",
			r.Name, ttlString, r.Priority, r.Content)
	case "NS":
		return fmt.Sprintf("%s.\t\t%s\tIN\tNS\t%s.",
			r.Name, ttlString, r.Content)
	case "SOA":
		soa := strings.SplitAfterN(r.Content, " ", 3)
		return fmt.Sprintf("%s.\t\t%s\tIN\tSOA\t%s %s ( %s )",
			r.Name, ttlString, soa[0], soa[1], soa[2])
	case "SRV":
		return fmt.Sprintf("%s.\t\t%s\tIN\tSRV\t%d\t%d\t%d\t%s.",
			r.Name, ttlString, r.Priority, r.Weight, r.Port, r.Target)
	case "TXT":
		return fmt.Sprintf("%s.\t\t%s\tIN\tTXT\t\"%s\"",
			r.Name, ttlString, r.Content)
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
