#!/bin/bash

endpoint="http://localhost:8000"
user="root"
pass="root"
ns="WITS"
db="AWIPS"

surreal import --endpoint $endpoint --username $user --password $pass --namespace $ns --database $db ./v2.surql
surreal import --endpoint $endpoint --username $user --password $pass --namespace $ns --database $db ./data/vtec.surql
surreal import --endpoint $endpoint --username $user --password $pass --namespace $ns --database $db ./data/state.surql
surreal import --endpoint $endpoint --username $user --password $pass --namespace $ns --database $db ./data/offshore.surql
surreal import --endpoint $endpoint --username $user --password $pass --namespace $ns --database $db ./data/office.surql
surreal import --endpoint $endpoint --username $user --password $pass --namespace $ns --database $db ./data/cwa.surql
surreal import --endpoint $endpoint --username $user --password $pass --namespace $ns --database $db ./data/counties.surql
surreal import --endpoint $endpoint --username $user --password $pass --namespace $ns --database $db ./data/zones.surql
surreal import --endpoint $endpoint --username $user --password $pass --namespace $ns --database $db ./data/marinezones.surql
surreal import --endpoint $endpoint --username $user --password $pass --namespace $ns --database $db ./data/firezones.surql

