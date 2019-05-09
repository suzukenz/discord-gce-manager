# discord-gce-manager

自分用。GCE インスタンスのステータスを調べて DiscordBot で通知する。

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

`k8s/configmap.yaml`に色々定義する

- gcp.project.id (PROJECT_ID) => GCP のプロジェクト ID
- discord.bot.token (DISCORD_TOKEN) => Discord Bot の Token
- discord.webhookurl (DISCORD_WEBHOOK) => Discord の WebhookURL (Cron 等でアプリ側から送信する用)

ローカルでプログラムを動かすときは `tools/configs.env`に直接環境変数を書いて、シェルで実行する。

### Build

go のコードは以下で build (dep は最初だけ)

```
dep ensure
make build
```

docker は docker-compose を build のためだけに使っている

```
cd docker
docker-compose build
```

## Appendix

[Google App Engine のフレキシブル環境を使おうと思ったけど高いので辞めたブランチ](https://github.com/suzukenz/discord-gce-manager/tree/gae-flexible)
