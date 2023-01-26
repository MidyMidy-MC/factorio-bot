package rconcontroller

type EventType string

const (
	EventEmpty             EventType = "empty"
	EventPlayerLeft        EventType = "player-left"
	EventPlayerJoin        EventType = "player-joined"
	EventPlayerKilled      EventType = "oops"
	EventConsoleChat       EventType = "console-chat"
	EventConsoleMe         EventType = "console-me"
	EventConsolePin        EventType = "console-pin"
	EventResourceDepleted  EventType = "resource-depleted"
	EventResearchCancelled EventType = "research-cancelled"
	EventResearchFinished  EventType = "research-finished"
	EventResearchReversed  EventType = "research-reversed"
	EventResearchStarted   EventType = "research-started"
	EventInvalid           EventType = "invalid"
)

type Munch map[string]any

func (m Munch) GetString(key string) string {
	t, ok := m[key].(string)
	if ok {
		return t
	}
	return ""
}

func (m Munch) GetInt(key string) int {
	t, ok := m[key].(int)
	if ok {
		return t
	}
	return 0
}

func (m Munch) GetBool(key string) bool {
	t, ok := m[key].(bool)
	if ok {
		return t
	}
	return false
}

func (m Munch) GetMunch(key string) Munch {
	t, ok := m[key].(Munch)
	if ok {
		return t
	}
	return nil
}

type Event struct {
	Munch
}

func (e *Event) EventType() EventType {
	t := e.GetString("type")
	if t == "" {
		return EventInvalid
	}
	return EventType(t)
}
