# What it is
Social media app for bookworms where book lovers can connect with likeminded readers. Minimal use of AI assistants during devlopment to preserve control of codebase.

## Tech stack
* Backend - _**Go**_
* Frontend - _**React & Javascript**_
* Database - _**MySQL**_
* Containerization - _**Docker**_

<p align="left">
  <img src="assets/gopher" alt="gopher" width="50%">
</p>

## Features
* User Profile
  * Personal _"library"_ of books read.
  * One-way Argon2 hashing of passwords for password integrity.
* Bookwise chatrooms to connect with likeminded readers.



### Upcoming Features
* Book clubs with configurable private/public chatrooms and invitation function.
  * Schedule coordinator for physical meetups in book clubs.
* Browsing of books based on genre and popularity.

# How to run it

To run the app, execute <br> <br>
`docker compose up` <br> <br> While standing in the root folder. This will:
* Build and run three Docker containers: mysql (DB), backend, and frontend
* Create database and initialize its schemas
* Build and run backend
* Install and run frontend

Then visit the frontend by navigating to http://localhost:3000/

# Frontend routing tree:

<p align="left">
  <img src="React_router_tree.png" alt="routing" width="50%">
</p>

