package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type fullStack struct {
	ID             int                 `json:"id"`
	Image          string              `json:"image"`
	Name           string              `json:"name"`
	Members        []string            `json:"members"`
	CreationDate   int                 `json:"creationDate"`
	FirstAlbum     string              `json:"firstAlbum"`
	Locations      []string            `json:"locations"`
	ConcertDates   []string            `json:"concertDates"`
	DatesLocations map[string][]string `json:"datesLocations"`
	WikiLink       []string
}

type MyArtist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type MyLocation struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type MyRelation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type MyDate struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type MyDates struct {
	Index []MyDate `json:"index"`
}

type MyLocations struct {
	Index []MyLocation `json:"index"`
}

type MyRelations struct {
	Index []MyRelation `json:"index"`
}

type MembersLinks struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

// Struct the data for this particular project
// variables
var (
	AllArtist []fullStack
	Artists   []MyArtist
	Dates     MyDates
	Locations MyLocations
	Relations MyRelations
	MemLinks  []MembersLinks
)

const rootC = "https://groupietrackers.herokuapp.com/api"

func main() {

	styling := http.FileServer(http.Dir("style"))
	http.Handle("/style/", http.StripPrefix("/style/", styling))

	http.HandleFunc("/", mainPage)
	http.HandleFunc("/concerts", concertPage)
	http.HandleFunc("/tours", tourPage)

	port := ":7000"
	fmt.Println("Server on", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("error", err)
	}
}



func Get_Artist() ([]MyArtist, error) {
	Artists := []MyArtist{}
	resp, err := http.Get(rootC + "/artists")
	if err != nil {
		return Artists, errors.New("error")
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Artists, errors.New("error reading all")
	}
	json.Unmarshal(bytes, &Artists)
	return Artists, nil
}

func GetDatesData() (MyDates, error) {
	Dates := MyDates{}
	resp, err := http.Get(rootC + "/dates")
	if err != nil {
		return Dates, errors.New("error by get")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Dates, errors.New("error by ReadAll")
	}
	json.Unmarshal(bytes, &Dates)
	return Dates, nil
}

func GetLocationsData() (MyLocations, error) {
	Locations := MyLocations{}
	resp, err := http.Get(rootC + "/locations")
	if err != nil {
		return Locations, errors.New("error by get")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Locations, errors.New("error by ReadAll")
	}
	json.Unmarshal(bytes, &Locations)
	return Locations, nil
}

func GetRelationsData() (MyRelations, error) {
	Relations := MyRelations{}
	resp, err := http.Get(rootC + "/relation")
	if err != nil {
		return Relations, errors.New("error by get")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Relations, errors.New("error by ReadAll")
	}
	json.Unmarshal(bytes, &Relations)
	fmt.Println(Relations.Index[0].DatesLocations)
	return Relations, nil
}

func GetData() error {
	if len(AllArtist) != 0 {
		return nil
	}
	Artists, err1 := Get_Artist()
	Locations, err2 := GetLocationsData()
	Dates, err3 := GetDatesData()
	Relations, err4 := GetRelationsData()
	err5 := links()
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		return errors.New("error by get data artists, locations, dates")
	}
	for i := range Artists {

		var tmpl fullStack
		var addMemLinks []string
		for j := 0; j < len(Artists[i].Members); j++ {
			for m := 0; m < len(MemLinks); m++ {
				if MemLinks[m].Name == Artists[i].Members[j] {
					addMemLinks = append(addMemLinks, MemLinks[m].Link)
				}
			}
		}
		tmpl.ID = i + 1
		tmpl.Image = Artists[i].Image
		tmpl.Name = Artists[i].Name
		tmpl.Members = Artists[i].Members
		tmpl.CreationDate = Artists[i].CreationDate
		tmpl.FirstAlbum = Artists[i].FirstAlbum
		tmpl.Locations = Locations.Index[i].Locations
		tmpl.ConcertDates = Dates.Index[i].Dates
		tmpl.DatesLocations = Relations.Index[i].DatesLocations
		tmpl.WikiLink = addMemLinks
		AllArtist = append(AllArtist, tmpl)
	}
	return nil
}

func GetArtistByID(id int) (MyArtist, error) {
	for _, artist := range Artists {
		if artist.ID == id {
			return artist, nil
		}
	}
	return MyArtist{}, errors.New("not found")
}

func GetDateByID(id int) (MyDate, error) {
	for _, date := range Dates.Index {
		if date.ID == id {
			return date, nil
		}
	}
	return MyDate{}, errors.New("not found")
}

func GetLocationByID(id int) (MyLocation, error) {
	for _, location := range Locations.Index {
		if location.ID == id {
			return location, nil
		}
	}
	return MyLocation{}, errors.New("not found")
}

func GetRelationByID(id int) (MyRelation, error) {
	for _, relation := range Relations.Index {
		if relation.ID == id {
			return relation, nil
		}
	}
	return MyRelation{}, errors.New("not found")
}

