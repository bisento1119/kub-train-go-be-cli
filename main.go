// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Person struct {
	Id           int    `json:"id"`
	Forename     string `json:"forename"`
	Lastname     string `json:"lastname"`
	ProfessionId int    `json:"professionId"`
}

type Profession struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type Professionals struct {
	Id         int        `json:"id"`
	Forename   string     `json:"forename"`
	Lastname   string     `json:"lastname"`
	Profession Profession `json:"profession,omitempty"`
}

type Config struct {
	GoBeAHost string
	GoBeAPort int
	GoBeAPath string
	GoBeBHost string
	GoBeBPort int
	GoBeBPath string
}

var myConfig Config

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func writeHttpStatus(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 - Internal Server Error"))
}

func returnAllProfessionals(w http.ResponseWriter, r *http.Request) {
	persons, err := getAllPersons()
	if err != nil {
		writeHttpStatus(w, err)
		return
	}

	professions, err2 := getAllProfessions()
	if err2 != nil {
		writeHttpStatus(w, err)
		return
	}

	professionsMap := make(map[int]Profession)
	for _, prof := range professions {
		professionsMap[prof.Id] = prof
	}

	retList := make([]Professionals, len(persons))

	for i, p := range persons {
		retList[i].Id = p.Id
		retList[i].Lastname = p.Lastname
		retList[i].Forename = p.Forename
		retList[i].Profession = professionsMap[p.ProfessionId]
	}
	fmt.Println("Endpoint Hit: returnAllProfessionals")
	json.NewEncoder(w).Encode(retList)
}

func getAllPersons() ([]Person, error) {
	var url = fmt.Sprintf("http://%v:%d/%v", myConfig.GoBeAHost, myConfig.GoBeAPort, myConfig.GoBeAPath)
	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		log.Error(err)
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var responseObject []Person
	json.Unmarshal(responseData, &responseObject)

	return responseObject, nil
}

func getAllProfessions() ([]Profession, error) {
	var url = fmt.Sprintf("http://%v:%d/%v", myConfig.GoBeBHost, myConfig.GoBeBPort, myConfig.GoBeBPath)
	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		log.Error(err)
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var responseObject []Profession
	json.Unmarshal(responseData, &responseObject)

	return responseObject, nil
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/professionals", returnAllProfessionals)
	log.Fatal(http.ListenAndServe(":4882", myRouter))
}

func readConfig() {
	myConfig = Config{
		GoBeAHost: "kub-train-go-be-a",
		GoBeAPort: 4880,
		GoBeAPath: "persons",
		GoBeBHost: "kub-train-go-be-b",
		GoBeBPort: 4881,
		GoBeBPath: "professions",
	}

	viper.SetConfigName("localConf")
	viper.AddConfigPath("./environments")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	myConfig.GoBeAHost = viper.GetString("be-go-a.host")
	myConfig.GoBeBHost = viper.GetString("be-go-b.host")
}

func main() {
	fmt.Println("Starting: kub-train-go-be-cli Endpoint")

	readConfig()

	handleRequests()
}
