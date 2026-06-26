# VPS Deployment with K3s (Kubernetes) - No Docker

**Production deployment on VPS: Liferay DXP + OpenAutonomyX Services on K3s**

---

## VPS Requirements

- **OS:** Ubuntu 20.04+ or CentOS 8+
- **CPU:** 4+ cores
- **RAM:** 8GB+ (16GB recommended)
- **Storage:** 50GB+
- **Network:** Public IP, ports 80, 443, 3000-3011 open

---

## Step 1: VPS Setup Commands

### SSH into VPS

```bash
ssh root@your-vps-ip
# or
ssh user@your-vps-ip
sudo su -
```

### Update System

```bash
apt-get update
apt-get upgrade -y
apt-get install -y \
  curl \
  wget \
  git \
  vim \
  htop \
  net-tools \
  build-essential \
  apt-transport-https \
  ca-certificates \
  gnupg \
  lsb-release
```

### Install PostgreSQL (for database)

```bash
# Add PostgreSQL repo
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list

# Install
apt-get update
apt-get install -y postgresql-15 postgresql-contrib-15

# Start PostgreSQL
systemctl start postgresql
systemctl enable postgresql

# Create database
sudo -u postgres psql -c "CREATE DATABASE openautonomyx;"
sudo -u postgres psql -c "CREATE USER openautonomyx WITH PASSWORD 'secure_password';"
sudo -u postgres psql -c "ALTER ROLE openautonomyx WITH CREATEDB;"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE openautonomyx TO openautonomyx;"

# Verify
sudo -u postgres psql -d openautonomyx -c "SELECT version();"
```

### Install Redis

```bash
apt-get install -y redis-server redis-tools

# Start Redis
systemctl start redis-server
systemctl enable redis-server

# Verify
redis-cli ping
# Should return: PONG
```

### Install Node.js (for services)

```bash
curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
apt-get install -y nodejs

# Verify
node --version
npm --version
```

---

## Step 2: Install K3s (Lightweight Kubernetes)

### Install K3s

```bash
# Download and install K3s
curl -sfL https://get.k3s.io | sh -

# Wait for installation (2-3 minutes)
sleep 10

# Verify K3s is running
kubectl get nodes
# Should show your node as Ready
```

### Setup kubeconfig

```bash
# Copy kubeconfig to home directory
mkdir -p ~/.kube
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chown $USER:$USER ~/.kube/config

# Test kubectl
kubectl cluster-info
kubectl get namespaces
```

### Check K3s Status

```bash
systemctl status k3s
sudo k3s kubectl get nodes
sudo k3s kubectl get all -A
```

---

## Step 3: Install Liferay on K3s

### Create Namespace

```bash
kubectl create namespace liferay
kubectl get namespaces
```

### Create ConfigMap for Portal Configuration

```bash
kubectl create configmap liferay-config \
  --from-literal=LIFERAY_JDBC_ONE_URL="jdbc:postgresql://postgres:5432/liferay" \
  --from-literal=LIFERAY_JDBC_ONE_DRIVER_CLASS_NAME="org.postgresql.Driver" \
  --from-literal=LIFERAY_JDBC_ONE_USERNAME="postgres" \
  --from-literal=LIFERAY_JDBC_ONE_PASSWORD="postgres" \
  --from-literal=LIFERAY_SETUP_WIZARD_ENABLED="false" \
  -n liferay
```

### Create PersistentVolume for Liferay Data

```bash
cat > /tmp/liferay-pv.yaml << 'EOF'
apiVersion: v1
kind: PersistentVolume
metadata:
  name: liferay-pv
spec:
  capacity:
    storage: 50Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: /mnt/liferay-data

---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: liferay-pvc
  namespace: liferay
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 50Gi
EOF

kubectl apply -f /tmp/liferay-pv.yaml
```

### Deploy Liferay

