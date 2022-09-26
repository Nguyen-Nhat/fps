### Prerequisite

* Before you begin, setting git to use go private repo

```
git config --global url."git@git.teko.vn:".insteadOf "https://git.teko.vn/"
```

* direnv: For manage and auto switch config on `.envrc` https://direnv.net/
* go 1.14+ with enable module support
* protobuf compiler: https://developers.google.com/protocol-buffers
* protoc-gen-go: https://github.com/golang/protobuf/protoc-gen-go
* protoc-gen-rpc-server: for gen rpc-runtime server https://git.teko.vn/shared/rpc-framework

You can install it by those command (not sure it can working properly)
```
## For macOS users
$ make install-osx

# for Linux user
$ sudo echo this command to promp password
$ make install-linux
```

### RPC Definition

Checkout at https://git.teko.vn/shared/rpc

### Application Config

* Set inside `.envrc`

```
$ cat .envrc
export LOG_LEVEL=debug
export LOG_MODE=development
```

* Auto env by command `direnv allow`

```
$ direnv allow
direnv: loading .envrc
direnv: export +LOG_LEVEL +LOG_MODE
```

### Project Config

Everything is inside a `Makefile`, there are something you could edit properly

```
PACKAGE = rpc.tekoapis.com

# List of implements server on rpc
SERVERS = health/health.proto sample/sample.proto

# List of depenencies client on rpc
CLIENTS =  # sample1/sample1.proto sample2/sample2.proto

# Enable grpc gateway
GATEWAY = true

```

### Run

* Update rpc && common lib to upstream

```
# Run this command if sometime local dev fail to update common
# go mod edit -droprequire go.tekoapis.com/kitchen
# go mod edit -droprequire rpc.tekoapis.com/rpc
# Then update common
$ make common-update
```

* Re-generate source code and run server

```
$ make run
```

* Just commpile all

```
$ make all
```

* Just generated depenencies source code

```
$ make generate
```


### Migration

* Every migration script is located at `migrations` folder, which follow pattern by https://github.com/golang-migrate/migrate
* Create a new migration script

```
$ make migrate name=[name of migration]

#example:

$ make migrate name="create person table"
2020-04-16T18:15:32.914+0700	INFO	migrate/migrate.go:103	create migration	{"name": "create-person-table"}
2020-04-16T18:15:32.914+0700	INFO	migrate/migrate.go:104	up script	{"up": "migrations/20200416181532_create-person-table.up.sql"}
2020-04-16T18:15:32.914+0700	INFO	migrate/migrate.go:105	down script	{"down": "migrations/20200416181532_create-person-table.up.sql"}
```

* Edit migration script by SQL (you should edit both up and down script)

* Run migration

```
make migrate-up
```

* Or sometime you want it migrate 1 step back

```
make migrate-down-1
```


### Common Lib

Should place shared with other team at: https://git.teko.vn/shared/kitchen
