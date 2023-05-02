package bigquery

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/navikt/nada-soda-service/pkg/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"google.golang.org/api/googleapi"
)

type BigQueryClient struct {
	client  *bigquery.Client
	project string
	dataset string
	table   string
	log     *logrus.Entry
}

func New(ctx context.Context, project, dataset, table string, log *logrus.Entry) (*BigQueryClient, error) {
	client, err := bigquery.NewClient(ctx, project)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	_, err = createTableIfNotExists(ctx, client, dataset, table)
	if err != nil {
		return nil, err
	}

	return &BigQueryClient{
		client: client,
		table:  table,
		log:    log,
	}, nil
}

func (b *BigQueryClient) StoreSodaResults(ctx context.Context, sodaTest models.SodaTest) error {
	client, err := bigquery.NewClient(ctx, b.project)
	if err != nil {
		return err
	}
	defer client.Close()

	table := client.Dataset(b.dataset).Table(b.table)
	inserter := table.Inserter()
	for _, r := range toBigqueryRows(sodaTest) {
		if err := inserter.Put(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

func createTableIfNotExists(ctx context.Context, bqClient *bigquery.Client, dataset, table string) (*bigquery.Table, error) {
	schema := bigquery.Schema{
		{Name: "id", Type: bigquery.StringFieldType, Required: true},
		{Name: "project", Type: bigquery.StringFieldType, Required: true},
		{Name: "dataset", Type: bigquery.StringFieldType, Required: true},
		{Name: "table", Type: bigquery.StringFieldType, Required: true},
		{Name: "test", Type: bigquery.StringFieldType, Required: true},
		{Name: "outcome", Type: bigquery.StringFieldType, Required: true},
		{Name: "definition", Type: bigquery.StringFieldType},
		{Name: "metrics", Type: bigquery.StringFieldType, Repeated: true},
		{Name: "resourceAttributes", Type: bigquery.StringFieldType, Repeated: true},
		{Name: "time", Type: bigquery.TimestampFieldType},
		{Name: "column", Type: bigquery.StringFieldType},
		{Name: "type", Type: bigquery.StringFieldType},
		{Name: "filter", Type: bigquery.StringFieldType, Repeated: true},
	}

	metadata := &bigquery.TableMetadata{
		Schema: schema,
	}

	tableRef := bqClient.Dataset(dataset).Table(table)
	if err := tableRef.Create(ctx, metadata); err != nil {
		var e *googleapi.Error
		if ok := xerrors.As(err, &e); ok {
			if e.Code == 409 {
				fmt.Println("already exists")
				return tableRef, nil
			}
		}
		return nil, err
	}

	return tableRef, nil
}

func toBigqueryRows(sodaTest models.SodaTest) []models.BigqueryRow {
	out := []models.BigqueryRow{}
	for _, r := range sodaTest.Results {
		out = append(out, models.BigqueryRow{
			ID:                 r.ID,
			Project:            sodaTest.GCPProject,
			Dataset:            sodaTest.Dataset,
			Table:              r.Table,
			Test:               r.Test,
			Outcome:            r.Outcome,
			Definition:         r.Definition,
			Metrics:            r.Metrics,
			ResourceAttributes: r.ResourceAttributes,
			Time:               r.Time,
			Column:             r.Column,
			Type:               r.Type,
			Filter:             r.Filter,
		})
	}

	return out
}
