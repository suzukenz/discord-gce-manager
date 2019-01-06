# discord-gce-manager

ほぼ自分用。GCE インスタンスのステータスを調べて DiscordBot で通知する。

## Usage

### BOT コマンド

```
/check            : 全部のゲームサーバーの状態を聞く。
/check <GameName> : 指定したゲームサーバーの状態を聞く。
/run <GameName>   : 指定したゲームサーバーを起動する。
/stop <GameName>  : 指定したゲームサーバーを停止する。
```

## Development

開発メモ。大体は Makefile に書いてあるのでそれ以外のことメモ

### Prerequisites

- Golang (1.9 以上でテスト）
- Docker と docker-compose
- Google Cloud SDK

### Configurations

`configs.env`に書いてある以下の設定が必要

- PROJECT_ID => GCP のプロジェクト ID
- DISCORD_TOKEN => Discord Bot の Token
- DISCORD_WEBHOOK => Discord の WebhookURL (Cron 等でアプリ側から送信する用)
