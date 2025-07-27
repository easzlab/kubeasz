import clickhouse_connect
from clickhouse_connect.driver.exceptions import ClickHouseError

CLICKHOUSE_HOST = 'localhost'
CLICKHOUSE_PORT = 8123 # for http port, use 8443 when TLS is enabled
CLICKHOUSE_USER = 'admin'
CLICKHOUSE_PASSWORD = 'password123'
CLICKHOUSE_DATABASE = 'default'

try:
  client = clickhouse_connect.get_client(
          host=CLICKHOUSE_HOST,
          port=CLICKHOUSE_PORT,
          user=CLICKHOUSE_USER,
          password=CLICKHOUSE_PASSWORD,
          database=CLICKHOUSE_DATABASE
      )
  print("connected to clickhouse")

  create_local_table_query = """
  CREATE TABLE my_table ON CLUSTER default (
    id UInt32,
    sku_id String,
    total_amount Decimal(16,2),
    create_time  Datetime
  ) ENGINE=ReplicatedMergeTree('/clickhouse/tables/{shard}/my_table', '{replica}')
    PARTITION BY toYYYYMMDD(create_time)
    PRIMARY KEY (id)
    ORDER BY (id, sku_id);
  """

  create_distributed_table_query = """
  CREATE TABLE my_table_dist ON CLUSTER default (
    id UInt32,
    sku_id String,
    total_amount Decimal(16,2),
    create_time  Datetime
  ) ENGINE=Distributed(default, default, my_table, hiveHash(sku_id));
  """

  # create table
  client.command(create_local_table_query)
  client.command(create_distributed_table_query)

  insert_query = """
  INSERT INTO my_table_dist (id, sku_id, total_amount, create_time) VALUES
  (1, 'SKU001', 100.00, '2023-10-01 10:00:00'),
  (2, 'SKU002', 150.50, '2023-10-01 10:05:00'),
  (3, 'SKU003', 200.75, '2023-10-01 10:10:00'),
  (4, 'SKU004', 250.00, '2023-10-01 10:15:00'),
  (5, 'SKU005', 300.25, '2023-10-01 10:20:00'),
  (6, 'SKU006', 350.50, '2023-10-01 10:25:00'),
  (7, 'SKU007', 400.00, '2023-10-01 10:30:00'),
  (8, 'SKU008', 450.75, '2023-10-01 10:35:00'),
  (9, 'SKU009', 500.00, '2023-10-01 10:40:00'),
  (10, 'SKU010', 550.25, '2023-10-01 10:45:00');
  """

  # inesert sql
  client.command(insert_query)

  # get data
  select_query = "SELECT * FROM my_table_dist ORDER BY id"
  result = client.query(select_query)

  # print result
  for row in result.result_rows:
      print(row)


  client.command("DROP TABLE IF EXISTS my_table_dist ON CLUSTER default sync")
  client.command("DROP TABLE IF EXISTS my_table ON CLUSTER default sync")

  # close connection
  client.close()

except ClickHouseError as e:
    print(f"An error occurred while connecting to ClickHouse: {e}")
except Exception as e:
    print(f"An unexpected error occurred: {e}")
finally:
    # drop tables
    if 'client' in locals():
        client.close()
