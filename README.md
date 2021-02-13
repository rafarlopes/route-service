# Route Service

Route service uses the ORSM API (http://project-osrm.org) to retrieve distance and time for the given routes based on Latitude and Longitude. Once we receive the routes, they are sorted by Time and Distance.

## Local env

For convenience, there's a Makefile available with actions to run tests, formatting, linting, security checks, imports check and create the binary.

To setup the pre-requisites, please run:
````bash
$ make setup
````

You can find the binary under the folder ```dist```, which is created after running make. The binary is named route-service
The default port is ```8080```, but you can easily switch using the flag ```-p 8000```. You can find the example below

````bash
$ ./dist/route-service -p 8000
````


## Dockerfile

In case you prefer to use docker, there's a Dockerfile which is a multistage docker build where first we build our solution using the Makefile and then we copy the binary under the ```/app``` folder.
An already built image is available on DockerHub if you would like to download:

````bash
$ docker pull hub.docker.com/rafarlopes/route-service
````

To run the image on your local environment:
````bash
$ docker run --rm -d -p 8080:8080 hub.docker.com/rafarlopes/route-service
````

To build the image on your computer, please use (in project root directory):
````bash
$ docker build -t route-service .
````

To  run the image you just built on your local environment:
````bash
$ docker run --rm -d -p 8080:8080 route-service
````

## Release

In order to create a new release, there's a make command for it. The usage is below.
It creates the tag and push to the origin main.
````bash
$ make release TAG_NAME=v0.1
````

## GitHub Actions
We have an action that will be build and test every single commit and pull request done against our Repository.
Also it's include another action to create a draft release and attach the binary every time a tag is pushed into our Repository.
After creating the release, another action kicks in and pushes the image to the DockerHub

## Design decisions

### Internal vs Pkg
I was considering that the package osrm could possible be located as pkg instead of internal.
After giving it some thought, I decided to keep in the internal. Since we have this as a service already, I believe having other apps referencing to the osrm package hardly makes sense.
The route package, I believe the best place is indeed internal because it's pretty much the API implementation itself and I don't believe this is also suitable for being reutilized somewhere else.

### Swagger
I was considering and reading about it, but for the sake simplicity I decided for not including any swagger/openapi definition here.

### Simplicity
My choice was to use mainly the standard library, with only 2 small exceptions for the ```golang.org/x/sync/errgroup``` and ```github.com/pkg/errors```.
The reason is those 2 packages are pretty much close to the standard library and the usage of it really made sense without adding complexity.
Regarding other frameworks and libraries, my preference was to avoid it and stick to the standard library. Basically because I believe if you are able to do deliver the code using only the standard library, you are also able to do it using other frameworks and libraries. So for the purpose of this task and to demonstrate my knowledge, I opted in this way.

### Tests
The tests are mainly integration tests, I didn't mock the http request which is calling the OSRM API.
I was considering maybe to create a struct that would have a HTTP Client and I could create it in a way that would be possible to override in the tests for a stub.
Again, due the simplicity I decided to keep it in the way it's right now.

### Middleware and logging
When it comes to logging, I'm logging a lot of information. I have logs in only a few places.
I could have added a middleware to log the incoming request, but since we only have 1 route it was just simpler to do in the RouteHandler.
This could be refactored later to move into a middleware. Same goes to another frameworks like Echo or Gorilla Mux. If this was a real production app, probably it would be better to use some more robust. Again, for simplicity I sticked to the standard library.
