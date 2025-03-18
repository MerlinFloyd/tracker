# My Fullstack App

This project is a fullstack web application built with React for the frontend, Go for the backend API, and PostgreSQL as the database. 

## Project Structure

```
my-fullstack-app
├── frontend          # React frontend application
│   ├── public
│   │   ├── index.html        # Main HTML file for the React app
│   │   └── favicon.ico       # Favicon for the web application
│   ├── src
│   │   ├── components
│   │   │   └── App.js        # Root component of the React application
│   │   ├── index.js          # Entry point for the React application
│   │   └── styles
│   │       └── index.css     # CSS styles for the React application
│   ├── package.json          # Configuration file for npm
│   └── README.md             # Documentation for the frontend
├── backend                 # Go backend API
│   ├── cmd
│   │   └── server
│   │       └── main.go      # Entry point for the Go API server
│   ├── internal
│   │   ├── api
│   │   │   └── handlers.go   # API request handlers
│   │   ├── database
│   │   │   └── postgres.go   # Database connection and interaction
│   │   └── models
│   │       └── models.go     # Data models for the application
│   ├── go.mod               # Go module file
│   └── README.md            # Documentation for the backend
├── database                # Database setup and migrations
│   ├── migrations
│   │   └── 001_initial_schema.sql  # Initial database schema
│   └── schema.sql          # Database schema definition
└── README.md               # Documentation for the entire project
```

## Getting Started

### Prerequisites

- Node.js and npm for the frontend
- Go for the backend
- PostgreSQL for the database

### Installation

1. **Frontend Setup**
   - Navigate to the `frontend` directory.
   - Run `npm install` to install the necessary dependencies.
   - Start the development server with `npm start`.

2. **Backend Setup**
   - Navigate to the `backend` directory.
   - Run `go mod tidy` to install the necessary Go dependencies.
   - Start the Go server with `go run cmd/server/main.go`.

3. **Database Setup**
   - Set up your PostgreSQL database using the SQL commands in `database/schema.sql` and `database/migrations/001_initial_schema.sql`.

### Usage

- Access the frontend application at `http://localhost:3000`.
- The backend API will be available at `http://localhost:8080`.

### Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

### License

This project is licensed under the MIT License.