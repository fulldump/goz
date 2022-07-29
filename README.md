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

## Download

Get your binaries from GitHub releases, precompiled for:

* Windows
* Linux
* Mac

