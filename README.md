# twSdrPower
This sensor monitors radio power by frequency and transmits it via syslog.

TWSNMP FCのための周波数別の無線電力センサーです。

[![Godoc Reference](https://godoc.org/github.com/twsnmp/twSdrPower?status.svg)](http://godoc.org/github.com/twsnmp/twSdrPower)
[![Go Report Card](https://goreportcard.com/badge/twsnmp/twSdrPower)](https://goreportcard.com/report/twsnmp/twSdrPower)

## Overview

RTL-SDRを利用して周辺の電波の強度をモニタし情報をTWSNMP FCなどへ  
syslogで送信するためのセンサープログラムです。  
現在のバージョンでは以下の情報を取得できます。

- 指定範囲(デフォルト24Mhz -1.67GHz)の指定単位(デフォルト1MHz）の電波の強度情報
- センサーのリソース
- 統計情報

取得した電波強度をグラフ出力することもできます。

## Status

最初のバージョン(v1.0.0)をリリース

## Build
### Env
ビルドするためには、以下が必要です。

- go 1.17
- librtlsdr
- docker(Linux版のビルド)
- make

RTL-SDRのライブラリはMac OSの場合
https://formulae.brew.sh/formula/librtlsdr
でインストールしました。
Linux版は、Docker環境の中でビルドするのでmakeとDokcerだけでビルドできます。

### Build
ビルドはmakeで行います。
```
$make
```
以下のターゲットが指定できます。
```
  all        全実行ファイルのビルド（省略可能）
  mac        Mac用の実行ファイルのビルド
  clean      ビルドした実行ファイルの削除
  zip        リリース用のZIPファイルを作成
```

```
$make
```
を実行すれば、MacOS,Windows,Linux(amd64),Linux(arm)用の実行ファイルが、  
`dist`のディレクトリに作成されます。


配布用のZIPファイルを作成するためには、
```
$make zip
```
を実行します。ZIPファイルが`dist/`ディレクトリに作成されます。

## Run
### 使用方法

```
Usage of ./dist/twSdrPower.app:
  -chart string
    	chart title
  -dark
    	dark mode chart
  -end string
    	end frequency (default "1667M")
  -folder string
    	chart folder (default "./")
  -gain int
    	RTL-SDR Tuner gain (0=auto)
  -interval int
    	syslog send interval(sec) (default 600)
  -list
    	List RTL-STR
  -once
    	Only once
  -sdr int
    	RTL-SDR Device Number
  -start string
    	start frequency (default "24M")
  -step string
    	step frequency (default "1M")
  -syslog string
    	syslog destnation list
```

syslogの送信先はカンマ区切りで複数指定できます。  
:に続けてポート番号を指定することもできます。

```
-syslog 192.168.1.1,192.168.1.2:5514
```

起動するためにはsyslogの送信先(-syslog)が必要です。

Mac OS,Windows,Linuxの環境では以下のコマンドで起動できます。  
（例はMac OS場合）

```
%twSdrPower.app -chart noise -gain 500  -dark  -folder /tmp -interval 300 -sdr 1 -syslog 192.168.1.250
```

### デバイス番号の確認
 -list  オプションを付けて起動でます。

```
 % twSdrPower.app -list
Device List count=1
0,Generic RTL2832U OEM,Realtek,RTL2838UHIDIR,00000001
``
先頭の0がデバイス番号です。

## Copyright

see ./LICENSE

```
Copyright 2022 Masayuki Yamai
```

