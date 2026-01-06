# BriefBuletin News Scraper

An automated news scraping and summarization system that fetches articles from multiple Bangladeshi and international news sources, extracts content, generates summaries using AI, and stores them in a PostgreSQL database.

## Features

- **Multi-source scraping**: Fetches news from ProthomAlo, The Daily Star, Dhaka Tribune, BBC, Al Jazeera, and more
- **AI-powered summarization**: Uses T5 transformer model to generate article summaries
- **Content extraction**: Intelligently extracts article titles, content, publication dates, and images
- **Database storage**: Stores articles in PostgreSQL with category mapping
- **Automated scheduling**: Runs every 30 minutes with automatic git version control
- **Category management**: Organizes articles by category (Politics, World, Sports, Science & Tech, etc.)
- **Duplicate detection**: Skips already-stored articles to avoid redundancy
- **Activity logging**: Maintains detailed logs of all operations in `activity.log`

## Requirements

- Python 3.8+
- PostgreSQL 12+
- Git (for version control)
- GNU Make
- PowerShell 7+ (Windows)

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/Shawan-Das/BriefBuletin_NewsScraper.git
cd BriefBuletin_NewsScraper
```

### 2. Install Python Dependencies

```bash
pip install feedparser requests beautifulsoup4 psycopg2-binary transformers torch python-dateutil
```

Or manually install:
```bash
pip install feedparser requests beautifulsoup4 psycopg2-binary transformers torch python-dateutil
```

### 3. Configure Database Connection

Edit `main.py` and update the database credentials:

```python
DB_HOST = 'your_db_host'
DB_NAME = 'your_database'
DB_USER = 'your_username'
DB_PASS = 'your_password'
DB_PORT = 'your_port'
```

### 4. Configure Git (for automated commits)

Ensure git is configured with your credentials:

```bash
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

## Usage

### One-Time Fetch

Run the scraper once:

```bash
python main.py
```

Or use the traditional make command:

```bash
make fetch
```

### Automated Background Loop (Every 30 Minutes)

**Start the background worker:**
```bash
make start
```

This launches a PowerShell background process that:
- Runs `main.py` every 30 minutes
- Automatically commits changes with timestamp: `git commit -m "last update YYYY-MM-DD HH:mm:ss"`
- Pushes changes to remote repository
- Logs all activity to console and `activity.log`

**Check if the loop is running:**
```bash
make status
```

Output:
```
Loop running with PID 12345
```

**Stop the background loop:**
```bash
make stop
```

**Run the loop in the foreground (for debugging):**
```bash
make run
```

This displays real-time console output from each run. Press `Ctrl+C` to stop.

## Architecture

### File Structure

```
.
├── main.py                 # Main scraper application
├── Makefile                # Automation targets (start, stop, status, run, fetch)
├── README.md               # This file
├── requirements.txt        # Python dependencies
├── activity.log            # Logs of all operations (auto-generated)
├── .run_loop.pid           # Process ID file (auto-generated when running)
└── scripts/
    ├── run_loop.ps1        # Controller script (start/stop/status)
    └── loop_worker.ps1     # Worker loop that runs every 30 minutes
```

### How It Works

1. **Scraping**: Fetches articles from configured news sources
2. **Extraction**: Parses HTML to extract:
   - Article title
   - Publication date
   - Article content
   - Featured image URL
3. **Summarization**: Uses T5 transformer to generate 50-100 word summaries
4. **Database Storage**: Inserts articles into PostgreSQL if not already present
5. **Logging**: Records cycle metrics (new articles, skipped, timestamps)
6. **Git Integration**: Commits and pushes changes after each cycle

## Configuration

### News Sources

The scraper monitors the following sources (configured in `SOURCES` dictionary):

| Source | Categories |
|--------|-----------|
| ProthomAlo | Bangladesh News, World, Jobs, Science & Tech |
| The Daily Star | Bangladesh, Asia |
| Dhaka Tribune | Bangladesh (Dhaka, Nation) |
| BBC News | International |
| Al Jazeera | International |
| Daily Sun | National, Science & Tech |

### Timing

- **Fetch interval**: 30 minutes (configurable as `FETCH_DELAY`)
- **News age limit**: 24 hours (configurable as `HOURS_LIMIT`)
- **URL fetch delay**: 1-10 seconds between requests (respectful scraping)

## Database Schema

The scraper expects a PostgreSQL database with articles table containing:

```sql
CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    title VARCHAR(500),
    content TEXT,
    summary TEXT,
    url VARCHAR(1000) UNIQUE,
    published_date TIMESTAMP,
    image_url VARCHAR(1000),
    source VARCHAR(200),
    category_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Logs

Activity logs are stored in `activity.log` with timestamps:

```
2025-12-05 14:30:00 - INFO - CYCLE COMPLETE
2025-12-05 14:30:00 - INFO -   Total new articles: 42
2025-12-05 14:30:00 - INFO -   Total skipped: 8
2025-12-05 14:30:00 - INFO -   Last Fetch: 2025-12-05 14:30:00
```

## Troubleshooting

### Logger Error
If you encounter `logger is not defined`:
- Ensure `setup_logging()` is called in `main()`
- Check that the global `logger` variable is properly initialized

### Database Connection Failed
- Verify PostgreSQL is running
- Check database credentials in `main.py`
- Ensure the database and tables exist

### Git Push Fails
- Verify git remote is configured: `git remote -v`
- Check git credentials/SSH keys
- Ensure you have push permissions

### Background Loop Not Starting
- Verify PowerShell execution policy: `Get-ExecutionPolicy`
- Check if port is already in use
- Run `make status` to verify process

## Performance Tips

- **Batch processing**: The scraper processes multiple sources sequentially
- **AI model**: First run downloads T5 model (~1GB), subsequent runs use cache
- **Database indexing**: Create indexes on `url` and `published_date` for faster queries

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-source`
3. Make changes and test locally
4. Commit with descriptive messages
5. Push to your fork and create a Pull Request

## License

This project is part of the Satcom Projects collection.

## Support

For issues, questions, or suggestions:
- Check existing GitHub issues
- Review `activity.log` for error details
- Test with `make run` in foreground mode for debugging

## Changelog

### v3 (Current - Direct Link Version)
- Added automated 30-minute scheduling with Make targets
- Implemented PowerShell background worker for Windows
- Added PID-based process management (start/stop/status)
- Integrated git auto-commit and push with timestamps
- Enhanced logging to file and console
- Improved article date extraction with multiple strategies
- Added content summarization with T5 model

### v2
- Multi-source news scraping
- Database storage with duplicate detection

### v1
- Initial RSS feed parsing version