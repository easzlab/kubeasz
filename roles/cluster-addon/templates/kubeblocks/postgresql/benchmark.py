import subprocess
import argparse
import psycopg2
from psycopg2 import sql

###############################################################################
###############################################################################
# README.md
# This is a simple test for PostgreSQL using sysbench
# Prerequisites
# - Sysbench installed
#   - `sudo apt-get install sysbench` on Ubuntu
#   - `brew install sysbench` on MacOS
# - Python 3.6+
# - PostgreSQL client libraries
#   - `pip install psycopg2-binary`
#
# Run
# - python3 benchmark.py --host <host> --user <user> --password <password>
###############################################################################
###############################################################################

# Configuration variables
PG_HOST = "127.0.0.1"
PG_PORT = "5432"
PG_USER = "<USERNAME>" #
PG_PASSWORD = "<PASSWD>"

SYSBENCH_DB_NAME = "sysbench_test"
SYSBENCH_USER = "sysbench_user"
SYSBENCH_PASSWORD = "SysbenchPass123!"

TABLE_SIZE = 100000
THREADS = [4, 8, 16]  # Different thread counts to test
DURATION = 60  # Test duration in seconds

def create_pg_user_and_db(host, port, root_user, root_password):
    try:
        connection = psycopg2.connect(
            host=host,
            port=port,
            user=root_user,
            password=root_password,
            dbname="postgres"
        )
        connection.autocommit = True
        cursor = connection.cursor()

        # Create database
        # Create database
        cursor.execute(
            sql.SQL("SELECT 1 FROM pg_database WHERE datname = {}")
            .format(sql.Literal(SYSBENCH_DB_NAME)))
        if not cursor.fetchone():
            cursor.execute(
                sql.SQL("CREATE DATABASE {}")
                .format(sql.Identifier(SYSBENCH_DB_NAME)))
            print(f"Created database {SYSBENCH_DB_NAME}")

        # Create user
        cursor.execute(
            sql.SQL("SELECT 1 FROM pg_roles WHERE rolname = {}")
            .format(sql.Literal(SYSBENCH_USER)))
        if not cursor.fetchone():
            cursor.execute(
                sql.SQL("CREATE USER {} WITH PASSWORD {}")
                .format(
                    sql.Identifier(SYSBENCH_USER),
                    sql.Literal(SYSBENCH_PASSWORD)
                ))
            print(f"Created user {SYSBENCH_USER}")

        # Grant privileges
        cursor.execute(
            sql.SQL("GRANT ALL PRIVILEGES ON DATABASE {} TO {}")
            .format(
                sql.Identifier(SYSBENCH_DB_NAME),
                sql.Identifier(SYSBENCH_USER)
            ))
        cursor.execute(
            sql.SQL("ALTER DATABASE {} OWNER TO {}")
            .format(
                sql.Identifier(SYSBENCH_DB_NAME),
                sql.Identifier(SYSBENCH_USER)
            ))

        cursor.close()
        connection.close()
        print("PostgreSQL user and database created successfully")

    except psycopg2.Error as e:
        print(f"Error creating PostgreSQL user/database: {e}")
        exit(1)

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
    # Create PostgreSQL user and database
    create_pg_user_and_db(PG_HOST, PG_PORT, PG_USER, PG_PASSWORD)

    # Prepare command
    prepare_cmd = (
        f"sysbench oltp_read_write "
        f"--db-driver=pgsql "
        f"--pgsql-host={PG_HOST} "
        f"--pgsql-port={PG_PORT} "
        f"--pgsql-user={SYSBENCH_USER} "
        f"--pgsql-password={SYSBENCH_PASSWORD} "
        f"--pgsql-db={SYSBENCH_DB_NAME} "
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
            f"--db-driver=pgsql "
            f"--pgsql-host={PG_HOST} "
            f"--pgsql-port={PG_PORT} "
            f"--pgsql-user={SYSBENCH_USER} "
            f"--pgsql-password={SYSBENCH_PASSWORD} "
            f"--pgsql-db={SYSBENCH_DB_NAME} "
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
        f"--db-driver=pgsql "
        f"--pgsql-host={PG_HOST} "
        f"--pgsql-port={PG_PORT} "
        f"--pgsql-user={SYSBENCH_USER} "
        f"--pgsql-password={SYSBENCH_PASSWORD} "
        f"--pgsql-db={SYSBENCH_DB_NAME} "
        f"--table-size={TABLE_SIZE} "
        f"--report-interval=1 "
        f"cleanup"
    )
    run_sysbench(cleanup_cmd, "Cleanup")

if __name__ == "__main__":
    # Verify dependencies
    check_dependencies()

    parser = argparse.ArgumentParser(description='PostgreSQL Sysbench Runner')
    parser.add_argument('--host', help='PostgreSQL Host', default=PG_HOST)
    parser.add_argument('--port', help='PostgreSQL Port', default=PG_PORT)
    parser.add_argument('--user', help='PostgreSQL User', default=PG_USER)
    parser.add_argument('--password', help='PostgreSQL Password', default=PG_PASSWORD)
    args = parser.parse_args()

    PG_HOST = args.host
    PG_PORT = args.port
    PG_USER = args.user
    PG_PASSWORD = args.password

    main()
