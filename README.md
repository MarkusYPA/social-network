# social-network

## Overview

A familiar style social network that includes key social features such as user profiles, followers, posts, groups, events,real-time chat, and notifications. The goal is to provide a functional and modular full-stack web application with authentication, data persistence, and real-time capabilities.

The frontend is developed using **Vue.js**, and the backend is written in **Go** with an **SQLite** database. Both frontend and backend are containerized with Docker.


## Features

### ✅ User System
- User registration with email, password, name, date of birth
- Optional fields: Avatar, Nickname, About Me
- Login and logout with session persistence (cookies)

### ✅ Profiles
- Public or private profiles
- User info, posts, and follower/following lists
- Toggle privacy setting end edit profile info

### ✅ Posts
- Create posts and comments
- Attach images
- Post visibility: public, followers-only, or selected followers

### ✅ Followers
- Follow/unfollow users
- Follow requests and approval for private profiles

### ✅ Groups
- Create and manage groups
- Group posts and comments
- Invite/request to join groups
- Create events with RSVP options

### ✅ Events
- Group specific events
- Members can choose attendance

### ✅ Chat
- Real-time private messages using WebSockets
- Emoji support
- Group chat rooms

### ✅ Notifications
- Follow request received
- Group invitation received
- Group join request (for group creator)
- Group event created (visible to members)


## Technologies Used

### Frontend
- [Vue.js](https://vuejs.org/)
- HTML, CSS

### Backend
- Go
- SQLite with [golang-migrate](https://github.com/golang-migrate/migrate)
- Gorilla WebSocket
- bcrypt (password hashing)
- uuid (user and object IDs)


## Dockerized Setup

Backend and frontend docker images are built using a docker-compose.yml file and a custom script (./dev.sh). 

### Backend
- Uses Go to build and serve the API
- Applies SQLite migrations on startup

### Frontend
- Built using Vue
- Served using serve package

## Authentication & Sessions
- Passwords are hashed using bcrypt
- Sessions are tracked via cookies and stored in database

## Notes
- WebSockets are used for real-time private and group chat
- Notifications are shown across all pages

## How to Run

### Without dockerizing
In the base directory, in two different terminals:
```
cd backend
go run .
```
```
cd frontend
nmp install
npm run dev
```

### With Docker
In the base directory
```
./dev.sh
```

### Visit site
Frontend: http://localhost:5173  
Backend API: http://localhost:8080

## Contributors
- [Markus Amberla](https://github.com/MarkusYPA)
- [Oleg Balandin](https://github.com/olegamobile)
- [Inka Säävuori](https://github.com/Inkasaa)
- [Mahdi Kheirkhah](https://github.com/mahdikheirkhah)
- [Kateryna Ovsiienko](https://github.com/mavka1207)
- [Fatemeh Kheirkhah](https://github.com/fatemekh78)