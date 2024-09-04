# Mengambil data dengan pagination filtering date range dan data di sort berdasarkan trx_date

curl --request GET \
  --url 'http://localhost:8080/transactions?page=1&pageSize=5&dateStart=2023-01-01&dateEnd=2024-01-31&sort=trx_date%20desc'

  curl --request GET \
    --url 'http://localhost:8080/transactions?page=1&pageSize=5&dateStart=2023-01-01&dateEnd=2024-01-31&sort=trx_date%20desc&trx_type=pemasukan&account_number=001901007760509'

curl --request GET \
  --url 'http://localhost:8080/transactions?page=1&pageSize=5&dateStart=2023-01-01&dateEnd=2024-01-31&sort=trx_date%20desc&trx_type=pengeluaran&account_number=001901007760509'