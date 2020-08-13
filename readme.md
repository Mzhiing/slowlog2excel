bin目录下的slowlog2excel可以直接运行

运行方式是：

chmod +x slowlog2excel  

例如：

slowlog2excel -s /usr/local/mysql/data/log/ -f slow.log -e 20200813.xlsx

其中：

-s 后跟slow.log放置的路径

-f 慢日志的文件名

-e 生成的xlsx的文件名

运行上述命令就会在slowlog2excel同级目录生成 excel文件

