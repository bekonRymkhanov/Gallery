# GalleryApp: Gallery Application

![GalleryApp Banner](/media/photo_2024-04-21_05-35-08.jpg)

## 🌟 Overview

GalleryApp is a gallery application that allows users to browse, search, and interact with a vast collection of photographs. Built with a tech stack (Golang backend, Angular frontend), GalleryApp delivers a seamless and responsive experience for photography enthusiasts.

**Live Demo:** [http://4.247.165.79](http://4.247.165.79)

**Video Demo:** 
<video controls width="600">
  <source src="media/Screencast from 2025-05-05 22-19-11.webm" type="video/webm">
</video>
## 📋 Features

- **Intuitive Gallery Interface**: Browse photos with a responsive grid layout
- **Powerful Search Capabilities**: Search photos by keywords, tags, or visual similarity
- **Interactive Elements**: Like and comment on favorite photographs
- **User Authentication**: Secure login and personalized experience
- **Responsive Design**: Optimal viewing on all device sizes

## 🛠️ Technology Stack

### Backend
- **Golang**: High-performance, concurrent server-side language
- **Postgres**: Robust relational database for structured data storage
- **RESTful API**: Clean architecture for client-server communication

### Frontend
- **Angular**: Component-based framework for building dynamic SPA
- **TypeScript**: Type-safe programming for improved code quality
- **SCSS**: Advanced styling with better organization capabilities

### Deployment
- **NGINX**: High-performance web server and reverse proxy
- **Azure**: Cloud hosting platform for reliable deployment
- **Docker**: Containerization for consistent deployment environments

## 🚀 Installation and Setup

### Prerequisites
- Go 1.18+
- Node.js 16+
- Angular CLI 14+
- PostgreSQL 14+
- Docker & Docker Compose (optional)

### Backend Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/photovista.git
cd photovista/backend

# Install dependencies
go mod download

# Configure database connection in .env file
cp .env.example .env
# Edit .env with your database credentials

# Run migrations
go run cmd/migrate/main.go

# Start the backend server
go run cmd/api/main.go
```

### Frontend Setup

```bash
# Navigate to frontend directory
cd ../frontend

# Install dependencies
npm install

# Start development server
ng serve

# Application will be available at http://localhost:4200
```

### Using Docker (Optional)

```bash
# From project root
docker-compose up -d
```

## 💭 Design Process & Development Journey

### Planning Phase
The development of PhotoVista began with a comprehensive planning phase where I mapped out the core functionalities and user experience. I created wireframes and mockups to visualize the user interface and established a clear architecture for both frontend and backend components.

### Technology Choices
I selected Golang for the backend due to its performance characteristics, strong concurrency model, and excellent support for building RESTful services. For the frontend, Angular was chosen for its robust framework features, TypeScript integration, and component-based architecture which allows for clean code organization and reusability.

The decision to use PostgreSQL was driven by the need for a reliable relational database with strong JSON support, which is ideal for storing varied metadata associated with photographs.

### Development Approach
I followed an iterative development process:
1. Built the core API endpoints and database schema
2. Implemented basic authentication and user management
3. Developed the gallery browsing experience
4. Added search functionality with both text and image-based queries
5. Implemented social features (likes, comments)
6. Performed optimization and refinement

Throughout the process, I maintained a focus on code quality, with comprehensive tests and documentation to ensure maintainability.

## 🧠 Unique Approaches & Methodologies

### Image Similarity Search
One of the most distinctive features of PhotoVista is its image similarity search. Rather than relying solely on metadata or tags, the application employs computer vision techniques to identify visually similar photographs. This was implemented using a combination of:

- Feature extraction using convolutional neural networks
- Vector embedding storage in PostgreSQL with the pgvector extension
- Optimized nearest-neighbor search algorithms

### Adaptive Image Loading
To improve performance and user experience, I implemented an adaptive image loading system that:

- Serves appropriately sized images based on device capabilities
- Implements progressive loading with blurred placeholders
- Utilizes intelligent caching strategies to reduce bandwidth usage

### Modular Architecture
The backend implements a clean, hexagonal architecture that separates business logic from external concerns, making the system highly testable and maintainable.

## ⚖️ Technical Trade-offs

### SQL vs NoSQL Database
While NoSQL solutions like MongoDB might offer more flexibility for document storage, I chose PostgreSQL for its:
- ACID compliance and data integrity guarantees
- Strong JSON support that still allows for schema flexibility when needed
- Better support for complex queries and relationships
- Integration with full-text search capabilities

The trade-off is slightly more rigid schema design, but the benefits in data consistency and query performance were worth it for this application.

### Monolith vs Microservices
I opted for a monolithic architecture initially rather than microservices because:
- It simplified the development process for a solo developer
- Reduced operational complexity during the initial launch
- Provided a solid foundation that can be broken into microservices later if needed

The application is designed with clean separation of concerns, making future transition to microservices possible if scaling requires it.

### Angular vs React
Angular was selected over React despite React's larger community because:
- Angular's opinionated structure enforced consistency across the codebase
- TypeScript integration was seamless
- The comprehensive ecosystem (routing, forms, HTTP client) reduced the need for third-party libraries

The trade-off is a steeper learning curve and slightly more boilerplate code, but the benefits in terms of maintainability and structure were valuable for this project.

## 🐛 Known Issues & Limitations

### Current Limitations
- **Mobile image upload**: The image upload feature may have inconsistent behavior on some mobile browsers
- **Search performance**: Very large image searches may experience latency on the free tier of Azure hosting
- **Browser compatibility**: Some advanced visual effects may not appear correctly in older browsers

### Future Improvements
- Implement WebSocket for real-time comment notifications
- Add user galleries and collections functionality
- Integrate with social media platforms for easy sharing
- Implement AI-powered automatic tagging of uploaded images
- Add support for video content

## 📊 Performance Considerations

PhotoVista was optimized for performance in several ways:
- Backend API response times kept under 100ms for most operations
- Frontend bundle size minimized through lazy loading and code splitting
- Image assets optimized and served through CDN
- Database queries optimized with appropriate indexes

## 👥 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgements

- [Unsplash](https://unsplash.com/) for the free high-quality images used in development
- The Go and Angular communities for their excellent documentation and support
- [Azure](https://azure.microsoft.com/) for cloud hosting services