```bash
cat > /tmp/liferay-deployment.yaml << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: liferay
  namespace: liferay
spec:
  replicas: 2
  selector:
    matchLabels:
      app: liferay
  template:
    metadata:
      labels:
        app: liferay
    spec:
      containers:
      - name: liferay
        image: liferay/dxp:latest
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 11311
          name: debug
        envFrom:
        - configMapRef:
            name: liferay-config
        volumeMounts:
        - name: liferay-data
          mountPath: /opt/liferay/data
        resources:
          requests:
            memory: "4Gi"
            cpu: "2"
          limits:
            memory: "8Gi"
            cpu: "4"
        livenessProbe:
          httpGet:
            path: /
            port: 8080
          initialDelaySeconds: 120
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /
            port: 8080
          initialDelaySeconds: 60
          periodSeconds: 10
      volumes:
      - name: liferay-data
        persistentVolumeClaim:
          claimName: liferay-pvc

---
apiVersion: v1
kind: Service
metadata:
  name: liferay-service
  namespace: liferay
spec:
  type: LoadBalancer
  selector:
    app: liferay
  ports:
  - name: http
    port: 8080
    targetPort: 8080
  - name: debug
    port: 11311
    targetPort: 11311
EOF

kubectl apply -f /tmp/liferay-deployment.yaml

# Wait for deployment
kubectl rollout status deployment/liferay -n liferay

# Check service
kubectl get svc -n liferay
```

---

## Step 4: Deploy OpenAutonomyX Services on K3s

### Create Namespace

```bash
kubectl create namespace openautonomyx
```

### Create ConfigMap for Backend Services

```bash
kubectl create configmap openautonomyx-config \
  --from-literal=DATABASE_URL="postgresql://openautonomyx:secure_password@postgres:5432/openautonomyx" \
  --from-literal=REDIS_URL="redis://redis:6379/0" \
  --from-literal=NODE_ENV="production" \
  -n openautonomyx
```

### Deploy PostgreSQL Service

```bash
cat > /tmp/postgres-deployment.yaml << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  namespace: openautonomyx
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:15-alpine
        ports:
        - containerPort: 5432
        env:
        - name: POSTGRES_DB
          value: openautonomyx
        - name: POSTGRES_USER
          value: openautonomyx
        - name: POSTGRES_PASSWORD
          value: secure_password
        volumeMounts:
        - name: postgres-data
          mountPath: /var/lib/postgresql/data
          subPath: postgres
      volumes:
      - name: postgres-data
        persistentVolumeClaim:
          claimName: postgres-pvc

---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-pvc
  namespace: openautonomyx
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 50Gi

---
apiVersion: v1
kind: Service
metadata:
  name: postgres
  namespace: openautonomyx
spec:
  selector:
    app: postgres
  ports:
  - port: 5432
    targetPort: 5432
EOF

kubectl apply -f /tmp/postgres-deployment.yaml
```

### Deploy Redis Service

```bash
cat > /tmp/redis-deployment.yaml << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis
  namespace: openautonomyx
spec:
  replicas: 1
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
        image: redis:7-alpine
        ports:
        - containerPort: 6379
        volumeMounts:
        - name: redis-data
          mountPath: /data
      volumes:
      - name: redis-data
        persistentVolumeClaim:
          claimName: redis-pvc

---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: redis-pvc
  namespace: openautonomyx
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi

---
apiVersion: v1
kind: Service
metadata:
  name: redis
  namespace: openautonomyx
spec:
  selector:
    app: redis
  ports:
  - port: 6379
    targetPort: 6379
EOF

kubectl apply -f /tmp/redis-deployment.yaml
```

### Deploy API Gateway Service

```bash
cat > /tmp/api-gateway-deployment.yaml << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-gateway
  namespace: openautonomyx
spec:
  replicas: 2
  selector:
    matchLabels:
      app: api-gateway
  template:
    metadata:
      labels:
        app: api-gateway
    spec:
      containers:
      - name: api-gateway
        image: api-gateway:latest
        imagePullPolicy: Never
        ports:
        - containerPort: 3000
        envFrom:
        - configMapRef:
            name: openautonomyx-config
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "1Gi"
            cpu: "1000m"
        livenessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 10
          periodSeconds: 5

---
apiVersion: v1
kind: Service
metadata:
  name: api-gateway
  namespace: openautonomyx
spec:
  type: LoadBalancer
  selector:
    app: api-gateway
  ports:
  - port: 3000
    targetPort: 3000
EOF

kubectl apply -f /tmp/api-gateway-deployment.yaml
```

### Deploy Blog Service

```bash
cat > /tmp/blog-deployment.yaml << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: blog
  namespace: openautonomyx
spec:
  replicas: 2
  selector:
    matchLabels:
      app: blog
  template:
    metadata:
      labels:
        app: blog
    spec:
      containers:
      - name: blog
        image: blog:latest
        imagePullPolicy: Never
        ports:
        - containerPort: 3009
        envFrom:
        - configMapRef:
            name: openautonomyx-config
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"

---
apiVersion: v1
kind: Service
metadata:
  name: blog
  namespace: openautonomyx
spec:
  selector:
    app: blog
  ports:
  - port: 3009
    targetPort: 3009
EOF

kubectl apply -f /tmp/blog-deployment.yaml
```

