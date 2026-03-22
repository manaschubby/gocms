- Accounts Keep track of user entities adding/subscribing to the CMS DB.

Package Domain has model definitions and easy to use DTOs / helpers if required.

Package Repo exposes interface implementations to DB interactions related to the domain models

Package handlers uses the service interfaces. We keep request payload validation in handlers and delegate business logic over to services

Imp Notes:
Numerical Columns are computed as float64. They are parsed into json.Number and used as float64. Reference: [internal/modules/cms/services/entry.go](./services/entry.go) and [internal/modules/cms/common/json.go](./common/json.go)
