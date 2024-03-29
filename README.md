<img src='https://user-images.githubusercontent.com/1423657/218816262-e0e8d7ad-44d0-4a7d-9497-0d383ed78b83.png' width=150 />


# PromCasa

 _Query, Aggregate and Publish_ anything from **ClickHouse** to **Prometheus** metrics in _zero time_
 
 Unlike other exporters limited to internal metrics, **PromCasa** is designed to unleash data from *any query* 

<br />

<img src="https://user-images.githubusercontent.com/1423657/153759412-bab0e246-4770-4fe4-b301-f48113c6b9d7.png" width=400 />


#### :star: Functionality

- Execute recurring Clickhouse `queries`
- Exctract mapped `labels` and `values`
- Aggregate results using `metric buckets`
- Publish as `prometheus` metrics

---

### Instructions
Download a [binary release](https://github.com/metrico/promcasa/releases/) or build from source


#### 📦 Download Binary
```
curl -fsSL github.com/metrico/promcasa/releases/latest/download/promcasa -O && chmod +x promcasa
```

#### :page_facing_up:	Configuration

**PromCasa** acts according to the query bucket parameters configured in `/etc/promcasa.json`

<br>

##### ▶️ Query Buckets
To provision and publish a new metrics bucket, extend the configuration with a query set:
```javascript
{
   "_info": "Custom Metrics from Clickhouse",
   "name": "my_status", // must be unique
   "help": "My Status",
   "query": "SELECT status, group, count(*) as counter FROM my_index FINAL PREWHERE (datetime >= toDateTime(now()-60)) AND (datetime < toDateTime(now()) ) group by status, group",
   "labels": ["status","group"], // must match columns
   "counter_name": "counter",
   "refresh": "60s", //  Refesh takes unit sign: (ns, ms, s, m, h)
   "type":"gauge" // gauge, histogram, counter
}
```

For a complete usage example, check out the [wiki](https://github.com/metrico/promcasa/wiki)

<br>

## License
 ©️ [qxip/metrico](https://metrico.in) Licensed under AGPLv3 as part of [qryn](https://qryn.dev)
