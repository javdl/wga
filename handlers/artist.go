package handlers

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

const (
	Yes           string = "yes"
	No            string = "no"
	NotApplicable string = "n/a"
)

func normalizedBioExcerpt(r *models.Record) string {
	s := []string{}

	yob := r.GetInt("year_of_birth")
	eyob := r.GetString("exact_year_of_birth")
	pob := r.GetString("place_of_birth")
	kpob := r.GetString("known_place_of_birth")

	yod := r.GetInt("year_of_death")
	eyod := r.GetString("exact_year_of_death")
	pod := r.GetString("place_of_death")
	kpod := r.GetString("known_place_of_death")

	if yob > 0 {

		c := []string{}

		prefix := "b."

		c = append(c, prefix)
		y := strconv.Itoa(yob)

		if eyob == No {
			y = "~" + y
		}

		c = append(c, y)

		if kpob == No {
			pob = pob + "?"
		}

		c = append(c, pob)

		s = append(s, strings.Join(c, " "))

	}

	if yod > 0 {

		c := []string{}

		prefix := "d."

		c = append(c, prefix)
		y := strconv.Itoa(yod)

		if eyod == No {
			y = "~" + y
		}

		c = append(c, y)

		if kpod == No {
			pod = pod + "?"
		}

		c = append(c, pod)

		s = append(s, strings.Join(c, " "))

	}

	return strings.Join(s, ", ")
}

func registerArtist(app *pocketbase.PocketBase) {

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {

		e.Router.GET("artists/:name", func(c echo.Context) error {
			slug := c.PathParam("name")
			cacheKey := "artist:" + slug

			if isHtmxRequest(c) {
				cacheKey = cacheKey + "-htmx"
			}

			html := ""

			found := app.Cache().Has(cacheKey)
			// found := false

			if found {
				html = app.Cache().Get(cacheKey).(string)
			} else {

				err := godotenv.Load()

				if err != nil {
					return apis.NewBadRequestError("Error loading .env file", err)
				}

				fullUrl := os.Getenv("WGA_PROTOCOL") + "://" + c.Request().Host + c.Request().URL.String()
				artist, err := app.Dao().FindRecordsByFilter("artists", "slug = '"+slug+"'", "+name", 1, 0)

				if err != nil {
					return apis.NewNotFoundError("", err)
				}

				works, err := app.Dao().FindRecordsByFilter("artworks", "author = '"+artist[0].GetString("id")+"'", "+title", 100, 0)

				if err != nil {
					return apis.NewNotFoundError("", err)
				}

				data := map[string]any{
					"Name":         artist[0].GetString("name"),
					"Bio":          artist[0].GetString("bio"),
					"Works":        []map[string]string{},
					"Slug":         slug,
					"BioExcerpt":   normalizedBioExcerpt(artist[0]),
					"CurrentUrl":   fullUrl,
					"Profession":   artist[0].GetString("profession"),
					"YearOfBirth":  artist[0].GetString("year_of_birth"),
					"YearOfDeath":  artist[0].GetString("year_of_death"),
					"PlaceOfBirth": artist[0].GetString("place_of_birth"),
					"PlaceOfDeath": artist[0].GetString("place_of_death"),
					"Jsonld": map[string]any{
						"@context":      "https://schema.org",
						"@type":         "Person",
						"name":          artist[0].GetString("name"),
						"url":           fullUrl,
						"hasOccupation": artist[0].GetString("profession"),
					},
				}

				if artist[0].GetInt("year_of_birth") > 0 {
					data["Jsonld"].(map[string]any)["birthDate"] = artist[0].GetString("year_of_birth")
				}

				if artist[0].GetInt("year_of_death") > 0 {
					data["Jsonld"].(map[string]any)["deathDate"] = artist[0].GetString("year_of_death")
				}

				if artist[0].GetString("place_of_birth") != "" {
					data["Jsonld"].(map[string]any)["birthPlace"] = map[string]string{
						"@type": "Place",
						"name":  artist[0].GetString("place_of_birth"),
					}
				}

				if artist[0].GetString("place_of_death") != "" {
					data["Jsonld"].(map[string]any)["deathPlace"] = map[string]string{
						"@type": "Place",
						"name":  artist[0].GetString("place_of_death"),
					}
				}

				for _, w := range works {
					data["Works"] = append(data["Works"].([]map[string]string), map[string]string{
						"Id":        w.GetId(),
						"Title":     w.GetString("title"),
						"Comment":   w.GetString("comment"),
						"Technique": w.GetString("technique"),
						"Image":     generateFileUrl(app, "artworks", w.GetString("id"), w.GetString("image")),
						"Thumb":     generateThumbUrl(app, "artworks", w.GetString("id"), w.GetString("image"), "320x240"),
					})
				}

				if isHtmxRequest(c) {
					html, err = renderBlock("artist:content", data)
				} else {
					html, err = renderPage("artist", data)
				}

				if err != nil {
					// or redirect to a dedicated 404 HTML page
					return apis.NewNotFoundError("", err)
				}

				app.Cache().Set(cacheKey, html)
			}

			c.Response().Header().Set("HX-Push-Url", "/artists/"+slug)

			return c.HTML(http.StatusOK, html)
		})

		return nil
	})
}
