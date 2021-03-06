package controllers

import (
	"net/http"
	"fmt"
	"strconv"
	"io/ioutil"
	"github.com/go-zoo/bone"

	"github.com/JacksonGariety/cetch/app/models"
	"github.com/JacksonGariety/cetch/app/utils"
)

func EntryShow(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(bone.GetValue(r, "id"))
	entry := models.Entry{}

	if !models.ExistsById(&entry, id) || err != nil {
		utils.NotFound(w, r)
		return
	}

	models.DB.Model(&entry).Related(&entry.Competition)
	models.DB.Model(&entry).Related(&entry.User)

	utils.Render(w, r, "entry.html", &utils.Props{ "entry": entry })
}

func EntryCreate(w http.ResponseWriter, r *http.Request) {
	comp, _ := (&models.Competition{}).Current()
	current_user := (*r.Context().Value("data").(*utils.Props))["current_user"]
	entry := &models.Entry{
		UserID: current_user.(models.User).ID,
		CompetitionID: comp.ID,
	}
	if models.Exists(&entry) {
		utils.BadRequest(w, r)
		return
	}

	// read the file
	reader, _ := r.MultipartReader()

	// language
	part, _ := reader.NextPart()
	language, _ := ioutil.ReadAll(part)
	languageString := string(language)

	if (languageString == "") {
		utils.BadRequest(w, r)
		return
	}

	// code
	part, _ = reader.NextPart()
	code, _ := ioutil.ReadAll(part)
	codeString := string(code)

	// pass it to the runner
	result, execTime, err := models.ProgramResultAndExecTime(codeString, languageString)

	if execTime == nil {
		utils.Render(w, r, "competition_show.html", &utils.Props{
			"competition": comp,
			"language": languageString,
			"current": true,
			"stderrError": true,
		})
	} else if *result == comp.Solution && err == nil {
		// run it 4 times and get the average
		n := 10
		avg := *execTime
		for i := 0; i < n; i++ {
			_, execTime, _ := models.ProgramResultAndExecTime(codeString, languageString)
			avg += *execTime
		}
		avg = avg / float64(n + 1)

		user := (*r.Context().Value("data").(*utils.Props))["current_user"]
		entry := models.Entry{
			CompetitionID: comp.ID,
			UserID: user.(models.User).ID,
			Language: "go",
			Code: codeString,
			ExecTime: avg,
		}

		models.Create(&entry)

		http.Redirect(w, r, fmt.Sprintf("/entry/%v", entry.ID), 307)
	} else {
		utils.Render(w, r, "competition_show.html", &utils.Props{
			"competition": comp,
			"language": languageString,
			"current": true,
			"stdoutError": true,
		})
	}
}
