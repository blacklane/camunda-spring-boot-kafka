# worker


### Template Usage

* clone repository
* replace template name with actual service name
```bash
find ./ -type f -exec sed -i '' 's/worker/<actual-service-name>/' {} \;
```
_(if sed gives you some errors, check if it worked anyway. Otherwise try using gnu-sed (gsed) instead)_
* (Optional) enable database config in `config.go`
* (Optional) uncomment sql-migrate dependency in dockerfile
* generate deployment files with blops as described in the [handbook](http://handbook.int.blacklane.io/devops/deploy.html#step-4-configure-kubernetes-and-set-up-some-environments)
```bash
blops generate -n <actual-service-name> --logs-enabled
```
* adjust generated `drone.yml` to your needs
* adjust generated `/deploy/overlays/...` with your required envs
* adjust generated `/deploy/base/config.yml` with your config envs
* adjust generated `/deploy/base/secrets.yml` with your secrets envs
* remove this section from README.md

--------

{{ Short description of what the service / library does. What is the scope of it? Why does it exist? }}

{{ Keep relevant Drone CI or Travis CI badge, replace `repo` (and `branch` if necessary) }}

![Build status badge](https://drone.blacklane.net/api/badges/blacklane/{{repo}}/status.svg?branch=master)

or

[![Build status badge](https://travis-ci.com/blacklane/elli.svg?token=eqEro8Uh7aLKHHx8ps1S&branch=master)](https://travis-ci.com/blacklane/{{repo}})

## Owners & contact

This service is owned by the [<Your-Team-Name>](https://blacklane.atlassian.net/wiki/home) team. You can get in touch with us via:
- Slack: `#your-team-slack`
- Email: `#your-team-email`

## Setup instructions

Install dependencies: `go mod download`

Configuration takes place via environment variables. Copy the sample file: `cp .env.sample .env`.
Take a look at the configuration in .env and adapt database connection or other configuration as needed.

To start the server: `go run main.go`
Access server on http://localhost:8000.

### Requirements

- go >= v1.13
- kafka (e.g. [bitnami docker-compose setup](https://github.com/bitnami/bitnami-docker-kafka))

### Tests

run `go test ./...` to execute all unit tests


## Deployment

{{ Describe how this service is being deployed. On wich environments is this service running? }}

This service is automatically being deployed via Kubernetes. TODO: link to Kubernetes handbook! It runs on the following stacks:

- Production
- Sandbox
- Auto
- Testing

## Contributing

{{ List any contribution guidelines in case they are different to the department's standard. }}

{{ Which Git flow is being used? Trunk-based development or develop / master? Keep the relevant paragraph. }}


This project is using [Gitflow](http://handbook.int.blacklane.io/git.html#gitflow) for development. Check the linked handbook for details.

or

This project is using [Simplified flow (GitHub flow)](http://handbook.int.blacklane.io/git.html#simplified-flow-github-flow) for development. Check the linked handbook for details.


## Links

{{ Add links that are helpful and relevant }}

- Documentation
  - [Confluence]()
  - [Web API](api-docs.int.blacklane.io/)
  - [Queues & Events](TODO)
- [Handbook](http://handbook.int.blacklane.io/git.html#simplified-flow-github-flow)

### Logs & dashboards

- Production
- - [](Logs)
- - [](Dashboard)
- Testing
- - [](Logs)
- - [](Dashboard)
