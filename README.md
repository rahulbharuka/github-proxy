# Github Proxy
This project acts as proxy for Github. It allows valid users to post comments and retrieve list of Github org members.


### Features
- Write/List/Delete comments for a given Github org.
- Retrieve list of **_public_** members of a given Github org.
---

### How to run ?
1. **Directly (without source code)**
   * Download _script_ directory.
   * cd to the directory and run
```
    docker-compose up
```
   * This will fetch docker images from Docker Hub.

2. **Using source code.**
   * Download the complete source code.
   * cd to the directory and run
```
    docker-compose up
```


Please note that by default, 
- _comment-app_ runs on port 6060.
- _member-app_ runs on port 7070.
- Refer _github-proxy.postman_collection_.
---

### Assumptions
- All operations are validated against only `public` membership of a user for given Github org. 
- User can post a comment only if he/she is a `public` member of given Github org.
- NOTE: _Private_ org membership is not considered because private membership can only be queried by existing member of the org.
---

### Architecture
This project will spin three containers:
- `db` - Hosts PostgreSQL DB that stores user comments against Github orgs.
- `comment-app` - Hosts a Golang service that serves comments related HTTP APIs.
- `member-app` - Hosts a Golang service that serves member retrieval HTTP API. 
---

### APIs
1. `POST /orgs/:org/comments`
  * Usage: To post comments against a given Github org.
  * Calls Github v3 API to validate that the user is a public member of given Github org.

```
    Request body:
    {
	    "author": "<user-name>",
	    "comment": "<comment>"
    }

    HTTP Response:
    200 - if comment is added successfully.
    400 - if request format is not correct.
    404 - if user is not a public member of given Github org.
    500 - if some error occured while validating user membership or saving comment in DB.
```  
2. `GET /orgs/:org/comments`
  * Usage: To retrieve list of all comments for given Github org.
  * Calls Github v3 API to validate Github org.
```
    HTTP Response:
    200 - on successful retrieval of all comments for given Github org.
    404 - if the given org does not exist on Github.
    500 - if some error occured while validating Github org or retrieving comments from DB.
```
3. `DELETE /orgs/:org/comments`
 * Usage: To (soft) delete all comments for given Github org.
 * Calls Github v3 API to validate Github org.

```
    HTTP Response:
    200 - on successful (soft) deletion of all comments against given Github org.
    204 - if no active comments exist for given Github org.
    404 - if the given org does not exist on Github.
    500 - if some error occured while validating Github org or deleting comments in DB.
```
4. `GET /orgs/:org/members`
  * Usage: To return a list of public members of a given Github org.
  * Calls Github v3 API to fetch list of public members of given Github org and then to fetch details of each user.
  * Respone is sorted in descending order of number of followers.

```
    HTTP Response:
    200 - on successful retrieval of all `public` members of given Github org.
    500 - if some error occured while retrieving Github org members.
```
---

### External Dependencies
- Uses `PostgreSQL` as a persistent data storage layer.
- Uses `go-pg` for creating PostgreSQL client and ORM.
- Uses `go-github` client library for accessing the GitHub API v3.
- Uses `Gin` web framework for routing.
---

### How to build Golang app binary for linux
- cd to app directory with `main.go` and run


    env GOOS=linux GOARCH=amd64 go build -o <binary-name> -v .
---

### Testing
- Added unit tests for logic layer.
