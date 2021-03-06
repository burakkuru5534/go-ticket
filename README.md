# Ticket Allocation Coding

## Introduction

In this project, it is aimed to allocate tickets. Basically, ticket creation, fetching information about the created ticket and purchasing the ticket are provided by rest APIs.



### Languages and frameworks

Technologies used in this project:

Golang,
postgresql,
Docker,
docker-compose

Test Environments:

postman,
jmeter

### Database

Postgresql was used as the database language.

Tables created:

ticket_options: In this table, the conditions of the tickets created for ticket sales, how many tickets are available for sale, and information about the event are kept.

purchases: In this table, the relations of the tickets that have been sold with the users who made the purchase are kept.

tickets: In this table, the relationship between the tickets sold and the ticket conditions is kept.


## Problem solution

The problem in this project was that the system was able to respond concurrently to incoming requests for the creation of tickets and the sale of the created tickets. We used golang's WaitGroup library to solve this problem. This library has three functions: Add, Done, Wait. When a request arrives, we inform the system that we are currently handling a request with the add function. With the Done function, we indicate that the processes related to this request are finished. The wait function also makes the system wait, which is necessary in the meantime.

### Create Ticket Option

Create Ticket Options request url example:

Method: POST

 http://localhost:8080/ticket_options
 
 request Body Example:
 ```json
 {
     "Name":"Fenerbahce vs Galatasaray",
     "Desc":"There are 10.000 available tickets.",
     "Allocation":10000
 }
 ```
 response example:
 
 200

### Get Ticket Option

Get Ticket Options request url example:

 Method: GET
 
  http://localhost:8080/ticket_options?id=9368ec0b-701a-4be8-a267-b5deb1063128
  
   id: this id should be one of the ticket_options's ids. 

  request Body: 
  
  response example:
  
 ```json
 {
    "ID":"297ac147-85cc-48bb-83e9-736077c22804",
    "Name":"Fenerbahce vs Galatasaray",
    "Desc":"There are 10.000 available tickets.",
    "Allocation":10000,
    "CreatedAt":"2022-01-05T17:41:08.039342Z",
    "UpdatedAt":"2022-01-05T17:41:08.039342Z"
}
  ```



### Purchase from Ticket Option


Purchases from Ticket Options request url example:

Method: POST

 http://localhost:8080/ticket_options/purchases?id=9368ec0b-701a-4be8-a267-b5deb1063128
 
 id: this id should be one of the ticket_options's ids. 
 request Body Example:
 ```json
 {
   "Quantity": 2,
   "UserID": "406c1d05-bbb2-4e94-b183-7d208c2692e1"
 }
 ```
 response example:
 
 200
 
 if quantity greater than allocation then our response:
 
 4xx kind of error and error message: there is not any available tickets.
 
 ### Project Run
 
 This project could build in your local also this is a dockerized project.
 
 $sudo docker build . -t go-ticket-docker
 
 $sudo docker-compose up -d --build
 
 $docker exec -it postgres psql -U postgres
 

### Test

I used postman and jmeter for testing.

You can use postman to create ticket options and get ticket options. But I needed to jmeter to test concurrent requests.
So that's why you can test purchases api with jmeter.

I also used golang's test library to test database connection and table controls.


