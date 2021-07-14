# ch2es

Data transfer from clickhouse to elasticsearch.

Clickhouse reader provide 4 types of cursor (set with `--ch-cursor`)

    0. offset
        transfer data with limit and offset (use for small tables)
        
    1. timestamp
        transfer data with timestamp offset (unix-timestamp <= <--ch-tsc-field> <= unix-timestamp + <--ch-tsc-step>).
        start at <--ch-tsc-min> end at <--ch-tsc-max>
        
    2. json file
        transfer data from file with clickhouse format JSONEachRow. line by line (use for big tables)
        
    3. stdin
        transfer data from stdin with clickhouse format JSONEachRow. line by line (use for big tables)

Elasticsearch writer provide 2 type of converter (set with `--es-converter`)
    
    0. null
        write data with unchanged clickhouse schema
        
    1. nested
        write data with nested schema
        
# elasticsearch additional info
1. At this stage of the project, all dots in the field name are rewritten to the value of the `--es-dot-replacer` flag, 
if it is set, otherwise the points remain in the fields, which can cause an elasticsearch error.

2. In elasticsearch by default use _id filed. You can rewrite this field with `--es-id-field` flag.

# the nested schema additional info 

Support [nested fields](https://www.elastic.co/guide/en/elasticsearch/reference/current/nested.html). 

Example:

Get from clickhouse :
```json
{
  "foo": 1,
  "bar": "string",
  "baz": [1,2,3,4],
  "baz2": ["one", "two"]
}
```
Insert to elasticsearch:
```json
{
  "foo": 1,
  "bar": "string",
  "<--es-nc-field>" : [
    {
      "baz": 1,
      "baz2": "one"
    },
    {
      "baz": 2,
      "baz2": "two"
    },
    {
      "baz": 3,
      "baz2": null // added if --es-nc-null=true
    },
    {
      "baz": 4,
      "baz2": null // added if --es-nc-null=true
    },
  ]
}
``` 


# client params

Use `-h` for the client params

# examples 

1. transfer from file

```bash
./ch2es.bin --es-bulksz=1000 --es-host=host --es-port=9201 --es-idx=idx --tn=16 --ch-cursor=2 --ch-jfc-file=./q.json --es-dot-replacer=_
```

2. transfer with timestamp cursor

```bash
./ch2es.bin --ch-fields=user_id,author_id,my_timestamp --ch-pass=xyz --ch-db=my_db --ch-host=host --ch-table=my_table --ch-query-timeout=60 --ch-conn-timeout=10 --es-bulksz=5000 --es-host=host --es-idx=my_index --tn=4 --ch-cursor=1 --ch-tsc-step=10 --ch-tsc-min=1622505600 --ch-tsc-max=1624924800 --ch-tsc-field=my_timestamp --es-dot-replacer=_
```
 
3. transfer with offset/limit cursor

```bash
./ch2es.bin --ch-fields=user_id,author_id,my_timestamp --ch-pass=xyz --ch-db=my_db --ch-host=host --ch-table=my_table --ch-query-timeout=60 --ch-conn-timeout=10 --es-bulksz=5000 --es-host=host --es-idx=my_index --tn=4 --ch-cursor=1 --ch-ofc-limit=10 --ch-ofc-max-offset=2000000 --ch-ofc-offset=20 --ch-ofc-order=user_id --es-dot-replacer=_
```

4. transfer from stdin

```bash
clickhouse-client -h 0.0.0.0 -q "select * from info limit 1000 format JSONEachRow" | ./ch2es.bin --es-bulksz=1000 --es-host=host --es-port=9200 --es-idx=test_idx --tn=1 --ch-cursor=3
```
