[![Build Status](https://travis-ci.org/jabley/ghcli.svg?branch=master)](https://travis-ci.org/jabley/ghcli)

A command-line client for interacting with the Github API.

## Installation
```shell
$ go get github.com/jabley/ghcli
```

## Usage

```shell
$ GH_OAUTH_TOKEN=a-token-here ghcli members -o alphagov
```

This will show all of the members of the organisation alphagov.

```shell
$ GH_OAUTH_TOKEN=a-token-here ghcli members add -o alphagov -u jabley
```

This will add the user jabley to the organisation alphagov.

```shell
$ GH_OAUTH_TOKEN=a-token-here ghcli members remove -o alphagov -u jabley
```

This will remove the user jabley from the organisation alphagov.


```shell
$ GH_OAUTH_TOKEN=a-token-here ghcli teams -o alphagov
```

This will show all of the teams of the organisation alphagov.
