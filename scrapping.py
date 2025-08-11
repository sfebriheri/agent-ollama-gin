import requests
from bs4 import BeautifulSoup
import re
from datetime import datetime
import json

class NewsArticleScraper:
    def __init__(self):
        self.headers = {
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
        }
        # Site-specific patterns for content extraction
        self.site_patterns = {
            'metrotvnews.com': {
                'content': {'tag': 'div', 'class_': 'detail-text'},
                'title': {'tag': 'h1', 'class_': 'title'}
            },
            'kompas.id': {
                'content': {'tag': 'div', 'class_': 'read-page--content'},
                'title': {'tag': 'h1', 'class_': 'read-page--title'}
            }
        }

    def clean_text(self, text):
        """Clean extracted text by removing extra whitespace and unwanted characters."""
        if not text:
            return ""
        # Remove extra whitespace and newlines
        text = re.sub(r'\s+', ' ', text.strip())
        # Remove special characters but keep punctuation
        text = re.sub(r'[^\w\s.,!?"-]', '', text)
        return text

    def extract_date(self, soup):
        """Extract article date using multiple methods."""
        # Try meta tags first
        date_meta = soup.find('meta', property=['article:published_time', 'article:modified_time', 'og:published_time'])
        if date_meta and date_meta.get('content'):
            return {
                'published_date': date_meta['content'],
                'extracted_from': 'meta tag'
            }
        
        # Try common date elements
        date_patterns = [
            r'\d{4}-\d{2}-\d{2}',  # YYYY-MM-DD
            r'\d{2}/\d{2}/\d{4}',  # DD/MM/YYYY
            r'\d{1,2}\s+(?:January|February|March|April|May|June|July|August|September|October|November|December)\s+\d{4}'
        ]
        
        # Look for dates in time tags, spans, and divs with date-related classes
        date_elements = soup.find_all(['time', 'span', 'div'], class_=lambda x: x and any(date_word in x.lower() for date_word in ['date', 'time', 'published', 'modified']))
        
        for element in date_elements:
            text = element.get_text()
            for pattern in date_patterns:
                match = re.search(pattern, text)
                if match:
                    return {
                        'published_date': match.group(),
                        'extracted_from': 'date element'
                    }
        
        return {
            'published_date': None,
            'extracted_from': None
        }

    def extract_author(self, soup):
        """Extract article author using multiple methods."""
        # Try meta tags first
        author_meta = soup.find('meta', {'name': ['author', 'article:author']}) or \
                     soup.find('meta', {'property': ['author', 'article:author']})
        if author_meta and author_meta.get('content'):
            return author_meta['content']

        # Try common author elements
        author_elements = soup.find_all(['a', 'span', 'div'], class_=lambda x: x and 'author' in x.lower())
        for element in author_elements:
            author = element.get_text().strip()
            if author and len(author) < 100:  # Basic validation
                return author

        return None

    def extract_main_content(self, soup, url):
        """Extract main article content while filtering out advertisements and irrelevant content."""
        # Get domain from URL to determine which site patterns to use
        domain = re.search(r'(?:https?://)?(?:www\.)?([^/]+)', url).group(1)
        
        # Common content container patterns
        content_patterns = [
            # Site-specific patterns first
            {'tag': 'div', 'class_': 'read-page--content'},  # Kompas new pattern
            {'tag': 'div', 'class_': 'detail-text'},  # Metrotvnews pattern
            {'tag': 'div', 'class_': 'read-page-content'},  # Kompas alternative
            # Generic patterns as fallback
            {'tag': 'div', 'class_': lambda x: x and any(c in x for c in ['kcm-read-text', 'read-text'])},
            {'tag': 'div', 'class_': lambda x: x and ('article' in x.lower() or 'content' in x.lower())},
            {'tag': 'article'},
            {'tag': 'div', 'itemprop': 'articleBody'},
            {'tag': 'div', 'role': 'main'}
        ]

        # Patterns for elements that typically contain ads
        ad_patterns = [
            'ads', 'advertisement', 'banner', 'promo', 'sponsor', 'widget',
            'subscription', 'subscribe', 'premium', 'recommended', 'related',
            'social-share', 'newsletter', 'popup', 'modal'
        ]

        main_content = None
        # First try to find the main content container
        for pattern in content_patterns:
            if 'tag' in pattern:
                if 'class_' in pattern:
                    elements = soup.find_all(pattern['tag'], class_=pattern['class_'])
                elif 'itemprop' in pattern:
                    elements = soup.find_all(pattern['tag'], itemprop=pattern['itemprop'])
                elif 'role' in pattern:
                    elements = soup.find_all(pattern['tag'], role=pattern['role'])
                else:
                    elements = soup.find_all(pattern['tag'])

                if elements:
                    # Filter out elements that might be ads
                    valid_elements = [
                        el for el in elements
                        if not any(ad_term in str(el.get('class', [])).lower() for ad_term in ad_patterns)
                    ]
                    if valid_elements:
                        # Find the element with the most substantial paragraphs
                        main_content = max(valid_elements, key=lambda x: sum(len(p.get_text()) for p in x.find_all('p')))
                        break

        if not main_content:
            return None

        # Remove unwanted elements from main content
        for el in main_content.find_all(['div', 'aside', 'iframe', 'script', 'style']):
            if any(ad_term in str(el.get('class', [])).lower() for ad_term in ad_patterns):
                el.decompose()

        # Extract paragraphs and clean them
        paragraphs = []
        seen_paragraphs = set()
        for p in main_content.find_all(['p', 'h2', 'h3']):
            text = self.clean_text(p.get_text())
            # More stringent filtering criteria
            if (text and len(text) > 20 and 
                text not in seen_paragraphs and
                not any(ad_term in text.lower() for ad_term in [
                    'subscribe', 'premium', 'download', 'install', 'click here', 'sign up',
                    'advertisement', 'sponsored', 'recommended', 'pre-order', 'bonus',
                    'kompas.id', 'akses', 'langganan', 'edisi khusus'
                ])):
                paragraphs.append(text)
                seen_paragraphs.add(text)

        return paragraphs

    def scrape_article(self, url):
        """Main method to scrape a news article."""
        try:
            response = requests.get(url, headers=self.headers, timeout=10)
            response.raise_for_status()
            
            soup = BeautifulSoup(response.content, 'html.parser')
            
            # Extract title
            title = soup.find('h1') or soup.find('meta', property='og:title')
            title_text = title.get_text() if hasattr(title, 'get_text') else title.get('content', '') if title else ''
            title_text = self.clean_text(title_text)

            # Extract date
            date_info = self.extract_date(soup)
            
            # Extract author
            author = self.extract_author(soup)
            
            # Extract content
            content = self.extract_main_content(soup, url)
            
            # Prepare the result
            article_data = {
                'url': url,
                'title': title_text,
                'date': date_info,
                'author': author,
                'content': content,
                'scraped_at': datetime.now().isoformat()
            }

            return article_data

        except requests.exceptions.RequestException as e:
            print(f"Error fetching the webpage: {e}")
            return None

