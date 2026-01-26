#!/usr/bin/env bash
set -euo pipefail

# DNSCrypt Stamp Parser
# Decodes DNS stamps from odoh-servers.md and extracts IP addresses and hostnames

readonly ODOH_URL="https://raw.githubusercontent.com/DNSCrypt/dnscrypt-resolvers/refs/heads/master/v3/odoh-servers.md"
readonly DOMAINS_FILE="domains.txt"
readonly IPV4_FILE="ipv4.txt"
readonly IPV6_FILE="ipv6.txt"

# Initialize output files
> "$DOMAINS_FILE"
> "$IPV4_FILE"
> "$IPV6_FILE"

# Function to convert base64url to standard base64
base64url_to_base64() {
    local input="$1"
    # Replace URL-safe characters with standard base64 characters
    input="${input//-/+}"
    input="${input//_//}"
    # Add padding if necessary
    local padding=$((4 - ${#input} % 4))
    if [ $padding -ne 4 ]; then
        printf '%s%*s' "$input" $padding '' | tr ' ' '='
    else
        printf '%s' "$input"
    fi
}

# Function to decode a DNS stamp and extract hostname/IP
decode_stamp() {
    local stamp="$1"

    # Remove sdns:// prefix
    stamp="${stamp#sdns://}"

    # Convert base64url to standard base64 and decode
    local base64_data
    base64_data=$(base64url_to_base64 "$stamp")

    # Decode base64 to binary and convert to hex
    local hex_data
    hex_data=$(echo "$base64_data" | base64 -d 2>/dev/null | hexdump -ve '1/1 "%.2x"')

    [ -z "$hex_data" ] && return 1

    # Parse the stamp format:
    # - First byte (2 hex chars): protocol ID
    # - Next 8 bytes (16 hex chars): props
    # - Next: length-prefixed hostname

    local protocol="${hex_data:0:2}"

    # Skip protocol (2 chars) + props (16 chars) = 18 chars
    local remaining="${hex_data:18}"

    # Get hostname length (first byte after props)
    local hostname_len_hex="${remaining:0:2}"
    local hostname_len=$((16#$hostname_len_hex))

    [ $hostname_len -eq 0 ] && return 1

    # Extract hostname hex data (2 hex chars per byte)
    local hostname_hex="${remaining:2:$((hostname_len * 2))}"

    # Convert hex to ASCII
    local hostname=""
    for ((i=0; i<${#hostname_hex}; i+=2)); do
        local byte="${hostname_hex:$i:2}"
        hostname+=$(printf "\\x$byte")
    done

    # Remove port if present (format: hostname:port)
    hostname="${hostname%%:*}"

    # Remove brackets from IPv6 addresses
    hostname="${hostname#[}"
    hostname="${hostname%]}"

    echo "$hostname"
}

# Function to classify address as IPv4, IPv6, or domain
classify_and_save() {
    local addr="$1"

    [ -z "$addr" ] && return

    # Check if it's an IPv6 address (contains colons and hex digits)
    if [[ "$addr" =~ ^[0-9a-fA-F:]+$ ]] && [[ "$addr" == *:* ]]; then
        # Validate IPv6 format
        if [[ "$addr" =~ ^([0-9a-fA-F]{0,4}:){2,7}[0-9a-fA-F]{0,4}$ ]]; then
            echo "$addr" >> "$IPV6_FILE"
            return
        fi
    fi

    # Check if it's an IPv4 address (four decimal numbers separated by dots)
    if [[ "$addr" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        # Validate each octet is 0-255
        local valid=1
        IFS='.' read -ra octets <<< "$addr"
        for octet in "${octets[@]}"; do
            if [ "$octet" -lt 0 ] || [ "$octet" -gt 255 ]; then
                valid=0
                break
            fi
        done
        if [ $valid -eq 1 ]; then
            echo "$addr" >> "$IPV4_FILE"
            return
        fi
    fi

    # Otherwise, treat it as a domain name
    echo "$addr" >> "$DOMAINS_FILE"
}

main() {
    echo "Downloading odoh-servers.md..."
    local content
    content=$(curl -sfL "$ODOH_URL")

    if [ -z "$content" ]; then
        echo "Error: Failed to download odoh-servers.md" >&2
        exit 1
    fi

    echo "Extracting and decoding DNS stamps..."

    # Extract all sdns:// stamps
    local stamps
    stamps=$(echo "$content" | grep -o 'sdns://[A-Za-z0-9_-]*' || true)

    if [ -z "$stamps" ]; then
        echo "Warning: No DNS stamps found in the file" >&2
        exit 0
    fi

    local count=0
    while IFS= read -r stamp; do
        [ -z "$stamp" ] && continue

        # Decode the stamp
        local hostname
        if hostname=$(decode_stamp "$stamp" 2>/dev/null); then
            [ -n "$hostname" ] && classify_and_save "$hostname"
            count=$((count + 1))
        fi
    done <<< "$stamps"

    echo "Processed $count DNS stamps"

    # Sort and deduplicate output files
    for file in "$DOMAINS_FILE" "$IPV4_FILE" "$IPV6_FILE"; do
        if [ -s "$file" ]; then
            sort -u "$file" -o "$file"
            local line_count
            line_count=$(wc -l < "$file")
            echo "  $(basename "$file"): $line_count entries"
        fi
    done

    echo "Done!"
}

main "$@"
