# RRInspector

is a powerful command-line utility for filtering and inspecting URLs based on various criteria. It is designed to identify vulnerabilities, such as IDORs, by analyzing and processing URLs and their responses.

## Features

- **URL Filtering**: Filter URLs by keywords or regex patterns.
- **Status Code Checking**: Check HTTP status codes and filter by specific codes.
- **Rate Limiting**: Control the rate of HTTP requests to avoid overloading.
- **Title Fetching**: Retrieve and display webpage titles.
- **Custom Headers**: Add custom headers, including User-Agent and cookies.
- **Case-Insensitive Matching**: Perform case-insensitive searches for more comprehensive results.
- **Threading**: Use multiple threads to speed up URL processing.
- **Response Body Inspection**: Search for specific patterns in the response body.
- **Exclusion Filters**: Exclude URLs based on file extensions or patterns.

## Installation

`go install -v github.com/Vulnpire/rrinspector@latest`

## Basic Usage

Let's suppose you have gathered urls, endpoints, params using tools like Waymore, Katana, Hakrawler, etc.

`cat urls.txt | rrinspector -word /admin,/login`

## Detailed Usage

    -word: Comma-separated list of words to filter URLs by.
    -wordL: File containing list of words to filter URLs by.
    -status: Check and return URLs with specified status codes.
    -rl: Rate limit for HTTP requests per second.
    -title: Fetch and return the title of the webpage.
    -mc: Specific status code to filter (e.g., 200, 301, 403).
    -x: Comma-separated list of file extensions to exclude (regex).
    -req: Regex pattern to match in the request URL.
    -reqL: File containing regex patterns to match in the request URL.
    -resp: Regex pattern to match in the response body.
    -respL: File containing regex patterns to match in the response body.
    -H: Comma-separated list of custom headers (e.g., 'Cookie: sessionid=abc; User-Agent: custom-agent').
    -t: Number of concurrent threads to use.
    -i: Case-insensitive matching for words, request, and response patterns.


## Example Commands

Check HTTP Status Codes and Fetch Titles:

`cat urls.txt | rrinspector -status -mc 200 -title`

Rate Limiting and Threading:

`cat urls.txt | rrinspector -rl 10 -t 5`

Inspect Request and Response Bodies with Regex Patterns:

`cat urls.txt | rrinspector -req "user/\d+" -resp "@gmail" -resp "user\sID:\s\d+" -i -t 5`

Excluding Specific File Types and Using Custom Headers:

`cat urls.txt | rrinspector -x .png,.jpeg -H "User-Agent: custom-agent" -t 5`

Using Regex Patterns from Files:

`cat urls.txt | rrinspector -reqL req_patterns.txt -respL resp_patterns.txt -i -t 10`

Add Custom Headers:

`cat urls.txt | rrinspector -H "User-Agent: custom-agent; Cookie: sessionid=abc"`
