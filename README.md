# 项目初始化
```
go mod init kubeants.io
go get -u github.com/gin-gonic/gin@v1.9.1
```

# 修改项目域名
1. 修改 go.mod 文件
```
module kubeants.io
# 改为目前域名
module kubeants.io
```
2. 批量替换代码中的 import 语句 运行以下命令批量替换 import "kubeant.cn 为 import "kubeants.io：
```
grep -rl 'kubeants.io' . | xargs sed -i 's|kubeants.io|kubeants.io|g'
# MacOS 用户请使用：
grep -rl 'kubeant.cn' . | xargs sed -i '' 's|kubeant.cn|kubeants.io|g'
```
3. 执行 go mod tidy 重新整理依赖
```
go mod tidy
go build
```


# 提交到仓库
…or create a new repository on the command line
```
echo "# kubeants" >> README.md
git init
git add README.md
git commit -m "first commit"
git branch -M main
git remote add origin https://github.com/geray-zsg/kubeants-apiserver.git
git push -u origin main
```
…or push an existing repository from the command line
```
git remote add origin https://github.com/geray-zsg/kubeants-apiserver.git
git branch -M main
git push -u origin main
```

# 有组名和无组名
获取某deploy，方法GET(普通用户)      
http://{{ks-apiserver}}/apis/clusters/<cluster>/apps/v1/namespaces/<namespace>/deployments/<deployment>
查看指定namespace下的pod，方法GET(普通用户)
http://{{ks-apiserver}}/api/clusters/<cluster>/v1/namespaces/<namespace>/pods

# http的几种请求方式
HEAD	类似 GET，但不返回响应体，只返回头部信息
OPTIONS	用于获取服务器支持的 HTTP 方法列表
PATCH	部分更新资源（区别于 PUT 的整体替换）
TRACE	服务器回显收到的请求（主要用于诊断）
CONNECT	建立隧道（通常用于 HTTPS 代理）

# 查看原生k8s接口信息
```
kubectl proxy

curl http://localhost:8001

# 无组名
curl http://localhost:8001/api/v1/

# 有组名
curl http://localhost:8001/apis/apps/

curl http://localhost:8001/api/v1/namespaces
curl http://localhost:8001/apis/apps/v1/deployments

```