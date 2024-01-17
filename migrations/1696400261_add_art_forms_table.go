package migrations

import (
	"encoding/json"

	"github.com/blackfyre/wga/assets"
	"github.com/blackfyre/wga/utils"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

type ArtForm struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func init() {
	tName := "Art_forms"
	tId := "art_forms"
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		collection := &models.Collection{}

		collection.Name = tName
		collection.Id = tId
		collection.Type = models.CollectionTypeBase
		collection.System = false
		collection.MarkAsNew()
		collection.Schema = schema.NewSchema(
			&schema.SchemaField{
				Id:          tId + "_name",
				Name:        "name",
				Type:        schema.FieldTypeText,
				Options:     &schema.TextOptions{},
				Presentable: true,
			},
			&schema.SchemaField{
				Id:      "schools_slug",
				Name:    "slug",
				Type:    schema.FieldTypeText,
				Options: &schema.TextOptions{},
			},
		)

		err := dao.SaveCollection(collection)

		if err != nil {
			return err
		}

		data, err := assets.InternalFiles.ReadFile("reference/forms.json")

		if err != nil {
			return err
		}

		var c []ArtForm

		err = json.Unmarshal(data, &c)

		if err != nil {
			return err
		}

		for _, g := range c {
			q := db.Insert(tId, dbx.Params{
				"id":   g.ID,
				"name": g.Name,
				"slug": utils.Slugify(g.Name),
			})

			_, err = q.Execute()

			if err != nil {
				return err
			}

		}

		return nil
	}, func(db dbx.Builder) error {
		q := db.DropTable(tId)
		_, err := q.Execute()

		return err
	})
}
