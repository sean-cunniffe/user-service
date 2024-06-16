message() {
    echo ""
    echo --------------$1-----------------
    echo ""
}

dir=src/generated/test_certs
ls ${dir} | xargs -I {} rm ${dir}/{}
message "create ca key"
openssl genpkey -algorithm RSA -out ${dir}/test_ca.key
message "create ca crt"
openssl req -x509 -new -nodes -key ${dir}/test_ca.key -sha256 -days 3650 -out ${dir}/test_ca.crt -config scripts/ca.cnf
message "create server key"
openssl genpkey -algorithm RSA -out ${dir}/test_server.key
message "create server csr"
openssl req -new -key ${dir}/test_server.key -out ${dir}/test_server.csr -config scripts/server.cnf
message "create server crt"
openssl x509 -req -in ${dir}/test_server.csr -CA ${dir}/test_ca.crt -CAkey ${dir}/test_ca.key -CAcreateserial -out ${dir}/test_server.crt -days 365 -sha256
message "verify server crt"
openssl verify -CAfile ${dir}/test_ca.crt ${dir}/test_server.crt
message "clean up"
rm -rf ${dir}/test_ca.key ${dir}/test_server.csr ${dir}/test_ca.srl
