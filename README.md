# GOZ

Goz is a tool to cypher/decyper files in place.

AES256 is the algorithm election for this implementation.

The binary is statically linked, no OS dependencies.

## How to use

Encrypt all files inside a directory (recursive):

```sh
goz --dir ./mydata/
```

Decrypt all files inside a directory (recursive):

```sh
goz --dir ./mydata/ --open
```

To operate with a single file use `--file` instead of `--dir`.

## Download

Get your binaries from GitHub releases, precompiled for:

* Windows
* Linux
* Mac

<!-- TODO: put links -->

## Disclaimer and drawbacks

Goz does not guarantee file corruption free.

It is designed to work with small files.

File access is not concurrent.

Goz does not ensure strong passwords.

