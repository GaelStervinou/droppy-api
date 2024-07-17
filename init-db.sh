#!/bin/bash
echo "Updating pg_hba.conf..."
echo "host all all 0.0.0.0/0 trust" > /var/lib/postgresql/data/pg_hba.conf
#pg_ctl restart -D /var/lib/postgresql/data
service postgresql restart