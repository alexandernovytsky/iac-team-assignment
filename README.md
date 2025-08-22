# iac-team-assignment

## Part 1

### Assumptions
- The Y in info logs message ("Counter: Y") represents the total number of objects in a successful data ingestion.
- The X in error logs message ("Error: X") represent the error code occurred during data ingestion.
(was considering to assume it to be the number of failed object but saw value 0, hence doesn't make sense)

### Alert

The requirement was to have the alert to be triggered occasionally.
As the logs count is similar between info and error, and there isn't much data in the log to rely on, I chose to create an alert that will track for periods of time of tilted ratio between infos and errors.
This alert will make aware of potential prolonged period which will indicate on unusual operation behaviour.  


## Part 2

- Implemented a simple SDK using a centralised approach for accessing clients.
- Clients are created though a ClientCreator which allows to reuse shared configuration between clients.
- Clients use HTTP requests to interact with BE.
- used openapi.yaml file to generate types via 3rd party [lib](https://github.com/oapi-codegen/oapi-codegen?tab=readme-ov-file#generating-api-clients)

generate types: 

    go generate ./...

export ENV vars before running the program:

    export CORALOGIX_API_KEY=
    export CORALOGIX_WEBHOOK_URL=

run app

    go run ./app/main.go

## Part 3

### API Strength

1. Versioning per resource (`v3/alerts-def`) allows dynamic resource usage and compatibility
2. Follows standard API error codes 2xx, 4xx, 5xx
3. Utilises HTTP methods for CRUD operations on resources with proper noun names
4. Standard authorization practise

### API Weakness

#### consistency with RESTful API practises

1. although versioning in URL is accepted practise, RESTful argues that versioning should be handled through headers (e.g., Accept: application/vnd.company.v3+json) rather than in the path, to keep the URL focused on resources.
2. 200 HTTP status code on creation instead of 201
3. PUT HTTP method URL to replace/full update doesn't use object identifier in URL, i.e. `v3/alerts-def/{id}`
4. No PATCH support
5. properties names should be plural if expected an array, e.g. `dayOfWeek: array` should be `daysOfWeek`

#### Complexity and Verbosity

payloads can cause confusion for end users and difficulty to maintain for developer.

1. payload object includes multiple models for different types in a single object, e.g. create alert/webhook payload that includes all possible types although only one can be used at a time. This makes the object large and not very user friendly. Additionally, make it harder to monitor and limit per type as it will require to read payload explicitly instead to implicitly know the type from the EP itself.
2. enums names are long and not friendly, e.g. `ALERT_DEF_PRIORITY_P5_OR_UNSPECIFIED` instead of just a simple `P5` or build a verbal scale like `lowest, low, medium, high, highest`.
3. Filtering 


### Suggested Improvements

1. Split endpoints by types, this will make request object to be clearer and will allow to separate type specific logic handling, storing and improve monitoring, limiting capabilities on API EP level. Possible to group by type of type e.g. `/v3/alert-defs/logs` or even more specific `/v3/alert-defs/logsImmediate` that excepts only it's specific request payload.
2. align with RESTful practises, e.g. PUT operations path in alerts `v3/alerts-def/{id}`, dedicated HTTP status codes `201 Created` for creation with POST, and `204` for PUT when no data is returned like in update webhook.
3. Consider addition of `304` status code for validating that large cached objects are still valid on GET requests like "list".
4. Change endpoint to RESTful style, remove operations from URL path, e.g. `v1/outgoing-webhooks:list` and use  `v1/outgoing-webhooks?q=...&f=...` with potentially query filtering for clarity. This will allow to create GET HTTP requests that can be easily cached in BE/CDN with the specific filters.