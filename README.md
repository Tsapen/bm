# Book Management System

The Book Management System is a simple command-line tool for managing books, collections, and their associations. It provides functionality to create, update, delete books and collections, as well as create and delete associations between books and collections.

## Table of Contents
- [Details](#details)
- [Prerequisites](#prerequisites)
- [Usage](#usage)
  - [Common Instructions](#common-instructions)
  - [Server](#server)
  - [Client](#client)
- [Contributing](#contributing)
- [License](#license)

## Details
BM is made of a server process that offers a REST API over both HTTP and a local UNIX socket.

## Prerequisites

Before running the commands, make sure you have the following installed:

- [Docker](https://www.docker.com/)

### Quick Start
1. Clone the repository:

    ```shell
    git clone https://github.com/Tsapen/bm.git
    ```
   
2. Navigate to the project directory:

    ```shell
    cd bm
    ```

3. Run tests:

    ```shell
    make test
    ```

3. Run the Docker container:
    ```shell
    make run-server
    ```
## Book Commands
### Create a book:
Using cli-server:
```shell
make create-book TITLE="The Great Gatsby" AUTHOR="F. Scott Fitzgerald" GENRE="Classic" PUBLISHED_DATE="2023-01-01"
```
or using http-server:
```shell
curl -X POST -H "Content-Type: application/json" -d '{"title":"The Great Gatsby","author":"F. Scott Fitzgerald","genre":"Classic","published_date":"2023-01-01"}' http://localhost:8080/api/v1/book

```
- TITLE (string, required): The title of the book.
- AUTHOR (string, required): The author of the book.
- GENRE (string, required): The genre of the book.
- PUBLISHED_DATE (date, optional): The published date of the book.
- Edition (string, optional): The edition of the book.
- Description (string, optional): The description of the book.  

### Get books:
Using cli-server:
```shell
make get-books TITLE='The Great Gatsby' AUTHOR='F. Scott Fitzgerald' GENRE='Classic' START_DATE='2020-01-01' ORDER_BY=title DESC=true PAGE=1 PAGE_SIZE=10 
```
or using http-server:
```shell
curl -X GET -H "Content-Type: application/json" 'http://localhost:8080/api/v1/books?order_by=title&desc=true&start_date=2020-01-01&page=1&page_size=10'
```
- ID (int64, optional): The id of book to retrieve.
- AUTHOR (string, optional): The author of the book to retrieve.
- GENRE (string, optional): The genre of the book to retrieve.
- COLLECTION_ID (int64, optional): The collection id of the book to retrieve.
- START_DATE (date, optional): The earliest possible published date of the book.
- FINISH_DATE (date, optional): The latest possible published date of the book.
- ORDER_BY (string optional): The field to order collections by: id|title|author|genre|published_date|edition
- DESC (bool, optional): Set to true for descending order.
- PAGE (int64, optional): The page number to retrieve.
- PAGE_SIZE (int64, optional): The number of collections per page, default 50.  

### Update a book:
Using cli-server:
```shell
make update-book ID=2 TITLE="Updated Title" AUTHOR="Updated Author" GENRE="Updated Genre"
```
or using http-server:
```shell
curl -X PUT -H "Content-Type: application/json" -d '{"id":2,"title":"Updated Title","author":"Updated Author","genre":"Updated Genre"}' http://localhost:8080/api/v1/book
```
- ID (int64, required): The id of the book to update.
- TITLE (string, optional): The updated title of the book.
- AUTHOR (string, optional): The updated author of the book.
- GENRE (string, optional): The updated genre of the book.
### Delete a book:
Using cli-server:
```shell
make delete-books IDS='1'
```
or using http-server:
```shell
curl -X DELETE -H "Content-Type: application/json" -d '{"ids":[2]}' http://localhost:8080/api/v1/books
```
- IDS (string, required): The IDs of the books to delete.
## Collection Commands
### Create a collection:
Using cli-server:
```shell
make create-collection NAME="New Collection" DESCRIPTION="A new collection"
```
or using http-server:
```shell
curl -X POST -H "Content-Type: application/json" -d '{"name":"New Collection","description":"A new collection"}' http://localhost:8080/api/v1/collection
```
- NAME (string, required): The name of the new collection.
- DESCRIPTION (string, optional): The description of the new collection.
### Get collections:
Using cli-server:
```shell
make get-collections IDS=1,2,3 ORDER_BY=name DESC=true PAGE=1 PAGE_SIZE=10
```
or using http-server:
```shell
curl -X GET -H "Content-Type: application/json" -d '{"ids":[1,2,3],"order_by":"name","desc":true,"page":1,"page_size"10}' http://localhost:8080/api/v1/collections
```
- IDS (string, optional): Ids of collections to retrieve.
- ORDER_BY (string, optional): The field to order collections by.
- DESC (bool, optional): Set to true for descending order.
- PAGE (int64, optional): The page number to retrieve.
- PAGE_SIZE (int64, optional): The number of collections per page.
### Update a collection:
Using cli-server:
```shell
make update-collection ID=1 NAME="Updated Collection" DESCRIPTION="An updated collection"
```
or using http-server:
```shell
curl -X PUT -H "Content-Type: application/json" -d '{"id": 1,"name":"Updated Collection","description":"An updated collection"}' http://localhost:8080/api/v1/collection
```
- ID (int64, required): The ID of the collection to update.
- NAME (string, optional): The updated name of the collection.
- DESCRIPTION (string, optional): The updated description of the collection.
### Delete a collection:
Using cli-server:
```shell
make delete-collection ID=1
```
or using http-server:
```
curl -X DELETE -d '{"collection_id": 1}'  http://localhost:8080/api/v1/collection
```
- ID (int64, required): The id of the collection to delete.
### Books-Collection Commands
Create a books-collection association:
Using cli-server:
```shell
make create-books-collection CID=1 BOOKS_ID='3,4'
```
or using http-server:
```shell
curl -X POST -H "Content-Type: application/json" -d '{"collection_id":2,"book_ids":[3,4]}' http://localhost:8080/api/v1/collection/books
```
- CID (int64, required): The id of the collection to associate books with.
- BOOK_IDS (string, required): Ids of books to associate with the collection.
### Delete a books-collection association:
Using cli-server:
```shell
make delete-books-collection CID=1 BOOKS_ID='3,4'
```
or using http-server:
```shell
curl -X DELETE -H "Content-Type: application/json" -d '{"collection_id":2,"books_ids":[3,4]}' http://localhost:8080/api/v1/collection/books
```
- CID (int64, required): The collection id to disassociate books from.
- BOOK_IDS (string, required): Ids of books to disassociate from the collection.

## HTTP Client
The Book Management System also provides an HTTP client for interacting with the API. You can use the client to make requests and receive responses programmatically.

### Example:
```shell
go get github.com/Tsapen/bm/pkg/http-client
go get github.com/Tsapen/bm/pkg/api
```

Use the API client to make requests:
```go
	client := httpclient.New(&httpclient.Config{
		Address: "runned server address",
		Timeout: 5 * time.Second,
	})

	req := &api.CreateBookReq{
		Title:         "title",
		Author:        "author",
		PublishedDate: time.Now(),
		Edition:       "edition",
		Description:   "description",
		Genre:         "genre",
	}

	resp, err := client.CreateBook(ctx, req)
```

