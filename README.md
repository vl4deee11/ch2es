# ch2es

Data transfer from clickhouse to elasticsearch.

# info

At this stage of the project, all points in the field name are rewritten to the value of the --ch-dot-replacer flag,
 if it is set, otherwise the points remain in the fields, which can cause an elasticsearch error. 

# client params
```
        Usage of ch2es.bin:
          -ch-cond string
                [Clickhouse] where condition
          -ch-conn-timeout int
                [Clickhouse] connect timeout in sec (default 20)
          -ch-cursor int
                [Clickhouse] cursor type. Available 0 (offset cursor), 1 (timestamp cursor), 2 (json file cursor)
          -ch-db string
                [Clickhouse] db name (default "default")
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
          -ch-table string
                [Clickhouse] table
          -ch-user string
                [Clickhouse] db username
          -ch-dot-replacer string
              	[Clickhouse] Replacer for dots in fields if need

          -ch-jfc-file string
                [Clickhouse json file cursor] path to file with data formatted JSONEachRow
          -ch-jfc-line int
                [Clickhouse json file cursor] start line in file with data formatted JSONEachRow


          -ch-ofc-limit int
                [Clickhouse offset cursor] limit (int). Use only if --ch-cursor=0 (by default) (default 100)
          -ch-ofc-max-offset int
                [Clickhouse offset cursor] max offset in clickhouse table. Use only if --ch-cursor=0 (by default)
          -ch-ofc-offset int
                [Clickhouse offset cursor] start offset (int). Use only if --ch-cursor=0 (by default)
          -ch-ofc-order string
                [Clickhouse offset cursor] order field (str). Use only if --ch-cursor=0 (by default)


          -ch-tsc-field string
                [Clickhouse timestamp cursor] field. Should be datetime type or timestamp. Use only if --ch-cursor=1
          -ch-tsc-max int
                [Clickhouse timestamp cursor] end time format unix timestamp. Use only if --ch-cursor=1
          -ch-tsc-min int
                [Clickhouse timestamp cursor] start time format unix timestamp. Use only if --ch-cursor=1
          -ch-tsc-step int
                [Clickhouse timestamp cursor] step in sec. Use only if --ch-cursor=1

          -es-id-field string
                [Elasticsearch]  id field
          -es-blksz int
                [Elasticsearch] search bulk insert size
          -es-host string
                [Elasticsearch] search host (default "0.0.0.0")
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

# Examples 

 1. Transfer from file
 
    `ch2es --es-blksz 1000 --es-host host --es-port 9201 --es-idx idx --tn 16 --ch-cursor 2 --ch-jfc-file ./q.json --ch-dot-replacer _`
    
 2. Transfer with timestamp cursor
 
    `ch2es --ch-fields user_id,author_id,my_timestamp --ch-pass xyz --ch-db my_db --ch-host host --ch-table my_table --ch-query-timeout 60 --ch-conn-timeout 10 --es-blksz 5000 --es-host host --es-idx my_index --tn 4 --ch-cursor 1 --ch-tsc-step 10 --ch-tsc-min 1622505600 --ch-tsc-max 1624924800 --ch-tsc-field my_timestamp --ch-dot-replacer _`
     
 3. Transfer with offset/limit cursor
 
    `ch2es --ch-fields user_id,author_id,my_timestamp --ch-pass xyz --ch-db my_db --ch-host host --ch-table my_table --ch-query-timeout 60 --ch-conn-timeout 10 --es-blksz 5000 --es-host host --es-idx my_index --tn 4 --ch-cursor 1 --ch-ofc-limit 10 --ch-ofc-max-offset 2000000 --ch-ofc-offset 20 --ch-ofc-order user_id --ch-dot-replacer _`