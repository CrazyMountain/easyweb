# easyweb
An simple implementation of web server with GO.

## note
测试前需要开启数据库服务，并保证配置文件里的数据库存在

把测试文件放到根目录是为了避免配置文件初始化的不兼容性。
Go test命令在执行时也会执行每个包下面的init函数，若测试文件路径与主程序不一致，会导致初始化的过程中因找不到配置文件而报错。