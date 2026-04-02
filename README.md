# 1行でDockerサーバ環境 + Goサンプル構築

サーバに`root`ログインし１行のコマンドを実行するだけでDocker環境とGo Webアプリ環境が構築できるスクリプトです。

## 対象OS

- Ubuntu 24

## ライセンス

[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)

# 内容

Ansibleのローカル実行でDocker環境を構築し、続けて最小のGo Webアプリをsystemdサービスとして導入します。

Goアプリは`playbooks/app/go-app`ディレクトリを単独プロジェクトとして管理し、playbook実行時にサーバへ配備されます。

## インストールモジュール

- `geerlingguy.docker` 8.0.0 (Ansible Galaxy ロール) で Docker をインストール
- `zip`, `unzip` をインストール
- `golang-go`, `sqlite3`, `curl`, `git` など Go サンプル実行に必要なパッケージをインストール
- `go-sample` サービスを作成し、`http://<server>:8080` で起動

# 使い方

新規にOSをインストールしたサーバに`root`でログインし、以下の１行のコマンドをそのままコピーして実行します。

## 実行コマンド

最新のリリースタグを使用して実行します。

```bash
curl -fsSL https://raw.githubusercontent.com/kdinstall/system-base/master/script/start.sh | REPO_USER=kdinstall REPO_NAME=system-base bash
```

> **注意:** デフォルトでは GitHub の最新リリースタグが自動的に取得・使用されます。  
> 開発中の最新コードを使いたい場合は、後述のテスト実行コマンドを使用してください。

オプション（`bash -s --` 経由で渡す）:

| オプション | 説明 |
|---|---------|
| `-test` | 最新の `master` ブランチを使用して実行 |
| `--help` | ヘルプを表示 |

## テスト実行

最新の master ブランチを使用してテスト実行する場合は、テスト用スクリプトを使用します。

```bash
curl -fsSL https://raw.githubusercontent.com/kdinstall/system-base/master/test/start.sh | REPO_USER=kdinstall REPO_NAME=system-base bash
```

> **注意:** `REPO_USER` と `REPO_NAME` の両方が必須です。未設定の場合はエラーで終了します。

## 導入後の確認

以下のコマンドで Docker と Go サンプルの導入状態を確認できます。

```bash
systemctl status docker --no-pager
systemctl status go-sample --no-pager
curl -fsSL http://localhost:8080/healthz
curl -fsSL http://localhost:8080/docker/ping
```

- `/healthz` は `{"status":"ok"}` を返します
- `/docker/ping` は Go アプリから `docker version` 実行結果を返します

### Webブラウザからのアクセス

デプロイ完了後、Webブラウザから以下のURLでアクセスして確認できます。

サーバのホスト名やIPアドレスが `example.com` または `192.168.1.100` の場合：

- **ヘルスチェック:** http://example.com:8080/healthz または http://192.168.1.100:8080/healthz
  - レスポンス例: `{"status":"ok"}`

- **dockerコマンド実行結果確認:** http://example.com:8080/docker/ping または http://192.168.1.100:8080/docker/ping
  - レスポンス例: Docker のバージョン情報を JSON で返します
  - エラーの場合: `{"ok":false, "error":"エラーメッセージ"}`

## Goアプリの管理

- Goアプリ本体は `playbooks/app/go-app` ディレクトリで管理します
- デプロイ時は `playbooks/app/main.yml` が `playbooks/app/go-app` を `/opt/go-sample-app/go-app` へコピーしてビルドします
- 変更を反映する場合は1行コマンドを再実行してください