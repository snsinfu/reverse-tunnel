package ports

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestParseNetPort_ParsesGoodInput(t *testing.T) {
	np, err := ParseNetPort("443/tcp")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if np.Protocol != "tcp" {
		t.Errorf("unexpected protocol: got %q, want %q", np.Protocol, "tcp")
	}

	if np.Port != 443 {
		t.Errorf("unexpected port number: got %d, want %d", np.Port, 443)
	}
}

func TestNetPort_String_IsInServiceForm(t *testing.T) {
	np := NetPort{
		Protocol: "udp",
		Port:     123,
	}

	expected := "123/udp"
	actual := np.String()

	if actual != expected {
		t.Errorf("unexpected string result: got %q, expected %q", actual, expected)
	}
}

func TestNetPort_IsYAMLMarshalable(t *testing.T) {
	data := map[string]NetPort{
		"ssh": NetPort{
			Protocol: "tcp",
			Port:     22,
		},
	}
	expected := "ssh: 22/tcp\n"

	doc, err := yaml.Marshal(data)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if string(doc) != expected {
		t.Errorf("unexpected document: got %q, expected %q", doc, expected)
	}
}

func TestNetPort_IsYAMLUnmarshalable(t *testing.T) {
	doc := []byte("[80/tcp, 123/udp]\n")

	expected := []NetPort{
		NetPort{Protocol: "tcp", Port: 80},
		NetPort{Protocol: "udp", Port: 123},
	}

	actual := []NetPort{}
	if err := yaml.Unmarshal(doc, &actual); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("unexpected result: got %v, expected %v", actual, expected)
	}
}
