test:
	go test -v ./...

run:
	go run cmd/server/main.go start

build:
	go build -o bin/server cmd/server/main.go

REPORTED_ISSUE_LINKS := "LOY-1297"

jira-test:
	rm -f testrun.*
	REPORTED_ISSUE_LINKS=${REPORTED_ISSUE_LINKS} JIRA_PWD=${PWD} \
		go test -count=1 -p 1 -covermode=count -coverprofile loyalty-file-processing-coverage.cov -tags integration \
		./...&& (cat testrun.tmp.json | jq -s "." > testrun.json)

jira-test-push:
	 /bin/zsh ./push-testcase.sh

test-n-coverage:
	go test ./... -coverprofile=fps.cov.tmp -coverpkg=./... -covermode count
	cat fps.cov.tmp | grep -s -v -e"internal/ent" -e"/cmd" -e"/configs" -e"api/server/common" > fps.cov \
    && go tool cover -func=fps.cov \
	&& go tool cover -func fps.cov | grep ^total

coverage:
	go tool cover -func loyalty-file-processing-coverage.cov | grep ^total

migrate:
	echo \# make migrate name="$(name)"
	go run cmd/server/main.go migrate create $(name)

migrate-up:
	go run cmd/server/main.go migrate up

migrate-down-1:
	go run cmd/server/main.go migrate down 1

jobs:
	go run cmd/server/main.go jobs

jobs-process-file-flatten:
	go run cmd/server/main.go jobs process-file flatten

jobs-process-file-execute-task:
	go run cmd/server/main.go jobs process-file execute-task

jobs-process-file-execute-row-group:
	go run cmd/server/main.go jobs process-file execute-row-group

jobs-process-file-update-status:
	go run cmd/server/main.go jobs process-file update-status