if __name__ == "__main__":
    # List of articles to scrape
    urls = [
        "https://www.metrotvnews.com/play/NnjCew9y-janjinya-19-juta-lapangan-kerja-tapi-di-manaf",
        "https://www.metrotvnews.com/read/2025/08/10/your-article",  # Replace with actual Metro TV article
        # Add more URLs as needed
    ]
    
    scraper = NewsArticleScraper()
    
    # Create a folder for storing article data if it doesn't exist
    import os
    if not os.path.exists('articles'):
        os.makedirs('articles')
    
    # Scrape each article
    for index, url in enumerate(urls, 1):
        print(f"\nScraping article {index} of {len(urls)}")
        print(f"URL: {url}")
        
        article = scraper.scrape_article(url)
        
        if article:
            # Create a unique filename for each article
            filename = f"articles/article_{index}_{datetime.now().strftime('%Y%m%d_%H%M%S')}.json"
            
            # Save article to its own JSON file
            with open(filename, 'w', encoding='utf-8') as f:
                json.dump(article, f, ensure_ascii=False, indent=2)
            
            print("\nArticle Details:")
            print(f"Title: {article['title']}")
            print(f"Date: {article['date']}")
            print(f"Author: {article['author']}")
            print(f"Saved to: {filename}")
            
            print("\nContent Preview (first 3 paragraphs):")
            if article['content']:
                for paragraph in article['content'][:3]:
                    print(f"\n{paragraph}")
                if len(article['content']) > 3:
                    print("\n... (full content in JSON file)")
            else:
                print("No content found")
        else:
            print("Failed to scrape the article")
        
        # Add a small delay between requests to be polite to the server
        import time
        time.sleep(2)  # 2 seconds delay
    
    print("\nScraping completed!")