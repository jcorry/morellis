package nlp

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jcorry/morellis/pkg/models/mysql"

	"google.golang.org/api/iterator"

	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"

	"cloud.google.com/go/dialogflow/apiv2"
)

type EntityTypesService struct {
	Client *dialogflow.EntityTypesClient
	DB     *sql.DB
}

var Svc EntityTypesService

func NewSessionEntityTypesService() (*EntityTypesService, error) {
	if Svc.Client != nil {
		return &Svc, nil
	}

	ctx := context.Background()
	c, err := dialogflow.NewEntityTypesClient(ctx)
	if err != nil {
		return nil, err
	}

	Svc := &EntityTypesService{
		Client: c,
	}

	return Svc, nil
}

func (s *EntityTypesService) ListEntityTypes() ([]*dialogflowpb.EntityType, error) {
	ctx := context.Background()
	req := &dialogflowpb.ListEntityTypesRequest{
		Parent:   `projects/morellis-api/agent`,
		PageSize: 120,
	}

	it := s.Client.ListEntityTypes(ctx, req)
	var results []*dialogflowpb.EntityType
	for {
		resp, err := it.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}
		results = append(results, resp)
	}
	return results, nil
}

func (s *EntityTypesService) AddEntities() error {
	fmt.Println("Adding entity types...")
	ctx := context.Background()

	flavorModel := &mysql.FlavorModel{DB: s.DB}

	flavors, err := flavorModel.List(100, 0, "")
	if err != nil {
		fmt.Println(fmt.Sprintf("%s", err))
		return err
	}

	var flavorEntities []*dialogflowpb.EntityType_Entity

	for _, f := range flavors {
		e := &dialogflowpb.EntityType_Entity{
			Value:    f.Name,
			Synonyms: []string{f.Name},
		}
		flavorEntities = append(flavorEntities, e)
	}

	fmt.Println(fmt.Sprintf("Flavor count: %d", len(flavorEntities)))

	req := &dialogflowpb.CreateEntityTypeRequest{
		Parent: `projects/morellis-api/agent`,
		EntityType: &dialogflowpb.EntityType{
			DisplayName:       `Flavor`,
			Kind:              dialogflowpb.EntityType_KIND_LIST,
			AutoExpansionMode: dialogflowpb.EntityType_AUTO_EXPANSION_MODE_DEFAULT,
			Entities:          flavorEntities,
		},
	}

	resp, err := s.Client.CreateEntityType(ctx, req)
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("resp: %+v", resp))

	return nil
}
