import requests

def fetch_simple_wikipedia_urls(total=10):
    urls = []
    while len(urls) < total:
        response = requests.get("https://simple.wikipedia.org/w/api.php",
                                params={
                                    "action": "query",
                                    "format": "json",
                                    "list": "random",
                                    "rnnamespace": "0",
                                    "rnlimit": min(total-len(urls), 500)  # Fetch in batches, API max is usually 500
                                })
        data = response.json()
        
        for item in data['query']['random']:
            url = f"https://simple.wikipedia.org/wiki/{item['title'].replace(' ', '_')}"
            print(url)

fetch_simple_wikipedia_urls(10)  # Example: Fetch 10 random simple English Wikipedia URLs
