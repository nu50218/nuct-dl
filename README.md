# nuct-dl

nuct-dlはnuctからwebdavを使って簡単にダウンロードしてこれるコマンドラインツールです。

## インストール

Goを[インストール](https://golang.org/doc/install)

`$ go get -u github.com/nu50218/nuct-dl`

## 使い方

`$ nuct-dl -id=[サイトID] -user=[ユーザー名] -pass=[パスワード]`

## オプション

### -output=[ダウンロード先ディレクトリ名]

デフォルトではサイトIDのディレクトリが作られてそこにダウンロードされます。

### -last_update=[時間（24h, 15dなど）]

指定した時間内に更新されたものだけダウンロードします。
