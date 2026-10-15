package tmpl

import (
	"os"
	"strings"
	"testing"
)

func TestDockerRegistryComposeMountsStorageRootdir(t *testing.T) {
	raw, err := os.ReadFile("../../../builtin/core/roles/image-registry/docker-registry/templates/docker-compose.yaml")
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name    string
		rootdir string
		want    bool
	}{
		{"filesystem rootdir is bind mounted", "/var/lib/docker-registry/data", true},
		{"no filesystem rootdir -> no data mount", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := map[string]any{
				"image_registry": map[string]any{
					"docker_registry_version": "2.8.3",
					"auth":                    map[string]any{"registry": "dockerhub.kubekey.local"},
					"docker_registry": map[string]any{"storage": map[string]any{
						"filesystem": map[string]any{"rootdir": tc.rootdir},
					}},
				},
				"groups":   map[string]any{},
				"hostvars": map[string]any{},
			}
			b, err := Parse(v, string(raw))
			if err != nil {
				t.Fatalf("render: %v", err)
			}
			out := string(b)
			mount := "source: " + tc.rootdir + "\n        target: " + tc.rootdir
			if got := strings.Contains(out, mount); got != tc.want {
				t.Fatalf("data mount present=%v, want %v\n%s", got, tc.want, out)
			}
		})
	}
}
