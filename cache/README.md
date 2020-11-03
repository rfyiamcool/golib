# go-cache

该库基于[ccache](github.com/karlseguin/ccache/v2)二次封装, 这个库的存储采用了分片方式的map，对于过期和maxSize的事件采用异步化处理，减少锁竞争及同步时延.

- 更简单的api
- 通过锁池解决了并发缓存击穿问题
- 实例化全局默认的cache对象

## usage

[单元测试](cache_test.go)
