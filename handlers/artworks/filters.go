package artworks

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
)

type filters struct {
	Title         string
	SchoolString  string
	ArtFormString string
	ArtTypeString string
	ArtistString  string
	Page          string
}

// AnyFilterActive checks if any filter is active.
// It returns true if any of the filter fields (Title, SchoolString, ArtFormString, ArtTypeString, ArtistString) is not empty.
func (f *filters) AnyFilterActive() bool {
	return f.Title != "" || f.SchoolString != "" || f.ArtFormString != "" || f.ArtTypeString != "" || f.ArtistString != ""
}

// FingerPrint returns a unique fingerprint string based on the filter values.
// The fingerprint is generated by concatenating the title, school, art form, art type, and artist strings.
func (f *filters) FingerPrint() string {
	return f.Title + ":" + f.SchoolString + ":" + f.ArtFormString + ":" + f.ArtTypeString + ":" + f.ArtistString + ":" + f.Page
}

// BuildFilter builds a filter string and parameters based on the values of the filters struct.
// The filter string is used to filter artworks based on various criteria such as title, school, art form, art type, and artist.
// The parameters map contains the values to be substituted in the filter string.
// The filter string and parameters are returned as a string and dbx.Params respectively.
func (f *filters) BuildFilter() (string, dbx.Params) {
	filterString := "published = true && author:length > 0"
	params := dbx.Params{}

	if f.Title != "" {
		filterString = filterString + " && title ~ {:title}"
		params["title"] = f.Title
	}

	if f.SchoolString != "" {
		filterString = filterString + " && school.slug = {:art_school}"
		params["art_school"] = f.SchoolString
	}

	if f.ArtFormString != "" {
		filterString = filterString + " && form.slug = {:art_form}"
		params["art_form"] = f.ArtFormString
	}

	if f.ArtTypeString != "" {
		filterString = filterString + " && type.slug = {:art_type}"
		params["art_type"] = f.ArtTypeString
	}

	if f.ArtistString != "" {
		filterString = filterString + " && author.name ~ {:artist}"
		params["artist"] = f.ArtistString
	}

	return filterString, params
}

// BuildFilterString builds a filter string based on the values of the filters struct.
// It concatenates the filter parameters with their corresponding values and returns the resulting filter string.
func (f *filters) BuildFilterString() string {
	filterString := ""

	if f.Title != "" {
		filterString = filterString + "&title=" + f.Title
	}

	if f.SchoolString != "" {
		filterString = filterString + "&art_school=" + f.SchoolString
	}

	if f.ArtFormString != "" {
		filterString = filterString + "&art_form=" + f.ArtFormString
	}

	if f.ArtTypeString != "" {
		filterString = filterString + "&art_type=" + f.ArtTypeString
	}

	if f.ArtistString != "" {
		filterString = filterString + "&artist=" + f.ArtistString
	}

	if f.Page != "" {
		filterString = filterString + "&page=" + f.Page
	}

	return filterString
}

func buildFilters(app *pocketbase.PocketBase, c echo.Context) *filters {
	f := &filters{
		Title:         c.QueryParamDefault("title", ""),
		SchoolString:  c.QueryParamDefault("art_school", ""),
		ArtFormString: c.QueryParamDefault("art_form", ""),
		ArtTypeString: c.QueryParamDefault("art_type", ""),
		ArtistString:  c.QueryParamDefault("artist", ""),
		Page:          c.QueryParamDefault("page", ""),
	}

	return f
}