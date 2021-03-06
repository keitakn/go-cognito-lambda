# go-cognito-lambda
![ci-master](https://github.com/keitakn/go-cognito-lambda/workflows/ci-master/badge.svg)
![cd-master](https://github.com/keitakn/go-cognito-lambda/workflows/cd-master/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/keitakn/go-cognito-lambda/badge.svg?branch=master)](https://coveralls.io/github/keitakn/go-cognito-lambda?branch=master)

CognitoUserPoolをトリガーとしたLambdaのサンプル色々

# Getting Started

## 環境変数の設定

以下の環境変数を設定して下さい。

[direnv/direnv](https://github.com/direnv/direnv) 等を利用するのがオススメです。

```
export DEPLOY_STAGE=デプロイターゲット（.eg. dev, stg, prod）
export TARGET_USER_POOL_ID=ターゲットとなるUserPoolのID
export TRIGGER_USER_POOL_NAME=ターゲットとなるUserPoolの名前
export REGION=AWSのリージョン（.eg. ap-northeast-1）
export API_DOMAIN_NAME=API Gatewayに設定するドメイン名を指定、予めRoute53にホストゾーンが設定されている必要があります
export CERTIFICATE_ARN=AWS Certificate ManagerのARNを指定、 "*.ドメイン名" で指定した証明書は利用出来ないので注意
export NEXT_IDAAS_SERVER_CLIENT_ID=クライアントシークレットを安全に保管出来るサーバーサイドアプリケーション用のUserPoolClientIDを指定
export DYNAMODB_TEST_ENDPOINT=テスト用のDynamoDBのエンドポイントを指定、ローカルで実行する時は http://localhost:58000 を指定
```

## AWSクレデンシャルの設定

従って以下のように [名前付きプロファイル](https://docs.aws.amazon.com/ja_jp/cli/latest/userguide/cli-configure-profiles.html) を作成して下さい。

`~/.aws/credentials`

```
[nekochans-dev]
aws_access_key_id=YOUR_AWS_ACCESS_KEY_ID
aws_secret_access_key=YOUR_AWS_SECRET_ACCESS_KEY
```

無論このプロファイル名は好きな名前に変えてもらって問題ありません。

その場合は `serverless.yml` 内の `custom.profiles` を全て修正して下さい。

## Goのインストール

`go1.15` をインストールします。

## Node.jsのインストール

最新安定版をインストールします。

## npm packageのインストール

`npm ci` を実行してpackageをインストールします。

## Dockerで必要なコンテナを生成する

`docker-compose up -d` で必要なコンテナを起動して下さい。

ソースコードのフォーマットやDynamoDBを使ったテストはDocker上でないと正常動作しません。

# デプロイ関連のコマンド

## Build & Deploy

`make deploy` を実行すると `build` , `deploy` が実行されます。

deployは [Serverless Framework](https://www.serverless.com/) を利用しています。

このツールを利用すると、既存のCognitoUserPoolに対してLambda関数をアタッチ出来るので、その機能を利用する事が主な目的です。

それ以外にも公式の [AWS SAM](https://docs.aws.amazon.com/ja_jp/serverless-application-model/latest/developerguide/serverless-sam-reference.html) と比較して痒いところに手が届くので、その点も良いと思います。

- （参考）[Serverless Frameworkの使い方まとめ](https://qiita.com/horike37/items/b295a91908fcfd4033a2)

## deployしたリソースを削除する

`make remove` を実行します。

# その他のコマンド

## テスト実行

Goのコンテナ内で `make test` を実行します。

`docker-compose exec go make test` でも大丈夫です。

## ソースコードのformat

Goのコンテナ内で `make format` を実行します。

`docker-compose exec go make format` でも大丈夫です。

## Lintの実行

`make format` では修正出来ないエラー内容を表示します。

Goのコンテナ内で `make lint` を実行します。

`docker-compose exec go make lint` でも大丈夫です。

# 開発を行う為の参考資料

Cognitoをカスタマイズする為のLambdaは以下の種類が存在します。

- カスタム認証フロー
- 認証イベント
- サインアップ
- メッセージ
- トークンの作成

詳しくは [こちら](https://docs.aws.amazon.com/ja_jp/cognito/latest/developerguide/cognito-user-identity-pools-working-with-aws-lambda-triggers.html) を見て下さい。

また `serverless.yml` にトリガーにCognitoのイベントを設定する必要があります。

それに関しては下記のドキュメントが参考になります。

- https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#cognito
- https://docs.aws.amazon.com/ja_jp/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpool-lambdaconfig.html

# APIの認証・認可について

本リポジトリで実装されているAPIは CognitoUserPool が発行するJWTトークンによって保護されています。

`serverless.yml` の `httpApi.authorizers` の設定次第ですが、ここでは `Client Credentials Grant(RFC 6749)` の仕組みでアクセストークンを発行する例を紹介します。

具体的な手順は下記の通りです。

## Client Credentials Grant(RFC 6749) でのアクセストークン発行方法

### 1. アプリクライアントIDとアプリクライアントのシークレットを `:` で繋いでBase64Encodeした値を生成する

`echo -n "【アプリクライアントID】:【アプリクライアントのシークレット】" | base64` で生成を行います。

仮に 【アプリクライアントID】が `aaa`, アプリクライアントのシークレットが `bbb` なら以下のようになります。

```
echo -n "aaa:bbb" | base64

# 実際にはもっと長い値が生成されます
YWFhOmJiYg==
```

### 2. 1で生成した値を使ってアクセストークンを発行する

以下のようにCognitoUserPoolのトークンエンドポイントに対してリクエストを行います。

```
curl -v \
-X POST \
-H "Content-Type: application/x-www-form-urlencoded" \
-H "Authorization: Basic {1で生成した値を指定する}" \
--data "grant_type=client_credentials" \
--data "scope={DEPLOY_STAGE}-cognito-admin-api.keitakn.de/admin" \
https://{CognitoUserPoolドメイン名}.auth.ap-northeast-1.amazoncognito.com/oauth2/token
```

- grant_typeは `client_credentials` 固定です
- scopeの `{DEPLOY_STAGE}` にはデプロイ時に利用している `{DEPLOY_STAGE}` の値を利用します。（`-cognito-admin-api.keitakn.de/admin` の部分は固定です）
- `CognitoUserPoolドメイン名` に関してはAWSコンソール上のCognitoUserPoolの「ドメイン名」からご確認下さい。

成功すると以下のようなリクエストが返ってきます。

```json
{
  "access_token": "JWT形式のトークン",
  "expires_in": 3600,
  "token_type": "Bearer"
}
```

### 3. アクセストークンを用いて、APIにリクエストを行う

以下のように `Authorization: Bearer + アクセストークン` を指定してリクエストを行います。

```
curl -v \
-X PATCH \
-H "Content-type: application/json" \
-H "Authorization: Bearer {2で取得したアクセストークンを指定}" \
-d \
'
{
  "userName": "対象のcognitoUsernameを指定",
  "newPassword": "新しいパスワード"
}
' \
https://${API_DOMAIN_NAME}/users/passwords | jq
```

トークンが有効な間は正常に200系のレスポンスが返ってきます。

トークンが無効、もしくは有効期限切れの場合は以下のようなレスポンスが返ってきます。

```
< HTTP/2 401
< date: Thu, 22 Oct 2020 03:10:39 GMT
< content-length: 26
< www-authenticate: Bearer scope="dev-cognito-admin-api.keitakn.de/admin" error="invalid_token" error_description="the token has expired"
< apigw-requestid: xxxxxxx=
<
{ [26 bytes data]
100   128  100    26  100   102    106    418 --:--:-- --:--:-- --:--:--   524
* Connection #0 to host xxxxx.execute-api.ap-northeast-1.amazonaws.com left intact
* Closing connection 0

{
  "message": "Unauthorized"
}
```

# コーディング規約

Goの標準的な慣例に従います。

- ファイル名はスネークケース
- ディレクトリ名は小文字のみ利用する形する
- 構造体、変数名はキャメルケース
- 引数、レシーバ名はなるべく1文字の単語を使用する
- package名は小文字のみを利用する

1つ例外があります。

`url` は `URL` を使うとか、 `api` は `API` を使う等のルールは採用していません。

理由としては下記の通りです。

- [golangci-lint](https://github.com/golangci/golangci-lint) でこのルールをチェック出来るが、カバーされている単語が少ない
- どれが略語なのかをその都度判断するのが難しい、人によっては出来たり出来なかったりするので、かえって統一感のないコードが出来上がってしまう可能性が高い

個人で書いた記事ですが [こちら](https://zenn.dev/keitakn/articles/go-naming-rules#%E9%A0%AD%E6%96%87%E5%AD%97%E3%82%84%E9%A0%AD%E6%96%87%E5%AD%97%E3%82%92%E3%81%A8%E3%81%A3%E3%81%9F%E5%90%8D%E5%89%8D%E3%81%AE%E5%8D%98%E8%AA%9E%E3%81%AE%E5%91%BD%E5%90%8D%E3%81%AB%E3%81%A4%E3%81%84%E3%81%A6) にさらに詳しい理由が書いてあります。

## `.go` 以外のファイル名に関するルール

- HTMLファイル（慣例に従いケバブケース）
- go buildで生成したファイル（package名と同様にしたいので、小文字のみを利用）
