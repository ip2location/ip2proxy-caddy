# IP2Proxy Caddy Middleware

[![Go](https://img.shields.io/badge/Go-1.25.1%2B-00ADD8?logo=go)](https://go.dev/) [![Caddy](https://img.shields.io/badge/Caddy-v2-1F88C0)](https://caddyserver.com/)

**IP2Proxy Caddy Middleware** enriches incoming requests with proxy and network intelligence through HTTP headers. It looks up the visitor's IP address using either a local **IP2Proxy BIN database** or the **IP2Location.io API**, then makes the results available to downstream applications such as PHP.

## Features

- **Two lookup modes:** local BIN database or remote IP2Location.io API.
- **15 request headers** covering proxy detection, network ownership, location and risk indicators.
- **Configurable header prefix** (default: `X-IP2Proxy`).
- **Caddyfile integration** for applications served through Caddy.
- **IPv4 and IPv6 lookups**, subject to data availability in the selected database or API plan.

## How It Works

```text
Client request → Caddy (ip2proxy middleware) → IP lookup
                                              ├─ Local IP2Proxy BIN
                                              └─ IP2Location.io API
               → Add X-IP2Proxy-* request headers → Application (e.g. PHP)
```

The middleware enriches **requests**, not HTTP responses. Your downstream application reads the injected request headers. If a lookup returns an error, the middleware continues to the next handler without applying new lookup headers.

## Requirements

- [Go](https://go.dev/doc/install) and [xcaddy](https://github.com/caddyserver/xcaddy) to build Caddy with this third-party module. The current `go.mod` declares Go 1.25.1.
- **Local mode:** An IP2Proxy `.BIN` database. Download a free LITE database from [IP2Location LITE](https://www.ip2location.com/database/lite), or use a compatible licensed BIN file.
- **Remote mode:** An [IP2Location.io API key](https://www.ip2location.io/sign-up) with access to the fields you need.

## Installation

Build Caddy with the middleware:

```bash
xcaddy build --with github.com/ip2location/ip2proxy-caddy
```

The resulting `caddy` executable includes the `ip2proxy` directive. Use this custom binary rather than a standard Caddy binary that does not contain the module.

## Configuration

Choose **one** of the following Caddyfile configurations. Both examples use PHP-FPM; adjust the domain, document root and PHP-FPM socket for your environment.

### Option 1: Local IP2Proxy BIN Database

Extract a compatible BIN file on the Caddy host and set its absolute path:

```caddyfile
{
    order ip2proxy before php_fastcgi
}

example.com {
    root * /var/www/html

    ip2proxy {
        mode local
        bin_path /path/to/IP2PROXY.BIN
        header_prefix X-IP2Proxy
    }

    php_fastcgi unix//run/php/php8.4-fpm.sock
    file_server
}
```

Local mode performs the lookup against the BIN database without calling the remote API for each request. The database file must be readable by the Caddy process.

### Option 2: Remote IP2Location.io API

Set the API key in the Caddy process environment, then configure:

```caddyfile
{
    order ip2proxy before php_fastcgi
}

example.com {
    root * /var/www/html

    ip2proxy {
        mode remote
        api_key {$IP2LOCATION_API_KEY}
        header_prefix X-IP2Proxy
    }

    php_fastcgi unix//run/php/php8.4-fpm.sock
    file_server
}
```

Remote mode makes IP2Location.io API lookups. Field availability and query limits depend on your subscription. Avoid committing API keys directly to your Caddyfile.

### Configuration Reference

| Option | Description | Default |
|---|---|---|
| `mode` | `local` for BIN lookup or `remote` for API lookup | `local` |
| `bin_path` | Path to the IP2Proxy BIN file; required in local mode | None |
| `api_key` | IP2Location.io API key; required in remote mode | None |
| `header_prefix` | Prefix for generated request headers | `X-IP2Proxy` |

## Usage

Start Caddy with the configuration file:

```bash
./caddy run --config /path/to/Caddyfile --adapter caddyfile
```

To check Caddyfile syntax before running:

```bash
./caddy validate --config /path/to/Caddyfile --adapter caddyfile
```

### Read Headers in PHP

Save this example as `index.php` in the configured website directory and open it in your browser. It displays any IP2Proxy headers received by PHP:

```php
<?php
$fields = [
    'IS_PROXY',
    'PROXY_TYPE',
    'COUNTRY_CODE',
    'COUNTRY_NAME',
    'REGION',
    'CITY',
    'ISP',
    'DOMAIN',
    'USAGE_TYPE',
    'ASN',
    'AS',
    'LAST_SEEN',
    'THREAT',
    'PROVIDER',
    'FRAUD_SCORE',
];

header('Content-Type: text/plain; charset=utf-8');

foreach ($fields as $field) {
    $key = 'HTTP_X_IP2PROXY_' . $field;
    if (isset($_SERVER[$key])) {
        echo $key . ': ' . $_SERVER[$key] . "\n";
    }
}
```

The example assumes the default `X-IP2Proxy` prefix. If you change `header_prefix`, update the PHP variable names accordingly. Values will depend on the IP address, BIN database fields or API subscription.

## Enrichment Headers

The middleware defines the following **15 request headers** in both lookup modes. Each uses the default `X-IP2Proxy` prefix.

| Header | What it represents |
|---|---|
| `X-IP2Proxy-Is-Proxy` | General proxy indicator (`True` or `False`) |
| `X-IP2Proxy-Proxy-Type` | Proxy or network classification |
| `X-IP2Proxy-Country-Code` | Two-letter country code |
| `X-IP2Proxy-Country-Name` | Country name |
| `X-IP2Proxy-Region` | Region or state |
| `X-IP2Proxy-City` | City |
| `X-IP2Proxy-ISP` | Internet service provider |
| `X-IP2Proxy-Domain` | Associated domain |
| `X-IP2Proxy-Usage-Type` | Network usage classification |
| `X-IP2Proxy-ASN` | Autonomous System Number |
| `X-IP2Proxy-AS` | Autonomous system name |
| `X-IP2Proxy-Last-Seen` | Last-seen indicator for proxy-related activity |
| `X-IP2Proxy-Threat` | Threat classification |
| `X-IP2Proxy-Provider` | Proxy provider information |
| `X-IP2Proxy-Fraud-Score` | IP fraud risk score |

**Note:** The code sets these headers on successful lookups, but meaningful field values depend on the selected BIN product or API subscription. The local `Is-Proxy` decision uses the library's classification logic, so do not treat it as interchangeable with every proxy type or as definitive proof of malicious activity.

## Understanding ASN Lookup

An **Autonomous System Number (ASN)** identifies a network operated under a common routing policy. The middleware exposes `X-IP2Proxy-ASN` and `X-IP2Proxy-AS`, helping applications understand which network is associated with a visitor's IP address.

ASN data can support network troubleshooting, traffic grouping, hosting-provider identification and security investigations. For offline IP intelligence resources, visit [IP2Location LITE](https://www.ip2location.com/database/lite).

## Troubleshooting

| Issue | Suggested check |
|---|---|
| `unrecognized directive: ip2proxy` | Use the Caddy binary built with `xcaddy --with github.com/ip2location/ip2proxy-caddy`. |
| `bin_path is required when mode is local` | Set `bin_path` to a readable BIN file. |
| BIN file cannot be opened | Verify the path, filesystem permissions and database format. |
| `api_key is required when mode is remote` | Make sure `IP2LOCATION_API_KEY` is available to the Caddy process. |
| Headers do not appear in PHP | Verify middleware ordering, the request path, PHP-FPM configuration and whether the lookup succeeds. |
| Some header values are missing or not useful | Verify that your BIN database or API subscription contains the corresponding fields. |
| Client IP appears incorrect | Check the proxy forwarding headers received by Caddy. |

## Security and Responsible Use

- **Trust the client IP carefully.** The current implementation checks `CF-Connecting-IP`, `X-Real-IP`, `X-Forwarded-For`, and `Forwarded` before the connection IP. Restrict and sanitize forwarding headers at your trusted proxy; otherwise clients may spoof them.
- **Protect credentials.** Keep the API key in your deployment environment or secret manager, not in committed Caddyfiles.
- **Use signals as context.** Proxy classifications, fraud scores and geolocation are risk indicators, not definitive evidence of fraud or a user's precise physical location.
- **Handle failed lookups explicitly.** A failed lookup does not block the request in the current implementation. Do not assume missing or previously supplied headers establish a clean result.

## Related IP2Location Products

| Product | Purpose |
|---|---|
| [IP2Location.io](https://www.ip2location.io/) | Hosted IP geolocation, ASN, and network intelligence APIs. |
| [IP2Location LITE](https://www.ip2location.com/database/lite) | Free downloadable IP geolocation and ASN databases for offline use. |
| [IP2Location Database](https://www.ip2location.com/database/ip2location) | Commercial offline IP geolocation databases. |
| [IP2Proxy](https://www.ip2location.com/database/proxy) | IP proxy, VPN, TOR and anonymizer intelligence. |
| [IP2WHOIS](https://www.ip2whois.com/) | Domain WHOIS lookup services. |

## Documentation and Support

- [IP2Location.io API documentation](https://www.ip2location.io/ip2location-documentation)
- [API plans and pricing](https://www.ip2location.io/pricing)
- [Report issues or request features](https://github.com/ip2location/ip2proxy-caddy/issues)


