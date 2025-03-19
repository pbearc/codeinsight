import requests
from requests.auth import HTTPBasicAuth

# Your GitHub Credentials
GITHUB_USERNAME = "pbearc"
GITHUB_TOKEN = "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"  # Make sure this is kept secret

# GitHub API Endpoint
GITHUB_API_URL = "https://api.github.com/search/code"

def search_github_code(query):
    """Search for code snippets in GitHub repositories using the GitHub API."""
    headers = {"Accept": "application/vnd.github+json"}
    auth = HTTPBasicAuth(GITHUB_USERNAME, GITHUB_TOKEN)

    # Define search query
    params = {
        "q": query,  # Search for the provided code (or part of it)
        "per_page": 5  # Limit results for testing
    }

    # Make API Request
    response = requests.get(GITHUB_API_URL, headers=headers, params=params, auth=auth)

    # Check if request is successful
    if response.status_code == 200:
        return response.json()
    else:
        print(f"Error {response.status_code}: {response.json()}")
        return None

# Example: Searching for code similar to the AI-generated Quick Sort code
search_query = """def quicksort(arr):
    if len(arr) <= 1:
        return arr
    pivot = arr[len(arr) // 2]
    left = [x for x in arr if x < pivot]
    middle = [x for x in arr if x == pivot]
    right = [x for x in arr if x > pivot]
    return quicksort(left) + middle + quicksort(right)"""

# Call the function with the query
search_results = search_github_code(search_query)

# Print out results
if search_results:
    for item in search_results["items"]:
        print(f"🔹 Found in: {item['repository']['full_name']}")
        print(f"📂 File Path: {item['path']}")
        print(f"🔗 Code URL: {item['html_url']}\n")
