# Technology Stack

## Architecture

レイヤードに近い構成を採用し、`validationcontext` 本体（検証 API）と `voauto`（バインド実行）、`internal/generator`（コード生成）を分離する。実行時責務と生成時責務を分け、利用者コードへの侵襲を最小化する。

## Core Technologies

- **Language**: Go 1.24.3
- **Framework**: なし（標準ライブラリ中心）
- **Runtime**: Go toolchain

## Key Libraries

- `github.com/google/uuid`: UUID 検証
- `github.com/gostaticanalysis/codegen`: コード生成補助
- `golang.org/x/tools`, `golang.org/x/mod`: 生成関連ユーティリティ

## Development Standards

### Type Safety

- `any` を使う箇所は `voauto` の反射境界のみに限定
- 公開 API は具体型シグネチャを優先し、破壊的変更を避ける

### Code Quality

- `go fmt ./...` と `go vet ./...` を最低限実行
- 生成コード（`registry_init.go`）は手編集しない
- エラー文言はフィールド名を含め、原因特定可能にする

### Testing

- 主要公開メソッドはテーブル駆動テストで網羅
- `ValidationAggregateError` のメッセージとスタックトレースを検証
- `voauto` は正常系・コンストラクタ未登録・検証失敗を最低カバー

## Development Environment

### Required Tools

- Go 1.24.3
- Task (`Taskfile.yml` 利用時)

### Common Commands

```bash
# Generate
go generate ./...

# Build
go build ./...

# Test
go test ./...

# Task-based workflow
task generate
task lint
```

## Key Technical Decisions

- **集約エラーモデル**: 複数エラーを `ValidationAggregateError` で返却
- **スタックトレース保持**: `AddError` で発生地点を記録しデバッグ容易化
- **生成コード戦略**: `//go:generate voauto-gen` コメントから登録コード生成
- **反射利用の境界化**: 反射は `voauto.BindAndValidate` 内に閉じる

---
_全依存の列挙ではなく、開発パターンに影響する技術判断のみを記載する。_
