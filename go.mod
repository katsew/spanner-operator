module github.com/katsew/spanner-operator

go 1.12

require (
	cloud.google.com/go v0.40.0
	github.com/evanphx/json-patch v4.4.0+incompatible // indirect
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/golang/groupcache v0.0.0-20190129154638-5b532d6fd5ef // indirect
	github.com/google/gofuzz v1.0.0 // indirect
	github.com/googleapis/gnostic v0.2.0 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/json-iterator/go v1.1.6 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/labstack/gommon v0.2.8
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/valyala/fasttemplate v1.0.1 // indirect
	go.opencensus.io v0.22.0 // indirect
	golang.org/x/exp v0.0.0-20190510132918-efd6b22b2522 // indirect
	golang.org/x/net v0.0.0-20190607181551-461777fb6f67 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sys v0.0.0-20190608050228-5b15430b70e3 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	golang.org/x/tools v0.0.0-20190608022120-eacb66d2a7c3 // indirect
	gonum.org/v1/gonum v0.0.0-20190608115022-c5f01565d866 // indirect
	google.golang.org/api v0.6.0
	google.golang.org/appengine v1.6.1 // indirect
	google.golang.org/genproto v0.0.0-20190605220351-eb0b1bdb6ae6
	google.golang.org/grpc v1.21.1
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.0.0-20190602205700-9b8cae951d65
	k8s.io/apimachinery v0.0.0-20190607205628-5fbcd19f360b
	k8s.io/client-go v0.0.0-20190531132438-d58e65e5f4b1
	k8s.io/code-generator v0.0.0-20190531131525-17d711082421
	k8s.io/gengo v0.0.0-20190327210449-e17681d19d3a // indirect
	k8s.io/klog v0.3.2
	k8s.io/kube-openapi v0.0.0-20190603182131-db7b694dc208 // indirect
	k8s.io/kubernetes v1.14.3
	k8s.io/utils v0.0.0-20190607212802-c55fbcfc754a // indirect
)

replace (
	golang.org/x/sync => golang.org/x/sync v0.0.0-20181108010431-42b317875d0f
	golang.org/x/sys => golang.org/x/sys v0.0.0-20190209173611-3b5209105503
	golang.org/x/tools => golang.org/x/tools v0.0.0-20190313210603-aa82965741a9
	k8s.io/api => k8s.io/api v0.0.0-20190531132109-d3f5f50bdd94
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190531131812-859a0ba5e71a
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190531132438-d58e65e5f4b1
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20190531131525-17d711082421
)

replace k8s.io/component-base => k8s.io/component-base v0.0.0-20190531133342-103ccccb7a11
