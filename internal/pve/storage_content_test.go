package pve

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestStorageContentItemCTime(t *testing.T) {
	body := []byte(`{"data":[
		{"volid":"local:iso/a.iso","content":"iso","ctime":1781588858},
		{"volid":"local-lvm:base-8109-disk-0","content":"images","format":"raw","size":1,"vmid":8109,"ctime":"1794483574"},
		{"volid":"local:iso/missing.iso","content":"iso"}
	]}`)
	var out struct {
		Data []StorageContentItem `json:"data"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out.Data) != 3 {
		t.Fatalf("items: got %d", len(out.Data))
	}
	if got := int64(out.Data[0].CTime); got != 1781588858 {
		t.Fatalf("directory ctime: got %d", got)
	}
	if got := int64(out.Data[1].CTime); got != 1794483574 {
		t.Fatalf("lvm ctime: got %d", got)
	}
	if got := int64(out.Data[2].CTime); got != 0 {
		t.Fatalf("missing ctime: got %d", got)
	}
	encoded, err := json.Marshal(out.Data[1])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(encoded), `"ctime":1794483574`) {
		t.Fatalf("marshal: %s", encoded)
	}

	var bad struct {
		Data []StorageContentItem `json:"data"`
	}
	err = json.Unmarshal([]byte(`{"data":[{"volid":"local-lvm:x","content":"images","ctime":"not-a-time"}]}`), &bad)
	if err == nil {
		t.Fatal("expected a non-numeric ctime to fail")
	}
}
