# fetch-selectel-domains

Fetch (dump) dns zones from https://selectel.ru/services/dns/

Linux x64 binary: https://github.com/nelsh/fetch-selectel-domains/releases/download/v1.0/fetch-selectel-domains

Example config `fetch-selectel-domains.yaml`

    APIURL:   "https://api.selectel.ru/domains/v1/"
    APItoken: "XXXXXX"    # get from https://support.selectel.ru/keys/
    TargetPath: "zones"   # path to store zones
    DefaultTTL: 3600      # default 86400
