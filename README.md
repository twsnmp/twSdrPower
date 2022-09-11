# twSdrPower
This sensor monitors radio power by frequency and transmits it via syslog.

TWSNMP FCのための周波数別の無線電力センサーです。

[![Godoc Reference](https://godoc.org/github.com/twsnmp/twSdrPower?status.svg)](http://godoc.org/github.com/twsnmp/twSdrPower)
[![Go Report Card](https://goreportcard.com/badge/twsnmp/twSdrPower)](https://goreportcard.com/report/twsnmp/twSdrPower)

## Overview

RTL-SDRを利用して周辺の電波の強度をモニタし情報をTWSNMP FCなどへ  
syslogで送信するためのセンサープログラムです。  
現在のバージョンでは以下の情報を取得できます。

- 24Mhz -1.67GHzの1MHz単位の電波の強度情報
- センサーのリソース

## Status

開発を始めたばかりです。

## Build

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
Usage of ./twSdrPower.app:
Usage of dist/twSdrPower.app:
  -sdr int
    	RTL-SDRのデバイス番号
  -interval int
    	syslog send interval(sec) (default 600)
  -syslog string
    	syslog destnation list
```

syslogの送信先はカンマ区切りで複数指定できます。  
:に続けてポート番号を指定することもできます。

```
-syslog 192.168.1.1,192.168.1.2:5514
```


### 起動方法

起動するためにはsyslogの送信先(-syslog)が必要です。

Mac OS,Windows,Linuxの環境では以下のコマンドで起動できます。  
（例はLinux場合）

```
#./twSdrPower  -syslog 192.168.1.1
```

## Copyright

see ./LICENSE

```
Copyright 2022 Masayuki Yamai
```

