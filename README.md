# JumpCloud Code Problem

## Instructions

#### Requirements
* Golang version 1.13

#### Running the service
To run the Hashing service, run the following command from the root of the project:
```go run cmd/main.go <port>```

#### Running the tests
To run the unit tests, run the following from the root of the project:
```go test ./...```

Using Postman:
* located in `/test` is a Postman collection with some simple requests to make calling the service
easier. The tests as committed assume the service is running on port 8088

## Project Layout

```
root/
 ├─ cmd/         Houses the entrypoint(s) into the software
 ├─ internal/
 |     ├─ app/   Application-specific code, such as construction of service, endpoints, etc
 |     └─ pkg/   Reusable modules of code, such as string parsing
 └─ test/        External testing, namely Postman collection
```

## Design Decisions
* Rather than only providing statistics for `POST`s against the `/hash` endpoint, my service will provide simple metrics for all endpoints
on the service.
    * Metrics captured at the router level, which has some interesting implications
        * Metrics are now available for insights into incorrect endpoint calls
            * Wrong method for given endpoint
            * Misspelled / incorrect endpoints
        * This could provide insight as to how users are trying to use the service not yet accounted for
* Hashes stored in-memory, though the service architecture will safely handle flushing to disc on shut-down if a different storage
mechanism were to be introduced.
* HTTP endpoint tests use the HTTP package directly running against an instance of the service

## Tools
* **go mod** - Intended to help portability of repository, currently only specifies golang version due to only using the standard library
* **golint** - Aid in code readability and consistency
* **Postman** - Used to test endpoints
* **https://md5decrypt.net/en/Sha512/** - Used as a third-party verification/sanity check for hashing
* **https://www.base64decode.org/** - Used as a third-party verification/sanity check for Base64 conversions

## Challenges

* Being limited to the standard library has made certain aspects of this project very challenging
    * Historically, I have worked with libraries such as `gorilla/mux`, which handle things like path parameters in a very clean and straightforward way
* Assertion libraries are helpful for making tests more readable, I have rolled a very primitive assertion library for this project.
    * It does not properly report the line the assertion took place on (instead it reports the line the within the assertion function)
* Windows system clock resolution
    * My tests indicate that the timing shown in the `/stats` endpoint are correct, but I dev on a Windows machine.
    * Timings show up as a zero-average, which I believe to be caused by the windows system clock resolution, which is less than that of linux.
It should appear correctly on a linux machine.