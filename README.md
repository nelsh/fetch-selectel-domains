# fetch-selectel-domains

Fetch (dump) dns zones from https://selectel.ru/services/dns/

Example config `fetch-selectel-domains.yaml`

    APIURL:   "https://api.selectel.ru/domains/v1/"
    APItoken: "XXXXXX"    # get from https://support.selectel.ru/keys/
    TargetPath: "zones"   # path to store zones
    DefaultTTL: 3600      # default 86400
