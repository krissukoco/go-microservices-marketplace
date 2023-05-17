rm auth-*.pem

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout auth-ca-key.pem -out auth-ca-cert.pem -subj "/C=ID/ST=Jakarta/L=Jakarta Timur/O=Kris Sukoco/OU=Computer/CN=*.krissukoco.dev/emailAddress=kristianto.sukoco@gmail.com"

echo "Auth CA's self-signed certificate generated"
# openssl x509 -in auth-ca-cert.pem -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout auth-server-key.pem -out auth-server-req.pem -subj "/C=ID/ST=Jakarta/L=Jakarta Timur/O=KrisSukoco/OU=Computer/CN=*.krissukoco.dev/emailAddress=kristianto.sukoco@gmail.com"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in auth-server-req.pem -days 60 -CA auth-ca-cert.pem -CAkey auth-ca-key.pem -CAcreateserial -out auth-server-cert.pem -extfile server-ext.cnf

echo "Auth Server's signed certificate generated"
# openssl x509 -in auth-server-cert.pem -noout -text

# 4. Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout auth-client-key.pem -out auth-client-req.pem -subj "/C=ID/ST=Jakarta/L=Jakarta Timur/O=KrisSukoco/OU=Computer/CN=*.krissukoco.dev/emailAddress=kristianto.sukoco@gmail.com"

# 5. Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -in auth-client-req.pem -days 60 -CA auth-ca-cert.pem -CAkey auth-ca-key.pem -CAcreateserial -out auth-client-cert.pem -extfile client-ext.cnf

echo "Auth Client's signed certificate generated"
# openssl x509 -in auth-client-cert.pem -noout -text