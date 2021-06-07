# ch2es
Data transfer from clickhouse to elasticsearch.

Ch2Es creates a ch2es.stats file to record the offset for each step. 
If you want to restart the program with an offset of 0, delete ch2es.stats.

# client params
```
  --ch-cond string
        Clickhouse clickhouse where condition (str) (default "1")
        
  --ch-conn-timeout int
        Clickhouse connect timeout in sec (int) (default 20)
        
  --ch-db string
        Clickhouse db name (str) (default "default")
        
  --ch-fields string
        Clickhouse clickhouse fields for transfer ex: f_1,f_2,f_3 (str) (default "*")
        
  --ch-host string
        Clickhouse host (str) (default "0.0.0.0")
        
  --ch-limit int
        Clickhouse limit (int) (default 100)
        
  --ch-order string
        Clickhouse order field (str)
        
  --ch-pass string
        Clickhouse db password (str)
        
  --ch-port int
        Clickhouse http host (int) (default 8123)
        
  --ch-protocol string
        Clickhouse protocol (str) (default "http")
        
  --ch-query-timeout int
        Clickhouse query timeout in sec (int) (default 60)
        
  --ch-table string
        Clickhouse table (str)
        
  --ch-user string
        Clickhouse db username (str)
        
  --es-blksz int
        Elastic search bulk insert size (int)
        
  --es-host string
        Elastic search host (str) (default "0.0.0.0")
        
  --es-idx string
        Elastic search index (str)
        
  --es-port int
        Elastic search port (int) (default 9200)
        
  --es-user string
        Elastic search username (str)
        
  --es-pass string
        Elastic search password (str)

  --es-protocol string
        Elastic search protocol (str) (default "http")
  
  --es-query-timeout string
        Elastic search query timeout in sec (int) (default 60)
  
  --max-offset int
        Max offset in clickhouse table (int)
        
  --tn int
        Threads number for parallel insert and read (int)


```

