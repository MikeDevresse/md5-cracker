[![Create and publish a Docker image](https://github.com/MikeDevresse/md5-cracker/actions/workflows/docker-image.yml/badge.svg?branch=main)](https://github.com/MikeDevresse/md5-cracker/actions/workflows/docker-image.yml)
[![Go build & tests](https://github.com/MikeDevresse/md5-cracker/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/MikeDevresse/md5-cracker/actions/workflows/go.yml)
[![React build & tests](https://github.com/MikeDevresse/md5-cracker/actions/workflows/react.yml/badge.svg?branch=main)](https://github.com/MikeDevresse/md5-cracker/actions/workflows/react.yml)

# MD5 Cracker

- [MD5 Cracker](#md5-cracker)
    * [Usage](#usage)
        + [Docker compose](#docker-compose)
        + [Docker images](#docker-images)
    * [Use the backend](#use-the-backend)
        + [Init](#init)
        + [Requests](#requests)
            - [search](#search)
        + [Responses](#responses)
            - [queue](#queue)
            - [slaves](#slaves)
            - [found](#found)
            - [Wrong hash given](#wrong-hash-given)
            - [Command commandName not found](#command-commandname-not-found)

This repository contains 2 projects: a backend written in Golang, and a frontend written with ReactJS. This project is a school project in order to learn how to use docker and scaling services in docker.
The goal of this project is to have a backend that communicates with "slaves" that will bruteforce in order to find the given md5 hash. In order to do that we split an alphabet in the number of slave connected
and tells to each slave what range of letters he should try.

The [ityt/hash-slave](https://hub.docker.com/r/itytophile/hash-slave) is used for this project has the md5 cracker slave.

## Usage

### Docker compose
A [docker-compose.yaml](https://github.com/MikeDevresse/md5-cracker/blob/main/docker-compose.yaml) file is given for a fast usage. In order to use it, clone the project and run

```docker-compose up -d```

You will then be able to navigate to the frontend page at [127.0.0.1](http://127.0.0.1)

### Docker images
This project has CI that generates two docker images, one for the frontend and one for the backend.

[Backend Image](https://github.com/MikeDevresse/md5-cracker/pkgs/container/md5-cracker-backend)

[Frontend Image](https://github.com/MikeDevresse/md5-cracker/pkgs/container/md5-cracker-frontend)

## Use the backend
In order to use the backend itself you need to connect to the websocket, the port with the image is `80` but it is mapped to `8080` with the docker-compose.

### Init
In order to initiate the client you must first send a message containing `client` in order to identify yourself

### Requests
#### search
```search md5_hash```

This function allows you to add a md5 hash in the queue 

### Responses
#### queue
```queue queue_size```

When the queue size gets updated, the websocket sends a queue message followed with the queue size

#### slaves
```slaves slave_count```

When the slave count gets updated, the websocket sends a slaves message followed by the slave count

#### found
```found md5_hash result```

When you send a search request, when it will be resolved a found response will be given with the md5 hash and the result that have been found 

#### Wrong hash given
```Wrong hash given```

Occurs when the search request has been sent with a second argument that is not a md5 hash

#### Command commandName not found
```Command commandName not found```

Occurs when a given request does not correspond to any known command