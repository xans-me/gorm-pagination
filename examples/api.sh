curl --request GET \
  --url 'http://localhost:8080/transactions?page=1&pageSize=5&dateStart=2023-01-01&dateEnd=2024-01-31&sort=trx_date%20desc'

  curl --request GET \
    --url 'http://localhost:8080/transactions?page=1&pageSize=5&dateStart=2023-01-01&dateEnd=2024-01-31&sort=trx_date%20desc&trx_type=income&account_number=001901007760509'

curl --request GET \
  --url 'http://localhost:8080/transactions?page=1&pageSize=5&dateStart=2023-01-01&dateEnd=2024-01-31&sort=trx_date%20desc&trx_type=expense&account_number=001901007760509&trx_amount=72608'

curl --request GET \
  --url 'http://localhost:8080/transactions?page=1&pageSize=5&dateStart=2023-01-01&dateEnd=2024-01-31&sort=trx_date%20desc&trx_type=expense&account_number=001901007760509&trx_amount=72608&search=TX'