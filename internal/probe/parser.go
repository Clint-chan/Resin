package probe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/netip"
	"strings"
)

func ParseEgressProbe(body []byte, format string) (netip.Addr, *string, error) {
	switch strings.TrimSpace(format) {
	case "", "cloudflare_trace":
		return ParseCloudflareTrace(body)
	case "plain_ip":
		return parsePlainIP(body)
	case "json_ip":
		return parseJSONIP(body)
	default:
		return netip.Addr{}, nil, fmt.Errorf("probe: unsupported egress probe format %q", format)
	}
}

func parsePlainIP(body []byte) (netip.Addr, *string, error) {
	raw := strings.TrimSpace(string(body))
	if raw == "" {
		return netip.Addr{}, nil, fmt.Errorf("probe: empty plain_ip response")
	}
	ip, err := netip.ParseAddr(raw)
	if err != nil {
		return netip.Addr{}, nil, err
	}
	return ip, nil, nil
}

func parseJSONIP(body []byte) (netip.Addr, *string, error) {
	var payload map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(body), &payload); err != nil {
		return netip.Addr{}, nil, err
	}
	for _, key := range []string{"ip", "origin", "query", "address"} {
		if raw, ok := payload[key]; ok {
			value, ok := raw.(string)
			if !ok {
				continue
			}
			value = strings.TrimSpace(value)
			if value == "" {
				continue
			}
			if strings.Contains(value, ",") {
				value = strings.TrimSpace(strings.Split(value, ",")[0])
			}
			ip, err := netip.ParseAddr(value)
			if err != nil {
				return netip.Addr{}, nil, err
			}
			return ip, nil, nil
		}
	}
	return netip.Addr{}, nil, fmt.Errorf("probe: ip field not found in json response")
}
