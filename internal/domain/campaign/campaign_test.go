package campaign

import "testing"

func TestNewCampaign(t *testing.T) {
	name := "Campaign Name"
	content := "Campaign Content"
	contacts := []string{"email@example.com", "email@example.com", "email@example.com"}

	c := NewCampaign(name, content, contacts)

	if c == nil {
		t.Error("Campaign is nil")
	} else if c.Name != name {
		t.Errorf("NewCampaign() Name = %v, want %v", c.Name, name)
	} else if c.Content != content {
		t.Errorf("NewCampaign() Content = %v, want %v", c.Content, content)
	} else if len(c.Contacts) != len(contacts) {
		t.Errorf("NewCampaign() Contacts = %v, want %v", c.Contacts, contacts)
	}
}
