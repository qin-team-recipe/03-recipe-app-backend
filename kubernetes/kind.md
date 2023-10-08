# kind（Kubernetes in Docker）編

課題：画像ファイルの置き場所をPersistentVolumeにする

## クラスタ作成

foo-cluster.yaml を作成

```foo-cluster.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
```

以下を実行

```shell
# クラスタ作成
kind create cluster --name foo-cluster --config=foo-cluster.yaml
```

## バックエンドのソースコード編集〜Dockerイメージ作成

以下を実行

```shell
git clone https://github.com/qin-team-recipe/03-recipe-app-backend
cd 03-recipe-app-backend
```

app.go 変更

```app.go
-        Addr:     "localhost:6379",
+        Addr:     os.Getenv("REDIS_ADDR"),
```

Dockerfile を作成

```Dockerfile
FROM aopontann/dev_go_opencv:latest as builder
WORKDIR /app
ADD ./ ./
RUN go mod download
RUN go build -o server app.go

FROM aopontann/dev_go_opencv:latest
WORKDIR /app
COPY --from=builder /app/server /app/
CMD ["/app/server"]
```

以下を実行

```shell
# Dockerビルド
docker build ./ -t 03-recipe-app-backend

# Dockerイメージをkindに登録
kind load docker-image 03-recipe-app-backend --name=foo-cluster

# 03-recipe-app-backendディレクトリから抜ける
cd ..
```

## ConfigMap & secret

configmap.yaml を作成

```configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-env
data:
  POSTGRES_USER: postgres
  POSTGRES_DB: postgres
  POSTGRES_HOSTNAME: db.default
  REDIS_ADDR: redis.default:6379
```

以下を実行

```shell
kubectl apply -f configmap.yaml
```

your_client_id と your_client_secret を書き換えて、以下を実行

```shell
kubectl create secret generic app-env-secret \
  --from-literal=POSTGRES_PASSWORD=password \
  --from-literal=POSTGRES_URL=postgres://postgres:password@db:5432/postgres?sslmode=disable \
  --from-literal=GOOGLE_OAUTH_CLIENT_ID=your_client_id \
  --from-literal=GOOGLE_OAUTH_CLIENT_SECRET=your_client_secret
```

## redis

redis.yaml を作成

```redis.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis
spec:
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
      - name: redis
        image: redis:latest
        ports:
        - containerPort: 6379
---
apiVersion: v1
kind: Service
metadata:
  name: redis
  labels:
    app: redis
spec:
  selector:
    app: redis
  ports:
  - name: redis
    port: 6379
    targetPort: 6379
    protocol: TCP
```

以下を実行

```shell
kubectl apply -f redis.yaml
```

## PostgreSQL

postgres.yaml を作成

```postgres.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-data
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
  storageClassName: standard
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: db
spec:
  selector:
    matchLabels:
      app: db
  template:
    metadata:
      labels:
        app: db
    spec:
      volumes:
      - name: postgres-data
        persistentVolumeClaim:
          claimName: postgres-data
      containers:
      - name: db
        image: groonga/pgroonga:3.0.8-alpine-15-slim
        envFrom:
        - configMapRef:
            name: app-env
        - secretRef:
            name: app-env-secret
        ports:
        - containerPort: 5432
---
apiVersion: v1
kind: Service
metadata:
  name: db
  labels:
    app: db
spec:
  selector:
    app: db
  ports:
  - name: postgresql
    port: 5432
    targetPort: 5432
    protocol: TCP
```

schema.yaml を作成（schema.sqlの実行のみに使用する作業用Pod）

```schema.yaml
apiVersion: v1
kind: Pod
metadata:
  name: schema
  labels:
    app: schema
spec:
  containers:
  - image: postgres:15-bullseye
    name: schema
    command:
    - "sh"
    - "-c"
    args:
    - |
      while true
      do
        sleep 5
      done
```

以下を実行

```shell
kubectl apply -f postgres.yaml
kubectl apply -f schema.yaml

# NAMEがdbから始まるPodと、schemaから始まるPodの起動を確認
kubectl get pods

# schemaのPodにschema.sqlをコピー
kubectl cp ./03-recipe-app-backend/db/schema.sql $(kubectl get --no-headers=true pods -l app=schema -o custom-columns=:metadata.name):/tmp/schema.sql

# schemaのPodのシェルにkubectl execで入る
kubectl exec -it $(kubectl get --no-headers=true pods -l app=schema -o custom-columns=:metadata.name) -- sh
```

schemaのPodのシェルで以下を実行

```shell
###### schema.sql全体をトランザクションにする

# schema.sqlがコピーされたことを確認
ls /tmp/schema.sql

# ファイル先頭にBEGIN;を挿入
sed -i '1s/^/BEGIN;\n\n/' /tmp/schema.sql

# ファイル末尾にCOMMIT;を挿入
sed -i '$s/$/\n\nCOMMIT;\n/' /tmp/schema.sql

# ファイル先頭のBEGIN;を確認
head /tmp/schema.sql

# ファイル末尾のCOMMIT;を確認
tail /tmp/schema.sql

###### psqlに入る。パスワードを聞かれるので入力する
psql -h db.default -p 5432 -d postgres -U postgres
```

psqlで以下を実行

```psql
# schema.sqlを実行してテーブル等を作成する
\i /tmp/schema.sql

# テーブルが作られたか確認する
\dt;
```

OSのシェルに戻って以下を実行

```shell
# 作業用Podを削除
kubectl delete -f schema.yaml
```

## app

service-app.yaml を作成

```service-app.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
  labels:
    app: app
spec:
  selector:
    matchLabels:
      app: app
  template:
    metadata:
      labels:
        app: app
    spec:
      volumes:
      - name: workspaces
        hostPath:
          path: ./workspaces
          type: DirectoryOrCreate
      containers:
      - name: app
        image: 03-recipe-app-backend:latest
        imagePullPolicy: Never
        envFrom:
        - configMapRef:
            name: app-env
        - secretRef:
            name: app-env-secret
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: app
spec:
  type: NodePort
  selector:
    app: app
  ports:
  - name: http
    port: 8888
    targetPort: 8080
    protocol: TCP
```

以下を実行

```shell
kubectl apply -f service-app.yaml

# NAMEがappから始まるPodの起動を確認
kubectl get pods

# appのPodのシェルにkubectl debugでアクセス
kubectl debug $(kubectl get --no-headers=true pods -l app=app -o custom-columns=:metadata.name) --image=curlimages/curl -it -- sh
```

デバッグに入ったら以下を実行

```shell
# APIを呼んでみる
curl http://app.default:8888/api/chefs/featured
```

{"data":[]} が返ったらOK。 OSのシェルに戻っておく。

## Ingress-Nginx Controllerをインストール

以下を実行

```shell
# kind用にインストール
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

# コントローラーの起動を待つ
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s

# 起動確認
kubectl get pods -n ingress-nginx
```

## Ingress

ingress-app.yaml を作成

```ingress-app.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: app
spec:
  ingressClassName: "nginx"
  rules:
  - http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: app
            port:
              number: 8888
```

以下を実行

```shell
kubectl apply -f ingress-app.yaml

# 起動確認
kubectl get pods

# APIを呼んでみる
curl http://localhost/api/chefs/featured
```

{"data":[]} が返ったらOK。

## 後始末

```shell
# クラスタ削除
kind delete clusters foo-cluster
```

