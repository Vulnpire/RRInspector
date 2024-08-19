package main

import (
    "bufio"
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "regexp"
    "strings"
    "sync"
    "time"
)

var (
    words        string
    wordFile     string
    status       bool
    rateLimit    int
    title        bool
    matchCode    int
    exclude      string
    reqPattern   string
    reqFile      string
    respPattern  string
    respFile     string
    headers      string
    threads      int
    caseInsensitive bool
)

func init() {
    flag.StringVar(&words, "word", "", "Comma-separated list of words to filter URLs by.")
    flag.StringVar(&wordFile, "wordL", "", "File containing list of words to filter URLs by.")
    flag.BoolVar(&status, "status", false, "Check and return URLs with specified status codes.")
    flag.IntVar(&rateLimit, "rl", 0, "Rate limit for HTTP requests per second.")
    flag.BoolVar(&title, "title", false, "Fetch and return the title of the webpage.")
    flag.IntVar(&matchCode, "mc", 0, "Specific status code to filter (e.g., 200, 301, 403).")
    flag.StringVar(&exclude, "x", "", "Comma-separated list of file extensions to exclude (regex).")
    flag.StringVar(&reqPattern, "req", "", "Regex pattern to match in the request URL.")
    flag.StringVar(&reqFile, "reqL", "", "File containing regex patterns to match in the request URL.")
    flag.StringVar(&respPattern, "resp", "", "Regex pattern to match in the response body.")
    flag.StringVar(&respFile, "respL", "", "File containing regex patterns to match in the response body.")
    flag.StringVar(&headers, "H", "", "Comma-separated list of custom headers (e.g., 'Cookie: sessionid=abc; User-Agent: custom-agent').")
    flag.IntVar(&threads, "t", 10, "Number of concurrent threads to use.")
    flag.BoolVar(&caseInsensitive, "i", false, "Case-insensitive matching for words, request, and response patterns.")
    flag.Parse()
}

func main() {
    if (words == "" && wordFile == "") && reqPattern == "" && reqFile == "" && respPattern == "" && respFile == "" {
        fmt.Println("Please provide words to filter with the -word flag, or specify a pattern with -req or -resp.")
        return
    }

    filters := loadFilters(words, wordFile)
    excludes := strings.Split(exclude, ",")
    excludeRegex := compileExclusionRegex(excludes)
    reqRegexes := loadRegexPatterns(reqPattern, reqFile)
    respRegexes := loadRegexPatterns(respPattern, respFile)
    customHeaders := parseHeaders(headers)

    var wg sync.WaitGroup
    urlChan := make(chan string, 100)
    resultChan := make(chan string, 100)

    // rate limiter
    var rateLimiter <-chan time.Time
    if rateLimit > 0 {
        rateLimiter = time.Tick(time.Second / time.Duration(rateLimit))
    }

    for i := 0; i < threads; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for url := range urlChan {
                if rateLimiter != nil {
                    <-rateLimiter
                }
                if reqRegexes != nil && !matchesAnyRegex(url, reqRegexes) {
                    continue
                }
                if status {
                    checkStatus(url, resultChan, excludeRegex, respRegexes, customHeaders)
                } else if !excludeURL(url, excludeRegex) {
                    resultChan <- url
                }
            }
        }()
    }

    go func() {
        wg.Wait()
        close(resultChan)
    }()

    go func() {
        scanner := bufio.NewScanner(os.Stdin)
        for scanner.Scan() {
            line := scanner.Text()
            if matchesFilter(line, filters) && !excludeURL(line, excludeRegex) {
                urlChan <- line
            }
        }
        close(urlChan)
    }()

    for result := range resultChan {
        fmt.Println(result)
    }
}

func matchesFilter(url string, filters []string) bool {
    for _, filter := range filters {
        if caseInsensitive {
            if strings.Contains(strings.ToLower(url), strings.ToLower(filter)) {
                return true
            }
        } else {
            if strings.Contains(url, filter) {
                return true
            }
        }
    }
    return false
}

