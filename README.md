# Forum


This is a web forum project created to enable communication between users through posts and comments. The forum allows users to associate categories to posts, like and dislike posts and comments, and filter posts based on categories, created posts, and liked posts.

## Objectives

The main objectives of this project are:

    Enable communication between users through posts and comments
    Allow associating categories to posts
    Enable liking and disliking of posts and comments
    Implement a filter mechanism for posts based on categories, created posts, and liked posts

Technologies Used

The following technologies were used in this project:

    Go programming language
    SQLite database library
    bcrypt and UUID packages
    Docker

Functionality

User Registration and Authentication

The forum allows users to register as new users by inputting their credentials, including their email, username, and password. The password is encrypted before storing it in the database. Upon registration, the user is logged in automatically and a session is created that allows them to access the forum.

To ensure that each user has only one opened session, cookies are used to store the session data. Each session contains an expiration date, and the cookie stays "alive" for a predetermined length of time.

Communication

Only registered users are able to create posts and comments. When a registered user creates a post, they can associate one or more categories to it. The implementation and choice of categories is left up to the developer.

Posts and comments are visible to all users, whether registered or not. Non-registered users are only able to view posts and comments.

Likes and Dislikes

Only registered users are able to like or dislike posts and comments. The number of likes and dislikes is visible to all users.

Filter

The forum includes a filtering mechanism that allows users to filter the displayed posts by categories, created posts, and liked posts. Filtering by categories is akin to subforums, where each subforum is a section of the online forum dedicated to a specific topic.

Usage

To run the project, you need to have Docker installed on your machine. Once you have Docker installed, you can clone the repository and run the following command in the project directory:

### Instructions


- Open download folder in terminal
- Type in terminal: 
- - `make build .` - Create image
- - `make d-run` - Run image
- - `make d-id` - View container ID
- - `docker start "CONTAINER ID"` - Start container
- go to link "http://localhost:8080"
- - `docker stop "CONTAINER ID"` - Stop container






Conclusion

This project provides a great opportunity to learn about web development, including HTML, HTTP, sessions and cookies, and encryption. It also allows you to learn about containerizing an application using Docker and using SQL to manipulate databases.


created by @rakhmeto, @dizdibay, @bshayakh