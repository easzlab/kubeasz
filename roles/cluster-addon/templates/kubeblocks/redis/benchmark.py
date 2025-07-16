import redis
import time
import argparse
from threading import Thread
from statistics import mean, median, pstdev

###############################################################################
###############################################################################
# README.md
# This is a simple bench mark test for redis.
# - Tests both SET and GET operations
# - Measures throughput and latency statistics
# - Supports concurrent client simulation
# - Configurable payload sizes
#
# Usage:
# python3 benchmark.py \
#   --host 127.0.0.1 \
#   --clients 100 \
#   --requests 1000000 \
#   --value-size 256
#   --username default \
#   --password <password>
#
# Installation Requirements:
# - pip install redis
#
# Output Example:
# ==================================================
# Total requests: 100,000
# Errors: 0
# Duration: 4.23s
# Throughput: 23,640.65 ops/sec
# Mean latency: 4.12ms
# Median latency: 3.98ms
# 99th percentile: 8.45ms
# Latency std dev: 1.23ms
# ==================================================
#
###############################################################################
###############################################################################
class RedisBenchmark:
    def __init__(self, host='localhost', port=6379, username=None, password=None, db=0):
        self.r = redis.Redis(
            host=host,
            port=port,
            username=username,
            password=password,
            db=db,
            decode_responses=True,

        )
        self.latencies = []
        self.results = {
            'total_requests': 0,
            'errors': 0,
            'start_time': None,
            'end_time': None
        }

    def flush_db(self):
        """Initialize test environment"""
        self.r.flushdb()
        print("✅ Database flushed")

    def _record_latency(self, start):
        latency = (time.perf_counter() - start) * 1000  # in milliseconds
        self.latencies.append(latency)
        self.results['total_requests'] += 1

    def _worker(self, op_type, key_size, value_size, num_ops):
        """Worker thread for benchmark operations"""
        key_prefix = f"key:{op_type}:{key_size}:"
        value = 'v' * value_size

        for i in range(num_ops):
            key = f"{key_prefix}{i}"
            try:
                start = time.perf_counter()
                if op_type == 'set':
                    self.r.set(key, value)
                elif op_type == 'get':
                    self.r.get(key)
                self._record_latency(start)
            except redis.RedisError as e:
                self.results['errors'] += 1

    def run_test(self, op_type, num_clients=50, total_requests=100000,
                key_size=32, value_size=128):
        """Run benchmark test with parameters"""
        self.latencies = []
        self.results.update({
            'total_requests': 0,
            'errors': 0,
            'start_time': time.time(),
            'op_type': op_type.upper()
        })

        ops_per_client = total_requests // num_clients
        threads = []

        print(f"\n🏁 Starting {op_type.upper()} test with {num_clients} clients...")

        # Create and start threads
        for _ in range(num_clients):
            t = Thread(target=self._worker,
                      args=(op_type, key_size, value_size, ops_per_client))
            threads.append(t)
            t.start()

        # Wait for all threads to complete
        for t in threads:
            t.join()

        self.results['end_time'] = time.time()
        self._calculate_stats()
        self._print_results()

    def _calculate_stats(self):
        """Calculate performance statistics"""
        total_time = self.results['end_time'] - self.results['start_time']
        self.results.update({
            'duration': total_time,
            'throughput': self.results['total_requests'] / total_time,
            'mean_latency': mean(self.latencies),
            'median_latency': median(self.latencies),
            'p99_latency': sorted(self.latencies)[int(len(self.latencies) * 0.99)],
            'stdev': pstdev(self.latencies)
        })

    def _print_results(self):
        """Print formatted benchmark results"""
        print(f"\n📊 {self.results['op_type']} Benchmark Results:")
        print("=" * 50)
        print(f"Total requests: {self.results['total_requests']:,}")
        print(f"Errors: {self.results['errors']}")
        print(f"Duration: {self.results['duration']:.2f}s")
        print(f"Throughput: {self.results['throughput']:,.2f} ops/sec")
        print(f"Mean latency: {self.results['mean_latency']:.2f}ms")
        print(f"Median latency: {self.results['median_latency']:.2f}ms")
        print(f"99th percentile: {self.results['p99_latency']:.2f}ms")
        print(f"Latency std dev: {self.results['stdev']:.2f}ms")
        print("=" * 50)

def main():
    parser = argparse.ArgumentParser(description='Redis Benchmark Tool')
    parser.add_argument('--host', default='localhost', help='Redis host')
    parser.add_argument('--port', type=int, default=6379, help='Redis port')
    parser.add_argument('--username',  default="default",help='Redis username')
    parser.add_argument('--password', help='Redis password')
    parser.add_argument('--clients', type=int, default=50,
                       help='Number of concurrent clients')
    parser.add_argument('--requests', type=int, default=100000,
                       help='Total number of requests')
    parser.add_argument('--key-size', type=int, default=32,
                       help='Size of keys in bytes')
    parser.add_argument('--value-size', type=int, default=128,
                       help='Size of values in bytes')

    args = parser.parse_args()

    benchmark = RedisBenchmark(
        host=args.host,
        port=args.port,
        username=args.username,
        password=args.password
    )

    try:
        # Prepare test environment
        benchmark.flush_db()

        # Run SET benchmark
        benchmark.run_test(
            op_type='set',
            num_clients=args.clients,
            total_requests=args.requests,
            key_size=args.key_size,
            value_size=args.value_size
        )

        # Run GET benchmark
        benchmark.run_test(
            op_type='get',
            num_clients=args.clients,
            total_requests=args.requests,
            key_size=args.key_size,
            value_size=args.value_size
        )

        # Cleanup
        benchmark.flush_db()
    except redis.ConnectionError as e:
        print(f"❌ Connection failed: {e}")
    except KeyboardInterrupt:
        print("\n🚫 Benchmark interrupted by user")

if __name__ == "__main__":
    main()
