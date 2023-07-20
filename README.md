## 使用技術
- [Go](https://github.com/golang/go)
- [PostgreSQL](https://www.postgresql.org/)
- [redis](https://redis.io/)
- [air](https://github.com/cosmtrek/air)
- [Gin](https://github.com/gin-gonic/gin)
- [sqlc](https://github.com/kyleconroy/sqlc)
- [golang-migrate](https://github.com/golang-migrate/migrate)

## セットアップ
docker compose -f .devcontainer/docker-compose.local.yml up

## APIサーバーの起動
```
air
```

## マイグレーション
スキーマ、テーブルなどを作成
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

## Git コミットルール

Conventional Commits に従う。

- [🐛] fix: コードベースのバグにパッチを当てる場合（セマンティックバージョン管理における PATCH に相当)
- [✨] feat: コードベースに新しい機能を追加した場合(セマンティックバージョン管理における MINOR に相当)
- [💥] BREAKING CHANGE: 本文または脚注に BREAKING CHANGE:が存在する、または型、範囲の直後に!が追加されているコミットは、API の破壊的変更を意味します。(セマンティックバージョン管理における MAJOR に相当) BREAKING CHANGE はあらゆる型のコミットに含めることができます。
- [📝] docs: ドキュメントの生成や修正を行う場合
- [♻️] refactor: ロジックの変化は行わず、内部構造を整理のみを行う場合
- [✅] test: テストの追加、及び修正を行う場合
- [👷] ci: CI ツールのファイルの変更を行う場合
