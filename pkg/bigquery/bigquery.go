package bigquery

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"google.golang.org/api/googleapi"
)

type SodaResult struct {
	ID                 string   `json:"id"`
	Project            string   `json:"project"`
	Dataset            string   `json:"dataset"`
	Table              string   `json:"table"`
	Test               string   `json:"test"`
	Definition         string   `json:"definition"`
	Metrics            []string `json:"metrics"`
	ResourceAttributes []string `json:"resourceAttributes"`
	Time               string   `json:"time"`
	Column             string   `json:"column"`
	Type               string   `json:"type"`
	Filter             []string `json:"filter"`
}

type NadaBigQuery struct {
	client *bigquery.Client
	table  *bigquery.Table
	log    *logrus.Entry
}

func New(ctx context.Context, project, dataset, tableName string, log *logrus.Entry) (*NadaBigQuery, error) {
	client, err := bigquery.NewClient(ctx, project)
	if err != nil {
		return nil, err
	}

	table, err := createTableIfNotExists(ctx, client, dataset, tableName)
	if err != nil {
		return nil, err
	}

	return &NadaBigQuery{
		client: client,
		table:  table,
		log:    log,
	}, nil
}

func (b *NadaBigQuery) StoreSodaResults(ctx context.Context, results []SodaResult) error {
	inserter := b.table.Inserter()
	for _, r := range results {
		fmt.Println("inserting row", r)
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
