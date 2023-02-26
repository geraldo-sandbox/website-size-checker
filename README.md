# Website Size Checker

Website Size Checker is a small utility for checking the response size of a request to a http/https website.

```bash
$ make build
$ ./build/bin/wsc -v -t 30 -c 10 https://www.1.de/ https://www.2.com/ https://www.3.de/ https://www.4.de https://www.5.de/ https://www.6.de/ https://www.7.de/ https://www.8.de/ https://www.9.de/ https://www.10.net/
https://www.1.com/ 14918 bytes
https://www.2.de/ 357742 bytes
https://www.3.de/ 491103 bytes
https://www.4.de/ 645932 bytes
https://www.5.de/ 866304 bytes
https://www.6.de/ 915400 bytes
https://www.7.de/ 997149 bytes
https://www.8.de/ 1021858 bytes
https://www.9.net/ 1183526 bytes
https://www.10.de/ 1386124 bytes
```

# Requirements

You have 3 options for building the project:

* Using Docker: in case you want to build and use it as a container. In that case, Go or Make are optional. For more information about how to install docker head to https://docs.docker.com/get-docker/.

```bash
cd <root folder of the project>
docker build . t wsc:latest
```

* Go: required for build the binary. The binary below is optimized with `-ldflags="-w -s"`, for more information about the linker head to https://pkg.go.dev/cmd/link. For more information about hot to install Go head to https://go.dev/doc/install. 

```bash
cd <root folder of the project>
go build -a -tags netgo -ldflags="-w -s" -o build/bin/wsc ./cmd/cli
``` 

* Go and Make: make is just an optional tool. For installing `make` on Ubuntu based systems https://linuxhint.com/install-make-ubuntu/. In MacOs with Homebrew: https://formulae.brew.sh/formula/make

```bash
cd <root folder of the project>
make build
``` 


# Usage

In this section you can find the description of the optional parameters and its basic usage. As follows: 

* Using no optional parameters, just the URLs for checking
```bash
# Using all default values
$ ./build/bin/wsc https://www.1.de/ https://www.2.com/ https://www.3.de/ https://www.4.de https://www.5.de/ https://www.6.de/ https://www.7.de/ https://www.8.de/ https://www.9.de/ https://www.10.net/
https://www.1.com/ 14918 bytes
https://www.2.de/ 357742 bytes
https://www.3.de/ 491103 bytes
https://www.4.de/ 645932 bytes
https://www.5.de/ 866304 bytes
https://www.6.de/ 915400 bytes
https://www.7.de/ 997149 bytes
https://www.8.de/ 1021858 bytes
https://www.9.net/ 1183526 bytes
https://www.10.de/ 1386124 bytes
```

* `-c` - integer - number of maximum of concurrent requests (default 1). Valid values are integers [1, 100].
```bash
# Process 10 concurrent requests with -c 10
$ ./build/bin/wsc -c 10 https://www.1.de/ https://www.2.com/ https://www.3.de/ https://www.4.de https://www.5.de/ https://www.6.de/ https://www.7.de/ https://www.8.de/ https://www.9.de/ https://www.10.net/
https://www.1.com/ 14918 bytes
https://www.2.de/ 357742 bytes
https://www.3.de/ 491103 bytes
https://www.4.de/ 645932 bytes
https://www.5.de/ 866304 bytes
https://www.6.de/ 915400 bytes
https://www.7.de/ 997149 bytes
https://www.8.de/ 1021858 bytes
https://www.9.net/ 1183526 bytes
https://www.10.de/ 1386124 bytes
```

* `-m` - string - http method to make the HTTP call (default "GET"). Valid values are [GET, HEAD, POST, PUT, PATCH, DELETE, CONNECT, OPTIONS, TRACE]
```bash
# Using POST http method
$ ./build/bin/wsc -m POST https://www.123.de/ https://www.123.com/
https://www.123.com/ 14918 bytes
https://www.123.de/ 357742 bytes
```

* `-t` - integer - request timeout for the http request in seconds (default 30). If a request fail to perform the byte size of the request is zero.
```bash
# Using custom request timeout, after 60 seconds with no response the timeout error is returned
$ ./build/bin/wsc -t 1 https://www.123.de/ https://www.123.com/
https://www.123.com/ 0 bytes  # this one failed, no error shown, no verbose flag
https://www.123.de/ 357742 bytes
```

* `-v` - if provided, enables the verbose mode, it outputs request errors.
```bash
# If -v param is provided the request error appear into the command line response
$ ./build/bin/wsc -v -t 1 https://www.123.com/
https://www.123.com/ 0 bytes (Post "https://www.123.com/": context deadline exceeded (Client.Timeout exceeded while awaiting headers)) # this one failed, with request error
```