func excludeURL(url string, excludeRegex *regexp.Regexp) bool {
    return excludeRegex != nil && excludeRegex.MatchString(url)
}

func compileExclusionRegex(excludes []string) *regexp.Regexp {
    if len(excludes) == 0 || (len(excludes) == 1 && excludes[0] == "") {
        return nil
    }
    pattern := `\.(` + strings.Join(excludes, "|") + `)$`
    return regexp.MustCompile(pattern)
}

func loadFilters(words string, wordFile string) []string {
    var filters []string
    if words != "" {
        filters = strings.Split(words, ",")
    }
    if wordFile != "" {
        fileFilters, err := loadLinesFromFile(wordFile)
        if err != nil {
            fmt.Println("Error loading words from file:", err)
            return filters
        }
        filters = append(filters, fileFilters...)
    }
    return filters
}

func loadRegexPatterns(pattern string, patternFile string) []*regexp.Regexp {
    var regexes []*regexp.Regexp
    if pattern != "" {
        regexes = append(regexes, compileRegex(pattern))
    }
    if patternFile != "" {
        filePatterns, err := loadLinesFromFile(patternFile)
        if err != nil {
            fmt.Println("Error loading patterns from file:", err)
            return regexes
        }
        for _, p := range filePatterns {
            regexes = append(regexes, compileRegex(p))
        }
    }
    return regexes
}

func loadLinesFromFile(filename string) ([]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    return lines, nil
}

func compileRegex(pattern string) *regexp.Regexp {
    if pattern == "" {
        return nil
    }
    if caseInsensitive {
        return regexp.MustCompile("(?i)" + pattern)
    }
    return regexp.MustCompile(pattern)
}

func matchesAnyRegex(text string, regexes []*regexp.Regexp) bool {
    for _, regex := range regexes {
        if regex.MatchString(text) {
            return true
        }
    }
    return false
}

func parseHeaders(headers string) map[string]string {
    if headers == "" {
        return nil
    }
    headerMap := make(map[string]string)
    pairs := strings.Split(headers, ";")
    for _, pair := range pairs {
        parts := strings.SplitN(pair, ":", 2)
        if len(parts) == 2 {
            headerMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
        }
    }
    return headerMap
}

func checkStatus(url string, resultChan chan<- string, excludeRegex *regexp.Regexp, respRegexes []*regexp.Regexp, customHeaders map[string]string) {
    client := &http.Client{
        Timeout: 5 * time.Second,
    }

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return
    }

    // add custom headers
    for key, value := range customHeaders {
        req.Header.Add(key, value)
    }

    resp, err := client.Do(req)
    if err != nil {
        return
    }
    defer resp.Body.Close()

    if excludeURL(url, excludeRegex) {
        return
    }

    // check if status code matches status code
    if (matchCode != 0 && resp.StatusCode == matchCode) ||
        (matchCode == 0 && (resp.StatusCode == http.StatusOK ||
            resp.StatusCode == http.StatusMovedPermanently ||
            resp.StatusCode == http.StatusForbidden)) {
        bodyBytes, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return
        }
        body := string(bodyBytes)

        if respRegexes != nil && !matchesAnyRegex(body, respRegexes) {
            return
        }

        if title {
            titleTag := fetchTitle(body)
            resultChan <- fmt.Sprintf("%s [%d] - %s", url, resp.StatusCode, titleTag)
        } else {
            resultChan <- fmt.Sprintf("%s [%d]", url, resp.StatusCode)
        }
    }
}

func fetchTitle(body string) string {
    // regex to extract the content of the <title> tag
    re := regexp.MustCompile("(?i)<title>(.*?)</title>")
    match := re.FindStringSubmatch(body)
    if len(match) > 1 {
        return match[1]
    }
    return "No Title"
}

