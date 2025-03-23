PingMe Threads Bot

A bot that posts using DeepSeek prompts.

Features
Auto-post content on Threads
AI generated responses

Future Features
Schedule posts
Support Image and Video Posts
Auto-reply to comments

Installation
Clone the repository
git clone https://github.com/maxcolliander/PingMe.git

Install dependencies
go mod tidy

Set up environment variables
THREADS_ACCESS_TOKEN=your_access_token
DEEPSEEK_API_KEY=your_api_key

In utils/deepseek.go you can modify user and system roles
In PingMe.go you can modify prompts

To post -> go run PingMe.go
