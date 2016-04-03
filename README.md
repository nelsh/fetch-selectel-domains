# fetch-selectel-domains

Fetch (dump) dns zones from https://selectel.ru/services/dns/

**Install and run**

1. Download Linux x64 binary https://github.com/nelsh/fetch-selectel-domains/releases/download/v1.01/fetch-selectel-domains to any location. For example: to `/usr/local/bin`

2. Create config `/etc/fetch-selectel-domains.yaml`. Example:

        APIURL:   "https://api.selectel.ru/domains/v1/"
        APItoken: "XXXXXX"    # get from https://support.selectel.ru/keys/
        TargetPath: "/home/username/selectel-zones"   # path to store zones
        DefaultTTL: 3600      # default 86400

3. Make `TargetDir`

        mkdir /home/username/selectel-zones
    
4. Run

        /usr/local/bin/fetch-selectel-domains
    
5. See files in `/home/username/selectel-zones`
