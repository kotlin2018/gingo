1、先下载该项目
````
git clone 
````

2、编译该项目 
````
go build 
````

3、将编译后生成的可执行命令文件 gingo  放到 GOPATH/bin目录下
````
mv gingo ${GOPATH}/bin
````

4、在任意目录下新建项目 (推荐: 在GOPATH/src目录下新建项目)
````
gingo new 项目名称 
````
5、在新建项目根目录下初始化
````
go mod init 项目名称
````
6、编译新项目
````
go build
````
