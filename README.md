# Reddit Gallery API

Minimal Golang [Gin](https://gin-gonic.com/) HTTP API for using Reddit API indirectly. The need arose from many browsers blocking cross-origin requests from [Reddit Gallery Client](https://github.com/pritoranta/reddit-gallery-client) to Reddit API directly. I decided against using a simple reverse proxy, because I only have a need for a couple endpoints, and I want to leave my options open for further customization. I decided to give Golang a try because this project seemed like an ideal candidate; small and simple.
