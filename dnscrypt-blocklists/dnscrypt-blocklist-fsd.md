# DNSCrypt blocklist

## Functional Specification Document (FSD)

**Project:** Shell script, that will run on GitHub CI (using "ubuntu-latest" runner) that will download and parse the DNSCrypt resolvers list which contains a list of DNS Stamp urls. Script will parse an IP addresses and/or a server hostnames from those DNS Stamps.

1. Downloads the DNSCrypt resolver list (see `DNSCrypt resolvers list` paragraph for lists)
2. Extracts all DNS stamps (format: `sdns://…`)
3. Decodes each DNS stamp per the official spec
4. Extracts hostnames/IP addresses embedded in the stamp (extracts both if both addr and hostname fields exist)
5. Saves all extracted `sdns://…` stamps in `sdns.txt` text file. Each line contains one stamp.
6. Saves all decoded hostnames into `domains.txt` text file. Each line contains one domain.
7. Saves all IPv4 addresses into `ipv4.txt`. Each line contains one IPv4 address.
8. Saves all IPv6 addresses into `ipv6.txt`. Each line contains one IPv6 address.
9. Sorts and deduplicates all output files.

At the end, convert all output text files to JSON format as well.

## DNSCrypt resolvers list

- https://raw.githubusercontent.com/DNSCrypt/dnscrypt-resolvers/refs/heads/master/v3/onion-services.md
- https://raw.githubusercontent.com/DNSCrypt/dnscrypt-resolvers/refs/heads/master/v3/opennic.md
- https://raw.githubusercontent.com/DNSCrypt/dnscrypt-resolvers/refs/heads/master/v3/odoh-relays.md
- https://raw.githubusercontent.com/DNSCrypt/dnscrypt-resolvers/refs/heads/master/v3/odoh-servers.md
- https://raw.githubusercontent.com/DNSCrypt/dnscrypt-resolvers/refs/heads/master/v3/parental-control.md
- https://raw.githubusercontent.com/DNSCrypt/dnscrypt-resolvers/refs/heads/master/v3/public-resolvers.md
- https://raw.githubusercontent.com/DNSCrypt/dnscrypt-resolvers/refs/heads/master/v3/relays.md

## Supported DNS Stamp Protocols

The script supports the following DNS stamp protocol types:
- 0x00: Plain DNS
- 0x01: DNSCrypt
- 0x02: DNS-over-HTTPS (DoH)
- 0x03: DNS-over-TLS (DoT)
- 0x04: DNS-over-QUIC (DoQ)
- 0x05: Oblivious DoH target
- 0x06: DNS over Oblivious HTTP target
- 0x81: Anonymized DNSCrypt relay
- 0x85: Oblivious DoH relay
- 0x86: DNS over Oblivious HTTP relay

Stamps that contain both an IP address field and a hostname field should extract both to their respective output files.

## Requirements:

- shell script that can run on GitHub CI (using `ubuntu-latest` runner) and can be easily run locally as well

## Usage

```bash
chmod +x parse-dnscrypt-stamps.sh
./parse-dnscrypt-stamps.sh
```

## Expected Output

Typical run processes ~1,100+ stamps from all 7 sources producing:
- `sdns.txt`: All unique stamps (~1,080 entries)
- `domains.txt`: Domain names (~80 entries)
- `ipv4.txt`: IPv4 addresses (~330 entries)
- `ipv6.txt`: IPv6 addresses (~290 entries)

Stamps with both IP address and hostname fields (e.g., DoH, DoT, DoQ protocols) will contribute entries to both the IP and domain output files.

## Notes & gotchas

Script:

- Handles hostnames and IPs
- Uses only standard Ubuntu tools available on GitHub's `ubuntu-latest` runner (`curl`, `grep`, `base64`, `hexdump`)
- Does not parse paths or public keys
- Assumes valid stamps, does not check them
- Silently skips stamps that fail to decode
- Continues processing if a source URL fails to download
- Properly handles IPv6 addresses in bracket notation (e.g., `[2001:db8::1]:443`)

## Other useful resources

DNS Stamp encoding protocol is documented at https://dnscrypt.info/stamps-specifications

## Conversion to JSON format

To convert `.txt` output file to JSON, use this one-liner:

```sh
awk_tool="$(command -v awk)"

"${awk_tool}" 'BEGIN{print "["} {printf "%s\"%s\"", (NR>1?",\n":""), $1} END{print "\n]"}' ./ipv4.txt > ./ipv4.json
```