### Deploy Formats Service

```bash
cat > /tmp/formats-deployment.yaml << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: formats
  namespace: openautonomyx
spec:
  replicas: 2
  selector:
    matchLabels:
      app: formats
  template:
    metadata:
      labels:
        app: formats
    spec:
      containers:
      - name: formats
        image: formats:latest
        imagePullPolicy: Never
        ports:
        - containerPort: 3011
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"

---
apiVersion: v1
kind: Service
metadata:
  name: formats
  namespace: openautonomyx
spec:
  selector:
    app: formats
  ports:
  - port: 3011
    targetPort: 3011
EOF

kubectl apply -f /tmp/formats-deployment.yaml
```

### Deploy Integrations Service

```bash
cat > /tmp/integrations-deployment.yaml << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: integrations
  namespace: openautonomyx
spec:
  replicas: 2
  selector:
    matchLabels:
      app: integrations
  template:
    metadata:
      labels:
        app: integrations
    spec:
      containers:
      - name: integrations
        image: integrations:latest
        imagePullPolicy: Never
        ports:
        - containerPort: 3010
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"

---
apiVersion: v1
kind: Service
metadata:
  name: integrations
  namespace: openautonomyx
spec:
  selector:
    app: integrations
  ports:
  - port: 3010
    targetPort: 3010
EOF

kubectl apply -f /tmp/integrations-deployment.yaml
```

---

## Step 5: Install Ingress Controller (for domain routing)

### Install Traefik (comes with K3s)

```bash
# K3s comes with Traefik by default
# Verify it's running
kubectl get deployment -n kube-system | grep traefik

# Check Traefik service
kubectl get svc -n kube-system
```

### Create Ingress for Liferay

```bash
cat > /tmp/liferay-ingress.yaml << 'EOF'
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: liferay-ingress
  namespace: liferay
  annotations:
    traefik.ingress.kubernetes.io/router.entrypoints: web,websecure
spec:
  rules:
  - host: publishing.openautonomyx.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: liferay-service
            port:
              number: 8080
EOF

kubectl apply -f /tmp/liferay-ingress.yaml
```

### Create Ingress for Backend Services

```bash
cat > /tmp/api-ingress.yaml << 'EOF'
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: api-ingress
  namespace: openautonomyx
  annotations:
    traefik.ingress.kubernetes.io/router.entrypoints: web,websecure
spec:
  rules:
  - host: api.publishing.openautonomyx.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: api-gateway
            port:
              number: 3000
  - host: blog.publishing.openautonomyx.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: blog
            port:
              number: 3009
  - host: formats.publishing.openautonomyx.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: formats
            port:
              number: 3011
  - host: integrations.publishing.openautonomyx.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: integrations
            port:
              number: 3010
EOF

kubectl apply -f /tmp/api-ingress.yaml
```

---

## Step 6: Build and Push Docker Images to K3s

### Build Services Locally

```bash
cd /Users/chinmaypanda/CustomApps

# Build each service
docker build -t api-gateway:latest ./services/api-gateway
docker build -t blog:latest ./services/blog
docker build -t formats:latest ./services/formats
docker build -t integrations:latest ./services/integrations

# Tag for VPS registry (optional)
# docker tag api-gateway:latest your-vps-ip:5000/api-gateway:latest
# docker push your-vps-ip:5000/api-gateway:latest
```

### OR Load Docker Images Directly to K3s

```bash
# Copy Docker images to VPS
scp api-gateway-latest.tar user@vps-ip:/tmp/
scp blog-latest.tar user@vps-ip:/tmp/
scp formats-latest.tar user@vps-ip:/tmp/
scp integrations-latest.tar user@vps-ip:/tmp/

# On VPS, load into k3s
ctr -n k8s.io image import /tmp/api-gateway-latest.tar
ctr -n k8s.io image import /tmp/blog-latest.tar
ctr -n k8s.io image import /tmp/formats-latest.tar
ctr -n k8s.io image import /tmp/integrations-latest.tar
```

---

## Step 7: Verify Deployment

### Check All Pods

