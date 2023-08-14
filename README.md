# rll-url-shortener
rll-url-shortener



## Running Application on Docker Container

This application you can run on docker container.

### Prerequisite :
* Docker
* OS Linux Docker Image, 
* Golang version 1.14 or latest

### Environment Variable available


### build docker image
```bash
$ docker build -t {IMAGE_NAME} -f deployment/dockerfiles/dockerfile--dev .
# example
$ docker build -t rll-url-shortener -f deployment/dockerfiles/dockerfile-dev .
```

### run docker container after build the image
```bash
# example run http serve:
$ docker run -i --name rll-url-shortener -p 8081:8081 -t rll-url-shortener

```


## Running on your local machine

Linux or MacOS

## Installation guide
#### 1. install go version 1.17 or above  
```bash
# please read this link installation guide of go
# https://golang.org/doc/install
```

#### 2. Create directory workspace    
```bash
run command below: 
mkdir $HOME/go
mkdir $HOME/go/src
mkdir $HOME/go/pkg
mkdir $HOME/go/bin
mkdir -p $HOME/go/src/github.com/ralali/rll-url-shortener
chmod -R 775 $HOME/go
cd $HOME/go/src/github.com/ralali
export GOPATH=$HOME/go
```    
```bash
# edit bash profile in $HOME/.bash_profile        
# add below to new line in file .bash_profile         
    PATH=$PATH:$HOME/bin:$HOME/go/bin
    export PATH  
    export GOPATH=$HOME/go 
# run command :
source $HOME/.bash_profile
```

#### 3. Build the application    
```bash
# run command :
cd $HOME/go/src/github.com/ralali/rll-url-shortener
git clone -b development https://github.com/ralali/rll-url-shortener .
cd $HOME/go/src/github.com/ralali/rll-url-shortener
go mod tidy && go mod download && go mod vendor
cp config/app.yaml.tpl config/app.yaml     
# edit config/app.yaml with server environtment
go build

# run application after build or create on supervisord 
./rll-url-shortener serve-http
```


### 4. Health check Route PATH
```bash
{{host}}/in/health
```

### Database Migration
migration up
```bash
make db.migrate.up
```

migration down
```bash
make db.migrate.down
```

migration reset
```bash
make db.migrate.reset
```

migration redo
```bash
make db.migrate.redo
```

migration status
```bash
make db.migrate.status
```

create migration table
```bash
make db.migrate.create name={table_name}

# example
make db.migrate.create name=user_profiles
```

to show all command
```bash
make db.migrate
```

## run docker compose on your local machine
Up containers
```bash
make docker.compose.up
```

Shutdown containers
```bash
make docker.compose.down
```

Restart containers
```bash
make docker.compose.restart
```