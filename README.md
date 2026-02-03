
<img src="theinfinitelibrary-frontend/src/logotype.png" alt="gopher" width="60%">

# What it is
Social media app for bookworms where book lovers can connect with likeminded readers. Minimal use of AI assistants during devlopment to preserve control of codebase.

**Status**: Beta - Usable, but significant features still under development.

**License**: MIT License (See file named LICENSE).

## Tech stack
* Backend - _**Go**_
* Frontend - _**React & Javascript**_
* Database - _**MySQL**_
* Containerization - _**Docker**_


  <img src="assets/gopher.svg" alt="gopher" width="10%">  &nbsp;&nbsp;&nbsp;  <img src="assets/golang.svg" alt="golang" width="10%"> &nbsp;&nbsp;&nbsp;   <img src="assets/react.svg" alt="react" width="10%">  &nbsp;&nbsp;&nbsp;  <img src="assets/typescript.svg" alt="typescript" width="10%">  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  <img src="assets/mysql.svg" alt="mysql" width="10%">      &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp; <img src="assets/docker.svg" alt="docker" width="10%">


## Features
* User Profile
  * Personal _"library"_, listing books that user has read, with entry points into book chatrooms.
  * One-way Argon2 hashing of passwords for password integrity.
* Bookwise chatrooms to connect with likeminded readers.



### Upcoming Features
* Book clubs with configurable private/public chatrooms and invitation function.
  * Schedule coordinator for physical meetups in book clubs.
* Browsing of books based on genre and popularity.

# How to run theinfinitelibrary

To run the app, execute <br> <br>
`docker compose up` <br> <br> While standing in the root folder. This will:
* Build and run three Docker containers: mysql (DB), backend, and frontend
* Create database and initialize its schemas
* Build and run backend
* Install and run frontend

Then visit the frontend by navigating to http://localhost:3000/


# High level architecture:

The below diagram visualizes the high-level architecture for The Infinite Library. Implementation details are abstracted away, and logical units below do not perfectly map onto request handlers in the app.

<p align="left">
  <img src="docs/HighLevel_Architecture.svg" alt="architecture" width="120%">
</p>

# Chatrooms architecture:

The below diagram visually describes the current implementation of seperate, concurrent, bookwise chatrooms, which leverages the concurrency support in Go with its Gochannels and Goroutines.

<p align="left">
  <img src="docs/Chatrooms_Architecture.svg" alt="chatrooms architecture" width="100%">
</p>

# Frontend routing tree:
The below diagram depicts the user flow in the frontend.
<p align="left">
  <img src="docs/React_router_tree.png" alt="routing" width="50%">
</p>

