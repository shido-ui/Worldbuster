# Events & Notifications

Phase 12 establishes the event/notification backbone.

Events are immutable facts emitted by authoritative domain services. Notifications are player-facing projections that may reference an event.

This allows future systems to consume the same facts for:
- activity feeds
- notifications
- achievements
- analytics
- world news
- bot memory
- missions
- economy history
- moderation/audit streams

AI will not create authoritative events. It may later turn validated events into readable summaries.
