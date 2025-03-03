# 项目初始化
```
go mod init kubeants.com
go get -u github.com/gin-gonic/gin@v1.9.1
```

# 修改项目域名
1. 修改 go.mod 文件
```
module kubeant.cn
# 改为目前域名
module kubeants.com
```
2. 批量替换代码中的 import 语句 运行以下命令批量替换 import "kubeant.cn 为 import "kubeants.com：
```
grep -rl 'kubeant.cn' . | xargs sed -i 's|kubeant.cn|kubeants.com|g'
# MacOS 用户请使用：
grep -rl 'kubeant.cn' . | xargs sed -i '' 's|kubeant.cn|kubeants.com|g'
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
git remote add origin https://github.com/geray-zsg/kubeants.git
git push -u origin main
```
…or push an existing repository from the command line
```
git remote add origin https://github.com/geray-zsg/kubeants.git
git branch -M main
git push -u origin main
```