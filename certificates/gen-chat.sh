rm chat-*.pem

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout chat-ca-key.pem -out chat-ca-cert.pem -subj "/C=ID/ST=Jakarta/L=Jakarta Timur/O=Kris Sukoco/OU=Computer/CN=*.krissukoco.dev/emailAddress=kristianto.sukoco@gmail.com"

echo "Chat CA's self-signed certificate generated"
# openssl x509 -in chat-ca-cert.pem -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout chat-server-key.pem -out chat-server-req.pem -subj "/C=ID/ST=Jakarta/L=Jakarta Timur/O=KrisSukoco/OU=Computer/CN=*.krissukoco.dev/emailAddress=kristianto.sukoco@gmail.com"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in chat-server-req.pem -days 60 -CA chat-ca-cert.pem -CAkey chat-ca-key.pem -CAcreateserial -out chat-server-cert.pem -extfile server-ext.cnf

echo "Chat Server's signed certificate generated"
# openssl x509 -in chat-server-cert.pem -noout -text

# 4. Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout chat-client-key.pem -out chat-client-req.pem -subj "/C=ID/ST=Jakarta/L=Jakarta Timur/O=KrisSukoco/OU=Computer/CN=*.krissukoco.dev/emailAddress=kristianto.sukoco@gmail.com"

# 5. Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -in chat-client-req.pem -days 60 -CA chat-ca-cert.pem -CAkey chat-ca-key.pem -CAcreateserial -out chat-client-cert.pem -extfile client-ext.cnf

echo "Chat Client's signed certificate generated"
# openssl x509 -in chat-client-cert.pem -noout -text