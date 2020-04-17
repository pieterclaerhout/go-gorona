package client

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/table"
	"github.com/pieterclaerhout/go-gorona/gorona"
)

const url = "https://corona.lmao.ninja/v2/countries/"

// GetCountry : Get by country
func GetCountry(country string) {

	today := getCountry(country, false)
	yesterday := getCountry(country, true)

	caseStates := gorona.CaseStates{today, yesterday}

	printCaseStates(caseStates, true)

}

func getCountry(country string, yesterday bool) gorona.CaseState {

	caseState := gorona.CaseState{}

	qs := ""
	if yesterday {
		qs = "?yesterday=1"
	}

	req, err := http.Get(url + country + qs)
	if err != nil {
		log.Fatalln(err)
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(body, &caseState)
	if err != nil {
		log.Fatalln(err)
	}

	caseState.Date = time.Now()
	if yesterday {
		caseState.Date = time.Now().Add(-24 * time.Hour)
	}

	return caseState

}

// GetCountries : GetCountries()
func GetCountries() {

	caseStates := gorona.CaseStates{}

	req, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(body, &caseStates)
	if err != nil {
		log.Fatalln(err)
	}

	printCaseStates(caseStates, false)

}

func printCaseState(caseState gorona.CaseState) {
	printCaseStates(gorona.CaseStates{caseState}, true)
}

func printCaseStates(caseStates gorona.CaseStates, printDates bool) {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)

	firstColumn := "Country"
	if printDates {
		firstColumn = "Date"
	}

	t.AppendHeader(
		table.Row{
			firstColumn, "Cases", "Today Cases", "Death", "Today Deaths", "Recovered", "Active", "Critical", "Cases Per Million", "Deaths Per Million",
		},
	)

	for _, caseState := range caseStates {

		firstValue := caseState.Country
		if printDates {
			firstValue = caseState.Date.Format("2006-01-02")
		}

		t.AppendRow(table.Row{
			firstValue,
			caseState.Cases,
			caseState.TodayCases,
			caseState.Deaths,
			caseState.TodayDeaths,
			caseState.Recovered,
			caseState.Active,
			caseState.Critical,
			caseState.CasesPerOneMillion,
			caseState.DeathsPerOneMillion,
		})

	}

	if !printDates {
		t.AppendFooter(
			table.Row{
				"Totals",
				caseStates.Cases(),
				caseStates.TodayCases(),
				caseStates.Deaths(),
				caseStates.TodayDeaths(),
				caseStates.Recovered(),
				caseStates.Active(),
				caseStates.Critical(),
				caseStates.CasesPerOneMillion(),
				caseStates.DeathsPerOneMillion(),
			},
		)
	}

	t.Render()

}
