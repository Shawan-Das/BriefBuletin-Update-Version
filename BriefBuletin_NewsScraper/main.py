# !pip install feedparser requests beautifulsoup4 psycopg2-binary transformers torch python-dateutil

### Version 3, Direct Link version (Test)
import requests
from bs4 import BeautifulSoup
import psycopg2
import time
from datetime import datetime, timedelta
from transformers import logging, pipeline
from urllib.parse import urljoin, urlparse
import re
import sys
import logging

# Database connection details
DB_HOST = 'localhost'
DB_NAME = 'db_name'
DB_USER = 'db_admin'
DB_PASS = 'db_password'
DB_PORT = '5432'
# Dictionary of news category page URLs mapped to category_ids
SOURCES = {
    'https://www.prothomalo.com':7,
    'https://www.prothomalo.com/collection/latest':7,
    'https://www.bd-pratidin.com': 7,
    'https://www.ittefaq.com.bd/country': 7,
    'https://www.kalerkantho.com/online/country-news': 7,
    # 'https://www.dailyvorerpata.com': 7,
    # 'https://www.somoynews.tv/read/recent': 7,
    'https://www.prothomalo.com/world': 7,
    'https://www.prothomalo.com/chakri': 7,
    # 'https://www.somoynews.tv': 7,
    # 'https://en.somoynews.tv/categories/politics': 6,
    'https://www.prothomalo.com/politics': 1,
    'https://en.prothomalo.com': 1,
    'https://en.prothomalo.com/bangladesh': 1,
    # 'https://en.somoynews.tv/read/recent': 1,
    # 'https://en.somoynews.tv/categories/bangladesh': 1,
    # 'https://en.somoynews.tv/categories/international': 6,
    'https://www.thedailystar.net/news/bangladesh': 1,
    'https://www.dhakatribune.com/bangladesh': 1,
    'https://www.dhakatribune.com/bangladesh/dhaka': 1,
    'https://www.dhakatribune.com/bangladesh/nation': 1,
    # 'https://www.daily-sun.com/topic/earthquake': 1,
    'https://www.daily-sun.com/national': 1,
    'https://en.prothomalo.com/international': 6,
    'https://www.thedailystar.net/news/asia': 6,
    'https://www.bbc.com/news': 6,
    'https://www.aljazeera.com/news': 6,
    'https://www.daily-sun.com/sci-tech':3,
    # 'https://en.somoynews.tv/categories/sports':4,
    # 'https://www.dhakatribune.com/sport':4,
    # 'https://en.prothomalo.com/sports': 4,
    # 'https://www.thedailystar.net/sports': 4,
    # 'https://www.bbc.com/sport': 4,
    # 'https://www.aljazeera.com/sports': 4,
    # 'https://www.daily-sun.com/sports': 4,
}

FETCH_DELAY = 1800  # Check every hour
HOURS_LIMIT = 24  # Only fetch news from last 24 hours

# Initialize summarizer once globally
summarizer = None

# Initialize logger at module level
logger = None

def setup_logging():
    """Setup logging configuration for activity.log file"""
    global logger
    log_file = 'activity.log'
    
    logger = logging.getLogger('news_scraper')
    logger.setLevel(logging.INFO)
    
    file_handler = logging.FileHandler(log_file, mode='a', encoding='utf-8')
    file_handler.setLevel(logging.INFO)
    
    formatter = logging.Formatter(
        '%(asctime)s - %(levelname)s - %(message)s',
        datefmt='%Y-%m-%d %H:%M:%S'
    )
    
    file_handler.setFormatter(formatter)
    
    logger.handlers.clear()
    logger.addHandler(file_handler)
    
    return logger

