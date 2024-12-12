# Forum

A simple forum application built with **Golang**, **SQLite3**, **JavaScript**, **HTML**, and **CSS**. This project allows users to create, read posts, as well as interact with comments, reactions, simulating the behavior of a modern forum.

## Features

- **User Authentication**
  - Register, login, and manage posts with user-specific data.
  
- **Forum Posts**
  - Users can create, read posts.

- **Comments Section**
  - Users can add comments to posts and read their own and others comments.

- **SQLite Database**
  - All posts, comments, and user data are stored in an SQLite3 database.

- **Responsive Frontend**
  - Built with HTML, CSS, and JavaScript for a responsive and user-friendly experience.

## Runing

1. To run this project you have mulitiple choices :
    1. cd Forum. 
   ```bash :
    "go run server.go"
2. ```Docker :
    "docker build -t <Name> ."
    "docker Run -d -p 8080:8080 <Name>"