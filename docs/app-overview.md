# App Overview

## Backend Overview

The backend is a Spring Boot application which is conceptually dividable into the following areas of responsibility:

* **Controller domain**
    * Handles Network traffic from API users. Exposes API endpoints. 
* **Service domain**
    * Handles business logic tied to APIs.
* **Security domain**
    * Handles authentication and authorization.
* **Entity domain**
    * Defines application objects.
* **Repository domain**
    * Handles database communication.

Partitioning the app into aforementioned areas of responsibility is motivated by the minimization of wasteful resource use in development, testing, and production. The divide also separates concerns during development, which reduces unwanted entanglement between features during app development, adding to the list of benefits.

### Data flow

Data flows thru the backend according to

**Controller**   ---> **Service** ---> **Repository**

where the areas relating to **Security** and **Entities** exert their effects in various parts of the above flow. They are still separate with respect to code repository structure, even if they appear at multiple stages in the above data flow.


### Controller domain

The app exposes the following API's over HTTP:

**Book**
* /api/addBook --- JSON body:
    * title
    * author
* /api/getBooks/{userID}

**User**
* /api/login   --- JSON body:
    * username
    * password
* /api/signup --- JSON body:
    * username
    * password
* /api/getUser/{userID}
* /api/addUser --- JSON body:
    * username
    * password

**Chatroom**
* /api/chatRoom/{chatRoomID}
* /api/postMessage/{chatRoomID} --- JSON body:
    * message, username


### Service domain

The app encompasses the following services, each responsible for handling some part of the business logic:

* BookService
    * addBookToUser(Long userID, String title, String author)
        * look up user
        * create Book entity
        * associate book with user
        * save (using BookRepository)
        * return response DTO
* UserService
    * signup(String username, String password)
        * check if username already exists
        * hash the password
        * Create the User Entity
        * Save it (using UserRepository)
        * Return response DTO or throw exception
    * login(String username, String password)
        * Check if username exists
        * Check password correctness (against stored hashcode)
        * Return respose DTO or throw exception
* ChatService
    * postMessage(Long chatRoomId, String username, String message)
        * Look up chatroom
        * Look up user
        * Create Message entity
        * Save it
        * Broadcast message to all Chatroom participants
        * Return ok or throw exception.


### Repository domain

The following repository classes enable database transactions within the app:

* BookRepository
    * save(Book)
    * retrieve(Book)
* UserRepository
    * save(User)
    * retrieve(user)
* MessageRepository
    * save(Message)
* ChatRepository
    * retrieve(Chat)


### Entities

The app houses the following entities:

* Book
* User
* Message
* Chat