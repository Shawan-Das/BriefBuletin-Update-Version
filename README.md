# BriefBuletin - Complete News Platform

BriefBuletin is a comprehensive news aggregation and management platform consisting of three main components: an Angular frontend for user interaction, a Go-based REST API for backend services, and a Python news scraper for automated content collection.

## Architecture

The platform is composed of three interconnected services:

- **BriefBuletin Frontend** (`BriefBuletin/`): Angular-based web application providing user interface for news browsing, user authentication, and content management.
- **BriefBuletin API** (`BriefBuletin_API/`): Go-based REST API server handling user authentication, article management, categories, comments, and database operations.
- **BriefBuletin News Scraper** (`BriefBuletin_NewsScraper/`): Python-based automated news scraper that fetches articles from multiple sources, generates AI-powered summaries, and stores them in the database.

## Features

- **User Authentication**: Registration, login, password reset with OTP verification
- **News Management**: Article creation, categorization, publishing, and viewing
- **Comment System**: User comments with admin approval workflow
- **Automated News Scraping**: Fetches from ProthomAlo, Daily Star, BBC, Al Jazeera, and more
- **AI Summarization**: Uses T5 transformer model to generate article summaries
- **Responsive Web Interface**: Modern Angular frontend with admin panel
- **PostgreSQL Database**: Robust data storage with proper schema design
- **JWT Authentication**: Secure API access with configurable bypass routes

## Prerequisites

- **Node.js** 16+ (for Angular frontend)
- **Go** 1.23+ (for API backend)
- **Python** 3.8+ (for news scraper)
- **PostgreSQL** 12+ (database)
- **Git** (version control)
- **GNU Make** (build automation)
- **PowerShell 7+** (Windows automation scripts)

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/Shawan-Das/BriefBuletin-Update-Version.git
cd BriefBuletin-Update-Version
```

### 2. Database Setup

Create a PostgreSQL database and run the schema from `BriefBuletin_API/config/sqlc/db_query/db-schema.sql`.

### 3. Configure Components

#### API Configuration
- Copy `BriefBuletin_API/config/connection/dev-config.json` and update database credentials
- Configure JWT key, SMTP settings (optional), and other parameters

#### Scraper Configuration
- Edit `BriefBuletin_NewsScraper/main.py` and update database connection variables

#### Frontend Configuration
- Update `BriefBuletin/src/environments/environment.ts` with API base URL
- Optionally configure proxy for local development

## Running the Project

### Quick Start (All Components)

Use the provided Makefile for easy management:

```bash
# Install dependencies for all components
make install-all

# Build all components
make build-all

# Run all services (API, Frontend, Scraper)
make run-all
```

### Individual Component Setup

#### 1. Start the API Backend

```bash
cd BriefBuletin_API
make dev
```

The API will start on `http://localhost:7070`

#### 2. Start the News Scraper

```bash
cd BriefBuletin_NewsScraper
make start  # Runs in background every 30 minutes
```

Or run once:
```bash
make fetch
```

#### 3. Start the Frontend

```bash
cd BriefBuletin
npm install
npm start  # Uses proxy to API at localhost:7070
```

Open `http://localhost:4200` in your browser.

### Production Deployment

#### Frontend (Vercel)
- Build: `npm run build -- --configuration production`
- Deploy to Vercel with rewrites to proxy `/api/*` to backend

#### API Backend
- Build: `make server` (creates Windows and Linux binaries)
- Deploy binary with config file

#### Scraper
- Schedule with cron/Windows Task Scheduler or use `make start` for background operation

## API Documentation

The API provides RESTful endpoints for:

- **Authentication**: `/api/auth/*` - User registration, login, password reset
- **Users**: `/api/auth/users` - User management
- **Categories**: `/api/category*` - News category management
- **Articles**: `/api/article*` - Article CRUD operations
- **Comments**: `/api/comment*` - Comment system with approval workflow

See `BriefBuletin_API/README.md` for detailed endpoint documentation and Postman examples.

## Database Schema

The system uses PostgreSQL with three main schemas:

- `news.users` - User accounts and authentication
- `news.categories` - News categories
- `news.articles` - Article content and metadata
- `news.comments` - User comments on articles

## Development

### Project Structure

```
BriefBuletin-Update-Version/
├── BriefBuletin/           # Angular frontend
├── BriefBuletin_API/       # Go backend API
├── BriefBuletin_NewsScraper/ # Python scraper
├── README.md              # This file
├── Makefile               # Build automation
└── .gitignore            # Git ignore rules
```

### Component-Specific Development

- **Frontend**: See `BriefBuletin/README.md` for Angular development setup
- **API**: See `BriefBuletin_API/README.md` for Go development and testing
- **Scraper**: See `BriefBuletin_NewsScraper/README.md` for Python development

## Testing

### API Testing
Use Postman collection or follow examples in `BriefBuletin_API/README.md`

### Frontend Testing
```bash
cd BriefBuletin
npm test
```

### Scraper Testing
```bash
cd BriefBuletin_NewsScraper
make run  # Foreground mode for debugging
```

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Verify PostgreSQL is running
   - Check credentials in config files
   - Ensure database and tables exist

2. **API CORS Errors**
   - Frontend proxy configuration
   - Backend CORS settings

3. **Scraper Not Running**
   - Check PowerShell execution policy
   - Verify Python dependencies installed

4. **Build Failures**
   - Ensure all prerequisites installed
   - Check Go modules: `go mod download`
   - Check Node modules: `npm install`

### Logs

- **API**: Console output or configured log files
- **Scraper**: `BriefBuletin_NewsScraper/activity.log`
- **Frontend**: Browser console

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes following component-specific guidelines
4. Test thoroughly
5. Submit a pull request

## License

This project is part of the Satcom Projects collection.

## Support

For issues and questions:
- Check component-specific READMEs
- Review logs for error details
- Test individual components in isolation