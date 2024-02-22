# filewatcher

一个基于fsnotify的只用于图片的文件系统事件监听库,使用rxgo实现事件过滤。

## 要解决的问题
- fsnotify 监听单一路径后，不会监听其子文件夹。若子文件夹发生变化，不会触发监听事件。
    - 分析：需要手动写一个方法去监听子目录，每次监听到新增文件都判断是否为文件夹，若是文件夹则加入到监听目录中。此时还有一种情况不好处理：剪切一个文件夹到目录中，该文件夹里面包含子文件夹，此时则不会监听到子文件夹（复制会监听到，剪切不会）。

- 当一个文件拷贝到监听路径后，偶尔会多次监听到事件 event:WRITE 和 modified file 。
  经多次测试 1M的图片会触发四次 event:WRITE 操作
  新增文件会触发多次监听，删除文件不会。
    - 分析：在fsnotify的issue中提到了这个问题，是因为一个“写”操作可能会进行多次写入导致的，建议的解决方法为写事件等待很短的时间，为每个新事件重新设置等待时间。

- 在尝试对子目录进行监听后，会发现监听到的时间有大量关于文件夹的事件，这些事件属于噪声，不需要处理。另外此时无法对父目录进行删除和重命名操作，因为子目录的监听进程占用了该目录，尝试删除父目录会报错。
    - 分析：如果要删除得从子目录删除，然后再逐级删除父目录，重命名是无法进行的。不过该项目中的监听目录属于内部目录，不会进行删除和重命名操作，所以不需要处理。

经过以上分析，发现直接采用fsnotify监听文件会存在子文件夹不触发监听以及漏掉文件的问题，这个在我的Java项目中也有存在，当时最终的解决方案是基于org.apache.commons.io.monitor包实现的监听，该包会每隔一段时间进行轮询操作，检测所监听文件夹下所有文件是否发生了变化，该方式可以参考。
所以在此项目中也采用了类似的方式，基于文件快照来实现文件变更监听。

- 基于文件快照也需要解决一个问题：如何过滤重复事件？每一次文件变更都会触发多次事件，如果不过滤会多次处理浪费计算资源。
    - 分析：设置延迟处理，并在监听到新事件时取消前面的任务执行，这样同时上传多个文件最终只会在上传结束后执行一次操作。

- 考虑极端情况，一张巨大的图片被拷贝到监听目录中，在刚上传时便会进行回调，后续等待两分钟之后才完全上传完，怎么处理这种情况
  - 分析：fsnotify在文件写完时会触发一次WRITE事件，所以在写完之前的那次WRITE事件可以不用处理，遇到异常直接忽略即可，文件在写完之后会触发一次WRITE事件，这时可正常处理。

- 基于文件快照可能会忽略文件修改事件，因为文件修改事件不会触发文件快照的更新。
    - 分析：采用文件快照和fsnotify双重监听，如果fsnotify监听到了修改但是文件快照没有监听到，说明是文件修改事件，此时依然进行回调处理。

## 特性

- 基于fsnotify和文件快照监听文件系统事件
- 基于文件快照实现自动监听子目录
- 基于rxgo注册回调接口处理文件变化
- 缓冲并批量处理文件事件
- 确保文件上传完成后有回调
- 完整的示例程序和测试代码
- 优雅退出