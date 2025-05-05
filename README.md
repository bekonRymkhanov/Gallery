# GalleryApp: Gallery Application

![GalleryApp Banner](/media/photo_2024-04-21_05-35-08.jpg)

## Overview

GalleryApp is a gallery application that allows users to browse, search, and interact with a vast collection of photographs. Built with a tech stack (Golang backend, Angular frontend), GalleryApp delivers a seamless and responsive experience for photography enthusiasts.

**Live Demo:** [http://4.247.165.79](http://4.247.165.79)

**Video Demo:** [Video](https://youtu.be/C4wKav5CzXY)

## Features

- **Gallery Interface**: Browse photos with a responsive grid layout
- **Search Capabilities**: Search photos by keywords, tags, or categories
- **Interactive Elements**: Like, Rate and comment on favorite photographs
- **User Authentication**: Secure login and registration
- **Design**: Optimal viewing on all device sizes

## Technology Stack

### Backend
- **Golang**: High-performance, concurrent tasks with goroutines and user friendly
- **Postgres**: The best relational database for my opinion, because of robust queries and for structured data storage
- **RESTful API**: Clean architecture for client-server communication

### Frontend
- **Angular**: Component based framework for building dynamic SPA and easiness for multi-platform usability ( you can test my application from tablet or even phone)
- **SCSS**: Simple and understandable styling tool

### Deployment
- **NGINX**: High-performance web server and reverse proxy
- **Azure**: Cloud hosting platform for reliable deployment
- **Docker**: Containerization for consistent deployment environments

## Installation and Setup

### Prerequisites
- Docker & Docker Compose (optional)
- make sure that port :80 is free

### Setup

```bash
git clone https://github.com/bekonRymkhanov/Gallery.git
cd Gallery
```
### Environment
 make ``.env`` file and fill it with your data
```bash
touch .env
nano .env
```
```env
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=isomething
POSTGRES_URL=postgres://user:password@postgres:5432/isomthing?sslmode=disable
POSTGRES_HOST=postgres
POSTGRES_PORT=5432

NODE_ENV=development
PORT=8080
```


### Start
```bash
docker compose up -d
```

## Process

### Planning Phase
The development of GalleryApp began with a planning phase where I mapped out the core functionalities and user experience. I established a clear architecture for both frontend and backend components. With postman testing app I predefined all API witch my system would have in the future

### Technology Choices
I selected Golang for the backend due to its performance characteristics, strong concurrency model, and excellent support for building RESTful services. For the frontend, Angular has been chosen for its battaries-included framework features, TypeScript integration, and component-based architecture which allows for clean code organization and reusability.
`In the future ts will be compiled in Golang language witch will increase its speed for nearly 10x`


The decision to use PostgreSQL was driven by the need for a reliable relational database with strong JSON support, which is ideal for storing varied metadata associated with photos.

### My Development Approach
I followed an iterative development process:
1. Built the core API endpoints and database schema
2. Implemented basic authentication and user management
3. Developed the gallery browsing experience
4. Added search functionality with both text and image-based queries
5. Implemented social features (likes, comments, ratings)
6. Though that adding ai agents and AI generative image needs more time on testing and preparing api's, Maybe, even requiring my oun trained model
7. I could've use open resoure api's for for photo retrival, but it will make any extentions to the code harder in the future and CRUD operations becomes impossible, so I decided to use my own photo CRUD.

## Technical trade-offs

### SQL vs NoSQL Database
even if noSQL solutions like MongoDB might offer more flexibility, I chose PostgreSQL for its:
- ACID and data integrity guarantees
- Strong JSON support
- Better support for complex queries and relationships
- Integration with full-text search capabilities

The benefits in data consistency and query performance were worth it for this application.

### Monolith vs Microservices
I choose a monolithic architecture rather than microservices because:
- It simplified the development process for a solo developer
- Reduced operational complexity during the initial launch`I had an only 2 days, witch is dealbreaker`
- Provided a solid foundation that can be broken into microservices later if needed

The application is designed of concerns, making future transition to microservices possible if scaling requires it.

### Angular vs React
Angular was selected over React despite React's larger community because:
- Opinionated, batteries-included
- TypeScript integration was seamless
- The very easy to understand ecosystem (routing, forms, HTTP client) reduced the need for third-party libraries

The bad side is a steeper learning curve and slightly more boilerplate code, but the benefits in terms of maintainability and structure were valuable for this project.

## Known issues & limitations

### Current Limitations
- **Search performance**: Very large image searches may experience latency on the free tier of Azure hosting
- **Browser compatibility**: Some advanced visual effects may not appear correctly in older browsers

### Future improvements
- Implement Websocket for real time comment notifications
- Implement redis cashing
- Integrate with social media platforms for easy sharing
- Implement AI-powered services

## Performance considerations

GalleryApp was optimized for performance in several ways:
- Backend API response times kept under 100ms for most operations, only photos list can exit it
- Frontend bundle size minimized through lazy loading and code splitting, and used html cashing for faster load of pages
- Database queries optimized with appropriate indexes


## License

This project is licensed under the MIT License - see the LICENSE file for details.
This project is designed and coded only to make nFactorial technical task.


## Acknowledgements

- The Go and Angular communities for their excellent documentation and support
- AI agents for their halp with problem solving tasks
- [Azure](https://azure.microsoft.com/) for cloud hosting services