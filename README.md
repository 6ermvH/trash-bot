#  Trash_bot
## Build
- Setup environment variable:
```bash
# ./.env
TELEGRAM_APITOKEN=<YOUR_BOT_TOKEN>

OPENROUTER_API_KEY=<OPENROUTE_API_KEY> # CAN BE DEPLOY WITHOUT

REDIS_USERNAME=<LOGIN_REDIS>
REDIS_PASSWORD=<PASSWORD_REDIS>

# ./config/local.yaml or ./config/server.yaml
CONFIG_PATH=<PATH_TO_CONFIGURATION>
```

- Setup config:
```yaml
redis:
  addr: <REDIS_ADDR> # default in local redis:6379
  db: 1 # count of clusters

telegram: {}

server:
  port: "8080" # server settings

```

- Deploy local:
```bash
# SETUP .env with all Setup variable
# REDIS_USERNAME = ""
# REDIS_PASSWORD = ""
# CONFIG_PATH = ./config/local.yaml
docker-compose up --build
```

If all correct you can check some command in bot:
![[Pasted image 20250522192700.png]]

- Deploy server:
```bash
# Need to be a redis server
# .env has all needed variable
docker build --env_file .env -t trash-bot .
docker run trash-bot
```