package probe

import (
	"net/netip"
	"testing"
)

func TestParseEgressProbe_CloudflareTrace(t *testing.T) {
	ip, loc, err := ParseEgressProbe([]byte("fl=29f\nip=203.0.113.10\nloc=US\n"), "cloudflare_trace")
	if err != nil {
		t.Fatalf("ParseEgressProbe: %v", err)
	}
	if ip != netip.MustParseAddr("203.0.113.10") {
		t.Fatalf("ip = %v, want 203.0.113.10", ip)
	}
	if loc == nil || *loc != "us" {
		t.Fatalf("loc = %v, want us", loc)
	}
}

func TestParseEgressProbe_PlainIP(t *testing.T) {
	ip, loc, err := ParseEgressProbe([]byte("198.51.100.5\n"), "plain_ip")
	if err != nil {
		t.Fatalf("ParseEgressProbe: %v", err)
	}
	if ip != netip.MustParseAddr("198.51.100.5") {
		t.Fatalf("ip = %v, want 198.51.100.5", ip)
	}
	if loc != nil {
		t.Fatalf("loc = %v, want nil", loc)
	}
}

func TestParseEgressProbe_JSONIP(t *testing.T) {
	ip, loc, err := ParseEgressProbe([]byte(`{"origin":"198.51.100.9, 198.51.100.10"}`), "json_ip")
	if err != nil {
		t.Fatalf("ParseEgressProbe: %v", err)
	}
	if ip != netip.MustParseAddr("198.51.100.9") {
		t.Fatalf("ip = %v, want 198.51.100.9", ip)
	}
	if loc != nil {
		t.Fatalf("loc = %v, want nil", loc)
	}
}

func TestParseEgressProbe_UnsupportedFormat(t *testing.T) {
	_, _, err := ParseEgressProbe([]byte("198.51.100.5"), "weird")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}
