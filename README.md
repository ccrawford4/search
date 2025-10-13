# Search Engine API

A high-performance web search engine built in Go, featuring intelligent web crawling, TF-IDF based relevance ranking, and a RESTful API. The engine respects robots.txt policies and supports distributed caching with Redis.

## Features

- **Intelligent Web Crawler**: Automatically crawls websites while respecting robots.txt policies and crawl delays
- **TF-IDF Ranking**: Uses Term Frequency-Inverse Document Frequency algorithm for search result relevance
- **Stemming Support**: Implements Snowball stemming algorithm for better search matching
- **Distributed Caching**: Redis integration for high-performance search operations
- **Multiple Database Support**: Compatible with MySQL, SQLite, and SQL Server via GORM
- **RESTful API**: Clean HTTP endpoints for searching and crawling operations
- **CORS Support**: Configurable cross-origin resource sharing for development
- **Concurrent Processing**: Efficient parallel crawling and indexing

## Architecture

The search engine consists of several key components:

- **Crawler**: Fetches and processes web pages concurrently
- **Indexer**: Builds and maintains the search index with word frequencies
- **TF-IDF Calculator**: Computes relevance scores for search results
- **Database Layer**: Persistent storage for indexed documents
- **Redis Cache**: Fast access to frequently searched terms
- **API Server**: Gin-based HTTP server handling search and crawl requests

## Prerequisites

- Go 1.23.0 or higher
- Redis server
- Database server (MySQL, SQLite, or SQL Server)
- Docker (optional, for containerized deployment)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd search
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables (see Configuration section)

## Configuration

Create a `.env` file in the project root:
`cp .example.env .env`

```env
MYSQL_USERNAME=<mysql-username>
MYSQL_PASSWORD=<mysql-password>
MYSQL_HOST=<mysql-host:port>
MYSQL_DATABASE=<database-name>
REDIS_HOST=<redis-host:port>
REDIS_PASSWORD=<redis-password>
ENV=development
```

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `MYSQL_USERNAME` | MySQL username | Yes | - |
| `MYSQL_PASSWORD` | MySQL password | Yes | - |
| `MYSQL_HOST` | MySQL host and port | No | localhost:3306 |
| `MYSQL_DATABASE` | MySQL database name | Yes | - |
| `REDIS_HOST` | Redis server host and port | Yes | - |
| `REDIS_PASSWORD` | Redis authentication password | Yes | - |
| `ENV` | Environment mode (development/production) | No | - |

**Note**: The application automatically constructs the MySQL DSN from the provided credentials in the format: `username:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local`

## Running the Application

### Local Development

```bash
ENV=development go run .
```

The server will start on `http://localhost:8080`

### Docker Deployment

#### Option 1: Using Docker Compose (Recommended)

```bash
docker compose up
```

#### Option 2: Manual Docker Setup

1. Build the Docker image:
```bash
docker build -t search-engine .
```

2. Run the container:
```bash
docker run -p 8080:8080 \
  -e MYSQL_USERNAME="your-mysql-username" \
  -e MYSQL_PASSWORD="your-mysql-password" \
  -e MYSQL_HOST="mysql:3306" \
  -e MYSQL_DATABASE="your-database-name" \
  -e REDIS_HOST="redis:6379" \
  -e REDIS_PASSWORD="your-redis-password" \
  search-engine
```

## API Endpoints

### Search Documents

Search the indexed documents for a given term.

**Endpoint**: `POST /search`

**Request Body**:
```json
{
  "SearchTerm": "golang programming"
}
```

**Response** (200 OK):
```json
{
  "HITS": [
    {
      "URL": "https://example.com/golang-guide",
      "TITLE": "Complete Guide to Golang",
      "DESCRIPTION": "Learn Go programming from basics to advanced",
      "TFIDF": 0.89
    }
  ],
  "TERM": "golang programming"
}
```

**Error Response** (404):
```json
{
  "error": "No results found"
}
```

### Start Crawling

Initiate a web crawl for a specified host.

**Endpoint**: `POST /crawl`

**Request Body**:
```json
{
  "Host": "https://example.com"
}
```

**Response** (200 OK):
```json
{
  "success": "true",
  "host": "https://example.com"
}
```

**Note**: Crawling runs asynchronously in the background. Results will be indexed and available for search as pages are processed.

### Serve Documents

Serves locally stored documents from the corpus (primarily for testing).

**Endpoint**: `GET /documents/top10/*any`

**Example**: `GET /documents/top10/index.html`

## How It Works

### Crawling Process

1. **Initialization**: The crawler fetches and parses `robots.txt` to determine crawl policies
2. **Queue Management**: URLs are queued and deduplicated using a hash set
3. **Concurrent Downloads**: Multiple pages are downloaded in parallel with respect to crawl delays
4. **Content Extraction**: HTML is parsed to extract text content, links, titles, and descriptions
5. **Stemming**: Words are stemmed using the Snowball algorithm for better matching
6. **Indexing**: Word frequencies and metadata are stored in the database and Redis cache

### Search Algorithm

1. **Query Processing**: Search terms are stemmed to match the indexed format
2. **Frequency Lookup**: Term frequencies are retrieved from the index
3. **TF-IDF Calculation**:
   - **TF** (Term Frequency): `termCount / totalWords`
   - **IDF** (Inverse Document Frequency): `log10(numDocs / (docsContainingWord + 1))`
   - **TF-IDF**: `TF * IDF`
4. **Ranking**: Results are sorted by TF-IDF score in descending order
5. **Response**: Top results are returned with URLs, titles, descriptions, and relevance scores

## Development

### Running Tests

```bash
go test -v ./...
```

### Code Structure

```
.
├── main.go           # Application entry point and routing
├── server.go         # HTTP handlers and template execution
├── crawl.go          # Web crawling logic
├── index.go          # Search index interface
├── db_index.go       # Database-backed index implementation
├── memory_index.go   # In-memory index implementation
├── tfidf.go          # TF-IDF calculation and search result ranking
├── extract.go        # HTML content extraction
├── download.go       # HTTP downloading with retry logic
├── clean.go          # URL normalization and validation
├── delay.go          # Crawl delay management
├── redis.go          # Redis client and caching
├── db.go             # Database models and operations
├── structs.go        # Data structures and types
└── stop_words.go     # Stop words filtering
```

### CORS Configuration

In development mode, CORS is enabled for:
- `http://localhost:3000`
- `http://127.0.0.1:3000`

For production, configure CORS according to your frontend domain requirements.

## Performance Considerations

- **Concurrent Crawling**: Respects crawl delays while maximizing throughput
- **Redis Caching**: Frequently accessed search terms are cached for fast retrieval
- **Database Indexing**: Proper indexes on URL and word columns for efficient lookups
- **Stemming**: Reduces index size and improves search recall
- **Batch Processing**: URLs are processed in batches to optimize database operations

## Limitations

- Crawling is limited to HTML content (does not process PDFs, images, etc.)
- Only crawls pages within the same host domain
- Respects robots.txt but does not handle JavaScript-heavy sites
- Search is currently single-term (no phrase matching or boolean operators)
