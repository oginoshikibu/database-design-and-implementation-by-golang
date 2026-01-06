# Database Design and Implementation by Golang

「Database Design and Implementation」（Springer出版）を読みながら、Golangでデータベースシステムを実装する輪読会プロジェクトです。

## 概要

このリポジトリは、データベース管理システム（DBMS）の内部構造を学ぶための教育的プロジェクトです。書籍の内容に沿って、SimpleDBをGolangで実装していきます。

## 参考

- 書籍: [Database Design and Implementation](https://link.springer.com/book/10.1007/978-3-030-33836-7)
- 参考実装: [yokomotod/database-design-and-implementation-go](https://github.com/yokomotod/database-design-and-implementation-go)

## 開発環境

### 必要なツール

- [mise](https://mise.jdx.dev/) - バージョン管理ツール
- Go 1.23以上

### セットアップ

```bash
# リポジトリをクローン
git clone https://github.com/oginoshikibu/database-design-and-implementation-by-golang.git
cd database-design-and-implementation-by-golang

# mise設定を信頼
mise trust

# 必要なツールをインストール
mise install

# 依存関係のインストール（必要に応じて）
go mod download
```

### コマンド

```bash
# テストの実行
go test ./...

# カバレッジ付きテスト
go test -v -race -coverprofile=coverage.out ./...

# リント
golangci-lint run

# フォーマット
gofmt -w .
```

## GitHub Copilot Code Reviewの設定

このリポジトリでは、GitHub Copilotによる自動コードレビューを推奨しています。

### 有効化手順

1. リポジトリの **Settings** > **Code security and analysis** に移動
2. **Copilot** セクションを探す
3. **Copilot Code Review** を有効化する

これにより、Pull Request作成時にCopilotが自動的にコードレビューを提供します。

レビュー時の重点ポイントは `.github/copilot-instructions.md` に記載されています。

## ディレクトリ構成

```
.
├── .github/
│   ├── workflows/               # GitHub Actions設定
│   ├── copilot-instructions.md  # Copilotへのレビュー指示
│   ├── copilot-workspace.yml    # Copilot Workspace設定
│   └── pull_request_template.md # PRテンプレート
├── .mise.toml                   # mise設定
├── .golangci.yml                # golangci-lint設定
├── go.mod                       # Go modules
└── README.md                    # このファイル
```

実装は今後追加していきます。

## ライセンス

MIT License
