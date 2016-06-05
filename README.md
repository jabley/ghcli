[![Build Status](https://travis-ci.org/jabley/ghcli.svg?branch=master)](https://travis-ci.org/jabley/ghcli)

A command-line client for interacting with the Github API.

## Installatation
```shell
$ go get github.com/jabley/ghcli
```

## Usage

```shell
$ GH_OAUTH_TOKEN=a-token-here ghcli members -o alphagov
```

This will show all of the members of the organisation alphagov.
