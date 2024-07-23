YouTube Downloader

A web application built with Go and the Gin framework for downloading YouTube videos in various formats. This application provides a simple interface for users to input a video URL, select a download format, and monitor the download progress.
Features

    Input a YouTube video URL.
    Choose from multiple download formats (e.g., 144p, 240p, 360p, 480p, 720p, 1080p, audio-only).
    View download progress with a real-time progress bar.
    Error handling with user-friendly feedback.

Technologies Used

    Go: Programming language for server-side logic.
    Gin: Web framework for building the HTTP server.
    Asynq: Task queue for managing background tasks.
    Redis: Database for storing download progress.
    Bootstrap: Front-end framework for styling.

Installation

    Clone the Repository




Install Dependencies

Make sure Go is installed, then run:

go mod tidy

Set Up Redis

Ensure Redis is installed and running on your local machine. Refer to the Redis installation guide if needed.

Configure the Application

Create a .env file in the root directory to configure Redis connection details. Example .env file:

REDIS_ADDR=localhost:6379

Run the Application

go

    go run main.go

    Access the Application

    Open your web browser and navigate to http://localhost:8080.

Usage

    Input Video URL: Enter the YouTube video URL in the provided input field.
    Select Format: Click the button to open a modal dialog where you can select the desired video format.
    Start Download: Click the "Download" button to initiate the download process.
    Monitor Progress: Track the download progress via the progress bar displayed on the interface.

Development
Frontend

    HTML/CSS: Provides the structure and design for the web interface.
    JavaScript: Manages AJAX requests for progress updates and user interactions.

Backend

    Gin Framework: Handles routing and HTTP requests.
    Asynq: Manages background tasks for processing downloads.
    Redis: Stores and retrieves download progress data.

