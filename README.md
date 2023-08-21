## 使用技術
- [Go](https://github.com/golang/go)
- [PostgreSQL](https://www.postgresql.org/)
- [redis](https://redis.io/)
- [air](https://github.com/cosmtrek/air)
- [Gin](https://github.com/gin-gonic/gin)
- [sqlc](https://github.com/kyleconroy/sqlc)
- [golang-migrate](https://github.com/golang-migrate/migrate)

## セットアップ
- 環境変数
  ```
  cp .env.sample .env
  ```
- docker compose で起動
  ```
  docker compose up --build
  ```
- DevContainer で起動
  1. 拡張機能 [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) をインストール
  2. VSCode のコマンドパレットを開く
  3. `>Devcontainers: Reopen in Container` を選択する
  4. コンテナを起動できたら、ターミナルで以下のコマンドを打ってAPIを起動する
     ```
     air
     ```

## APIサーバーの起動
```
air
```

## マイグレーション
1. app コンテナのターミナルを開く
2. 以下のコメントを実行してスキーマ、テーブルなどを作成
```
migrate -database ${POSTGRES_URL} -path db/migrations up
```
スキーマに変更を加える場合
1. 以下のコマンドを実行
   ```
   migrate create -ext sql -dir db/migrations -seq [変更内容]
   ```
2. db/migrationsの中に新規で下記のようなファイルが作成されるため、[変更内容].up.sqlに変更するSQLを書き、[変更内容].down.sqlには変更前に戻すSQLを書く
   ```
   000001_[変更内容].down.sql
   000001_[変更内容].up.sql
   ```

golang-migrateの[チュートリアル](https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md)も参考にしてください！

psqlを使ってDBにアクセスする場合
```
psql ${POSTGRES_URL}
```

## OpenAPI Specification

Stoplight Studioで各APIの仕様は次の[URL](https://project3-shimabu.stoplight.io/docs/project3-backend/branches/main/u788di80n61qf-team3-backend-api)をみてください。

yaml からコード生成
```
oapi-codegen -config docs/config.yaml docs/openapi.yaml > docs/openapi.gen.go
```

docs/openapi.gen.go に作成されるInterfaceに沿ってAPIを開発してください。

## Git コミットルール

Conventional Commits に従う。

- [🐛] fix: コードベースのバグにパッチを当てる場合（セマンティックバージョン管理における PATCH に相当)
- [✨] feat: コードベースに新しい機能を追加した場合(セマンティックバージョン管理における MINOR に相当)
- [💥] BREAKING CHANGE: 本文または脚注に BREAKING CHANGE:が存在する、または型、範囲の直後に!が追加されているコミットは、API の破壊的変更を意味します。(セマンティックバージョン管理における MAJOR に相当) BREAKING CHANGE はあらゆる型のコミットに含めることができます。
- [📝] docs: ドキュメントの生成や修正を行う場合
- [♻️] refactor: ロジックの変化は行わず、内部構造を整理のみを行う場合
- [✅] test: テストの追加、及び修正を行う場合
- [👷] ci: CI ツールのファイルの変更を行う場合
