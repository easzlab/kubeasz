import subprocess
import argparse
import mysql.connector
from mysql.connector import Error

###############################################################################
###############################################################################
# README.md
# This is a simple test for MySQL using sysbench
#
# Prerequisites
# - Sysbench installed
#   - `sudo apt-get install sysbench` on Ubuntu
#   - `brew install sysbench` on MacOS
# - Python 3.6+
# - MySQL client libraries (`python3 -m pip install mysql-connector-python`)
#
# Run test
# - python3 benchmark.py --host <host> --user <user> --password <password>
###############################################################################
###############################################################################

# Configuration variables
MYSQL_HOST = "127.0.0.1"
MYSQL_PORT = 3306
MYSQL_USER = "<USERNAME>" #
MYSQL_PASSWORD = "<PASSWD>"

SYSBENCH_DB_NAME = "sysbench_test"
SYSBENCH_USER = "sysbench_user"
SYSBENCH_PASSWORD = "SysbenchPass123!"

TABLE_SIZE = 100000
THREADS = [4, 8, 16]  # Different thread counts to test
DURATION = 60  # Test duration in seconds

def create_mysql_user_and_db(host, port, root_user, root_password):
    try:
        connection = mysql.connector.connect(
            host=host,
            user=root_user,
            password=root_password,
            port=port,
            auth_plugin='caching_sha2_password'
        )

        cursor = connection.cursor()

        # Create database
        cursor.execute(f"CREATE DATABASE IF NOT EXISTS {SYSBENCH_DB_NAME}")

        # Create user and grant privileges
        cursor.execute(f"""
            CREATE USER IF NOT EXISTS '{SYSBENCH_USER}'@'%'
            IDENTIFIED BY '{SYSBENCH_PASSWORD}'
        """)
        cursor.execute(f"""
            GRANT ALL PRIVILEGES ON {SYSBENCH_DB_NAME}.*
            TO '{SYSBENCH_USER}'@'%'
        """)
        connection.commit()

        print("MySQL user and database created successfully")

    except Error as e:
        print(f"Error creating MySQL user/database: {e}")
        exit(1)
    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()

def check_dependencies():
    try:
        subprocess.run(["sysbench", "--version"], check=True, capture_output=True)
    except (subprocess.CalledProcessError, FileNotFoundError):
        print("Error: sysbench is not installed. Please install it first.")
        print("For Ubuntu/Debian: sudo apt-get install sysbench")
        print("For CentOS/RHEL: sudo yum install sysbench")
        exit(1)

def run_sysbench(command, test_type):
    try:
        print(f"Running {test_type} test...")
        result = subprocess.run(
            command,
            check=True,
            shell=True,
            capture_output=True,
            text=True
        )
        print(result.stdout)
        print(f"✅ {test_type} test completed successfully")
        return True
    except subprocess.CalledProcessError as e:
        print(f"❌ Error during {test_type} test:")
        print(e.stderr)
        return False

def main():
    # Create database if not exists
    # Create MySQL user and database
    create_mysql_user_and_db(MYSQL_HOST, MYSQL_PORT, MYSQL_USER, MYSQL_PASSWORD)

    # Prepare command
    prepare_cmd = (
        f"sysbench oltp_read_write "
        f"--db-driver=mysql "
        f"--mysql-host={MYSQL_HOST} "
        f"--mysql-port={MYSQL_PORT} "
        f"--mysql-user={SYSBENCH_USER} "
        f"--mysql-password={SYSBENCH_PASSWORD} "
        f"--mysql-db={SYSBENCH_DB_NAME} "
        f"--table-size={TABLE_SIZE} "
        f"--report-interval=1 "
        f"prepare"
    )

    if not run_sysbench(prepare_cmd, "Prepare"):
        return

    # Run benchmark for different thread counts
    for threads in THREADS:
        print(f"\n🏁 Starting benchmark with {threads} threads")
        run_cmd = (
            f"sysbench oltp_read_write "
            f"--db-driver=mysql "
            f"--mysql-host={MYSQL_HOST} "
            f"--mysql-port={MYSQL_PORT} "
            f"--mysql-user={SYSBENCH_USER} "
            f"--mysql-password={SYSBENCH_PASSWORD} "
            f"--mysql-db={SYSBENCH_DB_NAME} "
            f"--table-size={TABLE_SIZE} "
            f"--threads={threads} "
            f"--time={DURATION} "
            f"--report-interval=1 "
            f"run"
        )

        if run_sysbench(run_cmd, f"Runtime ({threads} threads)"):
            print(f"📊 Results for {threads} threads:")
            print("--------------------------------")

    # Cleanup
    cleanup_cmd = (
        f"sysbench oltp_read_write "
        f"--db-driver=mysql "
        f"--mysql-host={MYSQL_HOST} "
        f"--mysql-port={MYSQL_PORT} "
        f"--mysql-user={SYSBENCH_USER} "
        f"--mysql-password={SYSBENCH_PASSWORD} "
        f"--mysql-db={SYSBENCH_DB_NAME} "
        f"--table-size={TABLE_SIZE} "
        f"--report-interval=1 "
        f"cleanup"
    )
    run_sysbench(cleanup_cmd, "Cleanup")

if __name__ == "__main__":
    # Verify dependencies
    check_dependencies()

    parser = argparse.ArgumentParser(description='MySQL Sysbench Runner')
    parser.add_argument('--host', help='MySQL Host', default=MYSQL_HOST)
    parser.add_argument('--port', help='MySQL Port', default=MYSQL_PORT)
    parser.add_argument('--user', help='MySQL User', default=MYSQL_USER)
    parser.add_argument('--password', help='MySQL Password', default=MYSQL_PASSWORD)
    args = parser.parse_args()

    MYSQL_HOST = args.host
    MYSQL_PORT = args.port
    MYSQL_USER = args.user
    MYSQL_PASSWORD = args.password

    main()
