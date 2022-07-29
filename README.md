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

Software released under The MIT License (MIT).

Goz does not guarantee file corruption free.

It is designed to work with small files (does not support streaming).

File access is not concurrent.

Goz does not ensure strong passwords.

## Thanks

This code is 100% based on these articles:

* https://levelup.gitconnected.com/a-short-guide-to-encryption-using-go-da97c928259f
* https://medium.com/@mertkimyonsen/encrypt-a-file-using-go-f1fe3bc7c635

