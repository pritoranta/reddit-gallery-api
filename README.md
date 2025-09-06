# Reddit Gallery API

Minimal Golang [Gin](https://gin-gonic.com/) HTTP API for using Reddit API indirectly. The need arose from many browsers blocking cross-origin requests from [Reddit Gallery Client](https://github.com/pritoranta/reddit-gallery-client) to Reddit API directly. I decided against using a simple reverse proxy, because I only have a need for a couple endpoints, and I wanted to leave my options open for further customization.

I decided to use Golang because this project seemed like an ideal candidate; small and simple. My first try of Golang, after years .NET and JS. I love the tools and ease of build and deployment processes, but the language is certainly quirky.
