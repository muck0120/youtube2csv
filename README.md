# YouTube2CSV (YouTube Data Export Tool)

このツールは、YouTube Data API を使って動画情報を取得し、CSV ファイルにエクスポートするためのアプリケーションです。主な機能は以下の通りです。

- 指定した YouTube チャンネル ID の動画情報一覧を取得。
- 動画のタイトル、URL、動画の長さを取得。
- 取得したデータを CSV ファイルとして保存。

## 📲 事前準備

YouTube2CSV を使用するために、以下の 2 点の事前準備を行なってください。

1. [**Go のインストール**](https://go.dev/)
2. [**YouTube API の認証情報を作成**](https://developers.google.com/youtube/v3/live/guides/auth/server-side-web-apps#creatingcred)

尚、「2.」で認証情報を作成した上で、シークレット情報の JSON ファイルをダウンロードしてください。このファイルは YouTube2CSV を使用するために必要です。

## 📝 使い方

### 1. 以下のコマンドを実行

引数は適宜変更してください（各フラグの詳細については [こちら](#フラグについて) を参照）。

```sh
go run github.com/muck0120/youtube2csv --secret="path/to/client_secret.json" --token="path/to/token.json" --channel-id="Target Channnel ID" --out="path/to/output.csv"
```

### 2. 表示された URL にアクセス

認証用の `token.json` ファイルが存在しない、またはトークンの期限が切れている場合は以下のように URL が表示されますので、表示された URL にアクセスしてください。

```sh
👇 Please access the displayed URL to obtain the authentication code and enter it.

https://accounts.google.com/o/oauth2/auth?access_type=offline&client_id=XXXXX&redirect_uri=xxxxx&response_type=code&scope=https://www.googleapis.com/auth/youtube.readonly&state=state-token
```

### 3. 認証トークンを取得する

上記の URL にアクセスした後は以下の手順に沿って、認証トークンを取得し、コンソールに貼り付けてください。

1. Google アカウントの選択画面になるので、「📲 事前準備」で認証情報を作成した Google アカウントを選択します。
2. 「このアプリは Google で確認されていません」という表示が出ますが、「続行」をクリックします。
3. さらに「youtube2csv が Google アカウントへのアクセスを求めています」と出るので「続行」をクリックします。
4. すると `http://localhost/?state=state-token&code=xxxxx&scope=https://www.googleapis.com/auth/youtube.readonly` という URL にリダイレクトされるので `code=xxxxx` の `xxxxx` の部分をコピーしてください。
5. 「4.」でコピーした `xxxxx` を `👇 Please access the displayed URL to obtain the authentication code and enter it.` に続いてコンソールにペーストして Enter をクリックします。

上手く認証が通れば `--out=` で指定したファイルに動画情報の一覧が書き出されます。

### フラグについて

以下のフラグを使用できます。

| フラグ       | 必須 | デフォルト                                | 内容                              |
| :----------- | :--: | :---------------------------------------- | :-------------------------------- |
| --secret     |      | ./client_secret.json                      | YouTube API の認証情報ファイル    |
| --token      |      | ./token.json                              | 認証トークンの保存ファイル        |
| --channel-id |  ✅  |                                           | 動画情報を取得したいチャンネル ID |
| --out        |      | ./output/{channel-id}\_yyyyMMddhhmmss.csv | エクスポート先ファイル            |

### 出力フォーマット

CSV には以下の項目が出力されます。

| 項目       | キー     | 内容                                                                                                   |
| :--------- | :------- | :----------------------------------------------------------------------------------------------------- |
| 番号       | no       | 動画のタイトル内に `#` と、それに続く数字が存在すれば、その数字をセットします。無ければ 0 が入ります。 |
| タイトル   | title    | 動画のタイトルが入ります。                                                                             |
| 動画の長さ | duration | 動画長さが分単位で入ります。分未満は四捨五入されます。                                                 |
| 動画 URL   | url      | 動画の URL が入ります。                                                                                |

## 🛠️ 開発方法

先に「📲 事前準備」を済ませた上で、リポジトリをクローンして自由にカスタマイズして使用してください。

尚、クローンした後に「📲 事前準備 - 2.」で取得した `client_secret.json` を `./input` ディレクトリ内に配置してください（`--secret` の引数を指定する場合は不要です）。

```sh
git clone https://github.com/muck0120/youtube2csv.git
cd youtube2csv
go mod download
```

### 各コマンド

| コマンド    | 引数       | 内容                                                               |
| :---------- | :--------- | :----------------------------------------------------------------- |
| `make run`  | CHANNEL_ID | YouTube2CSV を実行します。取得したいチャンネルの ID を指定します。 |
| `make gen`  |            | `go generate` を実行します。モックファイルの作成などに使用します。 |
| `make lint` |            | `golangci-lint` を実行します。lint のチェックに使用します。        |
| `make test` |            | `go test` を実行します。ユニットテストの確認に使用します。         |
| `make tidy` |            | `go mod tidy` を実行します。                                       |

## 🪪 ライセンス

このプロジェクトは [MIT ライセンス](./LICENSE) のもとで公開されています。
