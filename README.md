#collmz

## 介绍
COLL-MZ项目主要用于采集煎蛋、飞G、妹子图、Xiuren网站，以及本地类似图片、视频等文件，并展示到浏览器中。

## 特别申明
该项目主要是个人学习golang而开发的第一个试水程序，请勿将该项目用于非法用途。

## 特点
* 专为闷骚程序员提供；
* 采集各大妹子图片数据；
* 手动采集、定时采集（2小时进行一次）；
* 在浏览器快速浏览相关采集数据；
* 可整理本地文件、视频、漫画、文本等数据；
* sqlite3开放式数据库，可自行构建访问，方便二次开发；
* 可根据具体需求，构建其他网站的采集程序；
* 纯Golang实现。

## 界面预览
### 浏览界面

## 使用方法
1、下载项目到本地任意文件；

2、运行collmz-server-..exe文件；

3、通过浏览器访问http://localhost:8888

4、可以看到项目，可在./config/config.json文件内自行修改端口。

5、初始用户名：admin@admin.com，密码：adminadmin

## 代码编译环境搭建步骤
1、安装golang语言运行环境，配置好环境变量；

2、安装gcc编译环境，并配置好环境变量，推荐使用mingw，下载地址：https://sourceforge.net/projects/mingw-w64/

3、安装golang第三方库：

    * goquery
    github.com/PuerkitoBio/goquery
    * sqlite3
    github.com/mattn/go-sqlite3
    * session
    github.com/gorilla/sessions

4、下载该项目代码，到golang工作目录中任意目录，建议使用git克隆。

5、因为是在win10 x64下开发、编译的，所以只能保证该环境下运行良好，其他环境请自行排错。

## 项目地址
GIthub：https://github.com/fotomxq/coll-mz

OSchina：https://git.oschina.net/fotomxq/collmz

## 开发日志

2016.12.9 根据代码缺陷，着手开发2.0版本。

2016.12.8 完成coll-mz 1.0版本。

2016.11.5 立项开发。

## 项目协议
Apache License

Version 2.0, January 2004

http://www.apache.org/licenses/

## FAQ 使用问题

## FAQ 二次开发问题