# DNSCrypt blocklist

## Functional Specification Document (FSD)

**Project:** Shell script, that will run on GitHub CI (using "ubuntu-latest" runner) that will download and parse the DNSCrypt resolvers list which contains a list of DNS Stamp urls. Script will parse an IP addresses and/or a server hostnames from those DNS Stamps.

1. Downloads the DNSCrypt resolver list (see `DNSCrypt resolvers list` paragraph for lists)
2. Extracts all DNS stamps (format: `sdns://…`)
3. Decodes each DNS stamp per the official spec
4. Extracts hostnames/IP addresses from embedded in the stamp
5. Saves all extracted `sdns://…` stamps in `sdns.txt` text file. Each line contains one stamp.
6. Saves all decoded hostnames into `domains.txt` text file. Each line contains one domain.
7. Saves all IPv4 addresses into `ipv4.txt`. Each line contains one IPv4 address.
8. Saves all IPv6 addresses into `ipv6.txt`. Each line contains one IPv6 address.
9. Sorts and deduplicates all output files.

## DNSCrypt resolvers list

- https://raw.githubusercontent.com/DNSCrypt/dnscrypt-resolvers/refs/heads/master/v3/odoh-servers.md
- https://raw.githubusercontent.com/DNSCrypt/dnscrypt-resolvers/refs/heads/master/v3/public-resolvers.md
- https://raw.githubusercontent.com/DNSCrypt/dnscrypt-resolvers/refs/heads/master/v3/relays.md

## Requirements:

- shell script that can run on GitHub CI (using `ubuntu-latest` runner) and can be easily run locally as well

## Notes & gotchas

Script:

- Handles hostnames and IPs
- Uses only standard Ubuntu tools available on GitHub's `ubuntu-latest` runner (`curl`, `grep`, `base64`, `hexdump`)
- Does not parse paths or public keys
- Assumes valid stamps, does not check them

## Other useful resources

DNS Stamp encoding protocol is documented at https://dnscrypt.info/stamps-specifications
