# 快速体验 Kubebuilder
## Kubebuilder 安装
```shell
os=$(go env GOOS)
arch=$(go env GOARCH)

# 下载 kubebuilder 并解压到 tmp 目录中
curl -L https://go.kubebuilder.io/dl/2.3.1/${os}/${arch} | tar -xz -C /tmp/

# 将 kubebuilder 移动 PATH 路径中 
sudo mv /tmp/kubebuilder_2.3.1_${os}_${arch} /usr/local/kubebuilder
echo "export PATH=\$PATH:/usr/local/kubebuilder/bin" >> /root/.bashrc
# kubebuilder 压缩包中，包括 etcd、kube-apiserver、kubebuilder、kubeclt 这四个二进制文件

# 准备 kustomize 工具，kubebuilder 依赖 kustomize 工具构建 manifests，kubeclt 中的 kustomize 子命令功能补全
curl -LO https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv3.8.7/kustomize_v3.8.7_linux_amd64.tar.gz
tar -xzvf kustomize_v3.8.7_linux_amd64.tar.gz -C /usr/local/bin/
```

## 创建一个 hello world 项目
```shell
mkdir /root/project/kubernetesAPI/helloWorld && cd /root/project/kubernetesAPI/helloWorld
kubebuilder init --domain desistdaydream.ltd --owner DesistDaydream --repo helloWorld
# 创建完成后，会生成如下文件
[root@lichenhao helloWorld]# ls
config  Dockerfile  go.mod  go.sum  hack  main.go  Makefile  PROJECT
```
* config # 是各种 mainfests 文件
* Dockerfile # 用于构建该项目

## 创建一个 API
下面的命令将会创建一个新的 API，组/版本 为 `test/v1`，并在其上创建一个**新的 `Kind` 为 HelloWorld**。
>这所谓的新 Kind 其实就是 Custom Resource(自定义资源，简称 CR)
```shell
[root@lichenhao helloWorld]# kubebuilder create api --group test --version v1 --kind HelloWorld
Create Resource [y/n]
y
Create Controller [y/n]
y
Writing scaffold for you to edit...
api/v1/helloworld_types.go
controllers/helloworld_controller.go
Running make:
$ make
/opt/gopath/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
go build -o bin/manager main.go
# 创建完成后，会生成如下文件
[root@lichenhao helloWorld]# ls
api  bin  config  controllers  Dockerfile  go.mod  go.sum  hack  main.go  Makefile  PROJECT
```
* api 目录 # 由 Create Resource 操作操作，其中的 `v1/helloworld_types.go` 文件。用于定义 api，主要是定义这个 CR 创建出来的对象应该具有哪些字段，比如 spec 下应该有 containers、volumes 之类的。
* controllers 目录 # 由 Create Controller 操作会创建，其中的 `helloworld_controller.go` 文件。用于定义 Operator 的执行逻辑。具体 Operator 都干什么的代码一般都在 controllers 目录中

## 测试
首先需要准备 kubectl 命令连接集群所需的 config 文件，以确保可以正常访问集群。
使用 `make install` 命令即可将 CRD 安装到集群中
```shell
[root@lichenhao helloWorld]# make install
/opt/gopath/bin/controller-gen "crd:trivialVersions=true" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
kustomize build config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/helloworldren.test.desistdaydream.ltd created
# 此时 k8s 集群中，就多出了一个 CRD
# 注意，helloworld 项目创建的 CRD 名称后面多了个 ren，是因为 加上 ren 表示为 复数，自动添加的
# 可以修改这个 crd 的 yaml 中的 .spec.names.plural 字段，来改变复数的名字
[root@lichenhao helloWorld]# kubectl get crd | grep hello
helloworldren.test.desistdaydream.ltd       2020-11-23T14:51:07Z
```

然后使用 `make run` 命令可以运行这个 **自定义控制器**
```shell
[root@lichenhao helloWorld]# make run
/opt/gopath/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
/opt/gopath/bin/controller-gen "crd:trivialVersions=true" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
go run ./main.go
2020-11-23T22:55:43.197+0800	INFO	controller-runtime.metrics	metrics server is starting to listen	{"addr": ":8080"}
2020-11-23T22:55:43.198+0800	INFO	setup	starting manager
2020-11-23T22:55:43.198+0800	INFO	controller-runtime.manager	starting metrics server	{"path": "/metrics"}
2020-11-23T22:55:43.199+0800	INFO	controller-runtime.controller	Starting EventSource	{"controller": "helloworld", "source": "kind source: /, Kind="}
2020-11-23T22:55:43.299+0800	INFO	controller-runtime.controller	Starting Controller	{"controller": "helloworld"}
2020-11-23T22:55:43.299+0800	INFO	controller-runtime.controller	Starting workers	{"controller": "helloworld", "worker count": 1}
```

make run 会保持控制器在前台运行，并始终 watch 着 hellworld 这个 CR。在 `config/samples/test_v1_helloworld.yaml` 文件中，是一个创建该 CR 的 yaml 文件。
```shell
[root@lichenhao helloWorld]# kubectl apply -f config/samples/test_v1_helloworld.yaml
helloworld.test.desistdaydream.ltd/helloworld-sample created
# 这时候，就会创建了一个自定义资源
[root@lichenhao helloWorld]# kubectl get helloworld -A
NAMESPACE   NAME                AGE
default     helloworld-sample   50s
# 并且上面在前台运行的控制器会收到对应的消息，消息内容如下：
2020-11-23T23:01:32.265+0800	DEBUG	controller-runtime.controller	Successfully Reconciled	{"controller": "helloworld", "request": "default/helloworld-sample"}
```

这时候，可以构建并推送镜像到仓库，然后将 控制器 部署到集群中。
```shell
# 这里需要注意，Dockerfile 中的镜像被墙了，需要找别的办法提前 pull 下来
# 并且其中有一步 是 go mod download 由于是在容器中，无法读取到外部的环境变量，所以需要修改一下，添加 go 代理，否则超时
# RUN GOPROXY=https://goproxy.io go mod download
# 后面的多阶段构建所需的镜像也被墙了，需要提前下载，可以从 registry.aliyuncs.com/byteforce/distroless:nonroot 这里下载
make docker-build docker-push IMG=lchdzh/helloworld:operator-v1
# 这时直接使用 make 即可部署 控制器
[root@lichenhao helloWorld]# make deploy IMG=lchdzh/helloworld:operator-v1
/opt/gopath/bin/controller-gen "crd:trivialVersions=true" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
cd config/manager && kustomize edit set image controller=lchdzh/helloworld:operator-v1
kustomize build config/default | kubectl apply -f -
namespace/helloworld-system created
customresourcedefinition.apiextensions.k8s.io/helloworldren.test.desistdaydream.ltd configured
role.rbac.authorization.k8s.io/helloworld-leader-election-role created
clusterrole.rbac.authorization.k8s.io/helloworld-manager-role created
clusterrole.rbac.authorization.k8s.io/helloworld-proxy-role created
clusterrole.rbac.authorization.k8s.io/helloworld-metrics-reader created
rolebinding.rbac.authorization.k8s.io/helloworld-leader-election-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/helloworld-manager-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/helloworld-proxy-rolebinding created
service/helloworld-controller-manager-metrics-service created
deployment.apps/helloworld-controller-manager created
```

## 清理
使用 `make uninstall` 命令即可将 CRD 从集群中删除
使用 `kustomize build config/default | kubectl delete -f -` 命令即可将 控制器 以及 CRD从集群中删除

