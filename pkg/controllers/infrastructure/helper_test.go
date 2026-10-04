package infrastructure

import (
	"reflect"
	"testing"

	capkkinfrav1beta1 "github.com/kubesphere/kubekey/api/capkk/infrastructure/v1beta1"
	kkcorev1 "github.com/kubesphere/kubekey/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kubesphere/kubekey/v4/pkg/converter"
)

func TestInventoryHostsEqual(t *testing.T) {
	host := func(raw string) kkcorev1.InventoryHost {
		return kkcorev1.InventoryHost{"h1": runtime.RawExtension{Raw: []byte(raw)}}
	}
	testcases := []struct {
		name string
		a, b kkcorev1.InventoryHost
		want bool
	}{
		{name: "identical", a: host(`{"a":1}`), b: host(`{"a":1}`), want: true},
		{name: "key order and whitespace", a: host(`{"a":1,"connector":{"host":"x","type":"ssh"}}`), b: host(`{ "connector": {"type":"ssh", "host":"x"}, "a": 1 }`), want: true},
		{name: "different value", a: host(`{"a":1}`), b: host(`{"a":2}`), want: false},
		{name: "different host name", a: host(`{"a":1}`), b: kkcorev1.InventoryHost{"h2": {Raw: []byte(`{"a":1}`)}}, want: false},
		{name: "different length", a: host(`{"a":1}`), b: kkcorev1.InventoryHost{}, want: false},
		{name: "identical invalid json", a: host(`{`), b: host(`{`), want: true},
		{name: "different invalid json", a: host(`{`), b: host(`{"a":1}`), want: false},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if got := inventoryHostsEqual(tc.a, tc.b); got != tc.want {
				t.Fatalf("inventoryHostsEqual() = %v, want %v", got, tc.want)
			}
		})
	}
}

// The converter marshals the connector struct in field order; the stored payload below simulates
// an API server response with sorted keys.
func TestInventoryHostsEqualConverterRoundTrip(t *testing.T) {
	kkcluster := &capkkinfrav1beta1.KKCluster{Spec: capkkinfrav1beta1.KKClusterSpec{
		InventoryHosts: []capkkinfrav1beta1.InventoryHost{{
			Name:      "h1",
			Connector: capkkinfrav1beta1.InventoryHostConnector{Type: "ssh", Host: "172.30.0.11", User: "k8s", Password: "pw"},
		}},
	}}
	desired, err := converter.ConvertKKClusterToInventoryHost(kkcluster)
	if err != nil {
		t.Fatal(err)
	}
	stored := kkcorev1.InventoryHost{"h1": {Raw: []byte(`{"connector":{"host":"172.30.0.11","password":"pw","type":"ssh","user":"k8s"}}`)}}
	if reflect.DeepEqual(stored, desired) {
		t.Fatal("expected raw bytes to differ")
	}
	if !inventoryHostsEqual(stored, desired) {
		t.Fatal("expected semantically equal hosts")
	}
}
