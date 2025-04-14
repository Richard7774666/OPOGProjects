Project Structure Overview

cmd/ - Contains executable entry points

csv_generator/ - Tool to generate sample CSV data
import/ - Command-line tool for importing CSV data
benchmark/ - Performance testing tool


internal/csv/ - Core CSV processing logic

importer.go - Standard CSV importer
bulk_importer.go - Optimized batch importer
importer_test.go - Tests and benchmarks



The implementation offers three different import strategies:

Sequential import - Processes one record at a time
Parallel import - Uses multiple worker goroutines
Bulk import - Uses database transactions with batching