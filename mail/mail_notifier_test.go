package mail

import "testing"

func TestBuildEnvelop_RequiredFields(t *testing.T) {
	env := buildEnvelop(Data{
		To:      "to@example.com",
		Subject: "Hello",
		Body:    "<p>Hi</p>",
	})

	if len(env.To) != 1 || env.To[0] != "to@example.com" {
		t.Fatalf("unexpected To: %v", env.To)
	}
	if env.Subject != "Hello" {
		t.Fatalf("unexpected Subject: %q", env.Subject)
	}
	if env.HTML != "<p>Hi</p>" {
		t.Fatalf("unexpected HTML: %q", env.HTML)
	}
	// Optional fields must stay empty when not provided.
	if len(env.Cc) != 0 || len(env.Bcc) != 0 || len(env.ReplyTo) != 0 {
		t.Fatalf("expected optional fields empty, got Cc=%v Bcc=%v ReplyTo=%v", env.Cc, env.Bcc, env.ReplyTo)
	}
}

func TestBuildEnvelop_OptionalFields(t *testing.T) {
	env := buildEnvelop(Data{
		To:      "to@example.com",
		Cc:      "cc@example.com",
		Bcc:     "bcc@example.com",
		ReplyTo: "reply@example.com",
		Subject: "Hello",
		Body:    "body",
	})

	if len(env.Cc) != 1 || env.Cc[0] != "cc@example.com" {
		t.Fatalf("unexpected Cc: %v", env.Cc)
	}
	if len(env.Bcc) != 1 || env.Bcc[0] != "bcc@example.com" {
		t.Fatalf("unexpected Bcc: %v", env.Bcc)
	}
	if len(env.ReplyTo) != 1 || env.ReplyTo[0] != "reply@example.com" {
		t.Fatalf("unexpected ReplyTo: %v", env.ReplyTo)
	}
}
