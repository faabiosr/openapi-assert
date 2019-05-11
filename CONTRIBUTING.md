# Contributing

By participating to this project, you agree to abide our [code of conduct](/CODE_OF_CONDUCT.md).

## Setup your machine

`openapi-assert` is written in [Go](https://golang.org/).

Prerequisites:

* `make`
* [Go 1.11+](https://golang.org/doc/install)

Clone `openapi-assert` from source into `$GOPATH`:

```sh
$ mkdir -p $GOPATH/src/github.com/faabiosr
$ cd $GOPATH/src/github.com/faabiosr
$ git clone git@github.com:faabiosr/openapi-assert.git
$ cd openapi-assert
```

A good way of making sure everything is all right is running the test suite:
```console
$ make test
```

## Formatting the code
Format the code running:
```console
$ make fmt
```

## Create a commit

Commit messages should be well formatted.

You should give the message a title, starting with uppercase and ending without a dot.
Keep the width of the text at 72 chars.
The title must be followed with a newline, then a more detailed description.

Please reference any GitHub issues on the last line of the commit message (e.g. `See #123`, `Closes #123`, `Fixes #123`).

An example:

```
Add example for --release-notes flag

I added an example to the docs of the `--release-notes` flag to make
the usage more clear.  The example is an realistic use case and might
help others to generate their own changelog.

See #284
```

## Submit a pull request

Push your branch to your `openapi-assert` fork and open a pull request against the
master branch.
