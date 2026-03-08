# Project Structure

## Organization Philosophy

責務分離を重視し、以下 3 層で整理する。

- ライブラリ本体（検証 API）
- 自動バインド実行（`voauto`）
- コード生成実装（`internal/generator` と `cmd/voauto-gen`）

## Directory Patterns

### Core Validation
**Location**: `/` 直下の `validationcontext.go`, `validate_*.go`  
**Purpose**: 検証コンテキストと組み込み検証ルールの提供  
**Example**: `ValidateEmail`, `ValidateDateTime`, `Required`

### Auto Binding
**Location**: `/voauto/`  
**Purpose**: レイヤー非依存な構造体間の反射バインドと constructor 呼び出し  
**Example**: `BindAndValidate[T]`, `Register`

### Code Generator CLI
**Location**: `/cmd/voauto-gen/`  
**Purpose**: `go:generate` から実行される生成ツールのエントリポイント  
**Example**: `voauto-gen -dir ./domain`

### Generator Internals
**Location**: `/internal/generator/`  
**Purpose**: AST 解析・登録情報収集・`registry_init.go` 出力  
**Example**: `ScanDirectory`, `GenerateRegistries`

### Examples
**Location**: `/example/`  
**Purpose**: 利用サンプルと実行確認コード  
**Example**: `example/cmd/all_tests/main.go`

## Naming Conventions

- **Files**: `snake_case.go`（例: `validate_string.go`）
- **Components**: 公開型は `PascalCase`（例: `ValidationContext`）
- **Functions**: 公開関数は `PascalCase`、非公開は `camelCase`

## Import Organization

```go
import (
    "fmt"
    "time"

    "github.com/take0fit/validationcontext"
)
```

**Path Aliases**:
- なし（Go module path を直接利用）

## Code Organization Principles

- 追加検証メソッドは用途別ファイル（`validate_*.go`）へ配置
- 公開 API 変更時はテストと README を同時更新
- `internal/` は CLI 実装専用。外部公開 API を置かない
- 生成ファイルは再生成可能性を前提に管理し、手修正しない

---
_新規ファイルはこの責務分割と命名規約に従う。_
