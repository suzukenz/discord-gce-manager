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

開発メモ。大体 Makefile を書けばわかる

### Prerequisites

- Golang (1.9 以上でテスト）
- Docker
- Google Cloud SDK

### Configurations

`Makefile` と `app.yaml`に書いてある以下の設定が必要

- PROJECT_ID => GCP のプロジェクト ID
- DISCORD_TOKEN => Discord Bot の Token
- DISCORD_WEBHOOK => Discord の WebhookURL (Cron 等でアプリ側から送信する用)

### Installing

dep で go の依存モジュールをインストール

```
make deps
```

### build

```
make build
```

Docker のビルドは

```
make docker-build
```

### Deployment

```
make deploy
```

以下のリソースが作られる

- Google App Engine (Flexible)
- GAE Cron Job
- Cloud Storage (GAE Deploy の時に勝手にできる)
- (Datastore アプリの中で UPDATE 処理してるが、最初のデータは自分で登録する)
