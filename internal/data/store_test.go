package data

import "testing"

func TestStoreLookupSingle(t *testing.T) {
	s := NewStore()
	s.LoadRTP([]string{"021000021", "031100209"})
	s.LoadFedNow([]string{"021000021", "071000013"})
	s.LoadNames(map[string]string{
		"021000021": "JPMORGAN CHASE",
		"031100209": "HSBC BANK USA",
		"071000013": "FEDERAL RESERVE BANK",
	})

	r := s.Lookup("021000021")
	if !r.RTP || !r.FedNow {
		t.Errorf("021000021: RTP=%v FedNow=%v, want both true", r.RTP, r.FedNow)
	}
	if r.Name != "JPMORGAN CHASE" {
		t.Errorf("Name = %q, want JPMORGAN CHASE", r.Name)
	}

	r = s.Lookup("031100209")
	if !r.RTP || r.FedNow {
		t.Errorf("031100209: RTP=%v FedNow=%v, want RTP=true FedNow=false", r.RTP, r.FedNow)
	}

	r = s.Lookup("071000013")
	if r.RTP || !r.FedNow {
		t.Errorf("071000013: RTP=%v FedNow=%v, want RTP=false FedNow=true", r.RTP, r.FedNow)
	}

	r = s.Lookup("999999999")
	if r.RTP || r.FedNow {
		t.Errorf("999999999: RTP=%v FedNow=%v, want both false", r.RTP, r.FedNow)
	}
}

func TestStoreAll(t *testing.T) {
	s := NewStore()
	s.LoadRTP([]string{"021000021", "031100209"})
	s.LoadFedNow([]string{"021000021", "071000013"})
	s.LoadNames(map[string]string{
		"021000021": "JPMORGAN CHASE",
		"031100209": "HSBC BANK USA",
		"071000013": "FEDERAL RESERVE BANK",
	})

	all := s.All()
	if len(all) != 3 {
		t.Fatalf("All() returned %d entries, want 3", len(all))
	}
}
