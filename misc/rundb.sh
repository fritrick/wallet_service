docker run --name wallet-postgres -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword -d postgres
sleep 2
export PGPASSWORD=mysecretpassword
psql -hlocalhost -p5432 -Upostgres -c "CREATE DATABASE wallet"
psql -hlocalhost -p5432 -Upostgres -d wallet < misc/sql/initdb.sql