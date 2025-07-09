# Fender Digital Platform Engineering Challenge

## Description

Design and implement a RESTful web service to facilitate a user authentication system. The authentication mechanism should be *token based*. Requests and responses should be in **JSON**.

## Requirements

**Models**

The **User** model should have the following properties (at minimum):

1. name
2. email
3. password

You should determine what, *if any*, additional models you will need.

**Endpoints**

All of these endpoints should be written from a user's perspective.

1. **User** Registration
2. Login (*token based*) - should return a token, given *valid* credentials
3. Logout - logs a user out
4. Update a **User**'s Information
5. Delete a **User**

**README**

Please include:
- A readme file that explains your thinking
- How to set up and run the project
- If you choose to use a database, include instructions on how to set that up
- If you have tests, include instructions on how to run them
- A description of what enhancements you might make if you had more time.

**Additional Info**

- We expect this project to take a few hours to complete
- You can use Python, Go, Node.js, or shiny-new-framework X, as long as you tell us why you chose it and how it was a good fit for the challenge. 
- You can use whichever database you'd like. 
- Bonus points for security, specs, etc. 
- Do as little or as much as you like.

Please fork this repo and commit your code to it. Then, you can show your work and process through those commits.

**How to Run**

- Setup Database:
  - docker-compose up -d. This will create the docker container for a postgres database
- Run Http Server:
  - set same environment variables as the ones inside docker-compose file for database
  - include a logging level environment variable
  - move to cmd/api execute go mod tidy and then go run main.go

**Thinking Process**
- How could I provide authentication? investigate how to use JWT
- Decide to use a database. why? I think update and delete make sense from a storage perspective in terms of resource (REST)
  - Decide to use gorm why? flexibility to make the database querys
- Project Structure decide to set all as internal packages.
- I do create annotations for all the endpoints.
- I decide to use Gin webframework faster development for me.

**How do I test it**
- postman collection included

**Enhacements**
- generate openapi spec
- better logging to include where the message happens
- database connection pool suitable for production environment
- address gin warning messages for production
- mechanism to enable a user I could use the email I am collecting.