```bash
# Check Liferay
kubectl get pods -n liferay
kubectl logs -f deployment/liferay -n liferay

# Check OpenAutonomyX Services
kubectl get pods -n openautonomyx
kubectl logs -f deployment/api-gateway -n openautonomyx

# Check all services
kubectl get all -A
```

### Get LoadBalancer IPs

```bash
# Get external IPs
kubectl get svc -A | grep LoadBalancer

# Example output:
# liferay      liferay-service       LoadBalancer   10.43.1.1   <VPS-IP>   8080:30123/TCP
# openautonomyx api-gateway          LoadBalancer   10.43.1.2   <VPS-IP>   3000:30456/TCP
```

### Update DNS Records

```
# Point your domains to VPS IP:
publishing.openautonomyx.com    →  YOUR-VPS-IP
api.publishing.openautonomyx.com  →  YOUR-VPS-IP
blog.publishing.openautonomyx.com →  YOUR-VPS-IP
```

---

## Step 8: Setup SSL/TLS with Let's Encrypt

### Install Cert-Manager

```bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml

# Wait for deployment
kubectl rollout status deployment/cert-manager -n cert-manager
```

### Create ClusterIssuer

```bash
cat > /tmp/cert-issuer.yaml << 'EOF'
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: your-email@example.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: traefik
EOF

kubectl apply -f /tmp/cert-issuer.yaml
```

### Update Ingress with SSL

```bash
cat > /tmp/liferay-ingress-ssl.yaml << 'EOF'
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: liferay-ingress
  namespace: liferay
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    traefik.ingress.kubernetes.io/router.entrypoints: web,websecure
spec:
  tls:
  - hosts:
    - publishing.openautonomyx.com
    secretName: liferay-tls
  rules:
  - host: publishing.openautonomyx.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: liferay-service
            port:
              number: 8080
EOF

kubectl apply -f /tmp/liferay-ingress-ssl.yaml
```

---

## Monitoring & Management

### View Cluster Status

```bash
# Node status
kubectl get nodes
kubectl top nodes

# Pod status
kubectl get pods -A
kubectl top pods -A

# Events
kubectl get events -A

# Resource usage
kubectl describe node
```

### Scaling Services

```bash
# Scale API Gateway to 5 replicas
kubectl scale deployment/api-gateway --replicas=5 -n openautonomyx

# Check auto-scaling
kubectl get hpa -n openautonomyx
```

### Update Service

```bash
# Update image
kubectl set image deployment/api-gateway \
  api-gateway=api-gateway:v2.0 \
  -n openautonomyx

# Rollout status
kubectl rollout status deployment/api-gateway -n openautonomyx

# Rollback if needed
kubectl rollout undo deployment/api-gateway -n openautonomyx
```

---

## Useful Commands

```bash
# View logs
kubectl logs -f deployment/liferay -n liferay
kubectl logs -f pod/api-gateway-xxx -n openautonomyx

# Execute command in pod
kubectl exec -it pod/liferay-xxx -n liferay -- /bin/bash

# Port forward (for debugging)
kubectl port-forward service/liferay-service 8080:8080 -n liferay

# Delete everything
kubectl delete ns liferay openautonomyx

# Get detailed info
kubectl describe deployment liferay -n liferay
kubectl describe svc liferay-service -n liferay

# Watch status
kubectl get pods -w -n openautonomyx
```

---

## Complete Deployment Checklist

- ☐ VPS setup (Ubuntu 20.04+)
- ☐ PostgreSQL installed & running
- ☐ Redis installed & running
- ☐ K3s installed & running
- ☐ Liferay deployed on K3s
- ☐ API Gateway deployed
- ☐ Blog service deployed
- ☐ Formats service deployed
- ☐ Integrations service deployed
- ☐ Ingress configured
- ☐ DNS records updated
- ☐ SSL/TLS certificates installed
- ☐ Monitoring setup

---

## Access Your Platform

```
Liferay Portal:  https://publishing.openautonomyx.com:8080
API Gateway:     https://api.publishing.openautonomyx.com:3000
Blog Service:    https://blog.publishing.openautonomyx.com:3009
Formats Service: https://formats.publishing.openautonomyx.com:3011
Integrations:    https://integrations.publishing.openautonomyx.com:3010
```

---

**Production Platform Live!** 🚀

K3s + Liferay DXP + OpenAutonomyX Services = Complete Enterprise Publishing Platform
