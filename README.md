# fetch-selectel-domains

Fetch (dump) dns zones from https://selectel.ru/services/dns/

**Install and run**

1. Download Linux x64 binary https://github.com/nelsh/fetch-selectel-domains/releases/download/v1.0.2/fetch-selectel-domains to any location. For example: to `/usr/local/bin`

2. Create config `/etc/fetch-selectel-domains.yaml`. Example:

        APIURL:   "https://api.selectel.ru/domains/v1/"
        APItoken: "XXXXXX"    # get from https://support.selectel.ru/keys/
        TargetPath: "/home/username/selectel-zones"   # path to store zones
        DefaultTTL: 3600      # default 86400

3. Make `TargetDir`

        $ mkdir /home/username/selectel-zones
    
4. Run

        $ /usr/local/bin/fetch-selectel-domains -v
    
5. See files in `/home/username/selectel-zones`

**Example result**

`$ cat /home/username/selectel-zones/example.com.dns`

    $ORIGIN example.com.
    $TTL 3600

    example.com.		300	IN	SOA	ns1.selectel.org.  support.selectel.ru.  ( 2016031633 10800 3600 604800 300 )
    example.com.		86400	IN	NS	ns1.selectel.org.
    example.com.		86400	IN	NS	ns3.selectel.org.
    example.com.		86400	IN	NS	ns2.selectel.org.
    example.com.		86400	IN	NS	ns4.selectel.org.
    example.com.			IN	MX	10	mx.yandex.ru.
    example.com.			IN	A	123.123.123.123
