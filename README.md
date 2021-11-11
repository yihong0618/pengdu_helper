# pengdu_helper

网易云歌词自动生成《捧读》用的 epub, 帮助学习。（GitHub Actions + shortcuts + telegram）

# 如何使用

## 前提

1. 有访问自由互联网的环境，有 telegram 请自行查询如何写 telegram bot
2. 有《捧读》app
3. fork 或者 clone 这个项目，自行配置需要的 secrets

## 之后

1. 参考我的博客 -- [巧妙利用 iOS 的快捷指令配合 GitHub Actions 实现自动化](https://github.com/yihong0618/gitblog/issues/198) 拿到 action id
2. 配置 iOS 的捷径，捷径如图

![image](https://user-images.githubusercontent.com/15976103/132662709-aabc32cc-9629-4ac4-b9c3-67b2feba13ce.png)
![image](https://user-images.githubusercontent.com/15976103/132666328-3d5eaf91-15f2-49d2-aa40-e0d02f71dc1f.png)


3. 复制网易云链接如下图，点击复制链接

![image](https://user-images.githubusercontent.com/15976103/132662918-cd7b223d-031d-4e9b-aecd-463130323280.png)

4. 跑捷径，Actions 就自动运行了，之后会生成 epub, 用《捧读》打开即可

# 效果 --> 就可以用《捧读》学习你喜欢的音乐啦

![ce7f27caa5683901f90c9e164cfe0fe](https://user-images.githubusercontent.com/15976103/132663493-5132ef3c-b376-4db7-ac50-e9880ea97567.jpg)


# 如果网页端的网易云复制链接跑以下命令即可

e.g.
```shell
go run main.go -nestring https://music.163.com/#/song?id=112143
```

# 利用 srt or ass 生成 epub 交给捧读

```shell
go run main.go -filestring {$something}.srt 

# 带时间
go run main.go -filestring {$something}.srt -withtime true
```