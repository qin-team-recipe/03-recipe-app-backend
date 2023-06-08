### 使用技術
- [Go](https://github.com/golang/go)
- [PostgreSQL](https://www.postgresql.org/)
- [redis](https://redis.io/)
- [Gin](https://github.com/gin-gonic/gin)
- [sqlc](https://github.com/kyleconroy/sqlc)

## セットアップ
[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/qin-team-recipe/03-recipe-app-backend)
1. ↑のボタン押す
2. 設定はそのままで良いので`Create codespace`ボタンをクリックします
3. 開発環境のビルドが始まります（ビルド時間はかなり長いです）
4. ビルドが完了すると、エディタが表示されます

## Git コミットルール

Conventional Commits に従う。

- [🐛] fix: コードベースのバグにパッチを当てる場合（セマンティックバージョン管理における PATCH に相当)
- [✨] feat: コードベースに新しい機能を追加した場合(セマンティックバージョン管理における MINOR に相当)
- [💥] BREAKING CHANGE: 本文または脚注に BREAKING CHANGE:が存在する、または型、範囲の直後に!が追加されているコミットは、API の破壊的変更を意味します。(セマンティックバージョン管理における MAJOR に相当) BREAKING CHANGE はあらゆる型のコミットに含めることができます。
- [📝] docs: ドキュメントの生成や修正を行う場合
- [♻️] refactor: ロジックの変化は行わず、内部構造を整理のみを行う場合
- [✅] test: テストの追加、及び修正を行う場合
- [👷] ci: CI ツールのファイルの変更を行う場合