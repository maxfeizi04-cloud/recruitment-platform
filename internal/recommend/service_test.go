package recommend

import (
	"testing"
)

func TestJaccardSimilarity_Identical(t *testing.T) {
	a := []string{"Go", "React", "TypeScript"}
	b := []string{"Go", "React", "TypeScript"}
	result := jaccardSimilarity(a, b)
	if result != 1.0 {
		t.Errorf("expected 1.0, got %.2f", result)
	}
}

func TestJaccardSimilarity_HalfMatch(t *testing.T) {
	a := []string{"Go", "React", "TypeScript", "Docker"}
	b := []string{"Go", "React"}
	result := jaccardSimilarity(a, b)
	if result != 0.5 {
		t.Errorf("expected 0.5, got %.2f", result)
	}
}

func TestJaccardSimilarity_NoMatch(t *testing.T) {
	a := []string{"Go", "React"}
	b := []string{"Python", "Django"}
	result := jaccardSimilarity(a, b)
	if result != 0.0 {
		t.Errorf("expected 0.0, got %.2f", result)
	}
}

func TestJaccardSimilarity_OneEmpty(t *testing.T) {
	a := []string{"Go", "React"}
	b := []string{}
	result := jaccardSimilarity(a, b)
	if result != 0.0 {
		t.Errorf("expected 0.0 for empty set, got %.2f", result)
	}
}

func TestParseSkills_ObjectFormat(t *testing.T) {
	json := `{"skills":["Go","React","Docker"]}`
	skills := parseSkills(json)
	if len(skills) != 3 {
		t.Fatalf("expected 3 skills, got %d", len(skills))
	}
	if skills[0] != "Go" || skills[1] != "React" || skills[2] != "Docker" {
		t.Errorf("unexpected skills: %v", skills)
	}
}

func TestParseSkills_ArrayFormat(t *testing.T) {
	json := `["Java","Python","Rust"]`
	skills := parseSkills(json)
	if len(skills) != 3 {
		t.Fatalf("expected 3 skills, got %d", len(skills))
	}
}

func TestParseSkills_InvalidJSON(t *testing.T) {
	json := `not json`
	skills := parseSkills(json)
	if skills != nil {
		t.Errorf("expected nil for invalid JSON, got %v", skills)
	}
}

func TestParseSkills_NoSkillsField(t *testing.T) {
	json := `{"name":"test","age":30}`
	skills := parseSkills(json)
	if skills != nil {
		t.Errorf("expected nil for missing skills field, got %v", skills)
	}
}

func TestTimeDecay_Recent(t *testing.T) {
	result := timeDecay("2026-06-20T00:00:00Z") // 2 days ago (test date is June 22 2026)
	if result != 1.0 {
		t.Errorf("expected 1.0 for recent job, got %.2f", result)
	}
}

func TestTimeDecay_Old(t *testing.T) {
	result := timeDecay("2025-06-01T00:00:00Z") // over 1 year ago
	if result != 0.2 {
		t.Errorf("expected 0.2 for old job, got %.2f", result)
	}
}

func TestMaskPhone(t *testing.T) {
	result := maskPhone("13812345678")
	if result != "138****5678" {
		t.Errorf("expected 138****5678, got %s", result)
	}
}
