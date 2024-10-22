package bigquery

import (
	"context"
	"errors"

	"cloud.google.com/go/bigquery"
	"github.com/navikt/nada-soda-service/pkg/models"
	"google.golang.org/api/googleapi"
)

type Client struct {
	client  *bigquery.Client
	project string
	dataset string
	table   string
}

type BigQueryRow struct {
	ID                 string   `json:"id"`
	Project            string   `json:"project"`
	Dataset            string   `json:"dataset"`
	Table              string   `json:"table"`
	Test               string   `json:"test"`
	Outcome            string   `json:"outcome"`
	Definition         string   `json:"definition"`
	Metrics            []string `json:"metrics"`
	ResourceAttributes []string `json:"resourceAttributes"`
	Time               string   `json:"time"`
	Column             string   `json:"column"`
	Type               string   `json:"type"`
	Filter             []string `json:"filter"`
	Image              string   `json:"image"`
}

func New(ctx context.Context, project, dataset, table string) (*Client, error) {
	client, err := bigquery.NewClient(ctx, project)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	_, err = createTableIfNotExists(ctx, client, dataset, table)
	if err != nil {
		return nil, err
	}

	return &Client{
		client:  client,
		project: project,
		dataset: dataset,
		table:   table,
	}, nil
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
		{Name: "image", Type: bigquery.StringFieldType},
	}

	metadata := &bigquery.TableMetadata{
		Schema: schema,
	}

	tableRef := bqClient.Dataset(dataset).Table(table)
	if err := tableRef.Create(ctx, metadata); err != nil {
		var e *googleapi.Error
		if ok := errors.As(err, &e); ok {
			if e.Code == 409 {
				return tableRef, nil
			}
		}
		return nil, err
	}

	return tableRef, nil
}

func (b *Client) StoreResults(ctx context.Context, report models.SodaReport) error {
	client, err := bigquery.NewClient(ctx, b.project)
	if err != nil {
		return err
	}
	defer client.Close()

	table := client.Dataset(b.dataset).Table(b.table)
	inserter := table.Inserter()
	for _, r := range toBigQueryRows(report) {
		if err := inserter.Put(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

func toBigQueryRows(report models.SodaReport) []BigQueryRow {
	rows := []BigQueryRow{}
	for _, r := range report.Results {
		rows = append(rows, BigQueryRow{
			ID:                 r.ID,
			Project:            report.GCPProject,
			Dataset:            report.Dataset,
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
			Image:              report.DockerImage,
		})
	}

	return rows
}
