test:
	go test -v ./...

run:
	go run cmd/server/main.go start

build:
	go build -o bin/server cmd/server/main.go

REPORTED_ISSUE_LINKS := "LOY-1231"

jira-test:
	rm -f testrun.*
	REPORTED_ISSUE_LINKS=${REPORTED_ISSUE_LINKS} JIRA_PWD=${PWD} \
		go test -count=1 -p 1 -covermode=count -coverprofile loyalty-file-processing-coverage.cov -tags integration \
		./... && (cat testrun.tmp.json | jq -s "." > testrun.json)
jira-test-push:
	 /bin/zsh ./push-testcase.sh

coverage:
	go tool cover -func loyalty-file-processing-coverage.cov | grep ^total

migrate:
	go run cmd/migrate/main.go start
