# billing-engine

this project is made based on these assumptions:
1. only 1 currency is supported, which is IDR
2. loan is disbursed at the day it is requested (created in the system)
3. total users = 1M (from google play downloads: 500k+)
   1. DAU (1% total users): 10k
   2. get_outstanding traffic (5 requests per user per day): 50k/day = ~1RPS (peak: 10x = 10RPS)
   3. is_delinquent traffic (5 requests per user per day): 50k/day = ~1RPS (peak: 10x = 10RPS)
   4. make_payment traffic (10% DAU): 1k/day = ~1WPS (peak: 10x = 10WPS)
   5. create_loan_request traffic (10% DAU): 1k/day = ~1WPS (peak: 10x = 10WPS)
4. billings_tab total rows
   1. 1k loan/day & 50x payments/loan = 50k rows/day = 18.250.000 rows/year = 91.250.000 rows in 5 years
   2. sharding is not necessary. However, can consider archiving the old rows to keep the DB performance 
5. from the traffic estimation, caching can be optional, but depends on the traffic behavior (e.g. peaking on certain day or occasion)