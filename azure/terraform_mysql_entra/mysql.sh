#!/bin/bash

set -e

my_cnf() {
  cat <<EOF
[client]
host=mysql-entra.mysql.database.azure.com
port=3306
password=$(az account get-access-token --resource-type oss-rdbms --output tsv --query accessToken)
EOF
}

exec mysql --defaults-file=<(my_cnf) "$@"
