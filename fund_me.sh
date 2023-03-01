#!/bin/bash
for f in /root/test_accounts/*.json
do 
	cat $f  | sed -n 's/.*address\": \"//p' | sed -n 's/\".*//p' >> test_accounts.json
        echo ""  >> test_accounts.json
done


