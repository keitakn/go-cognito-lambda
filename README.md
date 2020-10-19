# go-cognito-lambda
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
export NEXT_IDAAS_SERVER_CLIENT_ID=クライアントシークレットを安全に保管出来るサーバーサイドアプリケーション用のUserPoolClientIDを指定
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

`make test`

## ソースコードのformat

`make format`

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
