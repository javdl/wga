package handlers

import (
    "fmt"
    "github.com/labstack/echo/v5"
    "github.com/pocketbase/pocketbase/models"
)

func normalizedBirthDeathActivity(record *models.Record) string {
	Start := record.GetInt("year_of_birth")
	End := record.GetInt("year_of_death")

	return fmt.Sprintf("%d-%d", Start, End)
}

func generateArtistSlug(artist *models.Record) string {
    if artist == nil {
        return ""
    }
    return artist.GetString("slug") + "-" + artist.GetString("id")
}

func generateCurrentPageUrl(c echo.Context) string {
    if c == nil || c.Request() == nil {
        return ""
    }
    return c.Scheme() + "://" + c.Request().Host + c.Request().URL.String()
}
