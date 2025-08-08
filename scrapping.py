import requests
from bs4 import BeautifulSoup

# URL of the article
url = "https://www.kompas.id/artikel/dua-jenis-bakteri-dalam-paket-mbg-di-ntt-badan-gizi-nasional-minta-maaf?open_from=Tagar_Page"

# Set headers to mimic a browser
headers = {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
}

try:
    # Send GET request to the webpage
    response = requests.get(url, headers=headers)
    
    # Check if the request was successful
    if response.status_code == 200:
        # Save the raw HTML for debugging
        with open("page_content.html", "w", encoding="utf-8") as f:
            f.write(response.text)
            
        # Parse the HTML content
        soup = BeautifulSoup(response.content, 'html.parser')
        
        print("DEBUG: Full response URL:", response.url)
        
        # Extract the article title using Tailwind classes
        title = soup.find('h1', class_='font-[\'Lora\']')
        if title:
            title_text = title.get_text(strip=True)
            print("DEBUG: Found title element:", title)
        else:
            print("DEBUG: Available h1 tags:", [h1.get('class', 'no-class') for h1 in soup.find_all('h1')])
            title_text = "Title not found"

        # Extract the publication date from meta tags
        published_date = soup.find('meta', property='article:published_time')
        modified_date = soup.find('meta', property='article:modified_time')
        if published_date and modified_date:
            published = published_date['content']
            modified = modified_date['content']
            date_text = f"Published: {published}, Last Modified: {modified}"
        else:
            date_text = "Date not found"

        # Extract the main article content (looking for article text container)
        article_body = soup.select_one('div[class*="article" i], div[class*="content" i]')
        paragraphs = article_body.find_all('p') if article_body else []
        article_text = "\n".join([p.get_text(strip=True) for p in paragraphs]) if paragraphs else "Content not found"

        # Print the extracted information
        print("Title:", title_text)
        print("Publication Date:", date_text)
        print("\nArticle Content:")
        print(article_text)

        # Optionally, save to a file
        with open("kompas_article.txt", "w", encoding="utf-8") as file:
            file.write(f"Title: {title_text}\n")
            file.write(f"Publication Date: {date_text}\n\n")
            file.write("Article Content:\n")
            file.write(article_text)

    else:
        print(f"Failed to retrieve the webpage. Status code: {response.status_code}")

except requests.exceptions.RequestException as e:
    print(f"Error fetching the webpage: {e}")