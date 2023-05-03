# NADA Soda Service
API for å samle inn soda metrikker til sentralt lager og varsle team på slack om feilende datakvalitetstester. Team kan sette opp SODA jobber ved å ta utgangspunkt i [eksempel repo](https://github.com/navikt/dp-nada-soda). Testresultatene sendes så til `nada-soda-service` som lagrer resultatene sentralt i BigQuery. 

- I dev lagres dette i `nada-dev-db2e.soda.historic`
- I prod lagres dette i `nada-prod-6977.soda.historic`

Videre vil `nada-soda-service` gå igjennom testresultatene og identifisere feil og varslinger og poste en rapport om disse avvikene til slack kanalen teamet angir for SODA jobben.

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

- `GCP_TEAM_PROJECT_ID`: GCP prosjektet hvor du ønsker å lagre soda testresultater
- `BIGQUERY_DATASET`: Datasett i GCP prosjektet hvor du ønsker å lagre soda testresultater
- `BIGQUERY_TABLE`: Tabell i datasett hvor du ønsker å lagre soda testresultater
- `SLACK_TOKEN`: Token for slack appen som skal brukes for å poste datakvalitetsavvik. 

`SLACK_TOKEN` kan settes med:
````bash
export SLACK_TOKEN=$(kubectl get secret --context=dev-gcp --namespace=nada slack-token -o jsonpath='{.data.SLACK_TOKEN}' | base64 -d)`
````

Kjør så appen med:
````bash
go run .
````

## Bygg og deploy til nais
````bash
docker build -t ghcr.io/navikt/nada-soda-service:<tag>
docker push ghcr.io/navikt/nada-soda-service:<tag>
````

Slack token mountes fra secreten `slack-token` som angitt i [nais.yaml](https://github.com/navikt/nada-soda-service/blob/main/.nais/nais.yaml#L26-L27).

````bash
k apply -f .nais/nais.yaml
````
