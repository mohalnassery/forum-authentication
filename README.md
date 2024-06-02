# Forum Application

This directory contains the code for a web-based forum application built using Go and HTML, CSS, JavaScript. The application allows users to create and participate in discussions on various topics.

## Overview

The application consists of two main components:

1. **Server**: A Go-based backend server that handles the application logic, database interactions, and API endpoints.
2. **Client**: A frontend client that provides the user interface for interacting with the forum.

## Server

The server component is responsible for managing the forum data and exposing a RESTful API for the client to interact with.

### Key Files

- `main.go`: The entry point of the server application.
- `routes/`: This directory contains the HTTP handlers for handling different API routes.
- `models/`: This directory defines the data models for the application, such as `User`, `comments`, `Post`, .
- `database/`: This directory contains the code for interacting with the SQLite database.

### Functionality

- **User Authentication**: Users can register, log in, and log out of the application.
- **Thread Management**: Users can create new threads, view existing threads, and post replies to threads.
- **Post Management**: Users can create new posts and delete their own posts.

## Client

The client component provides the user interface for interacting with the forum application.

### Key Files

- `client/index.html`: The main entry point of the application.
- `client/` contains all the pages of the application.
- `client/js`: This directory contains all front-end logic for the application.
- `client/css`: This directory contains the styling for the application.
- `client/assets` contains the images used in the application.

### Functionality

- **Thread Listing**: Users can view a list of all available threads.
- **Thread Details**: Users can view the details of a specific thread, including all posts and replies.
- **Post Creation**: Authenticated users can create new posts within a thread.
- **Post Deletion**: Authenticated users can delete their own posts.

## Setup and Installation

To run the application locally using go or docker, follow these steps:

1. Clone the repository: `git clone https://learn.reboot01.com/git/halmakan/forum.git`
2. Navigate to the project directory: `cd forum/server`
3. Run the go server: `go run main.go`
4. Navigate to `http://localhost:8080` to access the application.

### Docker
1. Navigate to the root directory of the project: `cd forum`
2. Build the Docker images: `docker-compose build`
3. Start the Docker containers: `docker-compose up`
4. Open your web browser and visit `http://localhost:8080` to access the application.

The `docker-compose.yml` file in the project root directory defines the services for the application:

- **app**: The Go server application, which is built from the Dockerfile in the project root.

The Dockerfile builds the Go server application and copies the client code into the Docker image.

