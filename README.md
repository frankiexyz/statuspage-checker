# statuspage-checker
Statuspage-checker

A simple golang program, it will get the PoP/item status for companies who use statuspage.io and export the values to prometheus and it will send a hipchat message to your hipchat room for non-green items.

I tested with Cloudflare and Incapsula's statuspage and it works fine.

Append the following in prometheus's config
```
  - job_name: statuspage-checker
    scrape_interval:     900s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
    static_configs:
      - targets: ['x.x.x.x:8888']

```