def sleep_with_progress(total_seconds):
    for elapsed in range(total_seconds):
        remaining = total_seconds - elapsed
        # Calculate a simple progress bar (30 chars wide)
        bar_length = 30
        filled_length = int(bar_length * elapsed // total_seconds)
        bar = '=' * filled_length + '-' * (bar_length - filled_length)
        mins, secs = divmod(remaining, 60)
        time_str = f"{mins:02d}:{secs:02d}"
        sys.stdout.write(f"\rWaiting for next fetch: [{bar}] {time_str} remaining")
        sys.stdout.flush()
        time.sleep(1)
    print() 

def get_summarizer():
    """Initialize summarizer once and reuse"""
    global summarizer
    if summarizer is None:
        summarizer = pipeline("summarization", model="t5-small")
    return summarizer

def get_headers():
    """Return headers to mimic a real browser"""
    return {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8',
        'Accept-Language': 'en-US,en;q=0.5',
        'Connection': 'keep-alive',
    }

def is_valid_article_url(url, base_url):
    """Check if URL is likely to be an actual article"""
    # Must be from same domain
    base_domain = urlparse(base_url).netloc
    url_domain = urlparse(url).netloc
    
    if base_domain not in url_domain:
        return False
    
    # Exclude patterns that are definitely not articles
    exclude_patterns = [
        '/tag/', '/tags/', '/category/', '/categories/', '/author/', '/authors/',
        '/page/', '/search', '/login', '/signup', '/subscribe', '/about',
        '/contact', '/privacy', '/terms', '/api/', '/oauth', '/auth',
        'javascript:', 'mailto:', '#', '.pdf', '.jpg', '.png', '.gif',
        '/gallery/', '/video/', '/photos/', '/images/',
    ]
    
    url_lower = url.lower()
    if any(pattern in url_lower for pattern in exclude_patterns):
        return False
    
    # Exclude URLs that look like category/section pages (too short or end with category name)
    path = urlparse(url).path.strip('/')
    if not path or len(path.split('/')) < 2:
        return False
    
    # Exclude if URL ends with common category names
    category_endings = ['news', 'sports', 'business', 'world', 'international', 
                       'bangladesh', 'asia', 'entertainment', 'lifestyle', 'opinion']
    last_segment = path.split('/')[-1].lower()
    if last_segment in category_endings and len(path.split('/')) <= 2:
        return False
    
    # Must have some content in the path (not just domain)
    if len(path) < 10:
        return False
    
    return True

def extract_article_links(category_url):
    """Extract all article links from a news category page"""
    try:
        response = requests.get(category_url, headers=get_headers(), timeout=15)
        response.raise_for_status()
        soup = BeautifulSoup(response.content, 'html.parser')
        
        links = set()
        
        # Strategy 1: Find article containers (most news sites use these)
        article_containers = soup.find_all(['article', 'div'], class_=re.compile(
            r'story|article|post|news|item|card|content-item', re.I
        ))
        
        for container in article_containers:
            a_tag = container.find('a', href=True)
            if a_tag and a_tag['href']:
                href = a_tag['href']
                if href.startswith('/'):
                    href = urljoin(category_url, href)
                
                # Clean URL
                clean_url = href.split('?')[0].split('#')[0]
                
                if is_valid_article_url(clean_url, category_url):
                    links.add(clean_url)
        
        # Strategy 2: If no article containers found, look for links with article-like patterns
        if len(links) < 5:
            for a_tag in soup.find_all('a', href=True):
                href = a_tag['href']
                
                if href.startswith('/'):
                    href = urljoin(category_url, href)
                
                clean_url = href.split('?')[0].split('#')[0]
                
                if is_valid_article_url(clean_url, category_url):
                    # Additional check: URL should have date pattern or be reasonably long
                    if len(clean_url) > 50 or re.search(r'/\d{4}/', clean_url):
                        links.add(clean_url)
        
        print(f"Found {len(links)} valid article links from {category_url}")
        return list(links)[:50]  # Limit to 50 articles per category
    
    except Exception as e:
        print(f"Error extracting links from {category_url}: {e}")
        return []

def extract_article_date(soup, url):
    """Extract publication date from article page"""
    try:
        # Strategy 1: Meta tags (most reliable)
        meta_tags = [
            ('property', 'article:published_time'),
            ('name', 'article:published_time'),
            ('property', 'og:published_time'),
            ('name', 'publishdate'),
            ('name', 'date'),
            ('itemprop', 'datePublished'),
        ]
        
        for attr, value in meta_tags:
            meta = soup.find('meta', {attr: value})
            if meta and meta.get('content'):
                try:
                    from dateutil import parser
                    return parser.parse(meta['content'])
                except:
                    pass
        
        # Strategy 2: Time tags with datetime attribute
        time_tags = soup.find_all('time', datetime=True)
        if time_tags:
            try:
                from dateutil import parser
                return parser.parse(time_tags[0]['datetime'])
            except:
                pass
        
        # Strategy 3: JSON-LD structured data
        json_ld_scripts = soup.find_all('script', type='application/ld+json')
        for script in json_ld_scripts:
            try:
                import json
                data = json.loads(script.string)
                
                # Handle both single object and array
                if isinstance(data, list):
                    data = data[0] if data else {}
                
                date_fields = ['datePublished', 'publishDate', 'dateCreated']
                for field in date_fields:
                    if field in data:
                        from dateutil import parser
                        return parser.parse(data[field])
            except:
                continue
        
        # Strategy 4: Look for date in URL (e.g., /2024/11/24/)
        date_match = re.search(r'/(\d{4})/(\d{1,2})/(\d{1,2})/', url)
        if date_match:
            year, month, day = date_match.groups()
            return datetime(int(year), int(month), int(day))
        
    except Exception as e:
        print(f"Error extracting date: {e}")
    
    return None

def scrape_article(url):
    """Scrape full article: title, content, featured image, and publish date"""
    try:
        response = requests.get(url, headers=get_headers(), timeout=15)
        response.raise_for_status()
        soup = BeautifulSoup(response.content, 'html.parser')

        # Extract title
        title = None
        
        # Try h1 first
        h1_tag = soup.find('h1')
        if h1_tag:
            title = h1_tag.get_text(strip=True)
        
        # Fallback to meta tags
        if not title or len(title) < 10:
            meta_title = soup.find('meta', property='og:title')
            if meta_title and meta_title.get('content'):
                title = meta_title['content']
        
        if not title or len(title) < 10:
            meta_title = soup.find('meta', attrs={'name': 'title'})
            if meta_title and meta_title.get('content'):
                title = meta_title['content']
        
        if not title or len(title) < 10:
            print(f"No valid title found for: {url}")
            return None, None, None, None

        # Extract publication date
        published_at = extract_article_date(soup, url)
        
        # Check if article is within time limit (24 hours)
        if published_at:
            # Make timezone-naive for comparison
            if published_at.tzinfo:
                published_at = published_at.replace(tzinfo=None)
            
            time_diff = datetime.now() - published_at
            if time_diff > timedelta(hours=HOURS_LIMIT):
                hours_old = time_diff.total_seconds() / 3600
                print(f"Article too old ({hours_old:.1f} hours): {title[:50]}...")
                return None, None, None, None
        else:
            print(f"Could not determine date for: {title[:50]}...")
            return None, None, None, None

        # Extract content - improved method
        content = ""
        
        # Remove unwanted elements
        for tag in soup.find_all(['script', 'style', 'nav', 'header', 'footer', 'aside', 'iframe', 'form']):
            tag.decompose()
        
        # Strategy 1: Find article tag
        article_tag = soup.find('article')
        if article_tag:
            # Look for paragraphs within article
            paragraphs = article_tag.find_all('p')
            if paragraphs and len(paragraphs) >= 3:
                content = ' '.join([p.get_text(strip=True) for p in paragraphs if len(p.get_text(strip=True)) > 20])
        
        # Strategy 2: Look for content divs
        if not content or len(content) < 200:
            content_divs = soup.find_all(['div'], class_=re.compile(
                r'article-body|story-body|post-content|entry-content|article-content|story-content', re.I
            ))
            
            for div in content_divs:
                paragraphs = div.find_all('p')
                if paragraphs and len(paragraphs) >= 3:
                    content = ' '.join([p.get_text(strip=True) for p in paragraphs if len(p.get_text(strip=True)) > 20])
                    break
        
        # Strategy 3: Find all paragraphs in body
        if not content or len(content) < 200:
            all_paragraphs = soup.find_all('p')
            valid_paragraphs = [p.get_text(strip=True) for p in all_paragraphs if len(p.get_text(strip=True)) > 30]
            if len(valid_paragraphs) >= 3:
                content = ' '.join(valid_paragraphs[:20])  # Limit to first 20 paragraphs
        
        if not content or len(content) < 100:
            print(f"No valid content found for: {title[:50]}...")
            return None, None, None, None

        # Extract featured image
        featured_image = None
        
        # Strategy 1: Open Graph image
        og_image = soup.find('meta', property='og:image')
        if og_image and og_image.get('content'):
            featured_image = og_image['content']
        
        # Strategy 2: Twitter card image
        if not featured_image:
            twitter_image = soup.find('meta', attrs={'name': 'twitter:image'})
            if twitter_image and twitter_image.get('content'):
                featured_image = twitter_image['content']
        
        # Strategy 3: Article tag images
        if not featured_image and article_tag:
            img = article_tag.find('img')
            if img:
                featured_image = img.get('src') or img.get('data-src') or img.get('data-lazy-src')
        
        # Strategy 4: Images with featured/main classes
        if not featured_image:
            img_candidates = soup.find_all('img', class_=re.compile(r'featured|hero|main|lead', re.I), limit=5)
            if img_candidates:
                featured_image = img_candidates[0].get('src') or img_candidates[0].get('data-src')
        
        # Make relative URLs absolute
        if featured_image:
            if featured_image.startswith('//'):
                featured_image = 'https:' + featured_image
            elif featured_image.startswith('/'):
                featured_image = urljoin(url, featured_image)
            elif not featured_image.startswith('http'):
                featured_image = urljoin(url, featured_image)

        return title, content, featured_image, published_at
    
    except requests.exceptions.RequestException as e:
        print(f"Request error for {url}: {e}")
        return None, None, None, None
    except Exception as e:
        print(f"Error scraping article {url}: {e}")
        return None, None, None, None

def generate_summary(text, max_length=150):
    """Generate summary using AI (T5 model)"""
    if not text or len(text) < 100:
        return text[:200] if text else ""
    try:
        # Truncate text for model (t5-small can handle ~512 tokens)
        truncated_text = text[:2000]
        
        sum_model = get_summarizer()
        summary = sum_model(truncated_text, max_length=max_length, min_length=30, do_sample=False)[0]['summary_text']
        return summary
    except Exception as e:
        print(f"Error generating summary: {e}")
        return text[:200] + "..." if len(text) > 200 else text

def main():
    """Main function to fetch and store news articles"""
    setup_logging()  # Initialize logging at the start
    conn = psycopg2.connect(host=DB_HOST, port=DB_PORT, dbname=DB_NAME, user=DB_USER, password=DB_PASS)
    cursor = conn.cursor()

    print(f"Starting news scraper... Fetching articles from last {HOURS_LIMIT} hours")
    print("Initializing AI model...")
    get_summarizer()  # Initialize once
    print("Model ready!\n")
    
    # while True:
    total_new = 0
    total_skipped = 0
    
    for category_url, category_id in SOURCES.items():
        try:
            print(f"\n{'='*80}")
            print(f"Processing: {category_url} (Category ID: {category_id})")
            print(f"{'='*80}")
            
            # Extract article links from category page
            article_links = extract_article_links(category_url)
            
            new_articles = 0
            skipped_articles = 0
            
            for article_url in article_links:
                try:
                    # Check if article already exists
                    cursor.execute("SELECT id FROM news.articles WHERE source_url = %s", (article_url,))
                    if cursor.fetchone():
                        skipped_articles += 1
                        continue
                    
                    # Scrape the article
                    title, content, featured_image, published_at = scrape_article(article_url)
                    
                    if not title or not content:
                        continue
                    
                    # Generate summary
                    summary = generate_summary(content)
                    
                    # Use current time if we couldn't extract publish date (shouldn't happen now)
                    if not published_at:
                        published_at = datetime.now()
                    
                    # Insert into database
                    if category_id == 7:
                        cursor.execute("""
                            INSERT INTO news.articles (title, summary, content, featured_image, category_id, status, published_at, source_url)
                            VALUES (%s, %s, %s, %s, %s, 'published', %s, %s)
                        """, (title, content, content, featured_image, category_id, published_at, article_url))
                    else:
                        cursor.execute("""
                            INSERT INTO news.articles (title, summary, content, featured_image, category_id, status, published_at, source_url)
                            VALUES (%s, %s, %s, %s, %s, 'published', %s, %s)
                        """, (title, summary, content, featured_image, category_id, published_at, article_url))
                    
                    conn.commit()
                    new_articles += 1
                    print(f"✓ [{new_articles}] {title[:70]}...")
                    print(f"  Published: {published_at} | Image: {'Yes' if featured_image else 'No'}")
                    
                    # Small delay between articles
                    time.sleep(3)
                
                except Exception as e:
                    print(f"Error processing article {article_url}: {e}")
                    continue
            
            total_new += new_articles
            total_skipped += skipped_articles
            
            print(f"\nCategory Summary:")
            print(f"  - New articles: {new_articles}")
            print(f"  - Skipped (already exists): {skipped_articles}")
            
            # Delay between category pages
            time.sleep(10)
        
        except Exception as e:
            print(f"Error processing category {category_url}: {e}")
    
    print(f"\n{'='*80}")
    print(f"CYCLE COMPLETE")
    print(f"  Total new articles: {total_new}")
    print(f"  Total skipped: {total_skipped}")
    print(f"  Waiting {FETCH_DELAY} seconds ({FETCH_DELAY//60} minutes) before next fetch...")
    # sleep_with_progress(FETCH_DELAY)
    print(f"{'='*80}\n")
    
    logger.info(f"{'='*80}")
    logger.info(f"CYCLE COMPLETE")
    logger.info(f"  Total new articles: {total_new}")
    logger.info(f"  Total skipped: {total_skipped}")
    logger.info(f"  Last Fetch: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    logger.info(f"{'='*80}")
    logger.info
        # break
        # time.sleep(FETCH_DELAY)

if __name__ == "__main__":
    main()