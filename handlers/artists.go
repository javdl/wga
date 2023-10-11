package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"blackfyre.ninja/wga/utils"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func registerArtists(app *pocketbase.PocketBase) {

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {

		e.Router.GET("/artists", func(c echo.Context) error {

			limit := 30
			page := 1
			searchExpression := ""
			searchExpressionPresent := false

			if c.QueryParam("page") != "" {
				err := error(nil)
				page, err = strconv.Atoi(c.QueryParam("page"))

				if err != nil {
					return apis.NewBadRequestError("Invalid page", err)
				}
			}

			if c.QueryParam("q") != "" {
				searchExpression = c.QueryParam("q")
			}

			if c.QueryParams().Has("q") {
				searchExpressionPresent = true
			}

			offset := (page - 1) * limit

			data := map[string]any{
				"Content": "",
			}

			filter := "published = true"

			if searchExpression != "" {
				filter = filter + " && name ~ {:searchExpression}"
			}

			records, err := app.Dao().FindRecordsByFilter(
				"artists",
				filter,
				"+name",
				limit,
				offset,
				dbx.Params{
					"searchExpression": searchExpression,
				},
			)

			if err != nil {
				return apis.NewBadRequestError("Invalid page", err)
			}

			totalRecords, err := app.Dao().FindRecordsByFilter(
				"artists",
				filter,
				"+name",
				0,
				0,
				dbx.Params{
					"searchExpression": searchExpression,
				},
			)

			if err != nil {
				return apis.NewBadRequestError("Invalid page", err)
			}

			recordsCount := len(totalRecords)

			preRendered := []map[string]string{}

			for _, m := range records {

				// TODO: handle aka

				school := m.GetStringSlice("school")

				schoolCollector := []string{}

				for _, s := range school {
					r, err := app.Dao().FindRecordById("schools", s)

					if err != nil {
						log.Print("school not found")
						continue
					}

					schoolCollector = append(schoolCollector, r.GetString("name"))

				}

				row := map[string]string{
					"Name":       m.GetString("name"),
					"Url":        artistUrl(m.GetString("slug")),
					"Profession": m.GetString("profession"),
					"BornDied":   normalizedBirthDeathActivity(m),
					"Schools":    strings.Join(schoolCollector, ", "),
				}

				preRendered = append(preRendered, row)
			}

			data["Content"] = preRendered
			data["Count"] = recordsCount

			pagination := utils.NewPagination(recordsCount, limit, page, "/artists?q="+searchExpression)

			data["Pagination"] = pagination.Render()

			html := ""

			if isHtmxRequest(c) {
				blockToRender := "artists:content"

				if searchExpression != "" || searchExpressionPresent {
					blockToRender = "artists:search-results"
				}

				html, err = renderBlock(blockToRender, data)
			} else {
				html, err = renderPage("artists", data)
			}

			if err != nil {
				// or redirect to a dedicated 404 HTML page
				return apis.NewNotFoundError("", err)
			}

			currentUrl := c.Request().URL.String()

			c.Response().Header().Set("HX-Push-Url", currentUrl)

			return c.HTML(http.StatusOK, html)
		})

		return nil
	})
}