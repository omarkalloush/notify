# Notify

A Slack application that monitors Github repositories and notifies specifi slack channels when new releases comes out.

### Firing it up!

1. Edit the `.env` file with the required values.

    - SLACK_TOKEN
    - SLACK_CHANNELS
    - repositoriesJSON

2. Building the image
```
docker build -t notify .
```

3. Starting the application
```
docker run -p 8080:8080 notify:latest
```
