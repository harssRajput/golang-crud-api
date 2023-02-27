# golang-crud-api
golang MongoDB crud API with all possible test cases. It is compatible with docker.

The main.go is the entry point of the application. It contains as a triggering point to the web-server and as a initiasization before running the server.

# Separation of Concerns, loose couplindg and high cohesion concepts used to make it extensible and scalable in future.
the app is structured in a way so that it can easily extendible later.
ex-
1. models dir contains modelling between Go struct and mongodb Doc. It again contains files corressponding to each model type.
2. app.go contains server logic and all business logic housed in different Methods. when app grows we can divide each method into its separate file easily and again later we can completely divide it into two layers. Controller + Service layers.
3. init.go contains all the initialisations thing. Again, it is divided into separate methods calls for each different service initialisation. i.e. init_mongo.go to initialise mongodb. this also can be extended easily by creating a new service code into init_xyz.go and just call it from init.go
4. app_test.go contains testing for each API and cover all the scenarios. 
5. config.go contains all the configuration of the application.

I was not able to establish connection to mongoDB container. Although, go-app is working fine as a local and docker container.

# Steps to Run
firstly, run mongodb service in a local. It is not able to connect to mongodb docker. This needs to fix. SO, run mongodb in a local
there are two ways.
## first way: run without docker.
1. pull the repo
2. install the docker and run the docker engine
3. go in repo directory
4. docker build -t myapp .       //to build the image
5. docker run -d -p8080:8080 --name myapp myapp        //to run the myapp as a docker container

## second way: run with docker
1. intall go
2. ttake pull of repo
3. go in repo directory
4. go mod download
5. go run .

## mongodb listens at 27017. app listens at 8080
https://user-images.githubusercontent.com/82873133/221620326-3663e57d-b31e-4c9a-965b-68bf1d734535.mp4

