# ch2es

Data transfer from clickhouse to elasticsearch.

# info

At this stage of the project, all points in the field name are rewritten to the value of the `--ch-dot-replacer` flag,
 if it is set, otherwise the points remain in the fields, which can cause an elasticsearch error. 
In elasticsearch by default use _id filed. You can rewrite this field with `--es-id-field` flag

# client params
```
  Usage of ./ch2es:
  -ch-cond string
        [Clickhouse] where condition
  -ch-conn-timeout int
        [Clickhouse] connect timeout in sec (default 20)
  -ch-cursor int
        [Clickhouse] cursor type. Available 0 (offset cursor), 1 (timestamp cursor), 2 (json file cursor), 3 (stdin cursor)
  -ch-db string
        [Clickhouse] db name (default "default")
  -ch-dot-replacer string
        [Clickhouse] Replacer for dots in fields if need
  -ch-fields string
        [Clickhouse] fields for transfer ex: f_1,f_2,f_3 (default "*")
  -ch-host string
        [Clickhouse] host (default "0.0.0.0")
  -ch-pass string
        [Clickhouse] db password
  -ch-port int
        [Clickhouse] http host (default 8123)
  -ch-protocol string
        [Clickhouse] protocol (default "http")
  -ch-query-timeout int
        [Clickhouse] query timeout in sec (default 60)              
  -ch-user string
        [Clickhouse] db username
  -ch-table string
        [Clickhouse] table        
        
  -ch-jfc-file string
        [Clickhouse json file cursor] path to file with data formatted JSONEachRow. Use only if --ch-cursor=2
  -ch-jfc-line int
        [Clickhouse json file cursor] start line in file with data formatted JSONEachRow. Use only if --ch-cursor=2
        
        
  -ch-ofc-limit int
        [Clickhouse offset cursor] limit. Use only if --ch-cursor=0 (by default) (default 100)
  -ch-ofc-max-offset int
        [Clickhouse offset cursor] max offset in clickhouse table. Use only if --ch-cursor=0 (by default)
  -ch-ofc-offset int
        [Clickhouse offset cursor] start offset. Use only if --ch-cursor=0 (by default)
  -ch-ofc-order string
        [Clickhouse offset cursor] order field. Use only if --ch-cursor=0 (by default)


  -ch-stdinc-line int
        [Clickhouse stdin cursor] start line in stdin with data formatted JSONEachRow. Use only if --ch-cursor=3
        
        
  -ch-tsc-field string
        [Clickhouse timestamp cursor] field. Should be datetime type or timestamp. Use only if --ch-cursor=1
  -ch-tsc-max int
        [Clickhouse timestamp cursor] end time format unix timestamp. Use only if --ch-cursor=1
  -ch-tsc-min int
        [Clickhouse timestamp cursor] start time format unix timestamp. Use only if --ch-cursor=1
  -ch-tsc-step int
        [Clickhouse timestamp cursor] step in sec. Use only if --ch-cursor=1
        
        
  -es-blksz int
        [Elasticsearch] search bulk insert size
  -es-host string
        [Elasticsearch] search host (default "0.0.0.0")
  -es-id-field string
        [Elasticsearch] id field
  -es-idx string
        [Elasticsearch] search index
  -es-pass string
        [Elasticsearch] search password
  -es-port int
        [Elasticsearch] search port (default 9200)
  -es-protocol string
        [Elasticsearch] protocol (default "http")
  -es-query-timeout int
        [Elasticsearch] search query timeout in sec (default 60)
  -es-user string
        [Elasticsearch] search username
        
        
  -tn int
        [Common] Threads number for parallel insert and read

```

# examples 

 1. transfer from file
 
    `ch2es --es-blksz 1000 --es-host host --es-port 9201 --es-idx idx --tn 16 --ch-cursor 2 --ch-jfc-file ./q.json --ch-dot-replacer _`
    
 2. transfer with timestamp cursor
 
    `ch2es --ch-fields user_id,author_id,my_timestamp --ch-pass xyz --ch-db my_db --ch-host host --ch-table my_table --ch-query-timeout 60 --ch-conn-timeout 10 --es-blksz 5000 --es-host host --es-idx my_index --tn 4 --ch-cursor 1 --ch-tsc-step 10 --ch-tsc-min 1622505600 --ch-tsc-max 1624924800 --ch-tsc-field my_timestamp --ch-dot-replacer _`
     
 3. transfer with offset/limit cursor
 
    `ch2es --ch-fields user_id,author_id,my_timestamp --ch-pass xyz --ch-db my_db --ch-host host --ch-table my_table --ch-query-timeout 60 --ch-conn-timeout 10 --es-blksz 5000 --es-host host --es-idx my_index --tn 4 --ch-cursor 1 --ch-ofc-limit 10 --ch-ofc-max-offset 2000000 --ch-ofc-offset 20 --ch-ofc-order user_id --ch-dot-replacer _`
 
 4. transfer from stdin
    `clickhouse-client -h 0.0.0.0 -q "select * from info limit 1000 format JSONEachRow" | ./ch2es.bin --es-blksz 1000 --es-host 0.0.0.0 --es-port 9200 --es-idx test_idx --tn 1 --ch-cursor 3`
