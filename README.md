
# GO-REC

Recently added song by your liked artists in Go.

## 1. Running Locally

- Make sure you have
    - [Go](http://golang.org/doc/install) version 1.16 or newer installed.
    - [Docker](https://docs.docker.com/engine/install/)
    - [Docker: Post-installation steps for Linux](https://docs.docker.com/engine/install/linux-postinstall/)
    - [Docker: Compose](https://docs.docker.com/compose/install/)

```sh
$ go mod vendor
$ go mod tidy
$ docker-compose up -d
```

### Environment

Copy .env.sample and update with environment settings\
Some envs values need to be requested to [hup.chat dev team](mailto:dev@hup.chat?subject=ENVS%20[GitHub]%20HUP.CHAT)

```sh
$ cp .env.sample .env
```

### 1.1 HEROKU MODE

Compile(Debug flags: https://github.com/go-delve) all files, listed in Procfile:
```sh
$ go build -gcflags="all=-N -l" -o bin/authenticate -v authcode/authenticate.go
$ go build -gcflags="all=-N -l" -o bin/spotify -v sṕotify.go
$ go build -gcflags="all=-N -l" -o bin/spotifyauth -v auth/auth.go
```

Procfile production:
```sh
web: bin/authenticate
...
```

Heroku local
```sh
$ heroku local
```

### 1.2 IDE MODE

Your app should now be running on [localhost:8080](http://localhost:8080/)

## Ngrok
Install and enable https://ngrok.com/\
Your app should now be running on [https://RANDOM.ngrok.io](https://RANDOM.ngrok.io)\
Request ngrok access(subdomain) when needs: [hup.chat dev team](mailto:dev@hup.chat?subject=ENVS%20[GitHub]%20HUP.CHAT)


## Deploying to Heroku

GitHub actions or

```sh
$ heroku create
$ git push heroku main
$ heroku open
```

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

### Extra documentation

#### Inspect HTTP Requests
- [Hookbin](https://hookbin.com/)