func GetFullDataById(id int) (fullStack, error) {
	for _, artist := range AllArtist {
		if artist.ID == id {
			return artist, nil
		}
	}
	return fullStack{}, errors.New("not found")
}

var data []fullStack

func mainPage(w http.ResponseWriter, r *http.Request) {
	err := GetData()
	if err != nil {
		errors.New("error by get data")
	}

	main := r.FormValue("main")
	search := r.FormValue("search")
	filterByCreationFrom := r.FormValue("startCD")
	filterByCreationTill := r.FormValue("endCD")
	filterByFA := r.FormValue("startFA")
	filterByFAend := r.FormValue("endFA")

	if !(search == "" && len(data) != 0) {
		data = Search(search)
	}

	if filterByCreationFrom != "" || filterByCreationTill != "" {
		if filterByCreationFrom == "" {
			filterByCreationFrom = "1900"
		}
		if filterByCreationTill == "" {
			filterByCreationTill = "2020"
		}

	}

	if filterByFA != "" || filterByFAend != "" {
		if filterByFA == "" {
			filterByFA = "1900-01-01"
		}
		if filterByFAend == "" {
			filterByFAend = "2020-03-03"
		}

	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		handle500(err, w)
		return
	}

	if main == "Main Page" {
		data = Search("a")
	}

	if err := tmpl.Execute(w, data); err != nil {
		handle500(err, w)
		return
	}
}

func concertPage(w http.ResponseWriter, r *http.Request) {
	listOfIds := r.URL.Query()["id"]
	id, err := strconv.Atoi(listOfIds[0])
	if err != nil {
		handle500(err, w)
		return
	}

	artist, err := GetFullDataById(id)
	if err != nil {
		http.Error(w, "Bad Request: 400", 400)
		return
	}

	tmpl, err := template.ParseFiles("concerts.html")
	if err != nil {
		handle500(err, w)
		return
	}

	if err := tmpl.Execute(w, artist); err != nil {
		handle500(err, w)
		return
	}
}

func tourPage(w http.ResponseWriter, r *http.Request) {
	listOfIds := r.URL.Query()["id"]
	id, err := strconv.Atoi(listOfIds[0])
	if err != nil {
		handle500(err, w)
		return
	}

	artist, err := GetFullDataById(id)
	if err != nil {
		http.Error(w, "Bad Request: 400", 400)
		return
	}

	tmpl, err := template.ParseFiles("tours.html")
	if err != nil {
		handle500(err, w)
		return
	}

	if err := tmpl.Execute(w, artist); err != nil {
		handle500(err, w)
		return
	}
}

func ConverterStructToString(AllArtist []fullStack) ([]string, error) {
	var data []string
	for i := 1; i <= len(AllArtist); i++ {
		artist, err1 := GetArtistByID(i)
		locations, err2 := GetLocationByID(i)
		date, err3 := GetDateByID(i)
		if err1 != nil || err2 != nil || err3 != nil {
			return data, errors.New("error by converter")
		}

		str := artist.Name + " "
		for _, member := range artist.Members {
			str += member + " "
		}
		str += strconv.Itoa(artist.CreationDate) + " "
		str += artist.FirstAlbum + " "
		for _, location := range locations.Locations {
			str += location + " "
		}
		for _, d := range date.Dates {
			str += d + " "
		}
		data = append(data, str)
	}
	return data, nil
}

func links() error {

	csvFile, err := os.Open("members.txt")

	if err != nil {
		fmt.Println(err)
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.LazyQuotes = true

	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var oneRecord MembersLinks
	var allRecords []MembersLinks

	for _, each := range csvData {
		oneRecord.Name = each[0]
		oneRecord.Link = each[1]
		allRecords = append(allRecords, oneRecord)
	}

	jsondata, err := json.Marshal(allRecords) // convert to JSON
	json.Unmarshal(jsondata, &MemLinks)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return nil
}

func Search(search string) []fullStack {
	if search == "" {
		return AllArtist
	}
	art, err := ConverterStructToString(AllArtist)
	if err != nil {
		errors.New("error by converter")
	}
	var search_artist []fullStack

	for i, artist := range art {
		lower_band := strings.ToLower(artist)
		for i_name, l_name := range []byte(lower_band) {
			lower_search := strings.ToLower(search)
			if lower_search[0] == l_name {
				length_name := 0
				indx := i_name
				for _, l := range []byte(lower_search) {
					if l == lower_band[indx] {
						if indx+1 == len(lower_band) {
							break
						}
						indx++
						length_name++
					} else {
						break
					}
				}
				if len(search) == length_name {
					band, err := GetFullDataById(i + 1)
					if err != nil {
						fmt.Println(err)
					}
					search_artist = append(search_artist, band)
					break
				}
			}
		}

	}
	return search_artist
}

func handle500(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["500"] = "Internal Server Error"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}
