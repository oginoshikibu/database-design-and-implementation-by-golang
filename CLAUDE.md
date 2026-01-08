# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

「Database Design and Implementation」（Springer出版）の書籍に基づき、SimpleDBをGolangで実装する輪読会プロジェクト。データベース管理システム（DBMS）の内部構造を学ぶための教育的プロジェクト。

## 開発コマンド

```bash
# テスト実行
go test ./...

# カバレッジ付きテスト
go test -v -race -coverprofile=coverage.out ./...

# 単一テスト実行
go test -v -run TestFileMgr ./src/file/

# リント
golangci-lint run

# フォーマット
gofmt -w .
```

## アーキテクチャ

書籍に沿ってDBMSをレイヤー構造で実装:

- **src/file/**: ディスク・ファイル管理層（Chapter 3）
  - `BlockID`: ファイル内のブロックを識別（Fig. 3.13）
  - `Page`: ディスクブロックの内容を保持（Fig. 3.14）
  - `FileMgr`: ブロックの読み書きを管理（Fig. 3.15）

今後、Buffer層、Transaction層などを追加予定。

## コーディング規約

- `gofmt`でフォーマット
- `golangci-lint`のチェックをパス
- パッケージ名は小文字・単数形
- エクスポートされる関数・型にはGoDocコメントを記述
- コードコメントは専門用語（BlockID、Page、mutex等）を除き日本語で記述
- エラーは`fmt.Errorf`でラップ
- 書籍のFigure番号をコメントに記載（例: `// Fig. 3.15`）
- 書籍の実装との差分がある場合はコメントに記載（例: Goの慣習に合わせた変更、不要なメソッドの削除理由など）

## パッケージドキュメント

各パッケージには `README.md` を作成し、以下の内容を含める:

- **概要**: パッケージの責務・役割
- **コンポーネント図**: Mermaidのクラス図で型間の関係を表現
- **シーケンス図**: Mermaidのシーケンス図で主要操作のデータフローを表現
- **型ごとの責務**: 各型（struct）のフィールド・メソッドと役割
- **書籍との実装差分**: Java版との違いを表形式で記載（理由も含む）
- **使用例**: 基本的な使い方を示すコード例

参考: `src/file/README.md`

## 参考

- 書籍: [Database Design and Implementation](https://link.springer.com/book/10.1007/978-3-030-33836-7)
- 参考実装: [yokomotod/database-design-and-implementation-go](https://github.com/yokomotod/database-design-and-implementation-go)
