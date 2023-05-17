rm product-*.pem

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout product-ca-key.pem -out product-ca-cert.pem -subj "/C=ID/ST=Jakarta/L=Jakarta Timur/O=Kris Sukoco/OU=Computer/CN=*.krissukoco.dev/emailAddress=kristianto.sukoco@gmail.com"

echo "Product CA's self-signed certificate generated"
# openssl x509 -in product-ca-cert.pem -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout product-server-key.pem -out product-server-req.pem -subj "/C=ID/ST=Jakarta/L=Jakarta Timur/O=KrisSukoco/OU=Computer/CN=*.krissukoco.dev/emailAddress=kristianto.sukoco@gmail.com"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in product-server-req.pem -days 60 -CA product-ca-cert.pem -CAkey product-ca-key.pem -CAcreateserial -out product-server-cert.pem -extfile server-ext.cnf

echo "Product Server's signed certificate generated"
# openssl x509 -in product-server-cert.pem -noout -text

# 4. Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout product-client-key.pem -out product-client-req.pem -subj "/C=ID/ST=Jakarta/L=Jakarta Timur/O=KrisSukoco/OU=Computer/CN=*.krissukoco.dev/emailAddress=kristianto.sukoco@gmail.com"

# 5. Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -in product-client-req.pem -days 60 -CA product-ca-cert.pem -CAkey product-ca-key.pem -CAcreateserial -out product-client-cert.pem -extfile client-ext.cnf

echo "Product Client's signed certificate generated"
# openssl x509 -in product-client-cert.pem -noout -text