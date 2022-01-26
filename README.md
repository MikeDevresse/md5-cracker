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
            - [stop-all](#stop-all)
            - [auto-scale](#auto-scale)
            - [max-search](#max-search)
            - [max-slaves-per-request](#max-slaves-per-request)
            - [slaves](#slaves)
        + [Responses](#responses)
            - [queue](#queue)
            - [slaves](#slaves-1)
            - [max-search](#max-search-1)
            - [max-slaves-per-request](#max-slaves-per-request-1)
            - [auto-scale](#auto-scale-1)
            - [found](#found)
            - [Command commandName with numberOfArgs arguments not found](#command-commandname-with-numberofargs-arguments-not-found)

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

Returns:
 - `Added Hash $hash to queue` The hash has been added to the queue
 - `Error Please provide a valid md5 hash` The provided hash is not a valid md5 hash

#### stop-all
```stop-all```

Stop all the slaves from working, empty the queue

#### auto-scale
```auto-scale true|false```

Tells if the server should scale the slaves automatically depending on the queue or not

#### max-search
```max-search limitForSearching```

Sets the limit to which the word can be possible, must follow the regex `^[9]{2,8}$`

Returns:
 - `Error max-search argument must follow the regex: ^[9]{2,8}$` The provided limit does not follow the required regex

#### max-slaves-per-request
```max-slaves-per-request numberOfSlaves```

Tells how many slaves at the max should be working a search request
for instance if we have 16 slaves but this option to 4, only 4 slaves will be working on the next search request.
But if we have only 2 slaves available and this option to 4, the 2 slaves available will work on the search request

Returns:
 - `Error max-slaves-per-request expects a number as second parameter` The given parameter was not a number
 - `Error max-slaves-per-request must be greater than 0` The given parameter is less or equal than 0

#### slaves
```slaves numberOfSlaves```

Allows you to scale the number of slaves that are working from 1 to 16

Returns:
- `Scaling` The application is getting scaled
- `Error An error occurred while trying to scale the application.` An unknown error occurred while scaling the application
- `Error slaves expects a number as second parameter` The given parameter was not a number
- `Error slaves must be between 1 and 16` The given parameter is not between 1 and 16


### Responses
#### queue
```queue queue_size search_request_being_handled```

When the queue size or the number of request being handled gets updated, the websocket sends a queue message followed with the queue size and the number of search request being handled

#### slaves
```slaves slave_count available_slaves working_slaves```

When the slave count gets updated, the websocket sends a slaves message followed by the slave count, the number of slaves not working and the number of slaves working

#### max-search
```max-search max_search_value```

When configuration is updated, the server sends the current state

#### max-slaves-per-request
```max-slaves-per-request number_of_slaves_per_request```

When configuration is updated, the server sends the current state

#### auto-scale
```auto-scale true|false```

When configuration is updated, the server sends the current state

#### found
```found md5_hash result```

When you send a search request, when it will be resolved a found response will be given with the md5 hash and the result that have been found

#### Command commandName with numberOfArgs arguments not found
```Command "commandName" with numberOfArgs arguments not found```

Occurs when a given request does not correspond to any known command