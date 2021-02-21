# slackops

slackの操作をシェルプロっぽくするツール。あとショートカットで文字送ったり。

## 提案

slack使ってるとネイティブアプリがWebアプリぽい操作感なのに、ターミナルがシェルなのでモヤる時ありませんか？

- メッセージを再編集するのに↑キーで履歴を選びたい
- 特定の人へのメンションを多用するので毎度コピペするのが面倒
- 同じ入力をするのにエイリアスみたいなので入力を減らしたい

特にslack→サーバーへコマンド投げ込みボットが居る場合は顕著に面倒になります。<br>
ボットを呼んでコマンド投げて、メンション入れて確認依頼して・・次のコマンド入れて・・をコピペして繰り返すことになります。

### というのをUXを解消してくれるのが、このツール！

## 使い方

```
git clone https://github.com/yasutakatou/slackops
cd slackops
go build slackops.go
```

バイナリをダウンロードして即使いたいなら[こっち](https://github.com/yasutakatou/slackops/releases)

### zlib1.dllが無い場合のエラー

使用しているrobotgoが依存しているのでzlib1.dllが無い場合、エラー落ちが予測されます。<br>
[こちらのリンク](https://sourceforge.net/projects/mingw-w64/files/External%20binary%20packages%20%28Win64%20hosted%29/Binaries%20%2864-bit%29/)からダウンロードしてパスが通ってる場所に配置すると動かせます

## コンフィグファイル

使うのにコンフィグファイルに動作を定義します。サンプルを参考にカスタマイズしてください。

- [TITLE]
	- 入力を投げ込むslackアプリのWindowタイトルを指定します。通常、 **Slack**のままで良いはず
		- Slack
- [SHORTCUT]
	- ショートカットで入力する文字列を**csv**で定義します。以下、例なら **ctrl+Q**でカンマの右が入力されます
		- Q,@admin @boss 
- [CONVERT]
	- 置き換え文字列を**csv**で定義します。以下、例なら **git**が入力された際にカンマの右に置き換えられます
		- git,https://github.com/yasutakatou/slackops
- [SEND]
	- slackの文字列書き込みショートカットです。以下、例なら **ctrl+S**で書き込みされます
		- S
		- ※slackアプリで**returnで書き込みする環境設定**にしてください。**ctr+return**の設定だと動きません！
- [ENTER]
	- slackのテキスト改行ショートカットです。以下、例なら **ctrl+X**でテキストボックス内で改行されます
		- X
- [WAIT]
	- 入力時のウェイトです。PCやネットワーク速度でうまく動かないときに微調整します。単位はミリsec
		- 100
- [SINGLELINE]
	- 一行モードです。**Y**の場合、[SEND]の入力無しで、slackで文字列が書き込まれます。ボットが居る場合などサクサク動かしたい場合に
		- Y
- [DELETE]
	- テキストボックス内の文字列削除ショートカットです。以下、例なら **ctrl+Z**でテキストボックス内の入力がクリアされます
		- Z

### コンフィグ記載例
```
[TITLE]
Slack
[SHORTCUT]
Q,@admin @boss 
W,@teamA @teamB 
E,CC: @teamA @teamB
[CONVERT]
git,https://github.com/yasutakatou/slackops
rm -rf,DONT USE!
[SEND]
S
[ENTER]
X
[WAIT]
100
[SINGLELINE]
Y
[DELETE]
Z
```

## 起動オプション

実行ファイルは以下、起動オプションがあります。

```
>slackops.exe -h
Usage of slackops.exe:
  -config string
        [-config=config file)] (default ".slackops")
  -debug
        [-debug=debug mode (true is enable)]
```

### -config
読み込むコンフィグファイルを指定します。デフォルトは実行ファイルのカレントにある **.slackops** です。

### -debug
デバッグモードで起動します。指定すると内部動作情報が色々出てきます。

## ライセンス
MIT License, Apache License 2.0, GNU General Public License v3.0


