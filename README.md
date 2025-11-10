# book-my-lab

Got it 👍 — here’s a clean **README section** you can drop into your repo to explain certificate generation.

---

## 🔐 TLS Certificate Setup

This project uses a self-signed Certificate Authority (CA) to sign server certificates for local development and testing.
Follow the steps below to generate the required certs.

---

### 1. Generate CA (root) certificate

Run once and reuse later:

```bash
openssl genrsa -out ca.key 4096
openssl req -x509 -new -nodes -key ca.key -sha256 -days 3650 \
  -out ca.crt \
  -subj "/C=IN/ST=Karnataka/L=Bangalore/O=MyCompany/OU=DevOps/CN=MyRootCA"
```

* `ca.key` → private key (keep safe, do not commit)
* `ca.crt` → root certificate (clients need this to trust your server)

---

### 2. Generate server key + CSR

```bash
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr \
  -subj "/C=IN/ST=Karnataka/L=Bangalore/O=MyCompany/OU=DevOps/CN=localhost"
```

* `server.key` → server private key
* `server.csr` → certificate signing request (temporary file)

---

### 3. Create SAN (Subject Alternative Name) config

Save as `server.ext`:

```bash
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
IP.1  = 127.0.0.1
```

---

### 4. Sign the server certificate with the CA

```bash
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
  -out server.crt -days 365 -sha256 -extfile server.ext
```

* `server.crt` → signed server certificate

---

### 📂 Final files you’ll use

* **Server**: `server.crt`, `server.key`
* **Client**: `ca.crt`
* **Keep safe**: `ca.key` (used to sign more certs later)
* **Can delete**: `server.csr`, `server.ext`

