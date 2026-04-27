package tfstate

import (
	"testing"
)

func scheduleResource(schedule string) Resource {
	attrs := map[string]interface{}{}
	if schedule != "" {
		attrs["schedule"] = schedule
	}
	return Resource{Type: "aws_scheduler_schedule", Attributes: attrs}
}

func TestFilterBySchedule_Match(t *testing.T) {
	resources := []Resource{
		scheduleResource("rate(5 minutes)"),
		scheduleResource("cron(0 12 * * ? *)"),
		scheduleResource(""),
	}
	got := FilterBySchedule(resources, "rate(5 minutes)", DefaultScheduleFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
}

func TestFilterBySchedule_EmptySchedule_ReturnsAll(t *testing.T) {
	resources := []Resource{scheduleResource("rate(1 hour)"), scheduleResource("cron(0 0 * * ? *)")}
	got := FilterBySchedule(resources, "", DefaultScheduleFilterOptions())
	if len(got) != len(resources) {
		t.Fatalf("expected %d results, got %d", len(resources), len(got))
	}
}

func TestFilterBySchedule_CaseInsensitive(t *testing.T) {
	resources := []Resource{scheduleResource("Rate(5 Minutes)")}
	got := FilterBySchedule(resources, "rate(5 minutes)", DefaultScheduleFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
}

func TestFilterBySchedule_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{scheduleResource("Rate(5 Minutes)")}
	got := FilterBySchedule(resources, "rate(5 minutes)", FilterOptions{CaseInsensitive: false})
	if len(got) != 0 {
		t.Fatalf("expected 0 results, got %d", len(got))
	}
}

func TestFilterBySchedules_ORSemantics(t *testing.T) {
	resources := []Resource{
		scheduleResource("rate(5 minutes)"),
		scheduleResource("cron(0 12 * * ? *)"),
		scheduleResource("rate(1 hour)"),
	}
	got := FilterBySchedules(resources, []string{"rate(5 minutes)", "rate(1 hour)"}, DefaultScheduleFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
}

func TestBuildScheduleIndex_Lookup(t *testing.T) {
	resources := []Resource{
		scheduleResource("rate(5 minutes)"),
		scheduleResource("cron(0 12 * * ? *)"),
	}
	idx := BuildScheduleIndex(resources)
	found := idx.Lookup("rate(5 minutes)")
	if len(found) != 1 {
		t.Fatalf("expected 1, got %d", len(found))
	}
}

func TestBuildScheduleIndex_LookupMissing(t *testing.T) {
	idx := BuildScheduleIndex([]Resource{scheduleResource("rate(1 hour)")})
	if idx.Lookup("rate(99 years)") != nil {
		t.Fatal("expected nil for missing schedule")
	}
}

func TestBuildScheduleIndex_Schedules(t *testing.T) {
	resources := []Resource{
		scheduleResource("rate(5 minutes)"),
		scheduleResource("cron(0 12 * * ? *)"),
	}
	idx := BuildScheduleIndex(resources)
	schedules := idx.Schedules()
	if len(schedules) != 2 {
		t.Fatalf("expected 2 schedules, got %d", len(schedules))
	}
}
