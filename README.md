<img src="https://user-images.githubusercontent.com/1423657/153759243-d950b5fa-d2a7-49b7-894c-cfd9c9531f82.png" width=100 />

# PromCasa

ClickHouse Custom Metrics Exporter for Prometheus Scrapers

<br />

<img src="https://user-images.githubusercontent.com/1423657/137605775-485de1af-20a1-47ae-933e-b91b9e08edb1.png" width=500 />

#### :star: ClickHeus Functionality

- Execute recurring Clickhouse `queries`
- Exctract mapped `labels` and `values`
- Aggregate results using `metric buckets`
- Publish as `prometheus` metrics

---

#### :page_facing_up:	Configuration

Clickheus acts according to the parameters configured in the `promcasa.json` file.


##### Query Buckets
To provision a metrics bucket, extend the configuration with a new query feed:
```
 {
      "_info": "Custom Metrics from Clickhouse",
      "name": "my_status", // must be unique
      "help": "My Status",
      "query": "SELECT status, group, count(*) as counter FROM my_index FINAL PREWHERE (datetime >= toDateTime(now()-60)) AND (datetime < toDateTime(now()) ) group by status, group",
      "labels": ["status","group"], // must match columns
      "counter_name": "counter",
      "refresh": "60s", //  Refesh takes unit sign: (ns, ms, s, m, h)
      "type":["g"] // gauge
    }
```
