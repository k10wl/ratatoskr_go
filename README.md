# Ratatoskr Telegram Bot

Ratatoskr is a Golang Telegram bot that redirects messages to a specified channel. This bot is intended to be used privately, but the code is available in an open repository for anyone to explore and use.

## Getting Started

To set up the Ratatoskr bot, follow these steps:

1. Clone the repository:
```bash
git clone https://github.com/k10wl/ratatoskr_go.git
```

2. Copy the `.env_bot.example` and `.env_webapp.example` files and rename it to `.env_bot` and `.env_webapp`. Fill in all missing information;

3. Start the webApp:
```bash
ngrok http $(PORT)
make dev_webapp
```

4. Start the bot:
Update `.env_bot` with ngrok address (telegram requires webApps to be https). Afterwards run bot
```bash
make dev_bot
```

Now the Ratatoskr bot should be up and running, ready to redirect messages to the specified channel.
