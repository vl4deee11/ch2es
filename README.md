# ch2es
Data transfer from clickhouse to elasticsearch.

Ch2Es creates a ch2es.stats file to record the offset for each step. 
If you want to restart the program with an offset of 0,delete ch2es.stats

# client params
```shell
  --ch-cond string
        Clickhouse where condition (str) (default "1")
        
  --ch-db string
        Clickhouse db name (str) (default "default")
        
  --ch-fields string
        Clickhouse fields for transfer ex: f_1,f_2,f_3 (str) (default "*")
        
  --ch-host string
        Clickhouse host (str) (default "0.0.0.0")
        
  --ch-limit int
        Clickhouse limit (int)
        
  --ch-order string
        Clickhouse order field (str)
        
  --ch-port int
        Clickhouse http host (int) (default 8123)
        
  --ch-table string
        Clickhouse table (str)
        
  --ch-timeout int
        Clickhouse connect timeout in ms (int)
        
  --es-blksz int
        Elastic search bulk insert size (int)
        
  --es-host string
        Elastic search host (str) (default "0.0.0.0")
        
  --es-idx string
        Elastic search index (str)
        
  --es-port int
        Elastic search port (int) (default 9200)
        
  --max-offset int
        Max offset in clickhouse table (int)
        
  --tn int
        Threads number for parallel bulk inserts (int)

```

