# PingMe Threads Bot 

**PingMe** is a bot that automates content posting on **Threads** using AI-generated responses from **DeepSeek**. 

## Features
- Auto-post content on Threads
- AI generated responses

## Future Features
- Schedule posts
- Support Image and Video Posts
- Auto-reply to comments

## Installation
### Prerequisites
- Go 1.18 or newer version installed
- A **DeepSeek** API key.
- A **Threads** Access token.
### Dependencies
This project relies on the following Go modules:
- **[deepseek](https://github.com/go-deepseek/deepseek)**
- **[godotenv](https://github.com/joho/godotenv)**

### Steps to Install
1. **Clone the repository**
```sh
git clone https://github.com/maxcolliander/PingMe.git 
cd PingMe
```
2. **Install dependencies**
```sh
go mod tidy
```
3. **Set up environment variables**
```env
THREADS_ACCESS_TOKEN=your_access_token
DEEPSEEK_API_KEY=your_api_key
```
4. **Modify AI Behavior / Prompt**
- In utils/deepseek.go you can modify user and system roles
- In PingMe.go you can modify prompts

5. **Run the bot**
```sh
go run PingMe.go
```

## Troubleshooting


## Licence
This project is licensed under the MIT License - see the LICENSE file for details