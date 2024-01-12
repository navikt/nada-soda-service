# NADA Soda Service

API for å samle inn Soda metrikker til sentralt lager og varsle team på Slack om feilende datakvalitetstester.
Team kan sette opp Soda-jobber ved å ta utgangspunkt i [navikt/dp-nada-soda](https://github.com/navikt/dp-nada-soda). 
Testresultatene sendes så til BigQuery-datasettet `nada-soda-service`.

- I dev lagres dette i `nada-dev-db2e.soda.historic`
- I prod lagres dette i `nada-prod-6977.soda.historic`

Videre vil `nada-soda-service` gå igjennom testresultatene, identifisere feil, og poste en rapport om disse avvikene til Slack kanalen teamet har angitt for Soda-jobben.

## Skisse

````mermaid
graph LR;
    subgraph " "
        subgraph "BigQuery Team A"
            bq1["datasett"]
            bq2["datasett"]
        end
        subgraph "BigQuery Team NADA"
            bq3["sentralt lager"]
        end
    end

    subgraph "Slack"
        teamslack["team slack kanal"]
    end

    subgraph "NAIS cluster"
        subgraph "Nada namespace"
            sodaservice["SODA service"]--"Lagre resultater sentralt"-->bq3
            sodaservice--"Rapporter datakvalitetsavvik"-->teamslack
        end
        subgraph "Team A namespace"
            sodajobb["SODA job"]--"Kjør tester"-->bq1
            sodajobb["SODA job"]--"Kjør tester"-->bq2
            sodajobb--"Send testresultater"-->sodaservice
        end
    end
````

## Kjør lokalt

Sett følgende miljøvariabler i terminalen du skal kjøre opp appen:

- `GCP_TEAM_PROJECT_ID`: GCP prosjektet hvor du ønsker å lagre Soda testresultater
- `BIGQUERY_DATASET`: Datasett i GCP prosjektet hvor du ønsker å lagre Soda testresultater
- `BIGQUERY_TABLE`: Tabell i datasett hvor du ønsker å lagre Soda testresultater
- `SLACK_TOKEN`: Token for Slack appen som skal brukes for å poste datakvalitetsavvik.

`SLACK_TOKEN` kan settes med:
````bash
export SLACK_TOKEN=$(kubectl get secret --context=dev-gcp --namespace=nada slack-token -o jsonpath='{.data.SLACK_TOKEN}' | base64 -d)`
````

Kjør så appen med:
````bash
go run .
````
