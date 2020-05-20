# JumpCloud Code Problem

## Instructions

#### Requirements
* Golang version 1.13

#### Running the service
To run the Hashing service:
```go run cmd/main.go <port>```

## Project Layout

```
root/
 ├─ cmd/        Houses entrypoint into the software
 └─ internal/
      ├─ app/   Application specific code, such as construction of service, endpoints, etc
      └─ pkg/   Reusable modules of code, such as string parsing
```

## Design Decisions
* Rather than only providing statistics for `POST`s against the `/hash` endpoint, my service will provide simple metrics for all endpoints
on the service.
    * Metrics captured at the router level, which has some interesting implications
        * Metrics are now available for insights into incorrect endpoint calls
            * Wrong method for given endpoint
            * Misspelled / incorrect endpoints
        * This could provide insight as to how users are trying to use the service that we haven't accounted for

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