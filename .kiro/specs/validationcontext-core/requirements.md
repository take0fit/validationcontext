# Requirements Document

## Introduction

本仕様は ValidationContext ライブラリの現行機能を、実装と整合する要件として明文化する。対象は「複数入力の検証結果を集約し、構造体変換時に一貫したエラー処理を行う」ためのコア機能である。

## Requirements

### Requirement 1: 検証エラー集約コンテキスト
**Objective:** As a Go バックエンド開発者, I want 検証結果を一箇所に蓄積・確認できること, so that ユースケース層で一括エラー処理を実行できる

#### Acceptance Criteria
1. When 開発者が `AddError(field, message)` を呼び出したとき, the system shall フィールド名・メッセージ・スタックトレースを 1 件の検証エラーとして保持する
2. If 検証エラーが 1 件以上存在するとき, then the system shall `HasErrors()` で `true` を返す
3. While 検証エラーが存在しない状態, the system shall `AggregateError()` で `nil` を返す
4. Where 複数エラーが蓄積されている場合, the system shall `AggregateError()` で全メッセージを連結した集約エラーを返す
5. The system shall `FormatErrors()` で各エラーのフィールド名と文言を人間可読文字列に整形できる

### Requirement 2: 組み込みバリデーションメソッド
**Objective:** As a ドメイン実装者, I want 文字列・数値・日時・ネットワーク・ファイルの標準検証を利用できること, so that アプリケーションごとの重複実装を減らせる

#### Acceptance Criteria
1. When 文字列が最小/最大文字数制約を満たさないとき, the system shall エラーを追加する
2. When メールアドレス・URL・UUID・IP アドレスが不正な形式のとき, the system shall エラーを追加する
3. While 数値が最小/最大制約を外れている状態, the system shall 該当フィールドにエラーを追加する
4. Where 日付・年月・年・月・日時・時刻フォーマット検証が指定された場合, the system shall 失敗時にエラーを追加する
5. The system shall ファイルパス・拡張子・ファイルサイズ検証を提供し、条件不一致時にエラーを追加する

### Requirement 3: 必須入力検証
**Objective:** As a ドメイン開発者, I want `Required` により nil/空値を統一判定できること, so that 必須項目ルールを一貫して実装できる

#### Acceptance Criteria
1. When 値が `nil` または空値のとき, the system shall 必須エラーを追加する
2. If `skipNil=true` かつ値が `nil` のとき, then the system shall エラーを追加しない
3. While ポインタ多重参照の値が渡された場合, the system shall 最終値まで間接参照して判定する
4. Where カスタムメッセージが指定されない場合, the system shall フィールド名を含むデフォルト文言を利用する
5. The system shall 文字列・配列・スライス・数値・bool を空値判定の対象として扱う

### Requirement 4: 構造体間の自動バインドと検証連携
**Objective:** As a アプリケーション開発者, I want `voauto.BindAndValidate` で任意レイヤーの構造体間変換を実行できること, so that 変換ロジックと検証呼び出しを共通化できる

#### Acceptance Criteria
1. When `vctag` が設定されているとき, the system shall タグ指定の constructor key と source field 名を使って変換する
2. When `vctag` が未設定のとき, the system shall `New<FieldName>` 規約で constructor key を解決する
3. If constructor が未登録のとき, then the system shall `constructor not found` エラーを返す
4. While constructor 実行中に検証エラーが追加された場合, the system shall `ValidationAggregateError` を返して結果オブジェクト返却を失敗させる
5. The system shall すべての対象フィールド変換成功時に、利用レイヤーを問わず型付き結果 `*T` を返す

### Requirement 5: コンストラクタ登録コード生成
**Objective:** As a ライブラリ利用者, I want `go:generate` コメントから登録コードを自動生成できること, so that 手動レジストリ管理のミスを防げる

#### Acceptance Criteria
1. When `voauto-gen` が対象ディレクトリを走査したとき, the system shall `//go:generate voauto-gen` コメントを検出する
2. When 対象メソッド情報を解析できたとき, the system shall パッケージ単位で `registry_init.go` を生成する
3. If 生成コメントが存在しないとき, then the system shall 生成対象なしとして処理を継続し利用方法を標準出力に表示する
4. While `-output` や `-methods` がコメント引数で指定された場合, the system shall 指定値を優先して生成設定へ反映する
5. The system shall 既存 `_test.go` と既存 `registry_init.go` を走査対象から除外する
