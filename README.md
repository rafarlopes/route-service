# Route Services

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
Once the release is approved, another action kicks in and pushes the image to the DockerHub